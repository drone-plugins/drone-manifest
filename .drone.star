def main(ctx):
  before = testing(ctx)

  stages = [
    linux(ctx, 'amd64'),
    linux(ctx, 'arm64'),
  ]

  after = manifest(ctx)

  for b in before:
    for s in stages:
      s['depends_on'].append(b['name'])

  for s in stages:
    for a in after:
      a['depends_on'].append(s['name'])

  return before + stages + after

def testing(ctx):
  return [{
    'kind': 'pipeline',
    'type': 'kubernetes',
    'name': 'testing',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'staticcheck',
        'image': 'golang:1.17',
        'pull': 'always',
        'commands': [
          'go get honnef.co/go/tools/cmd/staticcheck',
          'go run honnef.co/go/tools/cmd/staticcheck ./...',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/go',
          },
        ],
      },
      {
        'name': 'lint',
        'image': 'golang:1.17',
        'pull': 'always',
        'commands': [
          'go get golang.org/x/lint/golint',
          'go run golang.org/x/lint/golint -set_exit_status ./...',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/go',
          },
        ],
      },
      {
        'name': 'vet',
        'image': 'golang:1.17',
        'pull': 'always',
        'commands': [
          'go vet ./...',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/go',
          },
        ],
      },
      {
        'name': 'test',
        'image': 'golang:1.17',
        'pull': 'always',
        'commands': [
          'go test -cover ./...',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/go',
          },
        ],
      },
    ],
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
        'refs/pull/**',
      ],
    },
  }]

def linux(ctx, arch):
  docker = {
    'dockerfile': 'docker/Dockerfile.linux.%s' % (arch),
    'repo': 'lemontech/drone-manifest-ecr',
    'username': {
      'from_secret': 'docker_username',
    },
    'password': {
      'from_secret': 'docker_password',
    },
  }

  if ctx.build.event == 'pull_request':
    docker.update({
      'dry_run': True,
      'tags': 'linux-%s' % (arch),
    })
  else:
    docker.update({
      'auto_tag': True,
      'auto_tag_suffix': 'linux-%s' % (arch),
    })

  if ctx.build.event == 'tag':
    build = [
      'go build -v -ldflags "-X main.version=%s" -a -tags netgo -o release/linux/%s/drone-manifest-ecr ./cmd/drone-manifest-ecr' % (ctx.build.ref.replace("refs/tags/v", ""), arch),
    ]
  else:
    build = [
      'go build -v -ldflags "-X main.version=%s" -a -tags netgo -o release/linux/%s/drone-manifest-ecr ./cmd/drone-manifest-ecr' % (ctx.build.commit[0:8], arch),
    ]

  return {
    'kind': 'pipeline',
    'type': 'kubernetes',
    'name': 'linux-%s' % (arch),
    'platform': {
      'os': 'linux',
      'arch': arch,
    },
    'steps': [
      {
        'name': 'environment',
        'image': 'golang:1.17',
        'pull': 'always',
        'environment': {
          'CGO_ENABLED': '0',
        },
        'commands': [
          'go version',
          'go env',
        ],
      },
      {
        'name': 'build',
        'image': 'golang:1.17',
        'pull': 'always',
        'environment': {
          'CGO_ENABLED': '0',
        },
        'commands': build,
      },
      {
        'name': 'executable',
        'image': 'golang:1.17',
        'pull': 'always',
        'commands': [
          './release/linux/%s/drone-manifest-ecr --help' % (arch),
        ],
      },
      {
        'name': 'docker',
        'image': 'plugins/docker',
        'pull': 'always',
        'settings': docker,
      },
    ],
    'depends_on': [],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
        'refs/pull/**',
      ],
    },
    'node_selector': {
      'kubernetes.io/arch': arch,
    },
    'tolerations': [
      {
        'key': 'node/arch',
        'operator': 'Equal',
        'value': arch,
        'effect': 'NoSchedule'
      },
    ],
  }

def manifest(ctx):
  return [{
    'kind': 'pipeline',
    'type': 'kubernetes',
    'name': 'manifest',
    'steps': [
      {
        'name': 'manifest',
        'image': 'plugins/manifest',
        'pull': 'always',
        'settings': {
          'auto_tag': 'true',
          'username': {
            'from_secret': 'docker_username',
          },
          'password': {
            'from_secret': 'docker_password',
          },
          'spec': 'docker/manifest.tmpl',
          'ignore_missing': 'true',
        },
      },
    ],
    'depends_on': [],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
      ],
    },
  }]
