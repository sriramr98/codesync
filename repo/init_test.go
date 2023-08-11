package repo

import (
	"io/fs"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInitRepo(t *testing.T) {

	t.Run("Creates all files and directories", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		repo := Repo{fs: fs}

		err := repo.Init("/abcd/edgh")
		assert.NoError(t, err)

		gitDirExists, gitDirErr := afero.Exists(fs, "/abcd/edgh/.git")
		assert.NoError(t, gitDirErr)
		assert.True(t, gitDirExists)

		branchesDirExists, branchesDirErr := afero.Exists(fs, "/abcd/edgh/.git/branches")
		assert.NoError(t, branchesDirErr)
		assert.True(t, branchesDirExists)

		objectsDirExists, objectsDirErr := afero.Exists(fs, "/abcd/edgh/.git/objects")
		assert.NoError(t, objectsDirErr)
		assert.True(t, objectsDirExists)

		tagsDirExists, tagsDirErr := afero.Exists(fs, "/abcd/edgh/.git/refs/tags")
		assert.NoError(t, tagsDirErr)
		assert.True(t, tagsDirExists)

		headsDirExists, headsDirErr := afero.Exists(fs, "/abcd/edgh/.git/refs/heads")
		assert.NoError(t, headsDirErr)
		assert.True(t, headsDirExists)

		headFile, headFileErr := afero.ReadFile(fs, "/abcd/edgh/.git/HEAD")
		assert.NoError(t, headFileErr)
		assert.Equal(t, "ref: refs/heads/master\n", string(headFile))

		configFile, configFileErr := afero.ReadFile(fs, "/abcd/edgh/.git/config")
		assert.NoError(t, configFileErr)
		expectedConfig, configFileErr := repo.getConfig()
		assert.Equal(t, expectedConfig, configFile)
		assert.NoError(t, configFileErr)

	})

	t.Run("init with existing repo", func(t *testing.T) {
		// given
		memfs := afero.NewMemMapFs()
		if err := memfs.Mkdir("/abc/.git", fs.ModeDir); err != nil {
			t.Fatal(err)
		}

		repo := Repo{fs: memfs}

		err := repo.Init("/abc")

		assert.ErrorIs(t, err, ErrAlreadyInitialized)

	})

}
