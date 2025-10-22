//go:build !race

package infra

// Package infra manages Terraform-based infrastructure provisioning for WebKit projects.
//
// # Why This Package Exists
//
// WebKit projects define infrastructure needs in app.json (databases, apps, storage).
// This package transforms those definitions into Terraform variables and orchestrates
// the Terraform CLI to provision actual cloud resources.
//
// Instead of maintaining Terraform modules in every project, WebKit uses centralized
// modules. Projects declare WHAT they need; WebKit handles HOW to provision it.
//
// # State Management
//
// Terraform state lives in BackBlaze B2 at: {project-name}/{environment}/terraform.tfstate
// Backend config is generated dynamically with credentials from environment variables.
//
// # Typical Flow
//
//  1. Init() - Copy embedded Terraform templates to temp dir, configure remote backend
//  2. Plan() - Generate tfvars from app.json, show what will change
//  3. Apply() - Provision infrastructure
//  4. Cleanup() - Remove temporary files
//
// All operations use the hashicorp/terraform-exec library to interact with Terraform CLI.
