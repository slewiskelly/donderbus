# donderbus/setup-donderbus

Installs `donderbus` and adds it to the `PATH`.

## Usage

Specify this action in a GitHub Actions workflow:

```yaml
steps:
  - uses: slewiskelly/donderbus/.github/actions/setup-donderbus@v0
    with:
      version: latest # Default.
```

### Configuration

| Input     | Required? | Default  | Description                                                                                                              |
| --------- | --------- | -------- | ------------------------------------------------------------------------------------------------------------------------ |
| `version` | No        | `latest` | Version of `donderbus` to be installed; can be either: a semantic version (`v0`, `v0.1.0`), branch (`main`), or `latest` |
