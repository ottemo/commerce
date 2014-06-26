package fsmedia

const (
	MEDIA_DB_COLLECTION = "media"
	MEDIA_DEFAULT_FOLDER = "./media"
)

type FilesystemMediaStorage struct {
	storageFolder string
	setupWaitCnt int
}

type FilesystemMediaItem struct {
	id string
	entityModel string
	entityId string
	entityMedia []string

	FilesystemMediaStorage
}
