package model

import (
	"time"
)

type Good struct {
	ID        int        `json:"id" example:"1"`
	Name      string     `json:"name" example:"T-Shirt"`
	Size      string     `json:"size" example:"s"`
	CreatedAt time.Time  `json:"created_at" example:"2024-03-24T21:16:31Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2024-03-24T21:16:31Z"`
	DeletedAt *time.Time `json:"deleted_at" example:"null"`
}

type CreateGoodRequest struct {
	Name string `json:"name" example:"T-Shirt"`
	Size string `json:"size" example:"m"`
}

type UpdateGoodRequest struct {
	Name string `json:"name" example:"Shoes"`
	Size string `json:"size" example:"s"`
}

type GetGoodsResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"goods_received"`
	Data    []Good `json:"data"`
}

type GetGoodByIDResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"good_received"`
	Data    Good   `json:"data"`
}

type CreateGoodResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"good_created"`
	Data    Good   `json:"data"`
}

type UpdateGoodResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"good_updated"`
	Data    Good   `json:"data"`
}

type DeleteGoodResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"good_deleted"`
}
