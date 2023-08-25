package object

type GitObject interface {
	// GitType returns the type of the git object
	GitType() string
	// Encode converts the content of the object into string representation
	Encode() (string, error)
}
