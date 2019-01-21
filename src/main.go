package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    "github.com/mongodb/mongo-go-driver/mongo"

    "kartiz/auth"
)

type server struct {
    config map[string]interface{}
    auth *auth.Auth
    db *mongo.Database
    router *mux.Router
}

func (s *server) setConfig() {
    file, _ := os.Open("config.json")
    defer closeFile(file)
    if err := json.NewDecoder(file).Decode(&s.config); err != nil {
        log.Fatal(err)
    }
}

func (s *server) setDatabase() {
    host := s.config["mongo"].(map[string]interface{})["host"].(string)
    port := s.config["mongo"].(map[string]interface{})["port"].(string)
    mongoUrl := fmt.Sprintf("mongodb://%s:%s", host, port)
    log.Println(mongoUrl)

    client, err := mongo.NewClient(mongoUrl); if err != nil {
        log.Fatal(err)
    }

    ctx, _ := context.WithTimeout(context.Background(), time.Second)
    err = client.Connect(ctx); if err != nil {
        log.Fatal(err)
    }
    s.db = client.Database("kartiz")

}

func (s *server) setAuth() {
    s.auth = &auth.Auth{}
    s.auth.SetConfig(s.config["auth"].(map[string]interface{}))
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
    // create server object.
    s := server{}

    // set configure to server
    s.setConfig()

    // set database
    s.setDatabase()
    defer closeDB(s.db.Client())

    // set auth
    s.setAuth()

    // set routes
    s.setRoutes()

    log.Println("server starts.")
    log.Fatal(http.ListenAndServe(":8080", s.router))
}
