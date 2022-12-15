package dbgeneric_option_service

import (

	// get an object type
	"context"
	"log"

	"github.com/CastroEduardo/golang-api-rest/models/interfaces_public"
	"github.com/CastroEduardo/golang-api-rest/pkg/mongo_db"
	"github.com/CastroEduardo/golang-api-rest/pkg/setting"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ClientMongo *mongo.Client
var nameCollection = "genericOptions_sys"

// var client *mongo.Client
var collection *mongo.Collection

func settingsCollections() {
	var nameDb = setting.MongoDbSetting.Name
	ClientMongo = mongo_db.ClientMongo
	if ClientMongo != nil {
		//fmt.Println(os.Getenv("TOKEN_HASH")))
		collection = ClientMongo.Database(nameDb).Collection(nameCollection)
	}
}

func Add(genericOption interfaces_public.GenericList) string {
	settingsCollections()

	if collection != nil {
		insertResult, err := collection.InsertOne(context.TODO(), genericOption)
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

func GetListFromIdCompany(idCompany string) []interfaces_public.GenericList {
	settingsCollections()
	var list []interfaces_public.GenericList
	if collection != nil {

		doc, _ := collection.Find(context.TODO(), bson.M{"status": 1, "idCompany": idCompany})
		//doc.Decode(&hero)
		doc.All(context.Background(), &list)
		doc.Close(context.TODO())
		//hide all password this no send to User_sys
		// for i := range list {
		// 	//list[i].Password = ""
		// }
	}

	return list
}

func GetListFromIdCompany_identity(identity string, idCompany string) []interfaces_public.GenericList {
	settingsCollections()
	var list []interfaces_public.GenericList
	if collection != nil {

		doc, _ := collection.Find(context.TODO(), bson.M{"identity": identity, "idCompany": idCompany})
		//doc.Decode(&hero)
		doc.All(context.Background(), &list)
		doc.Close(context.TODO())
		//hide all password this no send to User_sys
		// for i := range list {
		// 	//list[i].Password = ""
		// }
	}

	return list
}

func FindToId(id string) interfaces_public.GenericList {
	//fmt.Println(id)
	settingsCollections()
	var modelSend interfaces_public.GenericList
	if collection != nil {
		//transform string _id to Object
		docID, _ := primitive.ObjectIDFromHex(id)
		doc := collection.FindOne(context.TODO(), bson.M{"_id": docID})
		doc.Decode(&modelSend)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}

	return modelSend
}

func UpdateOne(item interfaces_public.GenericList) bool {
	settingsCollections()

	if collection != nil {
		var id = item.ID
		item.ID = ""
		update2 := bson.M{
			"$set": item,
		}
		// update := bson.M{"$set": bson.M{}}
		docID, _ := primitive.ObjectIDFromHex(id)
		_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": docID}, update2)

		if err != nil {
			log.Fatalln("Error Update User ", err)
			return false
		}
	}

	return true
}
