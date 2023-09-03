package git

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type Ref struct {
	SHA  string
	Path string
}

func (g Git) FetchRefs() ([]Ref, error) {
	refRoot := path.Join("refs", "heads")

	files, err := os.ReadDir(path.Join(g.gitPath, refRoot))
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
		refContent, err := g.ResolveRef(filePath)
		if err != nil {
			// Corrupt ref, continue on
			continue
		}

		refs = append(refs, Ref{SHA: refContent, Path: filePath})
	}

	return refs, nil
}

func (g Git) ResolveRef(refPath string) (string, error) {
	contentBytes, err := os.ReadFile(path.Join(g.gitPath, refPath))
	if err != nil {
		fmt.Printf("%v\n", err)
		return "", err
	}

	contentLength := len(contentBytes)
	if strings.HasPrefix(string(contentBytes[:5]), "ref: ") {
		// ends with a \n, so we remove it
		newRefPath := string(contentBytes[5 : contentLength-1])
		return g.ResolveRef(newRefPath)
	} else {
		refContent := string(contentBytes)
		if strings.HasSuffix(refContent, "\n") {
			log.Printf("Trimming newline from ref %s", refContent)
			return strings.TrimSuffix(refContent, "\n"), nil
		}
		return refContent, nil
	}
}

func (g Git) FetchHeadRefPath() (string, error) {
	_, err := os.Stat(path.Join(g.gitPath, "HEAD"))
	if err == nil {
		return "", err
	}

	headRefBytes, err := os.ReadFile(path.Join(g.gitPath, "HEAD"))
	if err != nil {
		return "", err
	}

	log.Printf("HEAD Content: %s\n", string(headRefBytes))
	// Remove the last \n character
	return strings.Split(string(headRefBytes[:len(headRefBytes)-1]), "ref: ")[1], nil
}

func (g Git) UpdateHeadRef(id string) error {
	headRefPath, err := g.FetchHeadRefPath()
	if err != nil {
		//TODO: handle if err is os.ErrNotFound
		return err
	}
	fmt.Printf("Updating Ref %s\n", headRefPath)
	refPath := path.Join(g.gitPath, headRefPath)
	return os.WriteFile(refPath, []byte(id), 0644)
}

func (g Git) UpdateHeadFile(id string) error {
	headPath := path.Join(g.gitPath, "HEAD")
	if _, err := os.Stat(headPath); err != nil {
		//TODO: Handle err is os.ErrNotFound
		return err
	}

	return os.WriteFile(headPath, []byte(id), 0644)
}

func (g Git) GetHeadRef() (string, error) {
	// headRefPath, err := g.FetchHeadRefPath()
	// if err != nil {
	// 	return "", err
	// }

	// log.Printf("Fetching ref from path %s\n", headRefPath)
	return g.ResolveRef("HEAD")
}
