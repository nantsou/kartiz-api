package utils

import (
    "encoding/json"
    "errors"
    "github.com/gin-gonic/gin"
    "net/http"
    "net/http/httptest"
    "reflect"
    "testing"
)

func performRequest(r http.Handler, method, path string, headers ...map[string]string) *httptest.ResponseRecorder {
    req := httptest.NewRequest(method, path, nil)
    for _, h := range headers {
        req.Header.Add(h["key"], h["value"])
    }
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    return w
}

func TestRecovery(t *testing.T) {
    router := gin.New()
    router.Use(Recovery())

    type args struct {
        path string
        handler gin.HandlerFunc
    }

    type fields struct {
        Data interface{}    `json:"data"`
        Error string        `json:"error"`
        StatusCode int      `json:"statusCode"`
    }

    tests := []struct {
        name string
        args args
        want fields
    }{
        {
            name: "panic",
            args: args{
                path: "/panic",
                handler: func(c *gin.Context) {
                    panic("panic")
                },
            },
            want: fields{
                Data: nil,
                Error: errors.New("panic").Error(),
                StatusCode: http.StatusInternalServerError,
            },
        },
        {
            name: "runtime",
            args: args{
                path: "/runtime",
                handler: func(c *gin.Context) {
                    arr := []int{}
                    arr[0] = 0
                },
            },
            want: fields{
                Data: nil,
                Error: errors.New("runtime error: index out of range [0] with length 0").Error(),
                StatusCode: http.StatusInternalServerError,
            },
        },
    }

    for _, tc := range tests {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            router.GET(tc.args.path, tc.args.handler)
            w := performRequest(router, "GET", tc.args.path)
            var actual fields
            err := json.Unmarshal(w.Body.Bytes(), &actual); if err != nil {
                t.Errorf("Fail to unmarshal for test case %v", tc.name)
            }
            if !reflect.DeepEqual(actual, tc.want) {
                t.Errorf("actual: %v, want: %v", actual, tc.want)
            }
        })
    }
}