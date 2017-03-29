package models

import (
	"beego-mgo-rabbitmq-jwt/utilities/helper"
)

type PageResult struct {
	Result interface{}       `json:"result"`
	Page   *helper.Paginator `json:"page"`
}
