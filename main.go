package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Post struct {
	gorm.Model
	Title string
	Body  string
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "sqlite.db")

	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&Post{})
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", version)
	r.GET("/posts", getPosts)
	r.GET("/posts/:id", getPost)
	r.POST("/posts", createPost)
	r.DELETE("/posts/:id", deletePost)

	return r
}

func version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "1.0"})
}

func getPosts(c *gin.Context) {
	var posts []Post

	db.Find(&posts)

	if len(posts) > 0 {
		c.JSON(http.StatusOK, gin.H{"data": posts})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"data": "No posts found."})
}

func getPost(c *gin.Context) {
	var post Post
	postID := c.Params.ByName("id")

	db.First(&post, postID)

	if post.ID != 0 {
		c.JSON(http.StatusOK, gin.H{"data": post})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"data": "Post not found."})
}

func createPost(c *gin.Context) {
	postTitle := c.PostForm("title")
	postBody := c.PostForm("body")

	db.Create(&Post{Title: postTitle, Body: postBody})

	c.JSON(http.StatusCreated, gin.H{"data": "Post created."})
}

func deletePost(c *gin.Context) {
	var post Post
	postID := c.Params.ByName("id")

	db.First(&post, postID)

	if post.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"data": "Not found."})
		return
	}

	db.Delete(&post)

	c.JSON(http.StatusOK, gin.H{"data": "Post deleted."})
}
