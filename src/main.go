package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/mongodb/mongo-go-driver/mongo"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
)

type server struct {
    router *mux.Router
    db *mongo.Database
}

func closeDB(c *mongo.Client) {
    if err := c.Disconnect(nil); err != nil {
        log.Fatal(err)
    }
}

func closeFile(f *os.File) {
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
}

func getConfig(env string) map[string]interface{} {
    file, _ := os.Open("conf.json")
    defer closeFile(file)

    var data map[string]interface{}
    if err := json.NewDecoder(file).Decode(&data); err != nil {
        log.Fatal(err)
        return nil
    }
    if val, ok := data[env]; ok {
        return val.(map[string]interface{})
    } else {
        return nil
    }
}

func main() {
    var env string
    val, ok := os.LookupEnv("env"); if ok {
       env = val
    } else {
       env = "local"
    }

    conf := getConfig(env)

    host := conf["mongo"].(map[string]interface{})["host"].(string)
    port := conf["mongo"].(map[string]interface{})["port"].(string)
    mongoUrl := fmt.Sprintf("mongodb://%s:%s", host, port)
    log.Println(mongoUrl)

    client, err := mongo.NewClient(mongoUrl); if err != nil {
        log.Fatal(err)
    }

    ctx, _ := context.WithTimeout(context.Background(), time.Second)
    err = client.Connect(ctx); if err != nil {
        log.Fatal(err)
    }
    defer closeDB(client)

    db := client.Database("kartiz")

    s := server{mux.NewRouter(), db}

    // apply routes
    s.routes()

    log.Println("server starts.")
    log.Fatal(http.ListenAndServe(":8080", s.router))
}
