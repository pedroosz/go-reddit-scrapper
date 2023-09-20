package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pedroosz/go-reddit-scrapper/src/database"
	"github.com/pedroosz/go-reddit-scrapper/src/extractors"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
)

var wg sync.WaitGroup
var client *mongo.Client

func getForum() string {
	forum := os.Getenv("forum")
	if forum == "" {
		utils.Fatal("Forum deve ser forencedio", nil)
	}
	return forum
}

func getMinutes() int {
	minutes := os.Getenv("interval")
	if minutes == "" {
		utils.Fatal("Minutes deve ser fornecido", nil)
	}
	intMinutes, err := strconv.Atoi(minutes)

	if err != nil {
		utils.Fatal("Erro ao parsear minutos para inteiro", err)
	}
	return intMinutes
}

func config() {
	err := godotenv.Load(".env")
	if err != nil {
		utils.Log("Não foi possível carregar o arquivo de configuração ambiente")
		os.Exit(-1)
	}
}

func getParams() (int, string) {

	forum := getForum()
	intMinutes := getMinutes()
	return intMinutes, forum
}

func main() {
	config()
	client = database.PrepareDatabase()
	minutes, forum := getParams()
	interval := time.Duration(minutes) * time.Minute
	for {
		startTime := time.Now()
		extractors.ExtractPagesOfForum(forum, &wg, client)
		wg.Wait()
		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)
		utils.Log(fmt.Sprintf("Extração Concluída! Levou %s", elapsedTime.String()))
		utils.Log(fmt.Sprintf("Próxima Extração em %s", interval.String()))
		time.Sleep(interval)
	}
}
