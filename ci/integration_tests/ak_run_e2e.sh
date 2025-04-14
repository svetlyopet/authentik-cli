#!/usr/bin/env bash

set -e +o pipefail

GO_BIN=$(which go)
GO_ARCH=$(uname -m)

CI_TEST_DIR=$(pwd)/ci/integration_tests
CI_ENV_FILE="$CI_TEST_DIR/.env"

AK_CLI_PATH=$(pwd)/bin
AK_CLI_BIN="$AK_CLI_PATH/ak"

TEST_TENANT_NAME="example-tenant"
TEST_USER_NAME="example-user"
TEST_OIDC_APP_NAME="example-app"

build_ak_bin() {
  export GOARCH=$GO_ARCH
  $GO_BIN build -o bin/ak .
}

set_config() {
  source $CI_ENV_FILE
  $AK_CLI_BIN config 1&>/dev/null <<STDIN
$AUTHENTIK_URL
$AUTHENTIK_BOOTSTRAP_TOKEN
STDIN
}

create_tenant() {
  echo "Test case: Create tenant"
  $AK_CLI_BIN create tenant $TEST_TENANT_NAME
  echo
}

create_admin_user_for_tenant() {
  echo "Test case: Create admin user for tenant"
  $AK_CLI_BIN create user $TEST_USER_NAME \
  --name=$TEST_USER_NAME \
  --email=$TEST_USER_NAME@example.com \
  --tenant-admin=$TEST_TENANT_NAME
  echo
}

delete_user() {
  echo "Test case: Delete user"
  $AK_CLI_BIN delete user $TEST_USER_NAME
  echo
}

delete_tenant() {
  echo "Test case: Delete tenant"
  $AK_CLI_BIN delete tenant $TEST_TENANT_NAME
  echo
}

main() {
  build_ak_bin
  set_config
  create_tenant
  create_admin_user_for_tenant
  delete_user
  delete_tenant
}

main "$@"