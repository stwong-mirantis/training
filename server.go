package main

import (
	"encoding/json"
	"finalProject/user"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"log"
	"net/http"
)

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

	us.Route(us.GET("/users/{username}").To(us.getUser).
		Doc("get user by username"))

	us.Route(us.POST("/login").To(us.loginUser).
		Doc("login user"))

	us.Route(us.DELETE("/logout").To(us.logoutUser).Doc("logout user"))

	return &us
}

func (us *UserService) getAllUsers(request *restful.Request, response *restful.Response) {
	users := us.GetAllUsers()
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
