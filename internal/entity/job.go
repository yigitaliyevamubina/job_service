package entity

import "time"

type Job struct {
	Id          string    `protobuf:"1"`
	Title       string    `protobuf:"2"`
	Description string    `protobuf:"3"`
	OwnerId     string    `protobuf:"4"`
	Price       float32   `protobuf:"5"`
	FromDate    string    `protobuf:"6"`
	ToDate      string    `protobuf:"6"`
	CreatedAt   time.Time `protobuf:"6"`
	UpdatedAt   time.Time `protobuf:"7"`
	DeletedAt   time.Time `protobuf:"8"`
}

type GetListFilter struct {
	Page           int64  `json:"page" protobuf:"1"`
	Limit          int64  `json:"limit" protobuf:"2"`
	Search         string `json:"search" protobuf:"3"`
	OrderBy        string `json:"order_by" protobuf:"4"`
	IncludeDeleted bool   `json:"include_deleted" protobuf:"5"`
}
