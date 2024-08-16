package infrastructure

import (
	domain "testing_task-manager_api/Domain"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type validateTokenTestSuite struct {
	suite.Suite
	validRefreshToken string
	validToken   string
	expiredToken string
}

func (suite *validateTokenTestSuite) SetupSuite(){
	var err error
	
	user := domain.User{
		User_id: primitive.NewObjectID().Hex(),
		Username: ptr("new user"),
		Email: ptr("email@example.com"),
		User_type: "ADMIN",
	}
	
	suite.validToken, suite.validRefreshToken, err = GenerateJWTToken(user.User_id, *user.Username, *user.Email, user.User_type)
}

func (suite *validateTokenTestSuite)  TestValidateToken_Valid() {
	claims, err := ValidateToken(suite.validToken)
	suite.NoError(err, "Expected no error with valid token")
	suite.NotNil(claims, "Expected claims to be non-nil with valid token")
}

func ptr(s string) *string {
	return &s
}