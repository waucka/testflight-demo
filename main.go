package main

import (
	api "github.com/waucka/testflight-demo/internal"
	"labix.org/v2/mgo"
	"log"
	"os"
)

var (
	apiConfig *api.Config
)

func main() {
	dburl := "mongodb://" + os.Getenv("MONGO_HOST")
	session, err := mgo.Dial(dburl)
	if err != nil {
		log.Println("Couldn't connect to MongoDB!")
		log.Println(err)
		return
	}

	apiConfig = api.NewConfig(session, "demo")

	router := apiConfig.GetRouter()

	listenAddr := "0.0.0.0:8080"
	router.Run(listenAddr)
}
