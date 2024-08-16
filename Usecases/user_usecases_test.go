package usecases

import (
	"context"
	"log"
	"testing"
	domain "testing_task-manager_api/Domain"
	repositories "testing_task-manager_api/Repositories"
	"testing_task-manager_api/config"
	"testing_task-manager_api/Domain/mocks"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userUsecaseSuite struct{
	suite.Suite
	client    *mongo.Client
	db        *mongo.Database
	repository *mocks.UserRepository
	usecase domain.UserUsecase
	validate *validator.Validate

}

func (suite *userUsecaseSuite) SetupSuite(){
	configs := config.GetConfig()
	client, db := config.ConnectDB(configs)

	repository := new(mocks.UserRepository)
	usecase := NewUserUsecase(repository, 10)

	suite.client = client
	suite.db = db
	suite.repository = repository
	suite.usecase = usecase
	suite.validate = validator.New()
}

func (suite *userUsecaseSuite) TearDownSuite() {
	// Drop the test database after all tests are run
	if suite.db != nil {
		suite.db.Drop(context.TODO())
	}

	// Close the MongoDB connection
	config.CloseMongoDBConnection(suite.client)
}

func (suite *userUsecaseSuite) TearDownTest() {
	// List of collections you might want to clear after each test
	collections := []string{"users"}

	for _, collection := range collections {
		_, err := suite.db.Collection(collection).DeleteMany(context.TODO(), bson.D{})
		if err != nil {
			log.Fatalf("Failed to clear collection %s: %v", collection, err)
		}
	}
}
func (suite *userUsecaseSuite) SetupTest() {
	// Initialize the repository with the real database
	repository := repositories.NewUserRepository(suite.db, "users")
	suite.usecase = NewUserUsecase(repository, 10*time.Second)
}

// Create user test
func (suite *userUsecaseSuite) TestCreate_Positive() {
	user := domain.User{
		Name:     ptr("John Doe"),
		Username: ptr("johndoe"),
		Password: ptr("strongpassword"),
		Email:    ptr("john.doe@example.com"),
	}
	suite.repository.On("Create", mock.Anything, &user).Return(nil)

	err := suite.usecase.Create(context.TODO(), &user)
	suite.Nil(err)
}

func (suite *userUsecaseSuite) TestUserValidation() {
	tests := []struct {
		name    string
		user    domain.User
		wantErr bool
	}{
		{
			name: "Valid User",
			user: domain.User{
				Name:     ptr("John Doe"),
				Username: ptr("johndoe"),
				Password: ptr("strongpassword"),
				Email:    ptr("john.doe@example.com"),
			},
			wantErr: false,
		},
		{
			name: "Invalid Email",
			user: domain.User{
				Name:     ptr("John Doe"),
				Username: ptr("johndoe"),
				Password: ptr("strongpassword"),
				Email:    ptr("invalid-email"),
			},
			wantErr: true,
		},
		{
			name: "Missing Name",
			user: domain.User{
				Username: ptr("johndoe"),
				Password: ptr("strongpassword"),
				Email:    ptr("john.doe@example.com"),
			},
			wantErr: true,
		},
		{
			name: "Short Password",
			user: domain.User{
				Name:     ptr("John Doe"),
				Username: ptr("johndoe"),
				Password: ptr("123"),
				Email:    ptr("john.doe@example.com"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		if tt.wantErr {
			// Since we expect an error, we shouldn't mock the repository's Create call
			suite.repository.On("Create", mock.Anything, &tt.user).Return(nil).Maybe()
		} else {
			// Expect the repository's Create method to be called with a valid user
			suite.repository.On("Create", mock.Anything, &tt.user).Return(nil)
		}

		// Call the use case's Create method
		err := suite.usecase.Create(context.TODO(), &tt.user)

		// Assert error expectation
		if tt.wantErr {
			suite.Error(err, "Expected validation error, but got none")
		} else {
			suite.NoError(err, "Expected no validation error, but got one")
		}

		// Verify that all expected methods were called on the mock
		// suite.repository.AssertExpectations(suite.T())
		
	}
}

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(userUsecaseSuite))
}

// Helper function to create a pointer to a string
func ptr(s string) *string {
	return &s
}