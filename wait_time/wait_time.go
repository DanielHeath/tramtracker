package wait_time

import (
	"github.com/danielheath/tramtracker"
	"time"
)

type WaitTime struct {
	time.Duration
	tramtracker.Tram
}

var now func() time.Time

func init() {
	now = time.Now
}

func WaitTimes(tt tramtracker.TrackerResponse) []WaitTime {
	result := make([]WaitTime, len(tt.Upcoming))
	waitFrom := now()
	for i, tram := range tt.Upcoming {
		result[i] = WaitTime{
			tram.PredictedTime.Sub(waitFrom),
			tram,
		}
		waitFrom = tram.PredictedTime.Time
	}
	return result
}

type WaitTimeWarner struct {
	alreadyWarned    map[int]bool
	warnIfLongerThan time.Duration
	Warnings         chan WaitTime
}

func NewWaitTimeWarner(max time.Duration) WaitTimeWarner {
	return WaitTimeWarner{
		alreadyWarned:    make(map[int]bool),
		warnIfLongerThan: max,
		Warnings:         make(chan WaitTime),
	}
}

func (w WaitTimeWarner) Warn(tt tramtracker.TrackerResponse) {
	for _, wt := range WaitTimes(tt) {
		if wt.Duration > w.warnIfLongerThan {
			if w.alreadyWarned[wt.VehicleNo] {
				// No need to do anything
			} else {
				w.alreadyWarned[wt.VehicleNo] = true
				w.Warnings <- wt
			}
		}
	}

	// Trams go back & forth all day.
	// Once a tram disappears from the upcoming
	// list, it will return in an hour or so
	// and we want to notice it when that happens.
	// Therefore we remove any entries from
	// alreadyWarned except for the ones
	// on the current tram list.
	for k, _ := range w.alreadyWarned {
		del := true
		for _, tram := range tt.Upcoming {
			if k == tram.VehicleNo {
				del = false
			}
		}
		if del {
			delete(w.alreadyWarned, k)
		}
	}
}
