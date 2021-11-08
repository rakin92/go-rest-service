package mongo

import (
	"context"
	"time"

	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MDB struct {
	DB *mongo.Database
}

// Close use this method to close database connection
func (r *MDB) Close() {
	logger.Warn("Closing all db connections")
}

// Init initializes the mongo db connection
func Init(c *cfg.MongoDB) (*MDB, error) {
	logger.Info("[Mongo.Init] Connecting to Mongo DB %s", c.Database)
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(c.Host))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = mongoClient.Connect(ctx)
	if err != nil {
		return nil, err
	}
	logger.Info("[Mongo.Init] Connected to Mongo DB %s", c.Database)
	return &MDB{DB: mongoClient.Database(c.Database)}, nil
}
