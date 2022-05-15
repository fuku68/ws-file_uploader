package uploader

type Uploader interface {
	UploadStart(key string) error
	Upload(input []byte) error
	UploadEnd()
}
