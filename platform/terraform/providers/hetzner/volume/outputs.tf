output "id" {
  description = "Volume ID"
  value       = hcloud_volume.this.id
}

output "name" {
  description = "Volume name"
  value       = hcloud_volume.this.name
}

output "size" {
  description = "Volume size in GB"
  value       = hcloud_volume.this.size
}

output "linux_device" {
  description = "Device path on the server (e.g., /dev/disk/by-id/scsi-0HC_Volume_12345)"
  value       = hcloud_volume.this.linux_device
}

output "location" {
  description = "Volume location"
  value       = hcloud_volume.this.location
}
