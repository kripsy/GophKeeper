// Package models defines data structures used throughout the GophKeeper application.
// This file includes models for multipart file upload and download operations.
package models

// MultipartUploadFileData represents the data structure for uploading a file in multiple parts.
type MultipartUploadFileData struct {
	Content  []byte // Content holds the binary data of the file part.
	FileName string // FileName is the name of the file being uploaded.
	//nolint:revive,stylecheck
	Guid     string // Guid is a unique identifier for the upload session.
	Hash     string // Hash represents the hash of the entire file for integrity checks.
	Username string // Username is the name of the user performing the upload.
}

// ObjectPart represents a part of a multipart upload, including its number and ETag.
type ObjectPart struct {
	PartNumber int    // PartNumber is the sequential number of the part in the multipart upload.
	ETag       string // ETag is the entity tag associated with the part, used for validation.
}

// MultipartDownloadFileResponse represents the data structure for downloading a file part.
type MultipartDownloadFileResponse struct {
	Content  []byte // Content holds the binary data of the file part.
	FileName string // FileName is the name of the file being downloaded.
	//nolint:revive,stylecheck
	Guid string // Guid is a unique identifier for the download session.
	Hash string // Hash represents the hash of the entire file for integrity checks.
}

// MultipartDownloadFileRequest represents the request data structure for downloading a file part.
type MultipartDownloadFileRequest struct {
	FileName string // FileName is the name of the file to be downloaded.
	//nolint:revive,stylecheck
	Guid string // Guid is a unique identifier for the download session.
	Hash string // Hash represents the hash of the entire file, used to validate file integrity.
}
