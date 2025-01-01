package ilosrv

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

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

	r.POST("/api/flows/exec", func(c *gin.Context) {
		projectPath := c.Query("project")
		flowName := c.Query("flow")

		project, err := ilolib.ReadProjectDefinition(projectPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
			return
		}

		flow, exists := project.Flows[flowName]

		if exists {
			daemon.RunFlow(flow)
			c.Status(http.StatusNoContent)
		} else {
			c.JSON(http.StatusBadRequest, map[string]any{
				"error": fmt.Sprintf("flow '%s' does not exist", flowName),
			})
		}
	})

	r.POST("/api/schedules", func(c *gin.Context) {
		projectPath := c.Query("project")
		flowName := c.Query("flow")
		scheduleDay, _ := strconv.Atoi(c.Query("day"))
		scheduleHour, _ := strconv.Atoi(c.Query("hour"))
		scheduleMinute, _ := strconv.Atoi(c.Query("minute"))

		project, err := ilolib.ReadProjectDefinition(projectPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
			return
		}

		flow, exists := project.Flows[flowName]

		if exists {
			schedule := ilolib.Schedule{Minute: scheduleMinute, Hour: scheduleHour, Day: time.Weekday(scheduleDay)}

			daemon.ScheduleFlow(flow, schedule)
			c.JSON(http.StatusOK, map[string]any{
				"project":  projectPath,
				"flow":     flowName,
				"schedule": schedule.String(),
			})
		} else {
			c.JSON(http.StatusBadRequest, map[string]any{
				"error": fmt.Sprintf("flow '%s' does not exist", flowName),
			})
		}
	})

	return r
}
