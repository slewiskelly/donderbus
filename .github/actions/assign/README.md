# donerbus/assign

Assigns a pull request to a random set of individuals from a GitHub team.

## Usage

Specify this action in a GitHub Actions workflow:

```yaml
steps:
  - uses: slewiskelly/ock/.github/actions/assign@v0
    with:
      number: 123
      owner: acme
      repository: foo
      version: latest # Default.
```

### Configuration

| Input          | Required? | Default                       | Description                                                                                                              |
| -------------- | --------- | ----------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| `owner`        | Yes       |                               | Target pull request's owner                                                                                              |
| `repository`   | Yes       |                               | Target pull request's repository                                                                                         |
| `number`       | Yes       |                               | Target pull request's number                                                                                             |
| `version`      | No        | `latest`                      | Version of `donderbus` to be installed; can be either: a semantic version (`v0`, `v0.1.0`), branch (`main`), or `latest` |
| `github_token` | No        | `${{ secrets.GITHUB_TOKEN }}` | Token used to authenticate with GitHub                                                                                   |
