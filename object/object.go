package object

type GitObject interface {
	// GitType returns the type of the git object
	GitType() string
	// Encode converts the content of the object into string representation that can be written to GIT
	Encode() (string, error)
	// ToString converts the content to a string represntation that can be displayed to the use
	Print()
}
