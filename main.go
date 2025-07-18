package main

import (
	"database/sql"
	"github.com/PlamenPenkov123/gin-gonic-intro/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var db *sql.DB

func main() {
	db, err := DatabaseInit(db)
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/albums", func(c *gin.Context) {
		albums, err := models.GetAlbums(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve albums"})
			return
		}
		c.JSON(http.StatusOK, albums)
	})
	router.GET("/albums/:id", func(c *gin.Context) {
		var album models.Album
		id := c.Param("id")
		albumID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
			return
		}
		album, err = models.GetAlbumById(albumID, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve album"})
			return
		}
		if (album == models.Album{}) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Album not found"})
			return
		}
		c.JSON(http.StatusOK, album)
	})
	router.GET("/albums/artist", func(c *gin.Context) {
		name := c.Query("name")
		albums, err := models.GetAlbumsByArtist(name, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve albums by artist"})
			return
		}
		if len(albums) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No albums found for this artist"})
			return
		}
		c.JSON(http.StatusOK, albums)
	})
	router.POST("/albums", func(c *gin.Context) {
		var album models.Album
		if err := c.ShouldBindJSON(&album); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if _, err := models.AddAlbum(album, db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add album"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Album added successfully"})
	})
	router.PUT("/albums/:id", func(c *gin.Context) {
		id := c.Param("id")
		newAlbum := models.Album{}
		if err := c.ShouldBindJSON(&newAlbum); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		albumID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
			return
		}
		_, err = models.UpdateAlbum(albumID, newAlbum, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update album"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Album updated successfully"})
	})
	router.DELETE("/albums/:id", func(c *gin.Context) {
		id := c.Param("id")
		albumID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
			return
		}
		rowsAffected, err := models.DeleteAlbum(albumID, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete album"})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No album found with this ID"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
	})
	err = router.Run(":8080")
	if err != nil {
		return
	} // Run on port 8080
}
