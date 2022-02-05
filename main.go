package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

type Task struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Text      string             `bson:"text"`
	Completed bool               `bson:"completed"`
}

func init() {
	// Create MongoDB Client and connect
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("tasker").Collection("tasks")
}

func main() {
	app := &cli.App{
		Name:     "tasker",
		Usage:    "A simple CLI program to manage your tasks",
		Commands: []*cli.Command{
			{
				Name: "add",
				Aliases: []string{"a"},
				Usage: "add a task to the list",
				Action: func(c *cli.Context) error {
					str := c.Args().First()
					if str == "" {
						return errors.New("cannot add an empty task")
					}

					task := &Task{
						ID: primitive.NewObjectID(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Text: str,
						Completed: false,
					}

					return createTask(task)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}


func createTask(task *Task) error {
	_, err := collection.InsertOne(ctx, task)
	return err
}

func getAll() ([]*Task, error) {
	// Passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	return filterTasks(filter)
}

func filterTasks(filter interface{}) ([]*Task, error) {
	// A slice of tasks for storing the decoded documents
	var tasks []*Task

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}

	for cur.Next(ctx) {
		var t Task
		err := cur.Decode(&t)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, &t)
	}

	if err := cur.Err(); err != nil {
		return tasks, err
	}

	// Once exhausted, close the cursor
	cur.Close(ctx)

	if len(tasks) == 0 {
		return tasks, mongo.ErrNoDocuments
	}

	return tasks, nil
}