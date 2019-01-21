package auth

import (
    "context"
    "errors"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/bson/primitive"
    "time"
)

func (auth *Auth) buildClaims(payload map[string]interface{}) jwt.Claims {
    claims := jwt.MapClaims{
        "exp": time.Now().Add(time.Second*time.Duration(auth.exp)).Unix(),
    }
    for k, v := range payload {
        claims[k] = v
    }
    return claims
}

func (auth *Auth) checkBlackList(tokenString string) error {
    ctx, _ := context.WithTimeout(context.Background(), time.Second)
    cnt, _ := auth.db.Count(ctx, bson.D{{"token", tokenString}})
    if cnt > 0 {
        return errors.New(fmt.Sprintf("token: %s is already invalidated", tokenString))
    }
    return nil
}

func (auth *Auth) GetToken(payload map[string]interface{}) (string, error) {
    claims := auth.buildClaims(payload)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(auth.secret))
}

func (auth *Auth) VerifyToken(tokenString string) (jwt.MapClaims, error) {
    if err := auth.checkBlackList(tokenString); err != nil {
        return nil, err
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(auth.secret), nil
    }); if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(jwt.MapClaims); if ok && token.Valid {
        return claims, nil
    } else {
        return nil, err
    }
}

func (auth *Auth) Invalidate(tokenString string) error {
    if err := auth.checkBlackList(tokenString); err != nil {
        return err
    }
    ctx, _ := context.WithTimeout(context.Background(), time.Second)
    blacklists := BlackList{Token: tokenString, CreatedAt: primitive.DateTime(time.Now().Unix()*1000)}
    if _, err := auth.db.InsertOne(ctx, &blacklists); err != nil {
        return err
    }
    return nil
}