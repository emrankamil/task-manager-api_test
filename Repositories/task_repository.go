package repositories

import (
	"context"
	"errors"
	domain "testing_task-manager_api/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	database   *mongo.Database
	collection string
}

func NewTaskRepository(db *mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

func (tr *taskRepository) Create(c context.Context, task *domain.Task) error {

	if task == nil{
		return errors.New("task cannot be nil")
	}
	task.DueDate = time.Now()

	taskCollection := tr.database.Collection(tr.collection)
	_, err := taskCollection.InsertOne(c, task)

	return err
}

func (tr *taskRepository) FetchAll(c context.Context) (*[]domain.Task, error) {
	var tasks []domain.Task
	taskCollection := tr.database.Collection(tr.collection)

	cur, err := taskCollection.Find(c, bson.D{{}})
	if err != nil {
		return &[]domain.Task{}, err
	}

	for cur.Next(c) {
		var task *domain.Task

		val := cur.Decode(&task)
		if val != nil {
			return &[]domain.Task{}, val
		}

		tasks = append(tasks, *task)
	}

	if err := cur.Err(); err != nil {
		return &[]domain.Task{}, err
	}

	cur.Close(c)
	return &tasks, nil
}

func (tr *taskRepository) FetchByTaskID(c context.Context, taskID string) (*domain.Task, error) {
	var task *domain.Task
	taskCollection := tr.database.Collection(tr.collection)
	objID, err := primitive.ObjectIDFromHex(taskID)

	if err != nil {
        return &domain.Task{}, err
    }

	filter := bson.D{{Key: "_id", Value: objID}}
	result := taskCollection.FindOne(c, filter).Decode(&task)
	if result != nil {
		return &domain.Task{}, result
	}
	return task, result
}

func (tr *taskRepository) Update(c context.Context, taskID string, updatedTask domain.Task) error {
	taskCollection := tr.database.Collection(tr.collection)
	objID, err := primitive.ObjectIDFromHex(taskID)
    if err != nil {
        return err
    }

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: updatedTask.Title},
			{Key: "description", Value: updatedTask.Description},
			{Key: "due_date", Value: updatedTask.DueDate},
			{Key: "status", Value: updatedTask.Status},
		}},
	}
	updateResult, result := taskCollection.UpdateOne(c, filter, update)
	if result != nil {
		return result
	}
	if updateResult.MatchedCount == 0{
		return errors.New("TASK NOT FOUND")
	}
	return nil
}

func (tr *taskRepository) Delete(c context.Context, taskID string) error {
	taskCollection := tr.database.Collection(tr.collection)
	objID, err := primitive.ObjectIDFromHex(taskID)
    if err != nil {
        return errors.New("INVALID ID")
    }

    result, err := taskCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: objID}})
    if err != nil {
        return err 
    }

    if result.DeletedCount == 0 {
        return errors.New("task not found")
    }

    return nil
}
