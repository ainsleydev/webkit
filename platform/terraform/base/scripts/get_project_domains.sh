#!/bin/bash
# Script to fetch domain URNs from a DigitalOcean project
# Used by Terraform external data source to preserve manually-added domains
#
# Input (JSON via stdin): { "project_id": "...", "do_token": "..." }
# Output (JSON): { "domain_urns": "urn1,urn2,urn3" }

set -e

# Parse input JSON
eval "$(jq -r '@sh "PROJECT_ID=\(.project_id) DO_TOKEN=\(.do_token)"')"

# If project_id is empty, return empty domains (first apply scenario)
if [ -z "$PROJECT_ID" ] || [ "$PROJECT_ID" = "null" ]; then
  jq -n '{"domain_urns": ""}'
  exit 0
fi

# Query DigitalOcean API for project resources
RESPONSE=$(curl -s -X GET \
  -H "Authorization: Bearer $DO_TOKEN" \
  -H "Content-Type: application/json" \
  "https://api.digitalocean.com/v2/projects/$PROJECT_ID/resources")

# Extract domain URNs (those starting with "do:domain:")
# Join them with commas for Terraform
DOMAIN_URNS=$(echo "$RESPONSE" | jq -r '
  .resources[]?.urn // empty
  | select(startswith("do:domain:"))
' | tr '\n' ',' | sed 's/,$//')

# Return as JSON
jq -n --arg urns "$DOMAIN_URNS" '{"domain_urns": $urns}'
