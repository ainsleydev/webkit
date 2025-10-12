package testutil

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
)

type AferoErrCreateFs struct {
	afero.Fs
}

func (e *AferoErrCreateFs) Create(_ string) (afero.File, error) {
	return nil, fmt.Errorf("create error")
}

func (e *AferoErrCreateFs) OpenFile(_ string, _ int, _ os.FileMode) (afero.File, error) {
	return nil, fmt.Errorf("openfile error")
}
