package database

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/pedroosz/go-reddit-scrapper/src/entity"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		utils.Fatal("Erro ao conectar ao Mongo DB", err)
	}
	errPing := client.Ping(ctx, nil)
	if errPing != nil {
		utils.Fatal("Erro ao conectar ao Mongo DB", errPing)
	}
	return client
}

func UpdatePost(oldPost *entity.CompletePost, newPost *entity.CompletePost, client *mongo.Client) error {
	collection := client.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("forum"))
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
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		utils.Log("Post " + oldPost.Url + " não foi atualizado")
	} else {
		utils.Log("Post " + oldPost.Url + " foi atualizado")
	}
	return nil
}

func MapPostsOnDatabase(client *mongo.Client, callback func(post *entity.CompletePost)) {
	coll := client.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("forum"))
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		utils.Fatal("Erro ao tentar recuperar os arquivos do banco de dados", err)
	}
	defer cursor.Close(context.Background())
	wg := sync.WaitGroup{}
	for cursor.Next(context.Background()) {
		var post entity.CompletePost
		if err := cursor.Decode(&post); err != nil {
			utils.Log("Erro ao tentar recuperar um post para atualização")
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			callback(&post)
		}()
	}
	wg.Wait()
}

func createCollectionForForum(name string, client *mongo.Client) {
	err := client.Database(os.Getenv("DATABASE_NAME")).CreateCollection(context.Background(), name)
	if err != nil {
		utils.Fatal(fmt.Sprintf("Erro ao tentar criar a collection %s", name), err)
	}
}

func PrepareDatabase() *mongo.Client {
	client := connect()
	forum := os.Getenv("forum")
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
	coll := client.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("forum"))

	existingDoc := coll.FindOne(context.Background(), bson.M{
		"url": post.Url,
	})

	return existingDoc.Err() == nil
}

func InsertPostsOnCollection(post entity.CompletePost, client *mongo.Client) {
	coll := client.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("forum"))
	coll.InsertOne(context.Background(), post)
}
