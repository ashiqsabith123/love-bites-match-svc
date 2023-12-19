package intrefaces

type Utils interface {
	UploadPhotos(key string, image []byte) error
}
