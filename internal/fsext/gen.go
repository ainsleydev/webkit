package fsext

import (
	"github.com/spf13/afero"
)

//go:generate go tool go.uber.org/mock/mockgen -source=gen.go -destination ../mocks/fs.go -package=mocks

// FS is a stub for afero.Fs for testing.
type FS interface {
	afero.Fs
}
