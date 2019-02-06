package utils

import (
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/bson/primitive"
    "net/url"
)

func BuildOutput(data interface{}, err error, statusCode int) map[string]interface{} {
    output := map[string]interface{}{"data": data, "error": nil}
    if err != nil {
        output["error"] = map[string]interface{}{"message": err.Error(), "statusCode": statusCode}
    }
    return output
}

func BuildFilter(queries url.Values) bson.D {
    filter := bson.D{}
    for key, value := range queries {
        if len(value) == 1 && value[0] == "" {
            continue
        }
        array := bson.A{}
        for _, element := range value {
            array = append(array, element)
        }
        condition := primitive.E{Key:key, Value: bson.D{{Key: "$in", Value: array}}}
        filter = append(filter, condition)
    }
    return filter
}