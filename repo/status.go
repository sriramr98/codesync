package repo

func (r Repo) Status() (string, error) {
	return "Git Root at " + r.gitPath, nil
}
