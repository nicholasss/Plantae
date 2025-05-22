#!/usr/bin/env bash

echo ""
echo "Running Hurl tests..."

# moves up to root from
if [[ -f "../.env" ]]; then
  # echo "Please run this script from the project root."
  cd ..
fi

if [[ ! -f ".env" ]]; then
  echo "Unable to find '.env' file."
  exit 1
fi

# source .env variables
source .env

# create secrets and variabes
hurl test/user.hurl \
  --variable super-admin-token=$SUPER_ADMIN_TOKEN \
  --secret super-admin-token \
  --test
