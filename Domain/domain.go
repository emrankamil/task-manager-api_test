package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DatabaseName = "taskmanager"
	CollectionTask = "tasks"
	CollectionUser = "users"
)

type Task struct {
 ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
 Title       string    `json:"title"`
 Description string    `json:"description"`
 DueDate     time.Time `json:"due_date"`
 Status      string    `json:"status"`
}

type User struct{
	ID				primitive.ObjectID		`bson:"_id"`
	Name		    *string			`json:"name" validate:"required,min=2,max=100"`
	Username		*string			`json:"username" validate:"required,min=2,max=100"`
	Password		*string			`json:"password" validate:"required,min=6"`
	Email			*string			`json:"email" validate:"email,required"`
	User_type		string			`json:"user_type"`
	Created_at		time.Time		`json:"created_at"`
	Updated_at		time.Time		`json:"updated_at"`
	User_id			string			`json:"user_id"`
}

type Config struct {
    MongoDBURI string
    Port       string
    TimeZone   string
    SecretKey  string
	DatabaseName string
}

type TaskRepository interface {
	Create(c context.Context, task *Task) error
	FetchAll(c context.Context) (*[]Task, error)
	FetchByTaskID(c context.Context, taskID string) (*Task, error)
	Update(c context.Context, taskID string, updatedTask Task) error
	Delete(c context.Context, taskID string) error
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	FindByUsername(c context.Context, usrname string) (User, error)
	Update(c context.Context, userID string) error
}

type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	FetchAll(c context.Context) (*[]Task, error)
	FetchByTaskID(c context.Context, taskID string) (*Task, error)
	Update(c context.Context, taskID string, updatedTask Task) error
	Delete(c context.Context, taskID string) error
}

type UserUsecase interface {
	Create(c context.Context, user *User) error
	HandleLogin(c context.Context, username *User) (string, string, error)
	Update(c context.Context, userID string) error
}

type TaskController interface{
	Create(c *gin.Context)
	FetchAll(c *gin.Context)
	FetchByTaskID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type UserController interface{
	Signup(c *gin.Context)
	Login(c *gin.Context)
	PromoteUser(c *gin.Context)
}
type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}
