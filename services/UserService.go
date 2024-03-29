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

func GetAllUsersBy(keyword string, order string) ([]*users.User, *errors.RestErr) {
	var result []*users.User
	result, err := users.GetAllBy(keyword, order)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func CreateUser(user users.User) (*users.User, *errors.RestErr) {

	// Call middleware to sanitize and check if the fields are correct
	if err := user.EmailValidation(); err != nil {
		return nil, err
	}

	if err := user.PasswordValidation(); err != nil {
		return nil, err
	}

	if err := user.UsernameValidation(); err != nil {
		return nil, err
	}

	// DTO save user to DB
	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUser(userId int64) *errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func UpdateUser(partialUpdate bool, user users.User) (*users.User, *errors.RestErr) {
	//Get user from db
	current, err := GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	//if partialUpdate, verify all fields to find what must be updated
	if partialUpdate {
		if user.Email != "" {
			if err := user.EmailValidation(); err != nil {
				return nil, err
			} else {
				current.Email = user.Email
			}
		}
		if user.Username != "" {
			if err := user.UsernameValidation(); err != nil {
				return nil, err
			} else {
				current.Username = user.Username
			}
		}
		if user.DisplayName != "" {
			current.DisplayName = user.DisplayName
		}
		if user.Biography != "" {
			current.Biography = user.Biography
		}
		if user.Birthday != "" {
			current.Birthday = user.Birthday
		}
		if user.ProfilePic != "" {
			current.ProfilePic = user.ProfilePic
		}
		if user.Points != 0 { //Fix: check for zero points or check pointer null
			current.Points = user.Points
		}
	} else { //fullUpdate, update all to info in current user
		// Call middleware to sanitize and check if the fields are correct
		if err := user.EmailValidation(); err != nil {
			return nil, err
		}
		if err := user.UsernameValidation(); err != nil {
			return nil, err
		} else {
			current.Email = user.Email
			current.Username = user.Username
			current.DisplayName = user.DisplayName
			current.Biography = user.Biography
			current.Birthday = user.Birthday
			current.ProfilePic = user.ProfilePic
		}
	}

	// DTO save user to DB
	if err := current.Update(); err != nil {
		return nil, err
	}

	return current, nil

}

func UpdatePassword(user users.User) (*users.User, *errors.RestErr) {
	current, err := GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	if err := user.PasswordValidation(); err != nil {
		return nil, err
	}

	oldPassword := user.Passwords.OldPassword
	newPassword := user.Passwords.Password

	if err := current.UpdatePassword(oldPassword, newPassword); err != nil {
		return nil, err
	}

	return current, nil
}

func LoginUser(request users.UserLoginRequest) (*users.User, *errors.RestErr) {
	dao := &users.User{
		Email: request.Email,
		Passwords: users.Passwords{
			Password: request.Password,
		},
	}

	if err := dao.GetUserByEmailAndPassword(); err != nil {
		return nil, err
	}
	return dao, nil
}
