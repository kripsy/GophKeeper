package models

type MiltipartUploadFileData struct {
	Content  []byte
	FileName string
	Guid     string
	Hash     string
	Username string
}

type ObjectPart struct {
	PartNumber int
	ETag       string
}
