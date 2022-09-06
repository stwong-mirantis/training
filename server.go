package main

import (
	"encoding/json"
	"finalProject/message"
	"finalProject/user"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	Message string `json:"message"`
}

type Username struct {
	Username string `json:"username"`
}

type MessagingService struct {
	*restful.WebService
	user.UserRepository
	message.MessageRepository
}

func NewMessagingService(ur user.UserRepository, mr message.MessageRepository) *MessagingService {
	ms := MessagingService{
		WebService:        &restful.WebService{},
		UserRepository:    ur,
		MessageRepository: mr,
	}

	ms.Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ms.Route(ms.GET("/users").To(ms.getAllUsers).
		Doc("get all users").
		Writes([]user.User{}).Returns(http.StatusOK, "OK", []user.User{}).
		Returns(http.StatusInternalServerError, "Internal Server Error", nil))

	ms.Route(ms.GET("/users/{username}").To(ms.getUser).
		Doc("get user by username"))

	ms.Route(ms.POST("/login").To(ms.loginUser).
		Doc("login user"))

	ms.Route(ms.DELETE("/logout").To(ms.logoutUser).
		Doc("logout user"))

	ms.Route(ms.GET("/messages").To(ms.getMessages).
		Doc("get messages"))

	ms.Route(ms.POST("/messages").To(ms.createMessage).Doc("create message"))

	return &ms
}

func isRequestAndAuthTokenValid(request *restful.Request, response *restful.Response, ms *MessagingService) bool {
	token := request.Request.Header.Get("Authorization")
	var err error

	if token == "" {
		err = fmt.Errorf("auth token is not provided in request: %w", err)
		log.Println(err)
		response.WriteError(http.StatusUnauthorized, err)
		return false
	}

	if !ms.DoesAuthTokenExist(token) {
		err = fmt.Errorf("auth token does not exist: %w", err)
		log.Println(err)
		response.WriteError(http.StatusForbidden, err)
		return false
	}

	return true
}

