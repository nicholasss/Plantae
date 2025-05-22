#!/usr/bin/env bash

echo "Running Hurl tests..."

if [[ -f ".env" ]]; then
  hurl --test --variables-file .env .
elif [[ -f "../.env" ]]; then
  hurl --test --variables-file ../.env .
else
  echo "Could not find the '.env' file."
  echo "Please run from either the root or tests directory."
fi
