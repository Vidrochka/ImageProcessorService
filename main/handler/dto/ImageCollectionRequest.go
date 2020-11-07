package dto

import "encoding/json"

//ImageCollectionRequest - base64 request DTO fore deserealize
type ImageCollectionRequest struct {
	Images []Image `json:"images"`
}

//Image - image DTO
type Image struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Data      string `json:"data"`
}

//ToJSON - serealize image dto
func (image Image) ToJSON() string {
	jsonData, err := json.Marshal(image)

	if err != nil {
		panic(err)
	}

	return string(jsonData)
}
