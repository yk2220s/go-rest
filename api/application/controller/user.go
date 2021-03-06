package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yk2220s/go-rest-sample/api/domain/model"
	"github.com/yk2220s/go-rest-sample/api/domain/repository"
)

// UserController is present data to show.
type UserController struct {
	uRepository repository.UserRepository
}

// UserControllerFactory create UserController injected depency
func UserControllerFactory() UserController {
	var urepo repository.UserRepository
	urepo = repository.UserRepositoryFactory()

	return UserController{urepo}
}

// ListUser fetch List of users
func (controller UserController) ListUser(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	users := controller.uRepository.List(page)

	c.IndentedJSON(http.StatusOK, gin.H{"users": users})
}

// GetUser fetch List of users
func (controller UserController) GetUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Params.ByName("id"))

	user, err := controller.uRepository.GetByID(userID)

	if err != nil {
		c.AbortWithError(err.StatusCode(), err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"user": user})
}

type paramPostUser struct {
	User struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"omitempty,max=255,email"`
	} `json:"user"`
}

// CreateUser create User record.
func (controller UserController) CreateUser(c *gin.Context) {
	var puser paramPostUser

	if err := c.ShouldBindJSON(&puser); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user := model.User{
		Name:  puser.User.Name,
		Email: puser.User.Email,
	}

	newUser, cerr := controller.uRepository.Store(&user)

	if cerr != nil {
		c.AbortWithError(cerr.StatusCode(), cerr)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"user": newUser})
}

type paramPatchUser struct {
	User struct {
		Name  string `json:"name" binding:""`
		Email string `json:"email" binding:"omitempty,max=255,email"`
	} `json:"user"`
}

// UpdateUser create User record.
func (controller UserController) UpdateUser(c *gin.Context) {
	var puser paramPatchUser

	if err := c.ShouldBindJSON(&puser); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userID, _ := strconv.Atoi(c.Params.ByName("id"))
	user, gerr := controller.uRepository.GetByID(userID)

	if gerr != nil {
		c.AbortWithError(gerr.StatusCode(), gerr)
		return
	}

	user.Name = puser.User.Name
	user.Email = puser.User.Email

	newUser, uerr := controller.uRepository.Update(user)

	if uerr != nil {
		c.AbortWithError(uerr.StatusCode(), uerr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": newUser})
}

// DeleteUser delete User record.
func (controller UserController) DeleteUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Params.ByName("id"))

	isSuccess, err := controller.uRepository.Delete(userID)

	if err != nil {
		c.AbortWithError(err.StatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"isSuccess": isSuccess})
}
