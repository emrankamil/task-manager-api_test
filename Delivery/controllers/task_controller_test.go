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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskControllerSuite struct {
	// we need this to use the suite functionalities from testify
	suite.Suite
	// the mocked version of the usecase
	usecase *mocks.TaskUsecase
	// the functionalities we need to test
	controller domain.TaskController
	// testing server to be used the handler
	testingServer   *httptest.Server

}

func (suite *taskControllerSuite) SetupSuite() {
	// create a mocked version of usecase
	usecase := new(mocks.TaskUsecase)
	// inject the usecase to be used by handler
	controller := NewTaskController(usecase)

	// create default server using gin, then register all endpoints
	router := gin.Default()
	router.GET("/tasks", controller.FetchAll)
	router.GET("/tasks/:id", controller.FetchByTaskID)
	router.PUT("/tasks/:id", controller.Update)
	router.DELETE("/tasks/:id", controller.Delete)
	router.POST("/tasks", controller.Create)

	// create and run the testing server
	testingServer := httptest.NewServer(router)

	// assign the dependencies we need as the suite properties
	// we need this to run the tests
	suite.testingServer = testingServer
	suite.usecase = usecase
	suite.controller = controller
}

func (suite *taskControllerSuite) TearDownSuite() {
	defer suite.testingServer.Close()
}

func (suite *taskControllerSuite) TestCreateTask_Positive() {
	task := domain.Task{
		Title: "new title",
		Description: "New description",
		Status: "Pending",
	}

	suite.usecase.On("Create", mock.Anything, &task).Return(nil)

	// marshalling and some assertion
	requestBody, err := json.Marshal(&task)
	suite.NoError(err, "can not marshal struct to json")

	// calling the testing server given the provided request body
	response, err := http.Post(fmt.Sprintf("%s/tasks", suite.testingServer.URL), "application/json", bytes.NewBuffer(requestBody))
	suite.NoError(err, "no error when calling the endpoint")
	defer response.Body.Close()

	// unmarshalling the response
	responseBody := domain.SuccessResponse{}
	json.NewDecoder(response.Body).Decode(&responseBody)

	// running assertions to make sure that our method does the correct thing
	suite.Equal(http.StatusCreated, response.StatusCode)
	suite.Equal(responseBody.Message, "Task created successfully")
	suite.usecase.AssertExpectations(suite.T())
}

func (suite *taskControllerSuite) TestCreateTask_nilTask(){}

func (suite *taskControllerSuite) TestGetAllTasks_Positive() {
	tasks := []domain.Task{
		{
		Title: "new title",
		Description: "New description",
		Status: "Pending",
	},
		{
		Title: "new title",
		Description: "New description",
		Status: "Pending",
	},
		{
		Title: "new title",
		Description: "New description",
		Status: "Pending",
	},
	}

	suite.usecase.On("FetchAll", mock.Anything).Return(&tasks, nil)

	response, err := http.Get(fmt.Sprintf("%s/tasks", suite.testingServer.URL))
	suite.NoError(err, "no error when calling this endpoint")
	defer response.Body.Close()

	responseBody := domain.SuccessResponse{}
	json.NewDecoder(response.Body).Decode(&responseBody)

	suite.Equal(http.StatusOK, response.StatusCode)
	suite.Equal(responseBody.Message, "Success to get all tasks")
	suite.usecase.AssertExpectations(suite.T())
}

func (suite *taskControllerSuite) TestGetTaskByID_Positive() {
	taskID := primitive.NewObjectID()
	task := domain.Task{
		ID: taskID,
		Title: "new title",
		Description: "New description",
		Status: "Pending",
	}
	suite.usecase.On("FetchByTaskID",mock.Anything, taskID.Hex()).Return(&task, nil)
	response, err := http.Get(fmt.Sprintf("%s/tasks/%v", suite.testingServer.URL, taskID.Hex()))
	suite.NoError(err, "no error when calling this endpoint")
	defer response.Body.Close()

	responseBody := domain.SuccessResponse{}
	json.NewDecoder(response.Body).Decode(&responseBody)

	suite.Equal(http.StatusOK, response.StatusCode)
	suite.Equal(responseBody.Message, fmt.Sprintf("Success to get task with id %v", taskID.Hex()))
	suite.usecase.AssertExpectations(suite.T())
}

func TestTaskController(t *testing.T) {
	suite.Run(t, new(taskControllerSuite))
}