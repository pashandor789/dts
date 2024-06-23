package storage

import (
	"context"
	"dts/pkg/log"
	"dts/tabasco/storage/types"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Storage struct {
	client *mongo.Client
	logger log.Logger
}

type Config struct {
	URI string `yaml:"uri"`
}

/*

   schema: (collection)
       task:
           taskId: string
           testId: ui64
           testType: input/output (string)
           content: string
       task_meta:
           taskId: string
           batches: vector<ui64>
           batchSize: ui64
       builds:
           id: ui64
           executeScript: string
           initScript: string

*/

func (s Storage) GetTests(taskId string) ([]types.Test, error) {
	coll := s.client.Database("tabasco").Collection("tests")

	filter := bson.M{"task_id": taskId}

	c, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("can't get tests: %w", err)
	}

	var tests []types.Test
	err = c.All(context.TODO(), &tests)
	if err != nil {
		return nil, fmt.Errorf("can't get tests: %w", err)
	}

	return tests, nil
}

func (s Storage) PutTests(tests []types.Test) error {
	coll := s.client.Database("tabasco").Collection("tests")

	models := make([]mongo.WriteModel, len(tests))
	for i, test := range tests {
		filter := bson.M{"id": test.Id, "type": test.Type, "task_id": test.TaskId}
		update := bson.M{"$set": test}
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		models[i] = model
	}

	_, err := coll.BulkWrite(context.TODO(), models)
	if err != nil {
		return fmt.Errorf("can't put or update tests: %w", err)
	}

	return nil
}

func (s Storage) GetBuilds() ([]types.Build, error) {
	coll := s.client.Database("tabasco").Collection("builds")

	c, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("can't find builds collection: %w", err)
	}

	var builds []types.Build
	err = c.All(context.TODO(), &builds)
	if err != nil {
		return nil, fmt.Errorf("can't get builds: %w", err)
	}

	return builds, nil
}

func (s Storage) PutBuild(build *types.Build) error {
	coll := s.client.Database("tabasco").Collection("builds")

	filter := bson.M{"id": build.Id}
	update := bson.M{"$set": build}

	_, err := coll.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("can't put build: %w", err)
	}

	return nil
}

func New(cfg Config, logger log.Logger) (*Storage, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("can't open connection to mongoDB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("no ping response from mongoDB: %w", err)
	}

	s := &Storage{client: client, logger: logger}
	return s, nil
}
