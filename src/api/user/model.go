package user

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/pbkdf2"
)

type userName struct {
	First string `bson:"first" json:"first"`
	Last string `bson:"last" json:"last"`
}

type User struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name userName `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	IsAdmin bool `bson:"isAdmin" json:"isAdmin"`
	Salt string `bson:"salt" json:"salt"`
	Hash string `bson:"hash" json:"hash"`
}

type userProfile struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name userName `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	IsAdmin bool `bson:"isAdmin" json:"isAdmin"`
}

func (user *User) toProfile() userProfile {
	up := userProfile{user.Id, user.Name, user.Email, user.IsAdmin}
	return up
}

func randByte(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return []byte(""), err
	}
	return bytes, nil
}

func (user *User) SetPassWord(password string) {
	saltByte, _ := randByte(16)
	user.Salt = hex.EncodeToString(saltByte)
	user.Hash = hex.EncodeToString(pbkdf2.Key([]byte(password), saltByte, 10000, 512, sha512.New))
}

func (user *User) VerifyPassWord(password string) bool {
	saltByte, _ := hex.DecodeString(user.Salt)
	hash := hex.EncodeToString(pbkdf2.Key([]byte(password), saltByte, 10000, 512, sha512.New))
	return user.Hash == hash
}