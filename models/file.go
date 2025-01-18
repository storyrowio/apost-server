package models

import "mime/multipart"

const (
	CloudinaryDriver  = "cloudinary"
	GoogleDriveDriver = "google-drive"
)

type FileRequest struct {
	MultiFile []*multipart.FileHeader `form:"files"`
	Driver    string                  `form:"driver"`
	Token     string                  `form:"token"`
}

type FileMetadata struct {
	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Extension   string `json:"extension"`
}

type UploadResult struct {
	Url      string       `json:"url"`
	FileName string       `json:"fileName"`
	Path     string       `json:"path"`
	Metadata FileMetadata `json:"metadata"`
	Driver   string       `json:"driver" default:"cloudinary"`
}

type FileInfo struct {
	Id        string      `json:"id"`
	UserId    string      `json:"userId"`
	FileName  string      `json:"fileName"`
	Path      string      `json:"path"`
	Metadata  interface{} `json:"metadata"`
	Url       string      `json:"url"`
	Driver    string      `json:"driver" default:"cloudinary"`
	BasicDate `bson:",inline"`
}

type DestroyFileRequest struct {
	Ids []string `json:"ids"`
}

type CloudinaryInstanceSetting struct {
	CloudinaryCloudName string `json:"cloudinaryCloudName" bson:"cloudinaryCloudName"`
	CloudinaryKey       string `json:"cloudinaryKey" bson:"cloudinaryKey"`
	CloudinarySecret    string `json:"cloudinarySecret" bson:"cloudinarySecret"`
}
