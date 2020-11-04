package dto

import "encoding/json"

//Image - image DTO
type Image struct {
	Name      string
	Extension string
	Data      string
}

//ToJSON - serealize image dto
func (image Image) ToJSON() string {
	jsonData, err := json.Marshal(image)

	if err != nil {
		panic(err)
	}

	return string(jsonData)
}
