// MIT License
//
// Copyright (c) 2019 Thibault NORMAND
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package mongodb

import (
	"context"
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
		return nil, err
	}

	log.For(ctx).Info("Connected to MongoDB.")

	// Return session
	return client, nil
}
