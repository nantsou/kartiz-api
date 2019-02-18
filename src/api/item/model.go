package item

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Item struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	OrderId primitive.ObjectID `bson:"orderId json:orderId"`
	Name string `bson:"name" json:"name"`
	Link string `bson:"link" json:"link"`
	Price float32 `bson:"price" json:"price"`
	Quantity int32 `bson:"quantity" json:"quantity"`
	Note string `bson:"note" json:"note"`
}