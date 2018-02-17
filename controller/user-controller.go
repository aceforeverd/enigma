package controller

import (
	"database/sql"
	"errors"
	"github.com/aceforeverd/enigma/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserCon interface {
	GetAll(c *gin.Context)
	GetUser(c *gin.Context)
	Save(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type UserController struct {
	Repo repository.UserRepo
}

// GetUser GET /users
func (ctr *UserController) GetAll(c *gin.Context) {
	users, err := ctr.Repo.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetByID GET /user?id=ID or /user?name=NAME
func (ctr *UserController) GetUser(c *gin.Context) {
	idName := c.Query("id")
	if len(idName) > 0 {
		id, errID := strconv.Atoi(c.Query("id"))
		if errID != nil {
			c.String(http.StatusInternalServerError, errID.Error())
			return
		}
		if user, err := ctr.Repo.GetByID(id); err != nil {
			if err == sql.ErrNoRows {
				c.String(http.StatusAccepted, ErrUserNotFound.Error())
				return
			}
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, user)
		}
		return
	}

	name := c.Query("name")
	if len(name) > 0 {
		user, err := ctr.Repo.GetByUsername(name)
		if err != nil {
			if err == sql.ErrNoRows {
				c.String(http.StatusAccepted, ErrUserNotFound.Error())
				return
			}
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, user)
		return
	}

	c.String(http.StatusBadRequest, "id or name required")
}

// Save POST /user
func (ctr *UserController) Save(c *gin.Context) {
	var user repository.User
	if err := c.BindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	saved, err := ctr.Repo.Save(user)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, saved)
}

// Update PUT /user
func (ctr *UserController) Update(c *gin.Context) {
	var user repository.User
	if err := c.BindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	saved, err := ctr.Repo.Update(user)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, saved)
}

// Delete DELETE /user
func (ctr *UserController) Delete(c *gin.Context) {
	var user repository.User
	if err := c.BindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := ctr.Repo.Delete(user); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "deleted")
}
