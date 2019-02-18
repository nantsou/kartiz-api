package item

import (
	"context"
	"encoding/json"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"log"
	"time"
)

type service struct {
	c *mongo.Collection
}

func (s *service) find(filter interface{}) ([]Item, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	var items []Item
	cur, err := s.c.Find(ctx, filter)

	// allocate users into array
	for cur.Next(nil) {
		user := Item{}
		err = cur.Decode(&user); if err != nil {
			log.Fatal("Decode error ", err)
			return nil, err
		}
		items = append(items, user)
	}

	// close db cursor
	err = cur.Close(nil); if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *service) create(rawData *json.Decoder) (*Item, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	var data map[string]interface{}
	err := rawData.Decode(&data); if err != nil {
		return nil, err
	}
	var item Item
	d, _ := json.Marshal(data)
	err = json.Unmarshal(d, &item); if err != nil {
		return nil, err
	}
	res, err := s.c.InsertOne(ctx, &item); if err != nil {
		return nil, err
	}
	item.Id = res.InsertedID.(primitive.ObjectID)
	return &item, nil
}

func (s *service) get(objectId primitive.ObjectID) (*Item, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	var item Item
	if err := s.c.FindOne(ctx, bson.D{{"_id",objectId}}).Decode(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *service) update(objectId primitive.ObjectID, rawData *json.Decoder) (*Item, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	filter := bson.D{{"_id", objectId}}
	var after = options.ReturnDocument(1)
	option := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	var item Item
	if err := s.c.FindOne(ctx, filter).Decode(&item); err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err := rawData.Decode(&data); if err != nil {
		return nil, err
	}
	d, _ := json.Marshal(data)
	err = json.Unmarshal(d, &item); if err != nil {
		return nil, err
	}

	if err := s.c.FindOneAndUpdate(ctx, filter, bson.D{{"$set", &item}}, &option).Decode(&item); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &item, nil
}

func (s *service) delete(objectId primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	if res := s.c.FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: objectId}}); res.Err() != nil {
		return res.Err()
	}
	return nil
}