package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	domain "testing_task-manager_api/Domain"
	"testing_task-manager_api/Domain/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type userControllerSuite struct {
	// we need this to use the suite functionalities from testify
	suite.Suite
	// the mocked version of the usecase
	usecase *mocks.UserUsecase
	// the functionalities we need to test
	controller domain.UserController
	// testing server to be used the handler
	testingServer   *httptest.Server

}

func (suite *userControllerSuite) SetupSuite() {
	// create a mocked version of usecase
	usecase := new(mocks.UserUsecase)
	// inject the usecase to be used by handler
	controller := NewUserController(usecase)

	// create default server using gin, then register all endpoints
	router := gin.Default()
	router.POST("/register", controller.Signup)
	router.POST("/login", controller.Login)
	router.PUT("/promote/:id", controller.PromoteUser)

	// create and run the testing server
	testingServer := httptest.NewServer(router)

	// assign the dependencies we need as the suite properties
	// we need this to run the tests
	suite.testingServer = testingServer
	suite.usecase = usecase
	suite.controller = controller
}

func (suite *userControllerSuite) TearDownSuite() {
	defer suite.testingServer.Close()
}
func (suite *userControllerSuite) TestSignup() {
    tests := []struct {
        name         string
        requestBody  string
        mockResponse error
        expectedCode int
        expectedBody string
    }{
        {
            name: "Valid User",
            requestBody: `{
                "name": "John Doe",
                "username": "johndoe",
                "password": "strongpassword",
                "email": "john.doe@example.com"
            }`,
            mockResponse: nil,
            expectedCode: http.StatusOK,
            expectedBody: `{"message":"user registered successfully"}`,
        },
        {
            name: "Invalid Email",
            requestBody: `{
                "name": "John Doe",
                "username": "johndoe",
                "password": "strongpassword",
                "email": "invalid-email"
            }`,
            mockResponse: fmt.Errorf("invalid email"),
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"message":"invalid email"}`,
        },
        {
            name: "Missing Name",
            requestBody: `{
                "username": "johndoe",
                "password": "strongpassword",
                "email": "john.doe@example.com"
            }`,
            mockResponse: fmt.Errorf("name is required"),
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"message":"name is required"}`,
        },
        {
            name: "Short Password",
            requestBody: `{
                "name": "John Doe",
                "username": "johndoe",
                "password": "123",
                "email": "john.doe@example.com"
            }`,
            mockResponse: fmt.Errorf("password is too short"),
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"message":"password is too short"}`,
        },
    }

    for _, tt := range tests {
		suite.usecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(tt.mockResponse)

		// Make the HTTP request with the JSON string as the request body
		response, err := http.Post(fmt.Sprintf("%s/register", suite.testingServer.URL), "application/json", bytes.NewBufferString(tt.requestBody))
		suite.Require().NoError(err)
		defer response.Body.Close()

		responseBody := domain.SuccessResponse{}
		json.NewDecoder(response.Body).Decode(&responseBody)

		// running assertions to make sure that our method does the correct thing
		suite.Equal(http.StatusCreated, response.StatusCode)
		suite.Equal(responseBody.Message, "User registered successfully")
		suite.usecase.AssertExpectations(suite.T())
}
}
func TestUserController(t *testing.T) {
	suite.Run(t, new(userControllerSuite))
}