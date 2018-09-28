package server

import (
	"github.com/gin-gonic/gin"
	"kowalski.io/fieldnote/db"
)

func Init() {
	r := gin.Default()

	r.POST("/join", func(c *gin.Context) {
		var cr db.Credentials
		c.BindJSON(&cr)
		u, _ := cr.Create()
		c.JSON(200, &u)
	})

	r.GET("/notes", func(c *gin.Context) {
		n, _ := db.FetchNotes()
		c.JSON(200, &n)
	})

	r.POST("/notes", func(c *gin.Context) {
		var n db.Note
		c.BindJSON(&n)

		r, _ := n.Upsert()

		c.JSON(200, &r)
	})

	r.Run(":9090")
}
