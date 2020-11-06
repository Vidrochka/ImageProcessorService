package dto

import "encoding/json"

//ImageCollection - Response with image collection
type ImageCollection struct {
	File []SaveImageResponseFile
}

//SaveImageResponseFile - Base64 file
type SaveImageResponseFile struct {
	ID         int64
	Name       string
	Extension  string
	Status     int8
	ResMessage string
}

//ToJSON - convert Base64Response to json
func (resp ImageCollection) ToJSON() string {
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
