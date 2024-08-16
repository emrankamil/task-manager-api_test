package usecases

import (
	"context"
	"errors"
	domain "testing_task-manager_api/Domain"
	infrastructure "testing_task-manager_api/Infrastructure"
	"time"
)

type userUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (uu *userUsecase) Create(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.Create(ctx, user)
}

func (uu *userUsecase) HandleLogin(c context.Context, user *domain.User) (signedToken, signedRefreshToken string, err error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	foundUser, err := uu.userRepository.FindByUsername(ctx, *user.Username)
	if err != nil{
		return "", "", errors.New("user not found")
	}
	check, verifMsg := infrastructure.VerifyPassword(*user.Password, *foundUser.Password)
	if !check{
		return "", "", errors.New(verifMsg)
	}
	if foundUser.Username == nil || foundUser.Email == nil || foundUser.User_type == "" {
		return "", "", errors.New("invalid user data")
	}	
	token, refreshToken, err := infrastructure.GenerateJWTToken(foundUser.User_id, *foundUser.Username, *foundUser.Email,  foundUser.User_type)
	return token, refreshToken, err
}

func (uu *userUsecase) Update(c context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.Update(ctx, userID)
}


