package operations

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
)

type errCreateFs struct {
	afero.Fs
}

func (e *errCreateFs) Create(_ string) (afero.File, error) {
	return nil, fmt.Errorf("create error")
}

func (e *errCreateFs) OpenFile(_ string, _ int, _ os.FileMode) (afero.File, error) {
	return nil, fmt.Errorf("openfile error")
}
