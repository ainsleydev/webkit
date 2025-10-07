output "id" {
  description = "The ID of the S3 bucket"
  value       = digitalocean_spaces_bucket.this.id
}

output "urn" {
  description = "The URN of the S3 bucket"
  value       = digitalocean_spaces_bucket.this.urn
}

output "name" {
  description = "The name of the S3 bucket"
  value       = digitalocean_spaces_bucket.this.name
}

output "region" {
  description = "What region the bucket was crearted in"
  value = digitalocean_spaces_bucket.this.region
}

output "domain_name" {
  description = "The domain name of the S3 bucket"
  value       = digitalocean_spaces_bucket.this.bucket_domain_name
}

