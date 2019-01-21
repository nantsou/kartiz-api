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

type service struct {
	c *mongo.Collection
}

func (s *service) find() ([]User, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	var users []User
	cur, err := s.c.Find(ctx, nil)

	// allocate users into array
	for cur.Next(nil) {
		user := User{}
		err = cur.Decode(&user); if err != nil {
			log.Fatal("Decode error ", err)
			return nil, err
		}
		users = append(users, user)
	}

	// close db cursor
	err = cur.Close(nil); if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *service) create(rawData *json.Decoder) (*User, error) {
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
	res, err := s.c.InsertOne(ctx, &user); if err != nil {
		return nil, err
	}
	user.Id = res.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (s *service) get(objectId primitive.ObjectID) (*User, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	var user User
	if err := s.c.FindOne(ctx, bson.D{{"_id",objectId}}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) update(objectId primitive.ObjectID, rawData *json.Decoder) (*User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	filter := bson.D{{"_id", objectId}}
	var after = options.ReturnDocument(1)
	option := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	var user User
	if err := s.c.FindOne(ctx, filter).Decode(&user); err != nil {
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

	if err := s.c.FindOneAndUpdate(ctx, filter, bson.D{{"$set", &user}}, &option).Decode(&user); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &user, nil
}

func (s *service) delete(objectId primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	if res := s.c.FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: objectId}}); res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (s *service) login(decoder *json.Decoder) (*User, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	var data map[string]interface{}
	err := decoder.Decode(&data); if err != nil {
		return nil, err
	}
	email, ok := data["email"]; if !ok {
		return nil, errors.New("email is necessary")
	}
	password, ok := data["password"]; if !ok {
		return nil, errors.New("password is necessary")
	}
	var user User
	if err := s.c.FindOne(ctx, bson.D{{"email",email}}).Decode(&user); err != nil {
		return nil, err
	}
	if !user.VerifyPassWord(password.(string)) {
		return nil, errors.New("email or password is not correct")
	}
	return &user, nil
}

func (s *service) logout() {

}