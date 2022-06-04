package dbdepartamentuser_service

import (
	"context"
	"fmt"
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
var nameCollection = "departaments_users_sys"

//var client *mongo.Client
var collection *mongo.Collection

func settingsCollections() bool {

	var nameDb = setting.MongoDbSetting.Name
	ClientMongo = mongo_db.ClientMongo
	fmt.Println("==> " + nameDb)

	if ClientMongo != nil {
		//fmt.Println(os.Getenv("TOKEN_HASH")))
		collection = ClientMongo.Database(nameDb).Collection(nameCollection)
		return true
	}

	return false
}

func Add(Model authinterfaces.DepartamentUserSys) string {
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

func FindToId(id string) authinterfaces.DepartamentUserSys {
	settingsCollections()
	var modelSend authinterfaces.DepartamentUserSys
	if collection != nil {
		//transform string _id to Object
		docID, _ := primitive.ObjectIDFromHex(id)
		doc := collection.FindOne(context.TODO(), bson.M{"_id": docID})
		doc.Decode(&modelSend)

	}
	return modelSend
}

func GetListFromIdCompany(id string) []authinterfaces.DepartamentUserSys {
	settingsCollections()
	var list []authinterfaces.DepartamentUserSys
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		doc, _ := collection.Find(context.TODO(), bson.M{"idcompany": id})
		//doc.Decode(&hero)
		var hero authinterfaces.DepartamentUserSys
		for doc.Next(context.TODO()) {
			// Declare a result BSON object
			//var result bson.M
			err := doc.Decode(&hero)
			if err != nil {
				fmt.Println(hero)
			}
			list = append(list, hero)
		}
	}

	return list
}

func UpdateOne(model authinterfaces.DepartamentUserSys) bool {
	settingsCollections()

	//var modelSend authinterfaces.User
	if collection != nil {

		var id = model.ID
		model.ID = ""
		update2 := bson.M{
			"$set": model,
		}
		// update := bson.M{"$set": bson.M{}}
		docID, _ := primitive.ObjectIDFromHex(id)
		_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": docID}, update2)

		if err != nil {
			log.Fatalln("Error on inserting new departament", err)
			return false
		}

		return true

	}

	return false
}

func Delete(id string) bool {
	settingsCollections()

	//var modelSend authinterfaces.User
	if collection != nil {
		//transform string _id to Object
		docID, _ := primitive.ObjectIDFromHex(id)
		deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": docID})

		if err != nil {
			log.Fatalln("Error on inserting new Departament", err)
			return false
		}

		if deleteResult.DeletedCount > 0 {
			return true
		}

	}

	return false
}
