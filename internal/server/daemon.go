package server

import (
	"log/slog"
	"time"

	"github.com/fourls/ilo/internal/data"
	"github.com/fourls/ilo/internal/data/toolbox"
	"github.com/fourls/ilo/internal/exec"
	"github.com/fourls/ilo/internal/ilofile"
)

type scheduledFlow struct {
	flow     ilofile.Flow
	schedule data.Schedule
}

type IloDaemon struct {
	ticker        *time.Ticker
	log           *slog.Logger
	toolbox       toolbox.Toolbox
	flowSchedules []scheduledFlow
}

func (d *IloDaemon) Run() {
	d.ticker = time.NewTicker(time.Minute)
	go d.worker()
}

func (d *IloDaemon) RunFlow(flow ilofile.Flow) {
	observer := newObserver(flow.Project, d.log)
	go exec.RunFlow(flow, exec.RunStep, d.toolbox, &observer)
}

func (d *IloDaemon) ScheduleFlow(flow ilofile.Flow, schedule data.Schedule) {
	d.flowSchedules = append(d.flowSchedules, scheduledFlow{
		flow:     flow,
		schedule: schedule,
	})
}

func (d *IloDaemon) worker() {
	if d.ticker == nil {
		return
	}

	for time := range d.ticker.C {
		d.tick(time)
	}
}

func (d *IloDaemon) tick(now time.Time) {
	for _, entry := range d.flowSchedules {
		if entry.schedule.Match(now) {
			d.RunFlow(entry.flow)
		}
	}
}
