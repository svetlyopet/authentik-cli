#!/usr/bin/env bash

set -eo pipefail

: ${AUTHENTIK_TAG:="2025.2.4"}
: ${AUTHENTIK_URL:="http://localhost:9000"}

AUTHENTIK_SECRET_KEY=$(openssl rand -base64 60 | tr -d '\n')
AUTHENTIK_BOOTSTRAP_PASSWORD=$(openssl rand -base64 30 | tr -d '\n')
AUTHENTIK_BOOTSTRAP_TOKEN=$(openssl rand -base64 30 | tr -d '\n')

PG_PASS=$(openssl rand -base64 36 | tr -d '\n')

CI_TEST_DIR=$(pwd)/ci/integration_tests

WAIT_FOR_INIT=60

help() {
  echo "Usage: $0 [command]"
  echo "Commands:"
  echo "  create          Create an Authentik setup with Docker compose"
  echo "  destroy         Remove the Docker compose environment created by this script"
}

generate_env() {
  if [ -f $CI_TEST_DIR/.env ]; then
    rm $CI_TEST_DIR/.env
  fi

  echo "AUTHENTIK_URL=$AUTHENTIK_URL" >> $CI_TEST_DIR/.env
  echo "AUTHENTIK_TAG=$AUTHENTIK_TAG" >> $CI_TEST_DIR/.env
  echo "AUTHENTIK_SECRET_KEY=$AUTHENTIK_SECRET_KEY" >> $CI_TEST_DIR/.env
  echo "AUTHENTIK_BOOTSTRAP_TOKEN=$AUTHENTIK_BOOTSTRAP_TOKEN" >> $CI_TEST_DIR/.env
  echo "AUTHENTIK_BOOTSTRAP_PASSWORD=$AUTHENTIK_BOOTSTRAP_PASSWORD" >> $CI_TEST_DIR/.env
  echo "PG_PASS=$PG_PASS" >> $CI_TEST_DIR/.env
  echo "CI_TEST_DIR=$CI_TEST_DIR" >> $CI_TEST_DIR/.env
}

compose() {
  docker compose -f $CI_TEST_DIR/docker-compose.yml up -d
}

wait_for_init() {
  echo "Sleeping for $WAIT_FOR_INIT seconds until Authentik initializes..."
  sleep $WAIT_FOR_INIT
}

print_creds() {
  echo "##################### Credentials ####################"
  echo "login link: $AUTHENTIK_URL"
  echo "user: akadmin"
  echo "password: $AUTHENTIK_BOOTSTRAP_PASSWORD"
  echo "api-token: $AUTHENTIK_BOOTSTRAP_TOKEN"
}

cleanup() {
  docker compose -f $CI_TEST_DIR/docker-compose.yml down -v
  rm -f $CI_TEST_DIR/.env 2&>/dev/null
  rm -rf $CI_TEST_DIR/certs $CI_TEST_DIR/custom-templates $CI_TEST_DIR/media 2&>/dev/null
}

main() {
  case "$1" in
    create)
      generate_env
      compose
      wait_for_init
      print_creds
      ;;
    destroy)
      cleanup
      ;;
    *)
      help
      exit 1
      ;;
  esac
}

main "$@"
