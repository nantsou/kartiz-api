package utils

import (
    "errors"
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "net/http"
    "net/url"
    "reflect"
    "testing"
)

// common
func TestBuildOutput(t *testing.T) {
    data := map[string]string{"data": "test"}
    err := errors.New("test")
    statusCode := http.StatusOK

    expected := map[string]interface{} {
        "data": data,
        "statusCode": statusCode,
        "error": map[string]interface{} {"message": err.Error()},
    }
    actual := BuildOutput(data, err, statusCode)

    isEqual := reflect.DeepEqual(expected, actual); if !isEqual {
        t.Error("The expected output is different from actual output.")
    }
}

// mongo
func TestBuildFilter(t *testing.T) {
    queries := url.Values{}
    queries.Set("testKey", "testValue")

    expected := bson.D{primitive.E{Key:"testKey", Value: bson.D{{Key: "$in", Value: bson.A{"testValue"}}}}}
    actual := BuildFilter(queries)

    isEqual := reflect.DeepEqual(expected, actual); if !isEqual {
        fmt.Println(expected)
        fmt.Println(actual)
        t.Error("The expected output is different from actual output.")
    }
}