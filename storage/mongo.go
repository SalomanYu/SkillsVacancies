package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive" // В пакет заложены объекты документов MongoDB
	"go.mongodb.org/mongo-driver/mongo"          // Основной пакет mongoDB для создания БД, коллекций и подключения к БД
	"go.mongodb.org/mongo-driver/mongo/options"  // В опциях указываем местоположение БД
)

var collection *mongo.Collection
var ctx = context.TODO()

const NameDatabase = "hh_ru"
const NameCollection = "vacancies"

type Vacancy struct{
	IdMongo		primitive.ObjectID	`bson:"_id"`
	IdVacancy	string				`bson:"vacancy_id"`
	Title		string				`bson:"vacancy_title"`
	MatchTitle	string				`bson:"edwica_profession"`
	FilePath	string				`bson:"csvfile"`
}


func init(){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	checkErr(err)

	err = client.Ping(ctx, nil)
	checkErr(err)

	collection = client.Database(NameDatabase).Collection(NameCollection)
}

func checkErr(err error){
	if err != nil {
		panic(err)
	}
}

func GetVacanciesCsvFile(csvPath string) ([]*Vacancy, error){
	filter := bson.D{primitive.E{
		Key: "csvfile",
		Value: csvPath,
	}}	
	return filterVacancies(filter)
}

func filterVacancies(filter interface{}) ([]*Vacancy, error){
	var vacancies []*Vacancy

	cur, err := collection.Find(ctx, filter)
	defer cur.Close(ctx)
	checkErr(err)

	for cur.Next(ctx){
		var vacancy Vacancy
		err := cur.Decode(&vacancy)
		checkErr(err)
		vacancies = append(vacancies, &vacancy)
	}
	if err := cur.Err(); err != nil{
		return vacancies, err
	}
	if len(vacancies) == 0{
		return vacancies, mongo.ErrNoDocuments
	}
	return vacancies, nil
}