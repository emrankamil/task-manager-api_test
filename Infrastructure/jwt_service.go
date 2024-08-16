package infrastructure

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)


type UserClaim struct{
	User_id			string			
	Username		string
	Email			string
	User_type		string
	jwt.StandardClaims
}

func ValidateToken(signedToken string) (claims *UserClaim, err error){
	token, msg := jwt.ParseWithClaims(
		signedToken, 
		&UserClaim{}, 
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if msg != nil || !token.Valid{
		err = msg
		return
	}

	claims, ok:= token.Claims.(*UserClaim)
	if !ok{
		err = errors.New("the token is invalid")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix(){
		err = errors.New("token is expired")
		return
	}
	return claims, err
}

func GenerateJWTToken(user_id string, username string, email string, user_type string) (signedToken, signedRefreshToken string, err error){
	claims := &UserClaim{
		User_id: user_id,
		Username: username,
		Email: email,
		User_type: user_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &UserClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
    if err != nil {
        return "", "", err
    }

    signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
    if err != nil {
        return "", "", err
    }

    return
}
