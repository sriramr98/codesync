package repo

import (
	"bufio"
	"bytes"
	"io"

	"gopkg.in/ini.v1"
)

func (r *Repo) getConfig() ([]byte, error) {
	config := ini.Empty()
	core, err := config.NewSection("core")
	if err != nil {
		return nil, err
	}

	// the version of the gitdir format. 0 means the initial format, 1 the same with extensions. If > 1, git will panic; rs will only accept 0.
	core.NewKey("repositoryformatversion", "0")
	// disable tracking of file mode (permissions) changes in the work tree.
	core.NewKey("filemode", "true")
	// indicates that this repository has a worktree. Git supports an optional worktree key which indicates the location of the worktree, if not ..; rs doesn’t.
	core.NewKey("bare", "false")

	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	if _, err := config.WriteTo(writer); err != nil {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	return io.ReadAll(buf)

}
