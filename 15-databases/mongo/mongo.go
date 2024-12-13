package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

//https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/query-document/

// DB struct is used to store MongoDB client, a reference to a database, and a collection.
type DB struct {
	client   *mongo.Client
	database *mongo.Database
	coll     *mongo.Collection
}

// Person struct represents the schema of a document in the MongoDB collection.
type Person struct {
	FirstName string   `bson:"first_name"`
	Email     string   `bson:"email"`
	Age       int      `bson:"age"`
	Marks     int      `bson:"marks"`
	Hobbies   []string `bson:"hobbies"`
}

// NewDB function initializes a connection to MongoDB, sets up the database and collection.
func NewDB(collection string) (*DB, error) {
	const uri = "mongodb://root:example@localhost:27017"
	//Create client options and set the connection URI
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Connect to MongoDB using the provided context and options
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	database := client.Database("info")
	fmt.Printf("%+v\n", database)
	if database == nil {
		return nil, fmt.Errorf("failed to get DB: %w", err)
	}
	coll := database.Collection(collection)
	if coll == nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	return &DB{
		client:   client,
		database: database,
		coll:     coll,
	}, nil
}

// Ping method is used to test if the MongoDB client can communicate with the database.
func (db *DB) Ping(ctx context.Context) error {
	return db.client.Ping(ctx, nil)
}

func main() {

	db, err := NewDB("users")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.Ping(ctx)
	if err != nil {
		panic(err)
	}

	db.InsertOne(ctx)
	db.InsertMany(ctx)
	fmt.Println("Connected to MongoDB!")

}

// InsertOne method inserts a single document into the MongoDB collection
func (db *DB) InsertOne(ctx context.Context) {
	u := Person{
		FirstName: "John",
		Email:     "john@email.com",
		Age:       30,
		Hobbies:   []string{"Sports", "Cooking"},
		Marks:     50,
	}

	// Insert the document into the collection
	res, err := db.coll.InsertOne(ctx, u)
	if err != nil {
		log.Println(err)
		return
	}

	//inserted id
	fmt.Println(res.InsertedID)
}

// create a function the inserts multiple users record in one go

// create a function the inserts multiple users record in one go
func (db *DB) InsertMany(ctx context.Context) {
	u1 := Person{
		FirstName: "Jane",
		Email:     "jane@email.com",
		Age:       30,
		Hobbies:   []string{"Sports", "Cooking"},
		Marks:     50,
	}
	u2 := Person{
		FirstName: "Johnson",
		Email:     "johnson@email.com",
		Age:       20,
		Hobbies:   []string{"Hiking", "Surfing"},
		Marks:     60,
	}

	// users := []Person{u1, u2}
	users := []interface{}{u1, u2}
	res, err := db.coll.InsertMany(ctx, users)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Inserted many entries:", res.InsertedIDs)
}
