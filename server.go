package main

import (
	"encoding/json"
	"finalProject/user"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"io/ioutil"
	"log"
	"net/http"
)

type Authorization struct {
	Authorization string `json:"authorization"`
}

type Username struct {
	Username string `json:"username"`
}

type UserService struct {
	*restful.WebService
	user.UserRepository
}

func NewUserService(ur user.UserRepository) *UserService {
	us := UserService{
		WebService:     &restful.WebService{},
		UserRepository: ur,
	}

	us.Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	us.Route(us.GET("/users").To(us.getAllUsers).
		Doc("get all users").
		Writes([]user.User{}).Returns(http.StatusOK, "OK", []user.User{}).
		Returns(http.StatusInternalServerError, "Internal Server Error", nil))

	//us.Route(us.GET("/users/{username}").To(us.getUser).
	//	Doc("get user by username"))

	us.Route(us.POST("/login").To(us.loginUser).
		Doc("login user"))
	//
	us.Route(us.DELETE("/logout").To(us.logoutUser).
		Doc("logout user"))

	return &us
}

func isRequestAndAuthTokenValid(request *restful.Request, response *restful.Response, us *UserService) bool {
	reqBody, err := ioutil.ReadAll(request.Request.Body)
	if err != nil {
		err = fmt.Errorf("unable to read request body: %w", err)
		log.Println(err)
		response.WriteError(http.StatusInternalServerError, err)
		return false
	}

	reqBodyUnmarshal := Authorization{}
	err = json.Unmarshal(reqBody, &reqBodyUnmarshal)

	if err != nil {
		err = fmt.Errorf("unable to unmarshal request body: %w", err)
		log.Println(err)
		response.WriteError(http.StatusBadRequest, err)
		return false
	}

	if len(reqBodyUnmarshal.Authorization) == 0 {
		err = fmt.Errorf("auth token is not provided in request: %w", err)
		log.Println(err)
		response.WriteError(http.StatusUnauthorized, err)
		return false
	}

	if !us.DoesAuthTokenExist(reqBodyUnmarshal.Authorization) {
		err = fmt.Errorf("auth token does not exist: %w", err)
		log.Println(err)
		response.WriteError(http.StatusForbidden, err)
		return false
	}

	return true
}

func (us *UserService) logoutUser(request *restful.Request, response *restful.Response) {

	reqBody, err := ioutil.ReadAll(request.Request.Body)
	if err != nil {
		err = fmt.Errorf("unable to read request body: %w", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	reqBodyUnmarshal := Authorization{}
	err = json.Unmarshal(reqBody, &reqBodyUnmarshal)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal request body: %w", err)
		log.Println(err)
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	deletedUser, err := us.RemoveUser(reqBodyUnmarshal.Authorization)

	if err != nil {
		err = fmt.Errorf("unable to delete: %w", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	deletedUserJSON, err := json.Marshal(deletedUser)

	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.Write(deletedUserJSON)

}

func (us *UserService) loginUser(request *restful.Request, response *restful.Response) {

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

	user, err := us.AddUser(reqBodyUnmarshal.Username)

	if err != nil {
		err = fmt.Errorf("username is of inappropriate format or username already exists")
		log.Println(err)
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	userJSON, err := json.Marshal(user.UUID)

	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.Write(userJSON)

}

//func (us *UserService) getUser(request *restful.Request, response *restful.Response) {
//	if !isRequestAndAuthTokenValid(request, response, us) {
//		return
//	}
//}

func (us *UserService) getAllUsers(request *restful.Request, response *restful.Response) {
	if !isRequestAndAuthTokenValid(request, response, us) {
		return
	}

	users := us.GetAllOnlineUsers()
	usersJSON, err := json.Marshal(users)

	if err != nil {
		err = fmt.Errorf("unable to marshal products: %w", err)
		log.Println(err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.Write(usersJSON)
}

func main() {

	ur := user.NewUserResource()
	us := NewUserService(ur)

	restful.DefaultContainer.Add(us.WebService)

	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
