// Package fsmedia is a default implementation of InterfaceMediaStorage declared in
// "github.com/ottemo/foundation/media" package
package fsmedia

// Package global constants
const (
	ConstMediaDBCollection  = "media"    // database collection name to store media assignment information into
	ConstMediaDefaultFolder = "./media/" // filesystem folder path to store media files in there
)

// InterfaceMediaStorage implementer class
type FilesystemMediaStorage struct {
	storageFolder string
	setupWaitCnt  int
}
