package app

import "io"

// Closeable is a system that needs to be closed gracefully when the
// application is signalled to terminate
type Closeable interface {
	io.Closer
}
