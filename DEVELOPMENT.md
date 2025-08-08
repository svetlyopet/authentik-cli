### Setup environment for e2e tests

See the [README for running Authentik locally](./ci/integration_tests/README.md)

### Setup pre-commit git hooks

Install [gitleaks](https://github.com/gitleaks/gitleaks) locally.

Setup the pre-commit hook
```bash
make git-hooks
```

### Mock generation

Generate mocks the interfaces
```bash
make mocks
```

### Run unit tests

Execute tests and save results to a coverage file(optional)
```bash
make test
```

### Code coverage

Install the Cobertura code coverage tool
```bash
make coverage
```

Nice local UI view for looking at code coverage per package
```bash
make html-coverage
```