package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"log"
	"time"
)

type userService struct {
	c *mongo.Collection
}

func (us *userService) find() ([]profile, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	var profileArray []profile
	cur, err := us.c.Find(ctx, nil)

	// allocate users into array
	for cur.Next(nil) {
		var p profile
		err = cur.Decode(&p); if err != nil {
			log.Fatal("Decode error ", err)
			return nil, err
		}
		profileArray = append(profileArray, p)
	}

	// close db cursor
	err = cur.Close(nil); if err != nil {
		return nil, err
	}
	return profileArray, nil
}

func (us *userService) create(rawData *json.Decoder) (*profile, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	var data map[string]interface{}
	err := rawData.Decode(&data); if err != nil {
		return nil, err
	}
	if _, ok := data["email"]; !ok {
		return nil, errors.New("email is necessary")
	}
	var user User
	d, _ := json.Marshal(data)
	err = json.Unmarshal(d, &user); if err != nil {
		return nil, err
	}
	if val, ok := data["password"]; ok {
		user.SetPassWord(val.(string))
	} else {
		return nil, errors.New("password is necessary")
	}
	res, err := us.c.InsertOne(ctx, &user); if err != nil {
		return nil, err
	}
	var p profile
	if err := us.c.FindOne(ctx, bson.D{{Key: "_id", Value: res.InsertedID}}).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (us *userService) get(objectId primitive.ObjectID) (*profile, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	var p profile
	if err := us.c.FindOne(ctx, bson.D{{"_id",objectId}}).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (us *userService) update(objectId primitive.ObjectID, rawData *json.Decoder) (*profile, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	filter := bson.D{{"_id", objectId}}
	var after = options.ReturnDocument(1)
	option := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	var user User
	if err := us.c.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err := rawData.Decode(&data); if err != nil {
		return nil, err
	}
	d, _ := json.Marshal(data)
	err = json.Unmarshal(d, &user); if err != nil {
		return nil, err
	}
	if val, ok := data["password"]; ok {
		user.SetPassWord(val.(string))
	}
	var p profile
	if err := us.c.FindOneAndUpdate(ctx, filter, bson.D{{"$set", &user}}, &option).Decode(&p); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &p, nil
}

func (us *userService) delete(objectId primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	if res := us.c.FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: objectId}}); res.Err() != nil {
		return res.Err()
	}
	return nil
}