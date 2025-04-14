## Setup local Authentik
Creates a Postgres, Redis and Authentik deployments with Docker compose.

```bash
./ci/integration_tests/ak_bootstrap.sh create
```

You can use an environment variable to specify the Authentik version
```bash
AUTHENTIK_TAG=2025.2.3 ./ci/integration_tests/ak_bootstrap.sh create
```

## Remove local Authentik
Removes the Docker resources created by this script.

```bash
./ci/integration_tests/ak_bootstrap.sh destroy
```

## Run e2e tests
Runs the e2e tests against the local Authentik instance

```bash
./ci/integration_tests/ak_run_e2e.sh
```