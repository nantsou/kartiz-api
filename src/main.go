package main

import (
    "context"
    "encoding/json"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/mongo"

    "kartiz/auth"
)

type server struct {
    config map[string]interface{}
    auth *auth.Auth
    db *mongo.Database
    router *mux.Router
}

func getConfig() map[string]interface{} {
    var config map[string]interface{}

    file, _ := os.Open("config.json")
    defer closeFile(file)

    if err := json.NewDecoder(file).Decode(&config); err != nil {
        log.Fatal(err)
    }

    return config
}

func (s *server) setDatabase(config map[string]interface{}) {
    host := config["host"].(string)
    port := config["port"].(string)
    mongoUri := fmt.Sprintf("mongodb://%s:%s", host, port)
    log.Println(mongoUri)

    client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri)); if err != nil {
        log.Fatal(err)
    }

    ctx, _ := context.WithTimeout(context.Background(), time.Second)
    err = client.Connect(ctx); if err != nil {
        log.Fatal(err)
    }
    s.db = client.Database("kartiz")

}

func (s *server) setAuth(config map[string]interface{}) {
    s.auth = &auth.Auth{}
    s.auth.SetConfig(config)
    s.auth.SetDB(s.db.Collection("blacklist"))
}

func (s *server) setRoutes() {
    s.router = mux.NewRouter()
    s.routes()
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

func main() {
    // get configure
    config := getConfig()

    // create server object.
    s := server{}

    // set database
    s.setDatabase(config["db"].(map[string]interface{}))
    defer closeDB(s.db.Client())

    // set auth
    s.setAuth(config["auth"].(map[string]interface{}))

    // set routes
    s.setRoutes()

    log.Println("server starts.")
    log.Fatal(http.ListenAndServe(":8080", s.router))
}
