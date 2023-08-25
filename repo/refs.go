package repo

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type Ref struct {
	SHA  string
	Path string
}

func (r Repo) FetchRefs() ([]Ref, error) {
	refRoot := path.Join("refs", "heads")

	files, err := os.ReadDir(path.Join(r.gitPath, refRoot))
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return []Ref{}, nil
	}

	var refs []Ref

	for _, file := range files {
		if file.IsDir() {
			// No dirs expected, continue on to next file
			continue
		}

		filePath := path.Join(refRoot, file.Name())
		refContent, err := r.ResolveRef(filePath)
		if err != nil {
			// Corrupt ref, continue on
			continue
		}

		refs = append(refs, Ref{SHA: refContent, Path: filePath})
	}

	return refs, nil
}

func (r Repo) ResolveRef(refPath string) (string, error) {
	contentBytes, err := os.ReadFile(path.Join(r.gitPath, refPath))
	if err != nil {
		fmt.Printf("%v\n", err)
		return "", err
	}

	contentLength := len(contentBytes)
	if strings.HasPrefix(string(contentBytes[:5]), "ref: ") {
		// ends with a \n, so we remove it
		newRefPath := string(contentBytes[5 : contentLength-1])
		return r.ResolveRef(newRefPath)
	} else {
		// ends with a \n so we remove it
		return string(contentBytes[:contentLength-1]), nil
	}
}