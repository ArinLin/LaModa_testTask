package model

import "lamoda/pkg/web"

type ValidationErrors struct {
	Errors []web.ValidationError `json:"errors"`
}

type BadRequestInvalidBodyResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"invalid_body"`
}

type BadRequestInvalidIDResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"invalid_id"`
}

type GoodNotFoundResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"good_not_found"`
}

type WarehouseNotFoundResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"warehouse_not_found"`
}

type ValidationResponse struct {
	Status  string           `json:"status" example:"ERROR"`
	MsgCode string           `json:"msg_code" example:"go_validation"`
	Data    ValidationErrors `json:"data" example:"{\"errors\": [{\"tag\": \"<tag>\", \"field\": \"<field>\", \"param\": \"<param>\"}]}"`
}

type InternalResponse struct {
	Status  string `json:"status" example:"ERROR"`
	MsgCode string `json:"msg_code" example:"general_internal"`
}
