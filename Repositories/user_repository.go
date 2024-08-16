package repositories

import (
	"context"
	"errors"
	"fmt"
	domain "testing_task-manager_api/Domain"
	infrastructure "testing_task-manager_api/Infrastructure"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

type userRepository struct {
	database   *mongo.Database
	collection string
}

func (ur *userRepository) Create(c context.Context, user *domain.User) error {
	userCollection := ur.database.Collection(ur.collection)

	if validationErr := validate.Struct(user) ; validationErr != nil{
		return validationErr
	}

	count, err := userCollection.CountDocuments(c, bson.M{"username":user.Username})
	
	if err != nil {
		// log.Panic(err)
		return err
	}
	if count > 0{
		return errors.New("this username already exists")
	}

	count, err = userCollection.CountDocuments(c, bson.M{"email":user.Email})
	// defer cancel()
	if err!= nil {
		// log.Panic(err)
		return err
	}

	if count > 0{
		return errors.New("this email already exists. ")
	}

	password := infrastructure.HashPassword(*user.Password)
	user.Password = &password
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex() 

	countUsers, err := userCollection.CountDocuments(c, bson.M{})
	if err!= nil {
		// log.Panic(err)
		return err
	}
	if countUsers == 0{
		user.User_type = "ADMIN"
	}else{
	user.User_type = "USER"
	}

	_, insertionErr := userCollection.InsertOne(c, user)
	if insertionErr != nil{
		return insertionErr
	}

	return nil
}

func (ur *userRepository) FindByUsername(c context.Context, username string) (domain.User, error) {
	var foundUser domain.User
	userCollection := ur.database.Collection(ur.collection)

	filter := bson.M{"username":username}
	result := userCollection.FindOne(c, filter).Decode(&foundUser)
	if result != nil {
		return domain.User{}, result
	}
	return foundUser, result
}

func (ur *userRepository) Update(c context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	fmt.Println(objID, err)
	userCollection := ur.database.Collection(ur.collection)
    if err != nil {
        return err
    }

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_type", Value: "ADMIN"},
		}},
	}
	updateResult, msg := userCollection.UpdateOne(context.TODO(), filter, update)
	if msg != nil {
		return msg
	}
	if updateResult.MatchedCount == 0{
		return errors.New("USER NOT FOUND")
	}
	return nil
}

func NewUserRepository(db *mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}
