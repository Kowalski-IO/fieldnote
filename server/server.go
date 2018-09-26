package server

import (
	"github.com/gin-gonic/gin"
	"kowalski.io/fieldnote/db"
	"strconv"
)

func Init() {
	r := gin.Default()

	r.POST("/join", func(c *gin.Context) {
		var cr db.Credentials
		c.BindJSON(&cr)
		cr.Create()
		c.JSON(200, &cr)
	})

	r.GET("/notes", func(c *gin.Context) {
		n, _ := db.FetchNotes()
		c.JSON(200, &n)
	})

	r.POST("/notes", func(c *gin.Context) {
		var n db.Note
		c.BindJSON(&n)
		n.Create()
		c.JSON(200, &n)
	})

	r.PUT("/notes/:id", func(c *gin.Context) {
		var n db.Note
		c.BindJSON(&n)

		id := c.Param("id")
		n.ID, _ = strconv.ParseInt(id, 10, 64)

		n.Update()
		c.JSON(200, &n)
	})

	r.Run()
}
