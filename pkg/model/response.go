// Package model has data models.
package model

// BaseResponse is the common response structure of all endpoints.
type BaseResponse struct {
	Data        interface{} `json:"data"`
	Description string      `json:"description"`
}
