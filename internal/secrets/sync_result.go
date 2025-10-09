package secrets

type (
	// SyncResult represents the outcome of syncing one
	// SOPS YAML file.
	SyncResult struct {
		FilePath     string
		Added        int
		Skipped      int
		WasEncrypted bool
		WasMissing   bool
		Error        error
		AddedSecrets []SecretInfo
	}
	// SyncResults aggregates all file results together
	// for analytics.
	SyncResults struct {
		Files []SyncResult
	}
)

// TotalAdded returns total secrets added across all files
func (r *SyncResults) TotalAdded() int {
	return r.sumInts(func(f SyncResult) int { return f.Added })
}

// TotalSkipped returns total secrets skipped
func (r *SyncResults) TotalSkipped() int {
	return r.sumInts(func(f SyncResult) int { return f.Skipped })
}

// EncryptedCount returns count of encrypted files encountered
func (r *SyncResults) EncryptedCount() int {
	return r.countTrue(func(f SyncResult) bool { return f.WasEncrypted })
}

// MissingCount returns count of missing files
func (r *SyncResults) MissingCount() int {
	return r.countTrue(func(f SyncResult) bool { return f.WasMissing })
}

// HasErrors returns true if any file had an error
func (r *SyncResults) HasErrors() bool {
	return r.countTrue(func(f SyncResult) bool { return f.Error != nil }) > 0
}

// sumInts sums a specific integer field across all files
func (r *SyncResults) sumInts(getter func(SyncResult) int) int {
	total := 0
	for _, f := range r.Files {
		total += getter(f)
	}
	return total
}

// countTrue counts how many files satisfy the predicate
func (r *SyncResults) countTrue(predicate func(SyncResult) bool) int {
	count := 0
	for _, f := range r.Files {
		if predicate(f) {
			count++
		}
	}
	return count
}
