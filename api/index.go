package api

import (
	"context"
	"encoding/json"
	"github.com/rawnly/splash-cli-analytics/lib"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	connectionUrl := os.Getenv("MONGODB_URL")

	encoder := json.NewEncoder(w)
	w.Header().Set("content-type", "application-json")

	if r.Method != "POST" {
		w.WriteHeader(405)

		if err := encoder.Encode(lib.MethodNotAllowed); err != nil {
			log.Fatalln(err.Error())
		}

		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err.Error())

		if err2 := encoder.Encode(err); err2 != nil {
			log.Fatalln(err2.Error())
		}
	}

	var requestBody lib.RequestBody
	if err := json.Unmarshal(data, &requestBody); err != nil {
		w.WriteHeader(400)
		log.Println(err.Error())

		if err2 := encoder.Encode(lib.BadRequest); err2 != nil {
			log.Fatalln(err2.Error())
		}
	}

	//Database connection
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(connectionUrl).SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		w.WriteHeader(500)
		log.Println(err.Error())

		if err2 := encoder.Encode(lib.InternalServerError); err2 != nil {
			log.Fatalln(err2.Error())
		}

		return
	}

	ip := lib.GetIP(r)
	doc := lib.MapToDocument(requestBody, ip)

	collection := client.Database("analytics").Collection("splash-cli")
	document, err := collection.InsertOne(ctx, doc)

	if err != nil {
		log.Println(err.Error())

		if err2 := encoder.Encode(lib.InternalServerError); err != nil {
			log.Fatalln(err2.Error())
		}
		return
	}

	if err := encoder.Encode(document); err != nil {
		log.Fatalln(err.Error())
	}
}
