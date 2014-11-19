// Package media represents interfaces to access media storage services
package media

// interface to media storage service
type InterfaceMediaStorage interface {
	GetName() string

	Load(model string, objID string, mediaType string, mediaName string) ([]byte, error)
	Save(model string, objID string, mediaType string, mediaName string, mediaData []byte) error

	Remove(model string, objID string, mediaType string, mediaName string) error

	ListMedia(model string, objID string, mediaType string) ([]string, error)

	GetMediaPath(model string, objID string, mediaType string) (string, error)
}
