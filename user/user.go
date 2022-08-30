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
	DoesAuthTokenExist(authToken string) bool
	GetAllOnlineUsers() []User
	GetUser(authToken string) (User, error)
	AddUser(username string) (User, error)
	RemoveUser(authToken string) (User, error)
}

type UserResource struct {
	users map[string]User
}

func (ur *UserResource) DoesAuthTokenExist(authToken string) bool {
	if _, ok := ur.users[authToken]; ok {
		return true
	}
	return false
}

func (ur *UserResource) GetAllOnlineUsers() []User {
	var userArr []User
	for _, v := range ur.users {
		if *v.OnlineStatus {
			userArr = append(userArr, v)
		}
	}
	return userArr
}

func (ur *UserResource) GetUser(username string) (User, error) {
	for _, v := range ur.users {
		if v.Username == username {
			return v, nil
		}
	}
	return User{}, ErrUsernameDoesNotExist
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

func (ur *UserResource) RemoveUser(authToken string) (User, error) {
	if _, ok := ur.users[authToken]; ok {
		deletedUser := ur.users[authToken]
		delete(ur.users, authToken)
		return deletedUser, nil
	}
	return User{}, ErrUsernameDoesNotExist
}

func NewUserResource() *UserResource {
	return &UserResource{users: map[string]User{}}
}
