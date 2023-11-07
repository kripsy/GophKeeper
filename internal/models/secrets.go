package models

type MultipartUploadFileData struct {
	Content  []byte
	FileName string
	//nolint:revive,stylecheck
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
	//nolint:revive,stylecheck
	Guid string
	Hash string
}

type MultipartDownloadFileRequest struct {
	FileName string
	//nolint:revive,stylecheck
	Guid string
	Hash string
}
