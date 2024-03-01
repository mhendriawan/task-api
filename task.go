package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Task struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var (
	users []User
	tasks []Task
)

func main() {
	router := gin.Default()

	// Middleware for logging
	router.Use(gin.Logger())

	// Middleware for recovering from panics
	router.Use(gin.Recovery())

	// User endpoints
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", createUser)
		userGroup.GET("/", getUsers)
		userGroup.GET("/:id", getUserByID)
		userGroup.PUT("/:id", updateUser)
		userGroup.DELETE("/:id", deleteUser)
	}

	// Task endpoints
	// Secure task endpoints with OAuth 2
	taskGroup := router.Group("/tasks")
	taskGroup.Use(authMiddleware)
	{
		taskGroup.POST("", createTask)
		taskGroup.GET("", getTasks)
		taskGroup.GET("/:id", getTaskByID)
		taskGroup.PUT("/:id", updateTask)
		taskGroup.DELETE("/:id", deleteTask)
	}

	router.Run(":8080")
}

// Middleware to authenticate requests
func authMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Missing token"})
		c.Abort()
		return
	}

	// Validate OAuth 2 token
	userInfo, err := getUserInfoFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Unauthorized: %v", err)})
		c.Abort()
		return
	}

	// Set user info in the context for downstream handlers to access
	c.Set("userInfo", userInfo)

	c.Next()
}

// Dummy function to validate OAuth 2 token and fetch user info
func getUserInfoFromToken(token string) (*User, error) {
	dummyUser := &User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "",
	}
	return dummyUser, nil
}

// User handlers
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Assign a unique ID
	user.ID = uint(len(users) + 1)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	users = append(users, user)
	c.JSON(http.StatusCreated, user)
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 || id > len(users) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user := users[id-1]
	c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 || id > len(users) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedUser.ID = uint(id)
	updatedUser.CreatedAt = users[id-1].CreatedAt
	updatedUser.UpdatedAt = time.Now()
	users[id-1] = updatedUser
	c.JSON(http.StatusOK, updatedUser)
}

func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 || id > len(users) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	users = append(users[:id-1], users[id:]...)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// Task handlers
func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Assign a unique ID
	task.ID = uint(len(tasks) + 1)
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	tasks = append(tasks, task)
	c.JSON(http.StatusCreated, task)
}

func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

func getTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 || id > len(tasks) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	task := tasks[id-1]
	c.JSON(http.StatusOK, task)
}

func updateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 || id > len(tasks) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	var updatedTask Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedTask.ID = uint(id)
	updatedTask.CreatedAt = tasks[id-1].CreatedAt
	updatedTask.UpdatedAt = time.Now()
	tasks[id-1] = updatedTask
	c.JSON(http.StatusOK, updatedTask)
}

func deleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 || id > len(tasks) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	tasks = append(tasks[:id-1], tasks[id:]...)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
