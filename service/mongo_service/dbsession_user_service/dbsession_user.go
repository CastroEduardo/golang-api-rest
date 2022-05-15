package dbsession_user_service

import (
	"context"
	"fmt"
	"log"
	"os" // get an object type

	"strings"
	"time"

	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/mongo_db"
	"github.com/CastroEduardo/golang-api-rest/pkg/setting"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbcompany_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbprivilege_rol_user_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbrol_user_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbusers_service"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ClientMongo *mongo.Client
var nameCollection = "session_users"

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

func Add(Model authinterfaces.SessionUser) string {
	settingsCollections()

	LogoutSessionToIdUser(Model.IdUser) //Disable Active Session

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
	return ""
}

func LogoutSessionToIdUser(idUser string) bool {
	var sessionUser authinterfaces.SessionUser
	sessionUser = FindToActiveIdUser(idUser)
	if sessionUser.IdCompany != "" {
		if sessionUser.Active {
			sessionUser.Active = false
			sessionUser.DateLogout = time.Now()
			UpdateOne(sessionUser)
		}
	}

	return true
}

func UpdateOne(ModelUpdate authinterfaces.SessionUser) bool {
	settingsCollections()

	//var modelSend authinterfaces.User
	if collection != nil {
		var id = ModelUpdate.ID
		ModelUpdate.ID = ""
		update := bson.M{
			"$set": ModelUpdate,
		}
		// update := bson.M{"$set": bson.M{}}
		docID, _ := primitive.ObjectIDFromHex(id)
		_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": docID}, update)

		if err != nil {
			log.Fatalln("Error on inserting new "+nameCollection, err)
			return false
		}
	}

	return true
}

func FindToId(id string) authinterfaces.SessionUser {
	settingsCollections()
	var modelSend authinterfaces.SessionUser
	if collection != nil {
		//transform string _id to Object
		docID, _ := primitive.ObjectIDFromHex(id)
		doc := collection.FindOne(context.TODO(), bson.M{"_id": docID})
		doc.Decode(&modelSend)

	}
	return modelSend
}
func FindToToken(token string) authinterfaces.SessionUser {
	settingsCollections()
	var modelSend authinterfaces.SessionUser
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex(token)
		doc := collection.FindOne(context.TODO(), bson.M{"token": token})
		doc.Decode(&modelSend)

	}
	return modelSend
}

func FindToActiveIdUser(idUser string) authinterfaces.SessionUser {
	settingsCollections()
	var modelSend authinterfaces.SessionUser
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex(idUser)
		doc := collection.FindOne(context.TODO(), bson.M{"iduser": idUser, "active": true})
		doc.Decode(&modelSend)
		// if err != nil {
		// 	log.Fatal(err)
		// 	return modelSend
		// }
	}
	return modelSend
}

func GetList() []authinterfaces.SessionUser {
	settingsCollections()
	var list []authinterfaces.SessionUser
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		doc, _ := collection.Find(context.TODO(), bson.M{})
		//doc.Decode(&hero)
		var hero authinterfaces.SessionUser
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

func GetIdSessionToToken(tokenHeader string) string {
	splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	tokenPart := splitted[1]                    //Grab the token part, what we are truly interested in
	tk := &authinterfaces.Token{}

	//fmt.Println(tokenPart)
	_, errt := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_HASH")), nil
	})
	if errt != nil {
		return ""
	}
	return tk.IdSession

}

func GetClaimForToken(token string) authinterfaces.ClaimSession {

	SendModel := authinterfaces.ClaimSession{}

	// if tokenHeader == "" {
	// 	return SendModel
	// }

	settingsCollections()

	// splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	// tokenPart := splitted[1]                    //Grab the token part, what we are truly interested in
	// tk := &authinterfaces.Token{}

	// //fmt.Println(tokenPart)
	// _, errt := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte(os.Getenv("TOKEN_HASH")), nil
	// })

	// if errt != nil {
	// 	return SendModel
	// }

	if token == "" {
		return SendModel
	}

	var dataSession authinterfaces.SessionUser
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex(token)
		doc := collection.FindOne(context.TODO(), bson.M{"token": token})
		doc.Decode(&dataSession)
	}

	var dataCompany authinterfaces.Company
	dataCompany = dbcompany_service.FindToId(dataSession.IdCompany)
	SendModel.Company = dataCompany

	var dataUser authinterfaces.User
	dataUser = dbusers_service.FindToId(dataSession.IdUser)
	dataUser.Password = ""
	SendModel.User = dataUser

	var dataPrivilegesRol authinterfaces.UserPrivileges
	dataPrivilegesRol = dbprivilege_rol_user_service.FindToIdRol(dataUser.IdRol)
	SendModel.UserPrivileges = dataPrivilegesRol

	var dataRolUser authinterfaces.RolUser
	dataRolUser = dbrol_user_service.FindToIdRol(dataUser.IdRol)
	SendModel.RolUser = dataRolUser

	//fmt.Println(dataPrivilegesRol)

	return SendModel
}
