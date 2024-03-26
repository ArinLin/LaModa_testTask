package model

import (
	"time"
)

type Warehouse struct {
	ID          int        `json:"id" example:"1"`
	Name        string     `json:"name" example:"Sofino"`
	IsAvailable bool       `json:"is_available" example:"true"`
	CreatedAt   time.Time  `json:"created_at" example:"2024-03-24T21:16:31Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2024-03-24T21:16:31Z"`
	DeletedAt   *time.Time `json:"deleted_at" example:"null"`
}

type CreateWarehouseRequest struct {
	Name        string `json:"name" example:"Bojazbus"`
	IsAvailable bool   `json:"is_available" example:"true"`
}

type UpdateWarehouseRequest struct {
	Name        string `json:"name" example:"Tegespu"`
	IsAvailable bool   `json:"is_available" example:"false"`
}

type GetWarehousesResponse struct {
	Status  string      `json:"status" example:"OK"`
	MsgCode string      `json:"msg_code" example:"warehouses_received"`
	Data    []Warehouse `json:"data" example:"[{\"id\":1,\"name\":\"Bojazbus\",\"is_available\":\"true\",\"created_at\":\"2024-03-24T21:16:36Z\",\"updated_at\":\"2024-03-24T21:16:36Z\",\"deleted_at\":null},{\"id\":2,\"name\":\"Ratasit\",\"is_available\":\"true\",\"created_at\":\"2024-03-24T22:04:29Z\",\"updated_at\":\"2024-03-24T22:04:29Z\",\"deleted_at\":null}]"`
}

type GetWarehouseByIDResponse struct {
	Status  string    `json:"status" example:"OK"`
	MsgCode string    `json:"msg_code" example:"warehouse_received"`
	Data    Warehouse `json:"data" example:"{\"id\":1,\"name\":\"Bojazbus\",\"is_available\":\"true\",\"created_at\":\"2024-03-24T21:16:36Z\",\"updated_at\":\"2024-03-24T21:16:36Z\",\"deleted_at\":null}"`
}

type CreateWarehouseResponse struct {
	Status  string    `json:"status" example:"OK"`
	MsgCode string    `json:"msg_code" example:"warehouse_created"`
	Data    Warehouse `json:"data" example:"{\"id\":1,\"name\":\"Bojazbus\",\"is_available\":\"true\",\"created_at\":\"2024-03-24T21:16:36Z\",\"updated_at\":\"2024-03-24T21:16:36Z\",\"deleted_at\":null}"`
}

type UpdateWarehouseResponse struct {
	Status  string    `json:"status" example:"OK"`
	MsgCode string    `json:"msg_code" example:"warehouse_updated"`
	Data    Warehouse `json:"data" example:"{\"id\":1,\"name\":\"Tegespu\",\"is_available\":\"false\",\"created_at\":\"2024-03-24T21:16:36Z\",\"updated_at\":\"2024-03-24T21:46:36Z\",\"deleted_at\":null}"`
}

type DeleteWarehouseResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"warehouse_deleted"`
}
