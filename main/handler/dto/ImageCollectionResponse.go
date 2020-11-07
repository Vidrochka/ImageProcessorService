package dto

import "encoding/json"

//ImageCollectionResponse - Response with image collection
type ImageCollectionResponse struct {
	File []SaveImageResponseFile `json:"file"`
}

//SaveImageResponseFile - Base64 file
type SaveImageResponseFile struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Extension  string `json:"extension"`
	Status     int8   `json:"statis"`
	ResMessage string `json:"resMessage"`
}

//ToJSON - convert Base64Response to json
func (resp ImageCollectionResponse) ToJSON() string {
	jsonData, err := json.Marshal(resp)

	if err != nil {
		panic(err)
	}

	return string(jsonData)
}

//ToJSON - convert Base64Response to json
func (resp SaveImageResponseFile) ToJSON() string {
	jsonData, err := json.Marshal(resp)

	if err != nil {
		panic(err)
	}

	return string(jsonData)
}
