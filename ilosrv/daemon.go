package ilosrv

import (
	"log/slog"
	"time"

	"fourls.dev/ilo/ilolib"
)

type IloDaemon struct {
	ticker  *time.Ticker
	log     *slog.Logger
	toolbox ilolib.Toolbox
}

func (d *IloDaemon) Run() {
	d.ticker = time.NewTicker(5 * time.Second)
	go d.worker()
}

func (d *IloDaemon) RunFlow(project ilolib.ProjectDefinition, name string) {
	observer := newObserver(&project, d.log)

	go func() {
		_, err := ilolib.ProjectExecutor{
			Definition: project,
			Toolbox:    d.toolbox,
		}.RunFlow(name, &observer, ilolib.BuildDefaultExecutor)

		if err != nil {
			d.log.Error("Executor failed", "project", project.Path, "error", err)
		}
	}()
}

func (d *IloDaemon) worker() {
	if d.ticker == nil {
		return
	}

	for time := range d.ticker.C {
		d.tick(time)
	}
}

func (d *IloDaemon) tick(_ time.Time) {
	// todo: run scheduled jobs
}
