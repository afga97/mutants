package models

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MutantModel struct {
	Adn      []string `bson:"adn,omitempty"`
	IsMutant bool     `bson:"isMutant,omitempty"`
}

type ResponseMutant struct {
	CountMutant int64   `json:"count_mutant_dna"`
	CountHuman  int64   `json:"count_human_dna"`
	Ratio       float32 `json:"ratio"`
}

func dataBaseConnection() *mongo.Client {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}

var clientDb = dataBaseConnection()
var collection = clientDb.Database("mutants").Collection("adns")

func disconnect() {
	if err := clientDb.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (resp *ResponseMutant) countMutants(canal chan int64) {
	filter := bson.D{primitive.E{Key: "isMutant", Value: true}}
	result, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	canal <- result
}

func (resp *ResponseMutant) countHumans(canal chan int64) {
	filter := bson.D{primitive.E{Key: "isMutant", Value: false}}
	result, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	canal <- result
}

func (resp *ResponseMutant) calculateRatio() {
	if resp.CountMutant != 0 && resp.CountHuman != 0 {
		resp.Ratio = float32(resp.CountMutant) / float32(resp.CountHuman)
	}
}

func (mutModel *MutantModel) saveAdn() {
	if _, err := collection.InsertOne(context.TODO(), mutModel); err != nil {
		panic(err)
	}
}

func GetDataCollection() *ResponseMutant {
	canal := make(chan int64)
	response := &ResponseMutant{}
	go response.countMutants(canal)
	response.CountMutant = <-canal
	go response.countHumans(canal)
	response.CountHuman = <-canal
	disconnect()
	response.calculateRatio()
	return response
}

func SaveDna(adn []string, isMutant bool) {
	mutModel := &MutantModel{
		Adn: adn, IsMutant: isMutant,
	}
	mutModel.saveAdn()
	disconnect()
}
