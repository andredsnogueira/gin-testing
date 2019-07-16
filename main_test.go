package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestVersionRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"version\":\"1.0\"}", w.Body.String())
}

func TestGetPostsRoute(t *testing.T) {
	db, err := gorm.Open("sqlite3", "sqlite.db")

	if err != nil {
		panic("Failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&Post{})

	db.Create(&Post{Title: "postTitle 1", Body: "postBody 1"})
	db.Create(&Post{Title: "postTitle 2", Body: "postBody 2"})

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	os.Remove("sqlite_test.db")
}
