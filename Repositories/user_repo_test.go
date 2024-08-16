package repositories

import (
	"context"
	domain "testing_task-manager_api/Domain"
	"testing_task-manager_api/config"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type userRepositorySuite struct{
	suite.Suite
	repository domain.UserRepository
	client *mongo.Client
	collection  *mongo.Collection
	// cleanupExecutor utils.TruncateTableExecutor
}

func (suite *userRepositorySuite) SetupSuite(){

	configs := config.GetConfig()
	client, db := config.ConnectDB(configs)
	collection := db.Collection("test_user")
	repository := NewUserRepository(db, domain.CollectionUser)

	suite.repository = repository
	suite.client = client
	suite.collection = collection
}

// Create User Test
func (suite *userRepositorySuite) TestCreate_Success() {
	username := "new_user"
	password := "password123"
	user := domain.User{
		Username: &username,
		Password: &password,
	}

	// Test user creation
	err := suite.repository.Create(context.TODO(), &user)

	// Assertions
	suite.Nil(err, "err is a nil pointer so no error in this process")
	
	// Verify the user is actually in the database
	var result domain.User
	err = suite.collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&result)
	suite.Nil(err, "err is a nil pointer so no error in this process")
	suite.Equal(user.Username, result.Username, "should be equal between result and user")
}

func (suite *userRepositorySuite) TestCreate_UserExistsByUsername() {
	username :=  "existing_user"
	email := "existing_user@example.com"
	password := "password123"
	user := domain.User{
		Username: &username,
		Email: &email,
		Password: &password,
	}

	// Insert a user with the same username
	err := suite.repository.Create(context.TODO(), &user)
	suite.Nil(err)
 
	// Try to create a new user with the same username
	err = suite.repository.Create(context.TODO(), &user)

	// Assertions
	suite.NotNil(err)
	suite.Equal("this username already exists", err.Error())
}

func (suite *userRepositorySuite) TestCreate_UserExistsByEmail() {
	username :=  "existing_user"
	email := "existing_user@example.com"
	password := "password123"
	user := domain.User{
		Username: &username,
		Email: &email,
		Password: &password,
	}

	// Insert a user with the same email
	err := suite.repository.Create(context.TODO(), &user)
	suite.Nil(err)
 
	// Try to create a new user with the same email
	err = suite.repository.Create(context.TODO(), &user)

	// Assertions
	suite.NotNil(err)
	suite.Equal("this email already exists", err.Error())
}

//Find by username test
func (suite *userRepositorySuite) TestFindByUsername_UserExists() {
	username :=  "existing_user"
	email := "existing_user@example.com"
	password := "password123"
	user := domain.User{
		Username: &username,
		Email:    &email,
		Password: &password,
	}

	err := suite.repository.Create(context.TODO(), &user)
	suite.Nil(err)

	// Test finding the user by username
	foundUser, err := suite.repository.FindByUsername(context.TODO(), *user.Username)

	// Assertions
	suite.Nil(suite.T(), err)
	suite.Equal(user.Username, foundUser.Username)
}

func (suite *userRepositorySuite) TestFindByUsername_UserNotFound() {
	_, err := suite.repository.FindByUsername(context.TODO(), "non_existent_user")
 
	// Assertions
	suite.NotNil(err)
	suite.Equal(mongo.ErrNoDocuments, err)
}  

func (suite *userRepositorySuite) TestUpdate_UserExists() {
	username := "user_to_update"
	email := "user_to_update@example.com"
	user := domain.User{
		Username: &username,
		Email:    &email, 
	}
	err := suite.repository.Create(context.TODO(), &user)
	suite.Nil(err)

	// Find the user and get the user ID
	foundUser, err := suite.repository.FindByUsername(context.TODO(), *user.Username)
	suite.Nil(err)

	// Test updating the user 
	err = suite.repository.Update(context.TODO(), foundUser.User_id)

	// Assertions
	suite.Nil(err)

	// Verify the user was updated 
	var updatedUser domain.User
	err = suite.collection.FindOne(context.TODO(), bson.M{"_id": foundUser.ID}).Decode(&updatedUser)
	suite.Nil(err)
	suite.Equal("ADMIN", updatedUser.User_type)
}

func (suite *userRepositorySuite) TestUpdate_UserNotFound() {
	nonExistentUserID := primitive.NewObjectID().Hex()

	// Test updating a non-existent user
	err := suite.repository.Update(context.TODO(), nonExistentUserID)

	// Assertions
	suite.NotNil(err)
	suite.Equal("USER NOT FOUND", err.Error())
}