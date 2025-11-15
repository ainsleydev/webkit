#!/bin/bash
set -e

# Read input JSON from Terraform
eval "$(jq -r '@sh "DO_TOKEN=\(.do_token)"')"

# Query DigitalOcean API for all projects
PROJECTS=$(curl -s -X GET \
  -H "Authorization: Bearer $DO_TOKEN" \
  -H "Content-Type: application/json" \
  "https://api.digitalocean.com/v2/projects")

# Count the projects
COUNT=$(echo "$PROJECTS" | jq '.projects | length')

# Return count as JSON
jq -n --arg count "$COUNT" '{"count": $count}'
