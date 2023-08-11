package repo

import (
	"github.com/spf13/afero"
)

type Repo struct {
	fs afero.Fs
}
