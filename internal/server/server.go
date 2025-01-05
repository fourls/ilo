package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fourls/ilo/internal/data"
	"github.com/fourls/ilo/internal/data/provide"
	"github.com/fourls/ilo/internal/data/toolbox"
	"github.com/fourls/ilo/internal/ilofile/iloyml"
	"github.com/gin-gonic/gin"
)

func BuildServer(provider provide.Provider[toolbox.Toolbox]) *gin.Engine {
	r := gin.Default()

	toolbox, _ := provider.Load(
		"toolbox",
		provide.YamlUnmarshal[toolbox.Toolbox])

	daemon := IloDaemon{
		toolbox: *toolbox,
		log:     slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
	daemon.Run()

	r.POST("/api/flows/exec", func(c *gin.Context) {
		projectPath := c.Query("project")
		flowName := c.Query("flow")

		project, err := iloyml.New(projectPath)
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

		project, err := iloyml.New(projectPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
			return
		}

		flow, exists := project.Flows[flowName]

		if exists {
			schedule := data.Schedule{Minute: scheduleMinute, Hour: scheduleHour, Day: time.Weekday(scheduleDay)}

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
