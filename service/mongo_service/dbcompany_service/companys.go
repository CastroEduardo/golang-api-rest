package dbcompany_service

import (
	"context"
	"log"

	// get an object type
	//"rest-api-golang/src/connectDb"
	//"rest-api-golang/src/models/authinterfaces"

	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/mongo_db"
	"github.com/CastroEduardo/golang-api-rest/pkg/setting"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ClientMongo *mongo.Client
var nameCollection = "companys_sys"

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

func Add(Model authinterfaces.Company_sys) string {
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

func FindToId(id string) authinterfaces.Company_sys {
	settingsCollections()
	var modelSend authinterfaces.Company_sys
	if collection != nil {
		//transform string _id to Object
		docID, _ := primitive.ObjectIDFromHex(id)
		doc := collection.FindOne(context.TODO(), bson.M{"_id": docID})
		doc.Decode(&modelSend)

	}
	return modelSend
}

func GetList() []authinterfaces.Company_sys {
	settingsCollections()
	var list []authinterfaces.Company_sys
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		doc, _ := collection.Find(context.TODO(), bson.M{})
		doc.All(context.Background(), &list)
		doc.Close(context.TODO())
	}

	return list
}
