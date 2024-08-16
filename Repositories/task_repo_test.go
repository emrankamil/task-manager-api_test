package repositories

import (
	"context"
	"testing"
	domain "testing_task-manager_api/Domain"
	"testing_task-manager_api/config"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepositorySuite struct{
	suite.Suite
	repository domain.TaskRepository
	client *mongo.Client
	// cleanupExecutor utils.TruncateTableExecutor
}

func (suite *taskRepositorySuite) SetupSuite(){

	configs := config.GetConfig()
	client, db := config.ConnectDB(configs)
	repository := NewTaskRepository(db, domain.CollectionTask)

	suite.client = client
	suite.repository = repository
}

func (suite *taskRepositorySuite) TearDownTest() {
	// clean-up the used table to be used for another session
}

// Create Task test
func (suite *taskRepositorySuite) TestCreateTask_Positive(){
	task := domain.Task{
		Title: "new title",
		Description: "New description",
		Status: "Pending",
	}

	err := suite.repository.Create(context.TODO(), &task)
	suite.NoError(err, "no error when create task with valid input")
}

func (suite *taskRepositorySuite) TestCreateTask_NilPointer_Negative(){
	err := suite.repository.Create(context.TODO(), nil)
	suite.Error(err, "create error with nil input returns error")

}

func (suite *taskRepositorySuite) TestCreateTask_EmptyFields_Positive(){
	var task domain.Task
	err := suite.repository.Create(context.TODO(), &task)
	suite.NoError(err, "no error when create task with empty fields")
}


// FetchAll test
func (suite *taskRepositorySuite) TestGetAllTasks_EmptySlice_Positive() {
	tasks, err := suite.repository.FetchAll(context.TODO())
	suite.NoError(err, "no error when get all tasks when the table is empty")
	suite.Equal(len(*tasks), 0, "length of tasks should be 0, since it is empty slice")
	suite.Equal(tasks, &[]domain.Task{}, "tasks is an empty slice")
}

func (suite *taskRepositorySuite) TestGetAllTasks_FilledRecords_Positive() {
	task := domain.Task{
		Title: "new title",
		Description: "New description",
		Status: "Pending",
	}

	// inserting 3 tasks to be queried later
	err := suite.repository.Create(context.TODO(), &task)
	suite.NoError(err, "no error when create task with valid input")
	err = suite.repository.Create(context.TODO(), &task)
	suite.NoError(err, "no error when create task with valid input")
	err = suite.repository.Create(context.TODO(), &task)
	suite.NoError(err, "no error when create task with valid input")

	tasks, err := suite.repository.FetchAll(context.TODO())
	suite.NoError(err, "no error when get all tasks when the table is empty")
	suite.Equal(len(*tasks), 3, "insert 3 records before get all data, so it should contain three tasks")
}

// FetchByTaskID test
func (suite *taskRepositorySuite) TestGetTaskByID_NotFound_Negative(){
	id := primitive.NewObjectID().Hex()

	_, err := suite.repository.FetchByTaskID(context.TODO(), id)
	suite.Error(err, "error db not found")
	suite.Equal(err.Error(), "mongo: no documents in result")
}

func (suite *taskRepositorySuite) TestGetTaskByID_Exists_Positive(){	
	task := domain.Task{
		ID: primitive.NewObjectID(),
		Title: "new title",
		Description: "New description",
		Status: "Pending",
	}

	err := suite.repository.Create(context.TODO(), &task)
	suite.NoError(err, "no error when create task with valid input")

	result, err := suite.repository.FetchByTaskID(context.TODO(), task.ID.Hex())
	suite.NoError(err, "no error because task is found")
	suite.Equal(task.Title, (result).Title, "should be equal between result and task")
	suite.Equal(task.Description, (result).Description, "should be equal between result and task")

}	

// Update task test
func (suite *taskRepositorySuite) TestUpdate_Positive(){
	var err error
	taskID := primitive.NewObjectID()
	task := domain.Task{
		ID: taskID,
		Title: "new title",
		Description: "New description",
		Status: "Pending",
	}

	err = suite.repository.Create(context.TODO(), &task)
	suite.NoError(err, "no error when create task with valid input")

	updatedTask := domain.Task{
		ID: taskID,
		Title: "updated title",
		Description: "updated description",
		Status: "updated status",
	}

	err = suite.repository.Update(context.TODO(), taskID.Hex(), updatedTask)
	suite.NoError(err, "no error when updating task with valid input")

	result, err := suite.repository.FetchByTaskID(context.TODO(), taskID.Hex())
	suite.NoError(err, "no error because task is found")
	suite.Equal(updatedTask.Title, (result).Title, "should be equal between result and updatedTask")
	suite.Equal(updatedTask.Description, (result).Description, "should be equal between result and updatedTask")
	suite.Equal(updatedTask.Status, (result).Status, "should be equal between result and updatedTask")
}

func (suite *taskRepositorySuite) TestUpdate_InvalidID() {
	invalidID := "invalidID"
	updatedTask := domain.Task{}
	err := suite.repository.Update(context.TODO(), invalidID, updatedTask)
	suite.Error(err)
	suite.Equal(err.Error(), "the provided hex string is not a valid ObjectID")
}

func (suite *taskRepositorySuite) TestUpdate_TaskNotFound() {
	nonExistentID := primitive.NewObjectID().Hex()
	updatedTask := domain.Task{
		Title:       "Updated Title",
		Description: "Updated Description",
		Status:      "Completed",
	}
	err := suite.repository.Update(context.TODO(), nonExistentID, updatedTask)
	suite.EqualError(err, "TASK NOT FOUND")
}

//Delete Task 
func (suite *taskRepositorySuite) TestDelete_Positive(){
	var err error
	taskID := primitive.NewObjectID()
	task := domain.Task{
		ID: taskID,
		Title: "new title",
		Description: "New description",
		Status: "Pending",
	}

	err = suite.repository.Create(context.TODO(), &task)
	suite.NoError(err, "no error when create task with valid input")

	err = suite.repository.Delete(context.TODO(), taskID.Hex())
	suite.NoError(err)

	_, err = suite.repository.FetchByTaskID(context.TODO(), taskID.Hex())
	suite.Error(err, "error db not found")
	suite.Equal(err.Error(), "mongo: no documents in result")
	}

func TestTaskRepository(t *testing.T) {
	suite.Run(t, new(taskRepositorySuite))
}