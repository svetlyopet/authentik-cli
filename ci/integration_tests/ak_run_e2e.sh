#!/usr/bin/env bash

set -eo pipefail

YELLOW="\033[33m"
GREEN="\033[32m"
RED="\033[31m"
ENDCOLOR="\033[0m"

GO_BIN=$(which go)
GO_ARCH=$(uname -m)

CI_TEST_DIR=$(pwd)/ci/integration_tests
CI_ENV_FILE="$CI_TEST_DIR/.env"

AK_CLI_PATH=$(pwd)/bin
AK_CLI_BIN="$AK_CLI_PATH/ak"

TEST_TENANT_NAME="example-tenant"
TEST_USER_NAME="example-user"
TEST_GROUP_NAME="example-group"
TEST_OIDC_APP_NAME="example-app"

build_ak_bin() {
  echo -n "Building binary... "
  export GOARCH=$GO_ARCH
  $GO_BIN build -o bin/ak .
  echo -e "${GREEN}DONE${ENDCOLOR}"
  echo -e "Built binary ${YELLOW}$AK_CLI_BIN${ENDCOLOR}"
}

create_config() {
  echo -n "Creating config... "
  source $CI_ENV_FILE
  $AK_CLI_BIN config 1&>/dev/null <<STDIN
$AUTHENTIK_URL
$AUTHENTIK_BOOTSTRAP_TOKEN
STDIN
  echo -e "${GREEN}DONE${ENDCOLOR}"
  echo -e "Created config file in ${YELLOW}$HOME/.authentik-cli${ENDCOLOR}"
  echo
}

create_tenant() {
  TEST_NAME="Create tenant"
  test_start "$TEST_NAME"
  $AK_CLI_BIN create tenant $TEST_TENANT_NAME
  test_passed "$TEST_NAME"
}

create_admin_user_for_tenant() {
  TEST_NAME="Create admin user for tenant"
  test_start "$TEST_NAME"
  $AK_CLI_BIN create user $TEST_USER_NAME \
  --name=$TEST_USER_NAME \
  --email=$TEST_USER_NAME@example.com \
  --tenant-admin=$TEST_TENANT_NAME
  test_passed "$TEST_NAME"
}

create_group() {
  TEST_NAME="Create group"
  test_start "$TEST_NAME"
  $AK_CLI_BIN create group $TEST_GROUP_NAME
  test_passed "$TEST_NAME"
}

delete_group() {
  TEST_NAME="Delete group"
  test_start "$TEST_NAME"
  $AK_CLI_BIN delete group $TEST_GROUP_NAME
  test_passed "$TEST_NAME"
}

delete_user() {
  TEST_NAME="Delete user"
  test_start "$TEST_NAME"
  $AK_CLI_BIN delete user $TEST_USER_NAME
  test_passed "$TEST_NAME"
}

delete_tenant() {
  TEST_NAME="Delete tenant"
  test_start "$TEST_NAME"
  $AK_CLI_BIN delete tenant $TEST_TENANT_NAME
  test_passed "$TEST_NAME"
}


cleanup() {
  echo "Cleaning up..."
  echo -n "rm " && rm -v $HOME/.authentik-cli
}

test_start() {
  echo "Test case: $1"
}

test_passed() {
  echo -e "$1 ${GREEN}PASSED${ENDCOLOR}\n"
}

main() {
  build_ak_bin
  create_config
  create_tenant
  create_admin_user_for_tenant
  create_group
  delete_group
  delete_user
  delete_tenant
  cleanup
}

main "$@"