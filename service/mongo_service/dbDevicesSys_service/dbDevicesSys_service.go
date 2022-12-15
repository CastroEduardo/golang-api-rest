package dbDevicesSys_service

import (

	// get an object type
	//"rest-api-golang/src/connectDb"
	//"rest-api-golang/src/models/authinterfaces"

	"context"
	"log"

	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/mongo_db"
	"github.com/CastroEduardo/golang-api-rest/pkg/setting"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ClientMongo *mongo.Client
var nameCollection = "devices_sys"

//var client *mongo.Client
var collection *mongo.Collection
func settingsCollections() bool {

	var nameDb = setting.MongoDbSetting.Name
	ClientMongo = mongo_db.ClientMongo
	//fmt.Println("==> " + nameDb)

	if ClientMongo != nil {
		//fmt.Println(os.Getenv("TOKEN_HASH")))
		collection = ClientMongo.Database(nameDb).Collection(nameCollection)
		return true
	}

	return false
}

func Add(Model authinterfaces.DevicesSys_sys) string {
	result := settingsCollections()
	if result {
		if collection != nil {
			insertResult, err := collection.InsertOne(context.TODO(), Model)
			if err != nil {
				log.Fatalln("Error on inserting new "+nameCollection, err)
				return ""
			}
			//get id Add
			if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
				return oid.Hex()
			}
		}
	}
	return ""

}

func FindToIdDevice(idDevice string) authinterfaces.DevicesSys_sys {
	settingsCollections()
	var modelSend authinterfaces.DevicesSys_sys
	if collection != nil {
		
		doc := collection.FindOne(context.TODO(), bson.M{"key": idDevice})
		doc.Decode(&modelSend)

	}
	return modelSend
}
