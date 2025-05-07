### Setup environment for e2e tests

See the [README for running Authentik locally](./ci/integration_tests/README.md)

### mock generation
Install https://github.com/uber-go/mock
```bash
go install go.uber.org/mock/mockgen@latest
```

Generate mocks the interfaces
```bash
mockgen -destination=mocks/ak/ak.go -package=mock_ak -source=internal/ak/ak.go
```

### Unit tests

Execute tests and save results to a coverage file(optional)
```bash
go test -v -coverprofile=coverage.out ./...
```

### Code coverage

Install the Cobertura code coverage tool
```bash
go get github.com/t-yuki/gocover-cobertura
```

Nice local UI view for looking at code coverage per package
```bash
go tool cover -html=coverage.out
```