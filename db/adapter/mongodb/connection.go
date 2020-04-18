package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.zenithar.org/pkg/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Configuration repesents database connection configuration
type Configuration struct {
	AutoMigrate      bool
	ConnectionString string
	DatabaseName     string
	Username         string
	Password         string
}

// Connection provides Wire provider for a MongoDB database connection
func Connection(ctx context.Context, cfg *Configuration) (*mongo.Client, error) {

	log.For(ctx).Info("Trying to connect to MongoDB servers ...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Extract database name from connection string
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.ConnectionString))
	if err != nil {
		log.For(ctx).Error("Unable to connect to MongoDB", zap.Error(err))
		return nil, fmt.Errorf("mongodb: %w", err)
	}

	// Check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.For(ctx).Error("Unable to ping MongoDB database", zap.Error(err))
		return nil, fmt.Errorf("mongodb: %w", err)
	}

	log.For(ctx).Info("Connected to MongoDB.")

	// Return session
	return client, nil
}
