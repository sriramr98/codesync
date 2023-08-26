package git

func (g Git) Status() (string, error) {
	return "Git Root at " + g.gitPath, nil
}
