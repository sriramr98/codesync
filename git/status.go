package git

func (r Git) Status() (string, error) {
	return "Git Root at " + r.gitPath, nil
}
