package dto

import "encoding/json"

//Response - ResCode may be 0 - success / 1 - server fail / 2 - IncorrectRequest
type Response struct {
	Message string `json:"message"`
	ResCode int    `json:"resCode"`
}

//ToJSON - convert Response to json
func (resp Response) ToJSON() string {
	jsonData, err := json.Marshal(resp)

	if err != nil {
		panic(err)
	}

	return string(jsonData)
}

//MakeJSON - create Response and convert to json
func (resp Response) MakeJSON(message string, resCode int) string {
	jsonData, err := json.Marshal(Response{Message: message, ResCode: resCode})

	if err != nil {
		panic(err)
	}

	return string(jsonData)
}
