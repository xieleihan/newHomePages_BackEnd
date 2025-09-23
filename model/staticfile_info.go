package model

type StaticFileInfo struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	Ext          string `json:"ext"`
	Size         int64  `json:"size"`
	SizeFormatted string `json:"sizeFormatted"`
	LastModified string `json:"lastModified"`
}