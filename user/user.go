package user

import (
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

var (
	ErrUsernameAlreadyInUse = errors.New("username already exists")
	ErrUsernameDoesNotExist = errors.New("cannot find username")
)

type User struct {
	UUID         string
	Username     string    `json:"username"    description:"username of the user"`
	OnlineStatus *bool     `json:"online"      description:"online status of the user"`
	LastSeenTime time.Time `json:"lastSeen"      description:"when the user was last seen"`
}

type Message struct {
	Message string `json:"message"`
}

type UserRepository interface {
	SetUserOnlineStatus(onlineStatus *bool, authToken string)
	UpdateUserLastSeenTime(authToken string)
	DoesAuthTokenExist(authToken string) bool
	GetAllOnlineUsers() []User
	GetUserWithUsername(username string) (User, error)
	GetUserWithToken(token string) User
	AddUser(username string) (User, error)
	RemoveUser(authToken string) (Message, error)
	GetAllUsers() []User
	GetUserMap() map[string]User
	GetMutex() *sync.Mutex
}

type UserResource struct {
	users map[string]User
	mu    sync.Mutex
}

func (ur *UserResource) GetUserMap() map[string]User {
	return ur.users
}

func (ur *UserResource) GetMutex() *sync.Mutex {
	return &ur.mu
}

func (ur *UserResource) SetUserOnlineStatus(onlineStatus *bool, authToken string) {
	ur.mu.Lock()
	userUpdated := ur.users[authToken]
	userUpdated.OnlineStatus = onlineStatus
	ur.users[authToken] = userUpdated
	ur.mu.Unlock()
}

func (ur *UserResource) UpdateUserLastSeenTime(authToken string) {
	userUpdated := ur.users[authToken]
	userUpdated.LastSeenTime = time.Now()
	ur.mu.Lock()
	ur.users[authToken] = userUpdated
	ur.mu.Unlock()
}

func (ur *UserResource) DoesAuthTokenExist(authToken string) bool {
	ur.mu.Lock()
	_, ok := ur.users[authToken]
	ur.mu.Unlock()
	if ok {
		return true
	}
	return false
}

func (ur *UserResource) GetAllOnlineUsers() []User {
	ur.mu.Lock()
	userArr := make([]User, 0)
	for _, v := range ur.users {
		if v.OnlineStatus != nil && *v.OnlineStatus {
			userArr = append(userArr, v)
		}
	}
	ur.mu.Unlock()
	return userArr
}

func (ur *UserResource) GetAllUsers() []User {
	ur.mu.Lock()
	var userArr []User
	for _, v := range ur.users {
		userArr = append(userArr, v)
	}
	ur.mu.Unlock()
	return userArr
}

func (ur *UserResource) GetUserWithUsername(username string) (User, error) {
	ur.mu.Lock()
	for _, v := range ur.users {
		if v.Username == username {
			return v, nil
		}
	}
	ur.mu.Unlock()
	return User{}, ErrUsernameDoesNotExist
}

func (ur *UserResource) GetUserWithToken(token string) User {
	ur.mu.Lock()
	user := ur.users[token]
	ur.mu.Unlock()
	return user
}

func (ur *UserResource) AddUser(username string) (User, error) {
	ur.mu.Lock()
	for _, v := range ur.users {
		if username == v.Username {
			return User{}, ErrUsernameAlreadyInUse
		}
	}
	ur.mu.Unlock()
	id := uuid.New().String()
	onlineStatus := new(bool)
	*onlineStatus = true
	newUser := User{id, username, onlineStatus, time.Now()}
	ur.mu.Lock()
	ur.users[id] = newUser // here
	ur.mu.Unlock()
	return newUser, nil

}

func (ur *UserResource) RemoveUser(authToken string) (Message, error) {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	if _, ok := ur.users[authToken]; ok {
		delete(ur.users, authToken)
		return Message{Message: "bye!"}, nil
	}
	return Message{}, ErrUsernameDoesNotExist
}

func NewUserResource() *UserResource {
	return &UserResource{users: map[string]User{}}
}
