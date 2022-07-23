package dblogs_service

import (
	"context"
	"log"

	// get an object type
	"time"

	"github.com/CastroEduardo/golang-api-rest/models/interfaces_public"
	"github.com/CastroEduardo/golang-api-rest/pkg/mongo_db"
	"github.com/CastroEduardo/golang-api-rest/pkg/setting"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ClientMongo *mongo.Client
var nameCollection = "logs_sys"

//var client *mongo.Client
var collection *mongo.Collection

func settingsCollections() {
	var nameDb = setting.MongoDbSetting.Name
	ClientMongo = mongo_db.ClientMongo
	if ClientMongo != nil {
		//fmt.Println(os.Getenv("TOKEN_HASH")))
		collection = ClientMongo.Database(nameDb).Collection(nameCollection)
	}
}

func Add(level int, logSytem string, idAssociated string, ipRequest string) string {
	settingsCollections()

	newLog := interfaces_public.LogSystem{
		Log:          logSytem,
		Level:        level,
		IdAssociated: idAssociated,
		Status:       1,
		Date:         time.Now(),
		Ip:           ipRequest,
	}

	if collection != nil {
		insertResult, err := collection.InsertOne(context.TODO(), newLog)
		if err != nil {
			log.Fatalln("Error on inserting new Hero", err)
			return ""
		}
		//get id Add
		if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
			return oid.Hex()
		}
	}

	return ""
}

func GetList() []interfaces_public.LogSystem {
	settingsCollections()
	var list []interfaces_public.LogSystem
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		doc, _ := collection.Find(context.TODO(), bson.M{})
		doc.All(context.Background(), &list)
		doc.Close(context.TODO())
	}

	return list
}

func init() {
	//fmt.Println("init.....")
}

func Demo() string {

	return "Envio desde Servicio"
}