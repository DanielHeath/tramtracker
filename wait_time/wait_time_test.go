package wait_time

import (
	"fmt"
	"github.com/danielheath/tramtracker"
	_ "testing"
	"time"
)

const startTime = 1407550175

func init() {
	now = func() time.Time {
		return time.Unix(startTime, 0)
	}
}

func trams(trams ...tramtracker.Tram) tramtracker.TrackerResponse {
	return tramtracker.TrackerResponse{
		Upcoming: trams,
	}
}

func ExampleWaitTimeWarner() {
	// Warn me if there's more than 10 mins between trams.
	w := NewWaitTimeWarner(time.Minute * 10)
	done := make(chan bool)
	go func() {
		for waitTime := range w.Warnings {
			fmt.Println(waitTime.Duration, waitTime.VehicleNo)
		}
		done <- true
	}()
	w.Warn(trams(
		tramIn(time.Second*30, 1),
		tramIn(time.Minute*11, 2),
		tramIn(time.Minute*15, 3),
	))
	w.Warn(trams(
		tramIn(time.Minute*1, 2),
		tramIn(time.Minute*5, 3),
		tramIn(time.Minute*17, 4),
	))
	w.Warn(trams(
		tramIn(time.Minute*1, 3),
		tramIn(time.Minute*2, 4),
		tramIn(time.Minute*17, 1),
	))
	w.Warn(trams(
		tramIn(time.Minute*2, 4),
		tramIn(time.Minute*5, 1),
		tramIn(time.Minute*17, 2),
	))
	close(w.Warnings)
	<-done // wait for channel to be processed

	// Output:
	// 10m30s 2
	// 12m0s 4
	// 15m0s 1
	// 12m0s 2
}

func tramIn(t time.Duration, id int) tramtracker.Tram {
	return tramtracker.Tram{
		PredictedTime: tramtracker.TramTrackerDate{now().Add(t)},
		VehicleNo:     id,
	}
}
func ExampleWaitTimes() {
	tt := trams(
		tramIn(time.Second*35, 1),
		tramIn(time.Second*105, 2),
		tramIn(time.Second*185, 3),
	)
	times := WaitTimes(tt)
	fmt.Println(times[0])
	fmt.Println(times[1])
	fmt.Println(times[2])
	// Output:
	// 35s
	// 1m10s
	// 1m20s
}
