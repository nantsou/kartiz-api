package auth

import (
    "encoding/json"
    "errors"
    "kartiz/utils"
    "net/http"
    "strings"
)

func (auth *Auth) AuthMiddleWare(next http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        var output interface{}
        tokenString := r.Header.Get("Authorization")
        if len(tokenString) == 0 {
            w.WriteHeader(http.StatusUnauthorized)
            output = utils.BuildOutput(nil,
                errors.New("missing authorization header"),
                http.StatusUnauthorized)
            _ = json.NewEncoder(w).Encode(output)
            return
        }
        tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
        claims, err := auth.VerifyToken(tokenString)
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            output = utils.BuildOutput(nil,
                err,
                http.StatusUnauthorized)
            _ = json.NewEncoder(w).Encode(output)
            return
        }
        r.Header.Set("userId", claims["userId"].(string))
        r.Header.Set("token", tokenString)
        next.ServeHTTP(w, r)
    })
}