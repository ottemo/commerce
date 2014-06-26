package media

type IMediaStorage interface {

	GetName() string

	Load(model string, objId string, mediaType string, mediaName string) ([]byte, error)
	Save(model string, objId string, mediaType string, mediaName string, mediaData []byte) error

	ListMedia(model string, objId string, mediaType string) ([]string, error)
}
