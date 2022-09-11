package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Panicf("fail to connect database test.db: %s\n", err)
	}

	db.Migrator().DropTable(&User{})
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Println(err)
	}

	if r := db.Create(&User{Name: "Pallat", Email: "pallat@go.dev"}); r.Error != nil {
		log.Println(r.Error)
	}
	if r := db.Create(&User{Name: "Gopher", Email: "gopher@go.dev"}); r.Error != nil {
		log.Println(r.Error)
	}
	if r := db.Create(&User{Name: "Yod", Email: "yod@go.dev"}); r.Error != nil {
		log.Println(r.Error)
	}

	r := gin.Default()

	handler := UserHandler{db: db}

	r.POST("/users", handler.AddUsers)
	r.GET("/users", handler.Users)
	r.GET("/users/:id", handler.UsersID)
	r.DELETE("/users/:id", handler.DeleteUsers)

	r.Run()
}

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserHandler struct {
	db *gorm.DB
}

func (h UserHandler) AddUsers(c *gin.Context) {
	var user User

	if err := c.ShouldBind(&user); err != nil {
		c.Error(err)
		return
	}

	r := h.db.Create(&user)
	if err := r.Error; err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h UserHandler) Users(c *gin.Context) {
	var users []User

	r := h.db.Find(&users)
	if err := r.Error; err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h UserHandler) UsersID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Error(err)
		return
	}

	user := new(User)
	user.ID = uint(id)
	r := h.db.Take(&user)
	if err := r.Error; err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h UserHandler) DeleteUsers(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Error(err)
		return
	}

	user := new(User)
	user.ID = uint(id)
	r := h.db.Delete(&user)
	if err := r.Error; err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusAccepted)
}
