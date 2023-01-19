# DocMerge

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
