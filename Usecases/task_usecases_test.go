package usecases

import (
	"context"
	"errors"
	"testing"
	domain "testing_task-manager_api/Domain"
	"testing_task-manager_api/Domain/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskUsecaseSuite struct{
	suite.Suite
	repository *mocks.TaskRepository
	usecase domain.TaskUsecase
}

func (suite *taskUsecaseSuite) SetupTest(){
	
	repository := new(mocks.TaskRepository)
	usecase := NewTaskUsecase(repository, 10)

	suite.repository = repository
	suite.usecase = usecase
}

// create task test
func (suite *taskUsecaseSuite) TestCreateTask_Positive(){
	task := domain.Task{
		Title: "new title",
		Description: "New description",
		Status: "Pending",
	}

	suite.repository.On("Create",mock.Anything, &task).Return(nil)

	err := suite.usecase.Create(context.TODO(), &task)

	// assertions to make sure our operation does the right thing
	suite.Nil(err, "err is a nil pointer so no error in this process")
	suite.repository.AssertExpectations(suite.T())

}

func (suite *taskUsecaseSuite) TestCreateTask_EmptyFields_Positive(){
	var task domain.Task

	suite.repository.On("Create",mock.Anything, &task).Return(nil)

	err := suite.usecase.Create(context.TODO(), &task)

	// assertions to make sure our operation does the right thing
	suite.repository.AssertExpectations(suite.T())
	suite.NoError(err, "no error when create task with empty fields")
}

// Test FetchAll - Positive case
func (suite *taskUsecaseSuite) TestFetchAll_Positive() {
	tasks := &[]domain.Task{
		{
			Title:       "Task 1",
			Description: "Description 1",
			Status:      "Pending",
		},
		{
			Title:       "Task 2",
			Description: "Description 2",
			Status:      "Completed",
		},
	}

	suite.repository.On("FetchAll", mock.Anything).Return(tasks, nil)

	result, err := suite.usecase.FetchAll(context.TODO())

	// Assertions
	suite.NoError(err)
	suite.Equal(len(*tasks), len(*result))
	suite.repository.AssertExpectations(suite.T())
}

// Test FetchByTaskID - Positive case
func (suite *taskUsecaseSuite) TestFetchByTaskID_Positive() {
	taskID := primitive.NewObjectID()
    task := &domain.Task{
		ID: taskID,
        Title:       "new title",
        Description: "New description",
        Status:      "Pending",
    }

    suite.repository.On("Create", mock.Anything, task).Return(nil)

    err := suite.usecase.Create(context.TODO(), task)

    // Assertions to make sure our operation does the right thing
    suite.Nil(err, "task created successfully")
    suite.repository.AssertExpectations(suite.T())

    // Fetch the task ID after creation

    suite.repository.On("FetchByTaskID", mock.Anything, taskID.Hex()).Return(task, nil)

    result, err := suite.usecase.FetchByTaskID(context.TODO(), taskID.Hex())

    // Assertions
    suite.NoError(err)
    suite.Equal(task, result)
    suite.repository.AssertExpectations(suite.T())
}


// Test FetchByTaskID - Negative case (Task not found)
func (suite *taskUsecaseSuite) TestFetchByTaskID_TaskNotFound() {
	taskID := primitive.NewObjectID().Hex()

	suite.repository.On("FetchByTaskID", mock.Anything, taskID).Return(&domain.Task{}, errors.New("task not found"))

	_, err := suite.usecase.FetchByTaskID(context.TODO(), taskID)

	// Assertions
	suite.Error(err)
	suite.EqualError(err, "task not found")
	suite.repository.AssertExpectations(suite.T())
}

// Test Update - Positive case
func (suite *taskUsecaseSuite) TestUpdate_Positive() {
	taskID := primitive.NewObjectID()

    task := domain.Task{
		ID: taskID,
        Title:       "new title",
        Description: "New description",
        Status:      "Pending",
    }

    suite.repository.On("Create", mock.Anything, &task).Return(nil)

    err := suite.usecase.Create(context.TODO(), &task)

    // Assertions to make sure our operation does the right thing
    suite.Nil(err, "task created successfully")
    suite.repository.AssertExpectations(suite.T())

	updatedTask := domain.Task{
		Title:       "Updated Title",
		Description: "Updated Description",
		Status:      "Completed",
	}

	suite.repository.On("Update", mock.Anything, taskID.Hex(), updatedTask).Return(nil)

	err = suite.usecase.Update(context.TODO(), taskID.Hex(), updatedTask)

	// Assertions
	suite.NoError(err)
	suite.repository.AssertExpectations(suite.T())
}

// Test Update - Negative case (Task not found)
func (suite *taskUsecaseSuite) TestUpdate_TaskNotFound() {
	taskID := primitive.NewObjectID().Hex()
	updatedTask := domain.Task{
		Title:       "Updated Title",
		Description: "Updated Description",
		Status:      "Completed",
	}

	suite.repository.On("Update", mock.Anything, taskID, updatedTask).Return(errors.New("task not found"))

	err := suite.usecase.Update(context.TODO(), taskID, updatedTask)

	// Assertions
	suite.Error(err)
	suite.EqualError(err, "task not found")
	suite.repository.AssertExpectations(suite.T())
}

// Test Delete - Positive case
func (suite *taskUsecaseSuite) TestDelete_Positive() {
	taskID := primitive.NewObjectID()

    task := domain.Task{
		ID: taskID,
        Title:       "new title",
        Description: "New description",
        Status:      "Pending",
    }

    suite.repository.On("Create", mock.Anything, &task).Return(nil)

    err := suite.usecase.Create(context.TODO(), &task)

    // Assertions to make sure our operation does the right thing
    suite.Nil(err, "task created successfully")
    suite.repository.AssertExpectations(suite.T())

	suite.repository.On("Delete", mock.Anything, taskID.Hex()).Return(nil)

	err = suite.usecase.Delete(context.TODO(), taskID.Hex())

	// Assertions
	suite.NoError(err)
	suite.repository.AssertExpectations(suite.T())

	suite.repository.On("FetchByTaskID", mock.Anything, taskID.Hex()).Return(&domain.Task{}, errors.New("task not found"))

	_, err = suite.usecase.FetchByTaskID(context.TODO(), taskID.Hex())

	// Assertions
	suite.Error(err)
	suite.EqualError(err, "task not found")
	suite.repository.AssertExpectations(suite.T())
}

// Test Delete - Negative case (Task not found)
func (suite *taskUsecaseSuite) TestDelete_TaskNotFound() {
	taskID := primitive.NewObjectID().Hex()

	suite.repository.On("Delete", mock.Anything, taskID).Return(errors.New("task not found"))

	err := suite.usecase.Delete(context.TODO(), taskID)

	// Assertions
	suite.Error(err)
	suite.EqualError(err, "task not found")
	suite.repository.AssertExpectations(suite.T())

	
}


func TestTaskUsecase(t *testing.T) {
	suite.Run(t, new(taskUsecaseSuite))
}
