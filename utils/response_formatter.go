package utils

import (
	"encoding/json"
	"net/http"
)

type ResponseParams interface {
	StatusCode() int
	Json() []byte
}

type ResponseSuccessParams struct {
	Code   int         `json:"status_code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func (r *ResponseSuccessParams) StatusCode() int {
	return r.Code
}

func (r *ResponseSuccessParams) Json() []byte {
	responseJson, err := json.Marshal(r)
	if err != nil {
		return []byte("error while parsing response to json")
	}
	return responseJson
}

type ResponseErrorParams struct {
	Code   int      `json:"status_code"`
	Status string   `json:"status"`
	Errors []string `json:"errors"`
}

func (r *ResponseErrorParams) StatusCode() int {
	return r.Code
}

func (r *ResponseErrorParams) Json() []byte {
	responseJson, err := json.Marshal(r)
	if err != nil {
		return []byte("error while parsing response to json")
	}
	return responseJson
}

func FormatResponse(res http.ResponseWriter, params ResponseParams) {
	res.WriteHeader(params.StatusCode())
	res.Write(params.Json())
}
