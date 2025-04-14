package models

type Video struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}
