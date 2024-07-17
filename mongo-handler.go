package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func db() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	//filter := bson.D{{"blooddrives", bson.D{{"$exists", false}}}}
	centerCollection := client.Database("DonationCenters").Collection("centers")

	userCollection := client.Database("DonationCenters").Collection("users")

	cursor, _ := centerCollection.Find(context.TODO(), bson.D{})
	if err = cursor.All(context.TODO(), &c.Centers); err != nil {
		log.Fatal(err)

	}

	for _, v := range c.Centers {
		fmt.Println(v, "\n")
	}

	cur, e := userCollection.Find(context.TODO(), bson.D{})
	if e = cur.All(context.TODO(), &AS.accounts); err != nil {
		log.Fatal(e)

	}
	fmt.Println("UsersDB: ", AS)
}

func update(id int, center DonorCenter) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	//filter := bson.D{{"blooddrives", bson.D{{"$exists", false}}}}
	centerCollection := client.Database("DonationCenters").Collection("centers")
	// newCenter := DonorCenter{Name: "Bruce Wayne Donation Center", Addr: Address{Street: "11223-44556-0000", CityState: "Gotham,902654", Zip: "", Region: "Africa"}, Num: 23, Recruiters: []Recruiter{Recruiter{"Bruce", "Wayne", "", "dontworryBoutIt@myDickBiyatch!!.com"}}}
	// centerCollection.InsertOne(context.Background(), newCenter)

	opts := options.FindOneAndReplace().SetUpsert(true)
	filter := bson.D{{"num", center.Num}}

	var replacedDocument bson.M
	err = centerCollection.FindOneAndReplace(
		context.TODO(),
		filter,
		center,
		opts,
	).Decode(&replacedDocument)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("replaced document %v", replacedDocument)
	return
}

func updateUsers(id int, user Account) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	//filter := bson.D{{"blooddrives", bson.D{{"$exists", false}}}}
	userCollection := client.Database("DonationCenters").Collection("users")
	// newCenter := DonorCenter{Name: "Bruce Wayne Donation Center", Addr: Address{Street: "11223-44556-0000", CityState: "Gotham,902654", Zip: "", Region: "Africa"}, Num: 23, Recruiters: []Recruiter{Recruiter{"Bruce", "Wayne", "", "dontworryBoutIt@myDickBiyatch!!.com"}}}
	// centerCollection.InsertOne(context.Background(), newCenter)

	opts := options.FindOneAndReplace().SetUpsert(true)
	filter := bson.D{{"name", user.Name}}

	var replacedDocument bson.M
	err = userCollection.FindOneAndReplace(
		context.TODO(),
		filter,
		user,
		opts,
	).Decode(&replacedDocument)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("replaced document %v", replacedDocument)
	return
}
