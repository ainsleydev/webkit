package appdef

// TODO: Validate .Path on App and ensure it exists.
// TODO: Validate that terraform-managed VM apps (.Infra.Type == "vm" && .IsTerraformManaged()) must have at least one domain in .Domains array.
// TODO: Validate that domain names in .Domains should not contain protocol prefixes (e.g., "https://").
