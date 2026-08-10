// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/joho/godotenv"
// 	"go.mongodb.org/mongo-driver/v2/bson"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// 	"go.mongodb.org/mongo-driver/v2/mongo/options"
// )

// func main() {
// 	// Load .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Get MongoDB URI
// 	uri := os.Getenv("MONGODB_URI")

// 	if uri == "" {
// 		log.Fatal("MONGODB_URI is not set")
// 	}

// 	// Create MongoDB client
// 	clientOption := options.Client().ApplyURI(uri)

// 	client, err := mongo.Connect(clientOption)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Create context with timeout
// 	ctx, cancel := context.WithTimeout(
// 		context.Background(),
// 		10*time.Second,
// 	)
// 	defer cancel()

// 	// Check connection
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Connected to MongoDB Atlas successfully!")

// 	// Disconnect when main finishes
// 	defer client.Disconnect(ctx)

// 	// Select database
// 	database := client.Database("student")

// 	// Select collection
// 	collection := database.Collection("students")

// 	// Create student documents
// 	students := []interface{}{
// 		bson.D{
// 			{Key: "name", Value: "Abebe"},
// 			{Key: "age", Value: 20},
// 			{Key: "department", Value: "Computer Science"},
// 		},
// 		bson.D{
// 			{Key: "name", Value: "Kebede"},
// 			{Key: "age", Value: 21},
// 			{Key: "department", Value: "Software Engineering"},
// 		},
// 		bson.D{
// 			{Key: "name", Value: "Almaz"},
// 			{Key: "age", Value: 19},
// 			{Key: "department", Value: "Information Technology"},
// 		},
// 	}

// 	// Insert students
// 	result, err := collection.InsertMany(ctx, students)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Students inserted successfully!")
// 	fmt.Println("Number of students:", len(result.InsertedIDs))
// }

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

//steps to connect to mongodb database
//1.load godotenv returns error
//2.os.Getenv("MONGO_URI") returns uri
//3.create client option clientOption = options.Client().Apply(uri)
//4.connect to DB (client ,err := mongo.Connnect(clientOption))
//5.create context (ctx,cancel := context.withTimeOut(context.background,10*time.Second)
//6.test connection err =  client.Ping(ctx,nil)
//7.defer client.Disconnect(ctx)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get MongoDB URI
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	// Create MongoDB client
	clientOption := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOption)
	if err != nil {
		log.Fatal(err)
	}

	// Create context
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// Test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas successfully!")

	defer client.Disconnect(ctx)

	//***********this part is CRUD implementation**************
	//*****Steps to select a collection
	//1.db := client.Database(db name)
	//2.cl := db.Collection(col  name)

	// Select database
	database := client.Database("student")

	// Select collection
	collection := database.Collection("students")

	// =========================
	// CREATE
	// =========================
	//steps to Insert
	//1.create the data to be inserted bson.D{{key:,value:}}
	//2.res.err := collection.InsertOne(ctx,data)

	student := bson.D{
		{Key: "name", Value: "Abebe"},
		{Key: "age", Value: 20},
		{Key: "department", Value: "Computer Science"},
	}

	result, err := collection.InsertOne(ctx, student)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Student inserted successfully!")
	fmt.Println("Inserted ID:", result.InsertedID)

	// =========================
	// READ ONE
	// =========================
	//steps to read one 
	//1. var oneData bson.M
	//2.give the bson.m data to filter by
	//3.collect.FindOne(ctx,filter_data).Decode(&oneData)

	var foundStudent bson.M

	filter := bson.M{
		"name": "Abebe",
	}

	err = collection.FindOne(ctx, filter).Decode(&foundStudent)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nStudent found:")
	fmt.Println(foundStudent)

	// =========================
	// READ ALL
	// =========================
	//steps to read all
	//1cursor ,err := collection.Find(ctx,bson{})
	//2create a bson slice
	//3cursor.All(ctx,&student)
	

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	var students []bson.M

	err = cursor.All(ctx, &students)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nAll students:")

	for _, student := range students {
		fmt.Println(student)
	}

	// =========================
	// UPDATE
	// =========================
	//1. create the filter data
	//2.craete bsosn{bson{}} update val
	//


	filter = bson.M{
		"name": "Abebe",
	}

	update := bson.M{
		"$set": bson.M{
			"age":        21,
			"department": "Software Engineering",
		},
	}

	updateResult, err := collection.UpdateOne(
		ctx,
		filter,
		update,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nStudent updated successfully!")
	fmt.Println("Matched:", updateResult.MatchedCount)
	fmt.Println("Modified:", updateResult.ModifiedCount)

	// =========================
	// DELETE
	// =========================

	filter = bson.M{
		"name": "Abebe",
	}

	deleteResult, err := collection.DeleteOne(
		ctx,
		filter,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nStudent deleted successfully!")
	fmt.Println("Deleted:", deleteResult.DeletedCount)
}
