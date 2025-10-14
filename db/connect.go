package db

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/ksiezykm/FerretMate/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect(c model.Connection) error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/?directConnection=true",
		url.QueryEscape(c.Username),
		url.QueryEscape(c.Password),
		c.Host, c.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	Client = client
	return nil
}
