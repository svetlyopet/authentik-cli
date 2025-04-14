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

source $CI_ENV_FILE

build_ak_bin() {
  echo -n "Building binary... "
  export GOARCH=$GO_ARCH
  $GO_BIN build -o bin/ak .
  echo -e "${GREEN}DONE${ENDCOLOR}"
  echo -e "Built binary ${YELLOW}$AK_CLI_BIN${ENDCOLOR}"
}

create_config() {
  echo -n "Creating config... "
  $AK_CLI_BIN config 1&>/dev/null <<STDIN
$AUTHENTIK_URL
$AUTHENTIK_BOOTSTRAP_TOKEN
STDIN
  echo -e "${GREEN}DONE${ENDCOLOR}"
  echo -e "Created config file in ${YELLOW}$HOME/.authentik-cli${ENDCOLOR}"
  echo
}

create_tenant() {
  local test_name="Create tenant"
  test_start "$test_name"
  $AK_CLI_BIN create tenant $TEST_TENANT_NAME
  test_passed "$test_name"
}

create_admin_user_for_tenant() {
  local test_name="Create admin user for tenant"
  test_start "$test_name"
  $AK_CLI_BIN create user $TEST_USER_NAME \
  --name=$TEST_USER_NAME \
  --email=$TEST_USER_NAME@example.com \
  --tenant-admin=$TEST_TENANT_NAME
  test_passed "$test_name"
}

create_group() {
  local test_name="Create group"
  test_start "$test_name"
  $AK_CLI_BIN create group $TEST_GROUP_NAME
  test_passed "$test_name"
}

delete_group() {
  local test_name="Delete group"
  test_start "$test_name"
  $AK_CLI_BIN delete group $TEST_GROUP_NAME
  test_passed "$test_name"
}

delete_user() {
  local test_name="Delete user"
  test_start "$test_name"
  $AK_CLI_BIN delete user $TEST_USER_NAME
  test_passed "$test_name"
}

delete_tenant() {
  local test_name="Delete tenant"
  test_start "$test_name"
  $AK_CLI_BIN delete tenant $TEST_TENANT_NAME
  test_passed "$test_name"
}


cleanup() {
  echo "Cleaning up..."
  if [ -f $HOME/.authentik-cli  ]; then
    echo -n "rm " && rm -v $HOME/.authentik-cli
  fi
}
trap cleanup EXIT

test_start() {
  echo "Test case: $1"
}

test_passed() {
  echo -e "$1 ${GREEN}PASSED${ENDCOLOR}\n"
}

check_authentik_status() {
  local max_retries=10
  local retry_count=0
  local retry_interval=2

  echo "Checking Authentik target health... "

  while [[ $retry_count -lt $max_retries ]]; do
    response_code=$(curl -s -o /dev/null -w "%{http_code}" "${AUTHENTIK_URL}/-/health/ready/")

    if [[ $response_code -eq 200 ]]; then
      echo -e "Authentik is ${GREEN}up and running${ENDCOLOR} at $AUTHENTIK_URL"
      echo
      return 0
    else
      echo -e "${YELLOW}Authentik is not ready. Retrying in $retry_interval seconds...${ENDCOLOR}"
      sleep $retry_interval
      retry_count=$((retry_count + 1))
    fi
  done

  echo -e "${RED}Authentik is not ready after $max_retries retries. Aborting...${ENDCOLOR}"
  return 1
}

main() {
  check_authentik_status
  build_ak_bin
  create_config
  create_tenant
  create_admin_user_for_tenant
  create_group
  delete_group
  delete_user
  delete_tenant
}

main "$@"