#!/bin/bash
# waitForContainer.sh
#set -e

# Max query attempts before consider setup failed
MAX_TRIES=3

function serviceIsReady() {
  #docker-compose logs payments | grep "Starting HTTP service"
  #$(curl --output /dev/null --silent --head --fail http://localhost:7070/api/v1/healthcheck/)
  STATUS=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:7070/api/v1/healthcheck/)

  if [ $STATUS -eq 200 ]; then
    return 0
  fi

  return 1
}

function waitForDockerContainer() {
  attempt=1
  while [ $attempt -le $MAX_TRIES ]; do
    if "$@"; then
      echo "$2 container is up!"
      break
    fi
    echo "Waiting for $2 container... (attempt: $((attempt++)))"
    sleep 5
  done

  if [ $attempt -gt $MAX_TRIES ]; then
    echo "Error: $2 not responding, cancelling set up"
    exit 1
  fi
}

waitForDockerContainer serviceIsReady "Checkout"
