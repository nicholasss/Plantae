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

# run admin tests with admin token for testing
hurl \
  --variable lisa_email=lisa@gmail.com \
  --variable lisa_password=Growl1ng! \
  --secret super_admin_token=$SUPER_ADMIN_TOKEN \
  --test \
  test/admin.hurl

# run user tests with admin token for testing
hurl \
  --variable craig_email=craig@gmail.com \
  --variable craig_password=@ssword472 \
  --secret super_admin_token=$SUPER_ADMIN_TOKEN \
  --test \
  test/users.hurl
