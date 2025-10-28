#!/usr/bin/env bash

set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"

print_section() {
  echo
  echo "=== $1 ==="
}

cleanup() {
  if [[ -n "${USER_ID:-}" ]]; then
    echo
    echo "Deleting user ${USER_ID}"
    curl -sS -X DELETE "${API_BASE_URL}/users/${USER_ID}" \
      -H "Authorization: Bearer ${TOKEN}" \
      -o /dev/null -w "Status: %{http_code}\n"
  fi
}

trap cleanup EXIT

print_section "Register"
REGISTER_RESPONSE="$(curl -sS -X POST "${API_BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"name":"Smoke Test","email":"smoke@example.com","password":"changeme123"}')"
echo "${REGISTER_RESPONSE}" | jq '.'

TOKEN="$(echo "${REGISTER_RESPONSE}" | jq -r '.token')"
USER_ID="$(echo "${REGISTER_RESPONSE}" | jq -r '.user.id')"

if [[ -z "${TOKEN}" || -z "${USER_ID}" || "${TOKEN}" == "null" || "${USER_ID}" == "null" ]]; then
  echo "Failed to parse token or user id"
  exit 1
fi

print_section "Login"
curl -sS -X POST "${API_BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"smoke@example.com","password":"changeme123"}' | jq '.'

print_section "List Users"
curl -sS "${API_BASE_URL}/users" \
  -H "Authorization: Bearer ${TOKEN}" | jq '.'

print_section "Get User"
curl -sS "${API_BASE_URL}/users/${USER_ID}" \
  -H "Authorization: Bearer ${TOKEN}" | jq '.'

print_section "Update User"
curl -sS -X PATCH "${API_BASE_URL}/users/${USER_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"name":"Smoke Test Updated"}' | jq '.'

print_section "Delete User"
curl -sS -X DELETE "${API_BASE_URL}/users/${USER_ID}" \
  -H "Authorization: Bearer ${TOKEN}" -o /dev/null -w "Status: %{http_code}\n"

USER_ID=""
TOKEN=""
echo
echo "Smoke test completed successfully."
