#!/usr/bin/env bash

set -o pipefail

YELLOW="\033[33m"
GREEN="\033[32m"
RED="\033[31m"
ENDCOLOR="\033[0m"


ARCH=$(uname -m)
if [ $ARCH == "arm64" ]; then
  GO_ARCH="arm64"
elif [ $ARCH == "x86_64" ]; then
  GO_ARCH="amd64"
else
  echo "unsuported architecture $ARCH"
  exit 1
fi

GO_BIN=$(which go)

CI_TEST_DIR=$(pwd)/ci/integration_tests
CI_ENV_FILE="$CI_TEST_DIR/.env"

AK_CLI_PATH=$(pwd)/bin
AK_CLI_BIN="$AK_CLI_PATH/ak"

TEST_TENANT_NAME="example-tenant"
TEST_USER_NAME="example-user"
TEST_GROUP_NAME="example-group"
TEST_GROUP_RENAMED="example-group-renamed"
TEST_OIDC_APP_NAME="example-app-oidc"

source $CI_ENV_FILE

build_ak_bin() {
  echo -n "Building binary... "
  export GOARCH=$GO_ARCH
  $GO_BIN build -o bin/ak .
  if [ $? -ne 0 ]; then
    echo -e "Building binary ${RED}FAILED${ENDCOLOR}"
    exit 1
  fi
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
  test_result "$test_name" $?
}

create_admin_user_for_tenant() {
  local test_name="Create admin user for tenant"
  test_start "$test_name"
  $AK_CLI_BIN create user $TEST_USER_NAME \
  --name=$TEST_USER_NAME \
  --email=$TEST_USER_NAME@example.com \
  --tenant-admin=$TEST_TENANT_NAME
  test_result "$test_name" $?
}

create_group() {
  local test_name="Create group"
  test_start "$test_name"
  $AK_CLI_BIN create group $TEST_GROUP_NAME
  test_result "$test_name" $?
}

update_group() {
  local test_name="Update group"
  test_start "$test_name"
  $AK_CLI_BIN update group $TEST_GROUP_NAME --name=$TEST_GROUP_RENAMED
  test_result "$test_name" $?
}

delete_group() {
  local test_name="Delete group"
  test_start "$test_name"
  $AK_CLI_BIN delete group $TEST_GROUP_RENAMED
  test_result "$test_name" $?
}

update_user() {
  local test_name="Update user"
  test_start "$test_name"
  $AK_CLI_BIN update user $TEST_USER_NAME \
  --email=$TEST_USER_NAME-updated@example.com \
  --active=false
  test_result "$test_name" $?

  test_name="Update user keeps unchanged attributes"
  test_start "$test_name"
  $AK_CLI_BIN get user $TEST_USER_NAME | grep -q "\"name\": \"$TEST_USER_NAME\""
  test_result "$test_name" $?
}

delete_user() {
  local test_name="Delete user"
  test_start "$test_name"
  $AK_CLI_BIN delete user $TEST_USER_NAME
  test_result "$test_name" $?
}

delete_tenant() {
  local test_name="Delete tenant"
  test_start "$test_name"
  $AK_CLI_BIN delete tenant $TEST_TENANT_NAME
  test_result "$test_name" $?
}

create_app_oidc() {
  local test_name="Create OIDC app"
  test_start "$test_name"
  $AK_CLI_BIN create app $TEST_OIDC_APP_NAME \
  --provider-type=oidc \
  --oidc-client-type=confidential \
  --oidc-redirect-uri-strict=http://localhost:8000 \
  --oidc-redirect-uri-regex='http://*.local:9000'
  test_result "$test_name" $?
}

update_app_oidc() {
  local test_name="Update OIDC app"
  test_start "$test_name"
  $AK_CLI_BIN update app $TEST_OIDC_APP_NAME \
  --oidc-client-type=public \
  --oidc-redirect-uri-strict=http://localhost:8080
  test_result "$test_name" $?

  test_name="Update OIDC app keeps unchanged redirect URIs"
  test_start "$test_name"
  $AK_CLI_BIN get app $TEST_OIDC_APP_NAME | grep -q 'local:9000'
  test_result "$test_name" $?
}

delete_app() {
  local test_name="Delete app"
  test_start "$test_name"
  $AK_CLI_BIN delete app $TEST_OIDC_APP_NAME
  test_result "$test_name" $?
}


cleanup() {
  echo "Cleaning up..."
  if [ -f $HOME/.authentik-cli  ]; then
    echo "rm $HOME/.authentik-cli" && rm -f $HOME/.authentik-cli
  fi

  if [ -f $AK_CLI_BIN  ]; then
    echo "rm $AK_CLI_BIN" && rm -f $AK_CLI_BIN
  fi
}
trap cleanup EXIT

test_start() {
  echo "Test case: $1"
}

test_result() {
  if [ $2 -ne 0 ]; then
    echo -e "$1 ${RED}FAILED${ENDCOLOR}\n"
    exit 1
  fi
  echo -e "$1 ${GREEN}PASSED${ENDCOLOR}\n"
}

check_authentik_status() {
  local max_retries=30
  local retry_count=0
  local retry_interval=5

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
  exit 1
}

wait_for_default_flows() {
  # Since authentik 2025.8 the worker applies the default blueprints in the
  # background, so the API can be ready before the default flows (and the
  # bootstrap token) exist.
  local slugs=("default-authentication-flow" "default-provider-authorization-explicit-consent" "default-invalidation-flow")
  local max_retries=30
  local retry_interval=5

  echo "Waiting for the default flows... "

  for slug in "${slugs[@]}"; do
    local retry_count=0
    until curl -sf -H "Authorization: Bearer $AUTHENTIK_BOOTSTRAP_TOKEN" \
      "${AUTHENTIK_URL}/api/v3/flows/instances/${slug}/" | grep -Eq "\"slug\": ?\"${slug}\""; do
      retry_count=$((retry_count + 1))
      if [[ $retry_count -ge $max_retries ]]; then
        echo -e "${RED}Flow $slug not found after $max_retries retries. Aborting...${ENDCOLOR}"
        exit 1
      fi
      echo -e "${YELLOW}Flow $slug not found yet. Retrying in $retry_interval seconds...${ENDCOLOR}"
      sleep $retry_interval
    done
  done

  echo -e "Default flows are ${GREEN}available${ENDCOLOR}"
  echo
}

main() {
  check_authentik_status
  wait_for_default_flows
  build_ak_bin
  create_config
  create_tenant
  create_admin_user_for_tenant
  update_user
  delete_user
  delete_tenant
  create_group
  update_group
  delete_group
  create_app_oidc
  update_app_oidc
  delete_app
}

main "$@"