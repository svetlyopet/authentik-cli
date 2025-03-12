## Setup local Authentik
Creates a Postgres, Redis and Authentik deployments with Docker compose.

```bash
./ci/integration-tests/ak_bootstrap.sh create
```

## Remove local Authentik
Removes the Docker resources created by this script.

```bash
./ci/integration-tests/ak_bootstrap.sh destroy
```