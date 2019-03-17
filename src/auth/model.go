package auth

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type BlackList struct {
	Token string `bson:"token" json:"token"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

func (auth *Auth) setIndex() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	iv := auth.db.Indexes()
	cur, _ := iv.List(ctx)
	reset := true
	for cur.Next(ctx) {
		var index map[string]interface{}
		_ = cur.Decode(index)
		name, ok := index["name"]; if !ok {
			continue
		} else {
			if name.(string) == "createdAt" {
				expireAfterSeconds, ok := index["expireAfterSeconds"]; if ok {
					reset = expireAfterSeconds.(int32) != auth.exp
					break
				}
			}
		}
	}
	_ = cur.Close(ctx)

	if reset {
		_, _ = iv.DropOne(ctx, "createdAt")
		iob := mongo.NewIndexOptionsBuilder()
		iob.Name("createdAt").ExpireAfterSeconds(auth.exp)
		im := mongo.IndexModel{
			Keys: bsonx.Doc{{"createdAt", bsonx.Int32(1)}},
		}
		_, _ = iv.CreateOne(ctx, im)
	}
}

type Auth struct {
	secret string
	exp int32
	db *mongo.Collection
}

func (auth *Auth) SetConfig(config map[string]interface{}) {
	auth.secret = config["secret"].(string)
	auth.exp = int32(config["exp"].(float64))
}

func (auth *Auth) SetDB(db *mongo.Collection) {
	auth.db = db
	auth.setIndex()
}