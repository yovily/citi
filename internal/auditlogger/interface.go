package auditlogger

import (
    "context"
    "net/http"
    "go.mongodb.org/mongo-driver/mongo"
)

type AuditLoggerInterface interface {
    InitCtx(req *http.Request) context.Context
    Write(p []byte) (n int, err error)
}

func New(conn *mongo.Database) AuditLoggerInterface {
    return &MongoWriter{DB: conn}
}