package utils

import (
    "errors"
    "github.com/gin-gonic/gin"
    "net/http"
    "reflect"
    "testing"
)

func TestBuildOutput(t *testing.T) {
    type args struct {
        data interface{}
        err error
        statusCode int
    }

    tests := []struct {
        name string
        args args
        want gin.H
    }{
        {
            name: "no data no error 200",
            args: args{
                statusCode: http.StatusOK,
            },
            want: gin.H{
                "data": nil,
                "error": nil,
                "statusCode": http.StatusOK,
            },
        },
        {
            name: "has data no error 200",
            args: args{
                data: map[string]string{"data":"test"},
                statusCode: http.StatusOK,
            },
            want: gin.H{
                "data": map[string]string{"data":"test"},
                "error": nil,
                "statusCode": http.StatusOK,
            },
        },
        {
            name: "no data error 400",
            args: args{
                err: errors.New("bad request"),
                statusCode: http.StatusBadRequest,
            },
            want: gin.H{
                "data": nil,
                "error": errors.New("bad request").Error(),
                "statusCode": http.StatusBadRequest,
            },
        },
        {
            name: "has data error 400",
            args: args{
                data: map[string]string{"data":"test"},
                err: errors.New("bad request"),
                statusCode: http.StatusBadRequest,
            },
            want: gin.H{
                "data": map[string]string{"data":"test"},
                "error": errors.New("bad request").Error(),
                "statusCode": http.StatusBadRequest,
            },
        },
    }

    for _, tc := range tests {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            actual := BuildJsonOutput(tc.args.data, tc.args.err, tc.args.statusCode)
            if !reflect.DeepEqual(actual, tc.want) {
                t.Errorf("actual: %v, want: %v", actual, tc.want)
            }
        })
    }
}