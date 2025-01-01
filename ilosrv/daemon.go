package ilosrv

import (
	"log/slog"
	"time"

	"fourls.dev/ilo/ilolib"
)

type scheduledFlow struct {
	flow     ilolib.Flow
	schedule ilolib.Schedule
}

type IloDaemon struct {
	ticker        *time.Ticker
	log           *slog.Logger
	toolbox       ilolib.Toolbox
	flowSchedules []scheduledFlow
}

func (d *IloDaemon) Run() {
	d.ticker = time.NewTicker(time.Minute)
	go d.worker()
}

func (d *IloDaemon) RunFlow(flow ilolib.Flow) {
	observer := newObserver(flow.Project, d.log)
	go ilolib.FlowExecutor{
		Toolbox: d.toolbox,
	}.RunFlow(flow, &observer)
}

func (d *IloDaemon) ScheduleFlow(flow ilolib.Flow, schedule ilolib.Schedule) {
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