func (ms *MessagingService) createMessage(request *restful.Request, response *restful.Response) {
	if !isRequestAndAuthTokenValid(request, response, ms) {
		return
	}

	token := request.Request.Header.Get("Authorization")
	ms.UpdateUserLastSeenTime(token)
	onlineStatus := true
	ms.SetUserOnlineStatus(&onlineStatus, token)
	user := ms.GetUserWithToken(token)
	reqBody, err := ioutil.ReadAll(request.Request.Body)

	if err != nil {
		err = fmt.Errorf("unable to read request body: %w", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	reqBodyUnmarshal := Message{}
	err = json.Unmarshal(reqBody, &reqBodyUnmarshal)

	if err != nil {
		err = fmt.Errorf("unable to unmarshal request body: %w", err)
		log.Println(err)
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	message := ms.CreateMessage(reqBodyUnmarshal.Message, user.Username)
	messageJSON, err := json.Marshal(message)

	if err != nil {
		err = fmt.Errorf("unable to marshal message: %w", err)
		log.Println(err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.Write(messageJSON)
}

func (ms *MessagingService) getMessages(request *restful.Request, response *restful.Response) {
	if !isRequestAndAuthTokenValid(request, response, ms) {
		return
	}

	token := request.Request.Header.Get("Authorization")
	ms.UpdateUserLastSeenTime(token)
	onlineStatus := true
	ms.SetUserOnlineStatus(&onlineStatus, token)
	countQueryStr := request.Request.URL.Query()["count"]
	offsetQueryStr := request.Request.URL.Query()["offset"]
	count := 10
	offset := 0

	if len(countQueryStr) != 0 {
		countQueryConverted, err := strconv.Atoi(countQueryStr[0])
		if err == nil {
			count = countQueryConverted
		}
	}

	if len(offsetQueryStr) != 0 {
		offsetQueryConverted, err := strconv.Atoi(offsetQueryStr[0])
		if err == nil {
			offset = offsetQueryConverted
		}
	}

	messages := ms.GetPaginatedMessages(count, offset)
	messagesJSON, err := json.Marshal(messages)

	if err != nil {
		err = fmt.Errorf("unable to marshal messages: %w", err)
		log.Println(err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.Write(messagesJSON)
}

func (ms *MessagingService) logoutUser(request *restful.Request, response *restful.Response) {
	if !isRequestAndAuthTokenValid(request, response, ms) {
		return
	}

	token := request.Request.Header.Get("Authorization")
	message, err := ms.RemoveUser(token)

	if err != nil {
		err = fmt.Errorf("unable to delete: %w", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	messageJSON, err := json.Marshal(message)

	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.Write(messageJSON)
}

func (ms *MessagingService) loginUser(request *restful.Request, response *restful.Response) {

	reqBody, err := ioutil.ReadAll(request.Request.Body)

	if err != nil {
		err = fmt.Errorf("unable to read request body: %w", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	reqBodyUnmarshal := Username{}
	err = json.Unmarshal(reqBody, &reqBodyUnmarshal)

	if err != nil {
		err = fmt.Errorf("unable to unmarshal request body: %w", err)
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	if len(reqBodyUnmarshal.Username) == 0 {
		response.WriteError(http.StatusBadRequest, err)
	}

	user, err := ms.AddUser(reqBodyUnmarshal.Username)

	if err != nil {
		err = fmt.Errorf("username is already in use")
		log.Println(err)
		response.Header().Add("WWW-Authenticate", "Token realm='Username is already in use'")
		response.WriteError(http.StatusUnauthorized, err)
		return
	}

	userJSON, err := json.Marshal(user.UUID)

	if err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	response.Write(userJSON)
}

func (ms *MessagingService) getUser(request *restful.Request, response *restful.Response) {
	if !isRequestAndAuthTokenValid(request, response, ms) {
		return
	}

	token := request.Request.Header.Get("Authorization")
	ms.UpdateUserLastSeenTime(token)
	onlineStatus := true
	ms.SetUserOnlineStatus(&onlineStatus, token)
	username := request.PathParameter("username")
	user, err := ms.GetUserWithUsername(username)

	if err != nil {
		err = fmt.Errorf("unable to get user or user does not exist")
		log.Println(err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	userJSON, err := json.Marshal(user)

	if err != nil {
		err = fmt.Errorf("unable to marshal user: %w", err)
		log.Println(err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.Write(userJSON)
}

func (ms *MessagingService) getAllUsers(request *restful.Request, response *restful.Response) {
	if !isRequestAndAuthTokenValid(request, response, ms) {
		return
	}
	token := request.Request.Header.Get("Authorization")
	ms.UpdateUserLastSeenTime(token)
	onlineStatus := true
	ms.SetUserOnlineStatus(&onlineStatus, token)
	users := ms.GetAllOnlineUsers()
	usersJSON, err := json.Marshal(users)

	if err != nil {
		err = fmt.Errorf("unable to marshal users: %w", err)
		log.Println(err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.Write(usersJSON)
}

func closeSessionForInactiveUsers(ur user.UserRepository) {
	for {
		mu := ur.GetMutex()
		usersMap := ur.GetUserMap()
		currentTime := time.Now().Unix()
		mu.Lock()
		for k, v := range usersMap {
			if currentTime-v.LastSeenTime.Unix() > 10 && v.OnlineStatus != nil {
				fmt.Println("inside the inner loop!")
				ur.SetUserOnlineStatusNoLock(nil, k)
			}
		}
		mu.Unlock()
	}
}

func main() {
	ur := user.NewUserResource()
	mr := message.NewMessageResource()
	messagingService := NewMessagingService(ur, mr)
	restful.DefaultContainer.Add(messagingService.WebService)
	go closeSessionForInactiveUsers(ur)
	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
