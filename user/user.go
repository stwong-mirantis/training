package user

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrUsernameAlreadyInUse = errors.New("username already exists")
	ErrUsernameDoesNotExist = errors.New("cannot find username")
)

type User struct {
	UUID         string
	Username     string `json:"username"    description:"username of the user"`
	OnlineStatus *bool  `json:"online"      description:"online status of the user"`
}

type Message struct {
	Message string `json:"message"`
}

type UserRepository interface {
	DoesAuthTokenExist(authToken string) bool
	GetAllOnlineUsers() []User
	GetUserWithUsername(username string) (User, error)
	GetUserWithToken(token string) User
	AddUser(username string) (User, error)
	RemoveUser(authToken string) (Message, error)
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

func (ur *UserResource) GetUserWithUsername(username string) (User, error) {
	for _, v := range ur.users {
		if v.Username == username {
			return v, nil
		}
	}
	return User{}, ErrUsernameDoesNotExist
}

func (ur *UserResource) GetUserWithToken(token string) User {
	return ur.users[token]
}

func (ur *UserResource) AddUser(username string) (User, error) {

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

func (ur *UserResource) RemoveUser(authToken string) (Message, error) {
	if _, ok := ur.users[authToken]; ok {
		delete(ur.users, authToken)
		return Message{Message: "bye!"}, nil
	}
	return Message{}, ErrUsernameDoesNotExist
}

func NewUserResource() *UserResource {
	return &UserResource{users: map[string]User{}}
}
