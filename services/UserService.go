package services

import (
	"moku-moku/domain/users"
	"moku-moku/utils/errors"
)

func GetUser(userId int64) (*users.User, *errors.RestErr) {
	result := &users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func CreateUser(user users.User) (*users.User, *errors.RestErr) {

	// Call middleware to sanitize and check if the fields are correct
	if err := user.EmailValidation(); err != nil {
		return nil, err
	}

}
