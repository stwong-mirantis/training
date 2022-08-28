package user

import "errors"

var (
	ErrEmptyUsername        = errors.New("username cannot be empty")
	ErrUsernameAlreadyInUse = errors.New("username already exists")
	ErrUsernameDoesNotExist = errors.New("cannot find username")
)

type User struct {
	Username     string `json:"username"    description:"username of the user"`
	OnlineStatus *bool  `json:"online"      description:"online status of the user"`
}

type UserRepository interface {
	GetAllUsers() []User
	GetUser(authToken string) (User, error)
	AddUser(user User) (User, error)
	RemoveUser(user User) (User, error)
}

type UserResource struct {
	users map[string]User
}

func (ur *UserResource) GetAllUsers() []User {
	var userArr []User
	for _, v := range ur.users {
		userArr = append(userArr, v)
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

func (ur *UserResource) AddUser(user User) (User, error) {

}

func (ur *UserResource) RemoveUser(user User) (User, error) {

}

func NewUserResource() *UserResource {
	return &UserResource{users: map[string]User{}}
}
