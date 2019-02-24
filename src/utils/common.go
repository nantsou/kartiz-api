package utils

func BuildOutput(data interface{}, err error, statusCode int) map[string]interface{} {
    output := map[string]interface{}{"data": data, "statusCode": statusCode, "error": nil}
    if err != nil {
        output["error"] = map[string]interface{}{"message": err.Error()}
    }
    return output
}
