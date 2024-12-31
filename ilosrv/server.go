package ilosrv

import (
	"log/slog"
	"net/http"
	"os"

	"fourls.dev/ilo/ilolib"
	"github.com/gin-gonic/gin"
)

func BuildServer() *gin.Engine {
	r := gin.Default()

	toolbox, _ := ilolib.NewProdToolbox()

	daemon := IloDaemon{
		toolbox: *toolbox,
		log:     slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
	daemon.Run()

	r.POST("/projects/build", func(c *gin.Context) {
		projectPath := c.PostForm("project")
		flowNames := c.PostFormArray("flows")

		project, err := ilolib.ReadProjectDefinition(projectPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
			return
		}

		for _, flowName := range flowNames {
			daemon.RunFlow(*project, flowName)
		}

		c.Status(http.StatusNoContent)
	})

	return r
}
