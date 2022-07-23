package dbusers_service

import (
	"context"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	// get an object type
	//"rest-api-golang/src/connectDb"
	//"rest-api-golang/src/models/authinterfaces"

	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/mongo_db"
	"github.com/CastroEduardo/golang-api-rest/pkg/setting"
	"github.com/CastroEduardo/golang-api-rest/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ClientMongo *mongo.Client
var nameCollection = "users_sys"

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

func init() {
	//fmt.Println("init Service1")
}

func GetListFromIdCompany(idCompany string) []authinterfaces.User_sys {
	settingsCollections()
	var list []authinterfaces.User_sys
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		doc, _ := collection.Find(context.TODO(), bson.M{"public": 1, "idcompany": idCompany})
		//doc.Decode(&hero)
		doc.All(context.Background(), &list)
		doc.Close(context.TODO())
		//hide all password this no send to User_sys
		for i := range list {
			list[i].Password = ""
		}
	}

	return list
}

func CheckUserPasswordForEmail(User_sys string, password string) authinterfaces.User_sys {
	settingsCollections()
	var modelSend authinterfaces.User_sys
	if collection != nil {

		getUser := FindToEmail(User_sys)

		if getUser.NickName != "" {

			result := util.Descrypt(getUser.Password, []byte(password))
			//fmt.Println(getUser.Password)
			if result {
				getUser.Password = ""
				return getUser
			} else {
				return modelSend
			}

		}

	}

	return modelSend
}

func CheckUserPasswordForUser(User_sys string, password string) authinterfaces.User_sys {
	settingsCollections()
	var modelSend authinterfaces.User_sys
	if collection != nil {

		getUser := FindToNickName(User_sys)

		if getUser.NickName != "" {
			result := util.Descrypt(getUser.Password, []byte(password))
			//fmt.Println(getUser.Password)
			if result {
				getUser.Password = ""
				return getUser
			} else {
				return modelSend
			}

		}

	}

	return modelSend
}

func FindToNickName(nickName string) authinterfaces.User_sys {
	settingsCollections()
	var modelSend authinterfaces.User_sys
	if collection != nil {
		//transform string _id to Object
		// docID, _ := primitive.ObjectIDFromHex(id)
		nickName = strings.ToLower(nickName)
		doc := collection.FindOne(context.TODO(), bson.M{"nickname": nickName})
		doc.Decode(&modelSend)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}

	return modelSend
}

func FindToEmail(email string) authinterfaces.User_sys {
	settingsCollections()
	var modelSend authinterfaces.User_sys
	if collection != nil {
		//transform string _id to Object
		// docID, _ := primitive.ObjectIDFromHex(id)
		email = strings.ToLower(email)
		doc := collection.FindOne(context.TODO(), bson.M{"email": email})
		doc.Decode(&modelSend)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}

	return modelSend
}

func FindToId(id string) authinterfaces.User_sys {
	//fmt.Println(id)
	settingsCollections()
	var modelSend authinterfaces.User_sys
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

func Add(User_sys authinterfaces.User_sys) string {
	settingsCollections()

	if collection != nil {
		insertResult, err := collection.InsertOne(context.TODO(), User_sys)
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

func DeleteToId(id string) bool {
	settingsCollections()

	//var modelSend authinterfaces.User_sys
	if collection != nil {
		//transform string _id to Object
		docID, _ := primitive.ObjectIDFromHex(id)
		deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": docID})

		if err != nil {
			log.Fatalln("Error Finding user", err)
			return false
		}

		if deleteResult.DeletedCount > 0 {
			return true
		}

	}

	return false
}

func UpdateLastLogin(UserUpdate authinterfaces.User_sys) bool {
	settingsCollections()

	//var modelSend authinterfaces.User_sys
	if collection != nil {

		var id = UserUpdate.ID
		UserUpdate.ID = ""
		update2 := bson.M{
			"$set": UserUpdate,
		}
		// update := bson.M{"$set": bson.M{}}
		docID, _ := primitive.ObjectIDFromHex(id)
		_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": docID}, update2)

		if err != nil {
			log.Fatalln("Error on inserting new Hero", err)
			return false
		}
		return true
	}

	return false
}

func UpdateOne1(UserUpdate authinterfaces.User_sys) bool {
	settingsCollections()

	//var modelSend authinterfaces.User_sys
	if collection != nil {

		var id = UserUpdate.ID
		UserUpdate.ID = ""

		// fields := bson.M{}
		// fields["nickname"] = UserUpdate.NickName
		// fields["nickname2"] = "sd"

		update2 := bson.M{
			"$set": UserUpdate,
		}

		// update := bson.M{"$set": bson.M{}}
		docID, _ := primitive.ObjectIDFromHex(id)
		_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": docID}, update2)

		if err != nil {
			log.Fatalln("Error on inserting new Hero", err)
			return false
		}

		return true

		// if deleteResult.DeletedCount > 0 {
		// 	return true
		// }

	}

	return false
}

func UpdateOne(UserUpdate authinterfaces.User_sys) bool {
	settingsCollections()

	//var modelSend authinterfaces.User_sys
	if collection != nil {

		var id = UserUpdate.ID
		UserUpdate.ID = ""

		// fields := bson.M{}
		// fields["nickname"] = UserUpdate.NickName
		// fields["nickname2"] = "sd"

		update2 := bson.M{
			"$set": UserUpdate,
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

func IsAccount(nickName string) bool {
	var isAccount = false
	resultNickName := FindToNickName(nickName)
	if resultNickName.NickName == "" {
		//search email
		resultEmail := FindToEmail(nickName)
		if resultEmail.NickName != "" {
			//fmt.Println("** IS EMAIL")
			isAccount = true
		}
	} else {
		isAccount = true
	}

	return isAccount
}

func structToMap(i interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// You ca use tags here...
		// tag := typ.Field(i).Tag.Get("tagname")
		// Convert each type into a string for the url.Values string map
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		values.Set(typ.Field(i).Name, v)
	}
	return
}

func Demo() string {

	return "Envio desde Servicio"
}
