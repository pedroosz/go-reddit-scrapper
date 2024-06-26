package database

import (
	"context"
	"fmt"
	"os"

	"github.com/pedroosz/go-reddit-scrapper/src/entity"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect() *mongo.Client {
	poolSize := uint64(10)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")), &options.ClientOptions{
		MaxPoolSize: &poolSize,
	})
	if err != nil {
		utils.Fatal("Erro ao conectar ao Mongo DB", err)
	}
	errPing := client.Ping(context.Background(), nil)
	if errPing != nil {
		utils.Fatal("Erro ao conectar ao Mongo DB", errPing)
	}
	return client
}

func UpdatePost(oldPost *entity.CompletePost, newPost *entity.CompletePost, client *mongo.Client) error {
	collection := client.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("FORUM"))
	filter := bson.M{"url": oldPost.Url}
	update := bson.M{
		"$set": bson.M{
			"url":          newPost.Url,
			"rawText":      newPost.RawText,
			"title":        newPost.Title,
			"text":         newPost.Text,
			"up":           newPost.Up,
			"creationDate": newPost.CreationDate,
		},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func MapPostsOnDatabase(client *mongo.Client, callback func(post *entity.CompletePost)) {
	coll := client.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("FORUM"))
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		utils.Fatal("Erro ao tentar recuperar os arquivos do banco de dados", err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var post entity.CompletePost
		if err := cursor.Decode(&post); err != nil {
			utils.Log("Erro ao tentar recuperar um post para atualização")
			continue
		}
		callback(&post)
	}
}

func createCollectionForForum(name string, client *mongo.Client) {
	err := client.Database(os.Getenv("DATABASE_NAME")).CreateCollection(context.Background(), name)
	if err != nil {
		utils.Fatal(fmt.Sprintf("Erro ao tentar criar a collection %s", name), err)
	}
}

func PrepareDatabase() *mongo.Client {
	client := connect()
	forum := os.Getenv("FORUM")
	collection, err := client.Database(os.Getenv("DATABASE_NAME")).ListCollectionNames(context.Background(), bson.M{
		"name": forum,
	})
	if err != nil {
		utils.Fatal("Erro ao tentar recuperar a collection do Forum", err)
	}
	if len(collection) > 0 {
		return client
	}
	createCollectionForForum(forum, client)
	return client
}

func PostExistsOnCollection(post entity.Post, client *mongo.Client) bool {
	coll := client.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("FORUM"))

	existingDoc := coll.FindOne(context.Background(), bson.M{
		"url": post.Url,
	})

	return existingDoc.Err() == nil
}

func InsertPostsOnCollection(post entity.CompletePost, client *mongo.Client) {
	coll := client.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("FORUM"))
	coll.InsertOne(context.Background(), post)
}
