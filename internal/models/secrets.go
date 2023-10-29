package models

type MultipartUploadFileData struct {
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

type MultipartDownloadFileResponse struct {
	Content  []byte
	FileName string
	Guid     string
	Hash     string
}

type MultipartDownloadFileRequest struct {
	FileName string
	Guid     string
	Hash     string
}
