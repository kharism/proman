package repository

import (
	"go.mongodb.org/mongo-driver/bson"
)

// BaseParam basic parameter for fetching data
type BaseParam struct {
	Skip   int
	Limit  int
	Filter bson.M
	Orders []string
}
