<img align=right width="168" src="https://raw.githubusercontent.com/gouef/proxier/refs/heads/main/docs/gouef_logo.png">

# proxier
Lightweight proxy

[![Static Badge](https://img.shields.io/badge/Github-gouef%2Fproxier-blue?style=for-the-badge&logo=github&link=github.com%2Fgouef%2Fproxier)](https://github.com/gouef/proxier)

[![GoDoc](https://pkg.go.dev/badge/github.com/gouef/proxier.svg)](https://pkg.go.dev/github.com/gouef/proxier)
[![GitHub stars](https://img.shields.io/github/stars/gouef/proxier?style=social)](https://github.com/gouef/proxier/stargazers)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouef/proxier)](https://goreportcard.com/report/github.com/gouef/proxier)
[![codecov](https://codecov.io/github/gouef/proxier/branch/main/graph/badge.svg?token=YUG8EMH6Q8)](https://codecov.io/github/gouef/proxier)

## Versions
![Stable Version](https://img.shields.io/github/v/release/gouef/proxier?label=Stable&labelColor=green)
![GitHub Release](https://img.shields.io/github/v/release/gouef/proxier?label=RC&include_prereleases&filter=*rc*&logoSize=diago)
![GitHub Release](https://img.shields.io/github/v/release/gouef/proxier?label=Beta&include_prereleases&filter=*beta*&logoSize=diago)


## Configuration

```yaml
listen_http: ":80"
listen_https: ":443"

routes:
  - host: "proxier.gouef.local"
    path: "/"
    target: "http://host.docker.internal:8081"
  - host: "web-project.gouef.local"
    path: "/"
    target: "http://host.docker.internal:8081"
  - host: "gouef.local"
    path: "/"
    target: "http://host.docker.internal:8081"

  - host: "sub.example.com"
    path: "/"
    target: "http://localhost:4000"

# TLS is in testing
tls:
  use_lets_encrypt: false
  cert_file: ""
  key_file: ""
#  cache_dir: ".cache"
  cache_dir: ""
  email: ""
  hosts: []

```

## Run

```shell
docker run --rm -d -p 80:80 -p 443:443 \
-v ./config.yaml:/app/config.yaml \
--add-host host.docker.internal:host-gateway \
--name proxier gouef/proxier:latest
```

## Contributing

Read [Contributing](CONTRIBUTING.md)

## Contributors

<div>
<span>
  <a href="https://github.com/JanGalek"><img src="https://raw.githubusercontent.com/gouef/proxier/refs/heads/contributors-svg/.github/contributors/JanGalek.svg" alt="JanGalek" /></a>
</span>
<span>
  <a href="https://github.com/actions-user"><img src="https://raw.githubusercontent.com/gouef/proxier/refs/heads/contributors-svg/.github/contributors/actions-user.svg" alt="actions-user" /></a>
</span>
<span>
  <a href="https://github.com/apps/dependabot"><img src="https://raw.githubusercontent.com/gouef/proxier/refs/heads/contributors-svg/.github/contributors/dependabot[bot].svg" alt="dependabot[bot]" /></a>
</span>
</div>

## Join our Discord Community! 🎉

[![Discord](https://img.shields.io/discord/1334331501462163509?style=for-the-badge&logo=discord&logoColor=white&logoSize=auto&label=Community%20discord&labelColor=blue&link=https%3A%2F%2Fdiscord.gg%2FwjGqeWFnqK
)](https://discord.gg/wjGqeWFnqK)

Click above to join our community on Discord!
