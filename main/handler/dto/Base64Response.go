package dto

import "encoding/json"

//Base64Response - Base64 Response
type Base64Response struct {
	File []SaveImageResponseFile
}

//SaveImageResponseFile - Base64 file
type SaveImageResponseFile struct {
	ID   int64
	Name string
}

//ToJSON - convert Base64Response to json
func (resp Base64Response) ToJSON() string {
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
