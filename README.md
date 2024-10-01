# DocMerge

[![Releases](https://img.shields.io/github/v/release/stenic/docmerge?style=for-the-badge)](https://github.com/stenic/docmerge/releases)
[![Build status](https://img.shields.io/github/actions/workflow/status/stenic/docmerge/release.yaml?style=for-the-badge)](https://github.com/stenic/docmerge/actions/workflows/ci.yml)
[![GitHub License](https://img.shields.io/github/license/stenic/docmerge?style=for-the-badge)](https://github.com/stenic/docmerge/blob/main/LICENSE)
[![Powered by Stenic](https://img.shields.io/badge/powered--by-stenic.io-blue?style=for-the-badge&logoColor=blue)](https://stenic.io)

Crawl your repositories for `docs` folders and create a central documentation source.

## Installation

```shell
# homebrew
brew install stenic/tap/docmerge

# gofish
gofish rig add https://github.com/stenic/fish-food
gofish install github.com/stenic/fish-food/docmerge

# scoop
scoop bucket add docmerge https://github.com/stenic/scoop-bucket.git
scoop install docmerge

# go
go install github.com/stenic/docmerge@latest

# docker
docker pull ghcr.io/stenic/docmerge:latest

# dockerfile
COPY --from=ghcr.io/stenic/docmerge:latest /docmerge /usr/local/bin/
```

> For even more options, check the [releases page](https://github.com/stenic/docmerge/releases).
