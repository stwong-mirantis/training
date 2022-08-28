package user

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrEmptyUsername        = errors.New("username cannot be empty")
	ErrUsernameAlreadyInUse = errors.New("username already exists")
	ErrUsernameDoesNotExist = errors.New("cannot find username")
)

type User struct {
	UUID         string
	Username     string `json:"username"    description:"username of the user"`
	OnlineStatus *bool  `json:"online"      description:"online status of the user"`
}

type UserRepository interface {
	GetAllUsers() []User
	GetUser(authToken string) (User, error)
	AddUser(username string) (User, error)
	RemoveUser(username string) (User, error)
}

type UserResource struct {
	users map[string]User
}

func (ur *UserResource) GetAllUsers() []User {
	var userArr []User
	for _, v := range ur.users {
		if *v.OnlineStatus {
			userArr = append(userArr, v)
		}
	}
	return userArr
}

func (ur *UserResource) GetUser(id string) (User, error) {
	u, ok := ur.users[id]
	if !ok {
		return User{}, ErrUsernameDoesNotExist
	}
	return u, nil
}

func (ur *UserResource) AddUser(username string) (User, error) {
	if len(username) == 0 {
		return User{}, ErrEmptyUsername
	}

	for _, v := range ur.users {
		if username == v.Username {
			return User{}, ErrUsernameAlreadyInUse
		}
	}
	id := uuid.New().String()
	onlineStatus := new(bool)
	*onlineStatus = true
	newUser := User{id, username, onlineStatus}
	ur.users[id] = newUser
	return newUser, nil

}

func (ur *UserResource) RemoveUser(username string) (User, error) {
	for k, v := range ur.users {
		if v.Username == username {
			delete(ur.users, k)
			return v, nil
		}
	}
	return User{}, ErrUsernameDoesNotExist
}

func NewUserResource() *UserResource {
	return &UserResource{users: map[string]User{}}
}
