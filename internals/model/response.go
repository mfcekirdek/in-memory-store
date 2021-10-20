package model

type BaseResponse struct {
	Data        interface{} `json:"data"`
	Description string      `json:"description"`
}
