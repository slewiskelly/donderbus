# Donderbus

Assigns a pull request to a random set of individuals from a GitHub team.

> [!NOTE]
> Where possible, it is recommended to instead use GitHub's [auto-assignment][about-auto-assignment]
> feature.

> [!WARNING]
> ___This is a work in progress and is considered experimental.___

## Overview

### Assignment

Pull requests are assigned to individuals from currently assigned teams.

Individuals are selected at random.

## Installation

### Docker

```shell
docker pull ghcr.io/slewiskelly/donderbus:latest
```

### Go

```shell
go install github.com/slewiskelly/donderbus/cmd/donderbus@latest
```

### Homebrew

```shell
brew install slewiskelly/tap/donderbus
```

## Usage

> [!TIP]
> For full usage and examples use `donderbus help` or `donderbus <command> --help`.

### Assignment

To assign individuals to a pull request:

> [!NOTE]
> The pull request must already be assigned to one or more teams.

```shell
donderbus assign [flags] <url>
```

[about-auto-assignment]: https://docs.github.com/en/organizations/organizing-members-into-teams/managing-code-review-settings-for-your-team#about-auto-assignment
