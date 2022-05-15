package mongo_db

import (
	"fmt"
	
	"context"
	"log"

	"github.com/CastroEduardo/golang-api-rest/pkg/setting"
	// "github.com/CastroEduardo/golang-api-rest/pkg/util"
	"github.com/CastroEduardo/golang-api-rest/service/logger_service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ClientMongo *mongo.Client

func Setup() bool {

	//var passwordStr = "123456"

	var Host string = setting.MongoDbSetting.Host
	var User = setting.MongoDbSetting.User
	//descripPwd, _ := util.DecryptAES(setting.MongoDbSetting.Password)
	var Password string =setting.MongoDbSetting.Password// descripPwd //util.Descrypt(setting.MongoDbSetting.Password, passwordStr)
	var Name string = setting.MongoDbSetting.Name

	URL_MONGO := "mongodb://" + User + ":" + Password + "@" + Host + "/" + Name + "?clickshield?replicaSet=rs0&authSource=" + Name
	 fmt.Println(URL_MONGO)
	
	//"mongodb://castro:231154Admin11@127.0.0.1:30001,127.0.0.1:30002,127.0.0.1:30003/admin?clickshield?replicaSet=0&connect=mongo_rsdirect"
	 
	
	//mongodb://castro:555555@172.16.18:27017
	urlDb := URL_MONGO //os.Getenv("URL_MONGODB")
	//fmt.Println(urlDb)

	

	// Rest of the code will go here
	// Set client options
	clientOptions := options.Client().ApplyURI(urlDb)
	

	//Context = context.TODO()
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	ClientMongo = client
	if err != nil {

		logger_service.Add(err.Error())
		log.Fatal(err)
	}


	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		logger_service.Add(err.Error())
		log.Fatal(err)

		return false
	}
	logger_service.Add("Success ==> connect to mongo db...")
	//fmt.Println("--#### Connected to MongoDB! #####")

	return true

}

func Status() *mongo.Client {

	// Check the connection
	err := ClientMongo.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		logger_service.Add(err.Error())
		result := Setup()
		if result {
			return ClientMongo
		}
		return ClientMongo
	}

	return ClientMongo

}
