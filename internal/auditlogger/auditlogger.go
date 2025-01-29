package auditlogger

import (
    "context"
    "encoding/json"
    "log/slog"
    "net/http"

    "go.mongodb.org/mongo-driver/mongo"
)

// MongoWriter is a struct that implements the io.Writer interface to a MongoDB collection
type MongoWriter struct {
    DB *mongo.Database
}

// Values struct holds the context values for logging
type Values struct {
    Uuid   string
    Soeid  string
    Logger *slog.Logger
}

// InitCtx creates a context object to store uuid, soeid and logger pointer
func (mw *MongoWriter) InitCtx(r *http.Request) context.Context {
    opts := slog.HandlerOptions{
        AddSource: true,
    }
    v := Values{
        Uuid:   r.Header.Get("UUID"),
        Soeid:  r.Header.Get("SOEID"),
        Logger: slog.New(slog.NewJSONHandler(mw, &opts)),
    }
    c := context.WithValue(r.Context(), "auditValues", v)
    return c
}

// Write implements the io.Writer interface
func (mw *MongoWriter) Write(p []byte) (n int, err error) {
    c := mw.DB.Collection("log")

    // Parse the incoming byte array into key: value pairs
    var v interface{}
    err = json.Unmarshal(p, &v)
    if err != nil {
        c.InsertOne(context.TODO(), p)
        return len(p), err
    }

    c.InsertOne(context.TODO(), v)

    return len(p), nil
}