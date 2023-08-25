package object

type IndexEntry struct {
}

type Index struct {
	entries []IndexEntry
}

//func NewIndex(content []byte) (Index, error) {
//	contentLength := len(content)
//
//	// Convert everything - the last 20 bytes to sha1 checksum
//	digest := sha1.Sum(content[:contentLength-20])
//	if !bytes.Equal(digest[:], content[contentLength-20:]) {
//		return Index{}, errors.New("checksum mismatch")
//	}
//
//	signature := content[:4]
//	if !bytes.Equal(signature, []byte("DIRC")) {
//		return Index{}, errors.New("invalid signature")
//	}
//
//	version := binary.BigEndian.Uint64(content[4:8])
//	if version != 2 {
//		// for simplicity, we only support version 2
//		return Index{}, errors.New("invalid version")
//	}
//
//	entriesCount := binary.BigEndian.Uint64(content[8:12])
//	entriesContent := content[12:]
//
//	entries := make([]IndexEntry, entriesCount)
//
//	//idx := 0
//	//for i := range entriesCount {
//	//
//	//}
//}
