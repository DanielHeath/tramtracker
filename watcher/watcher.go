package tramtracker

import (
	"errors"
	"github.com/danielheath/tramtracker"
	"time"
)

type Gap struct {
	StartsIn time.Duration
	EndsIn   time.Duration
}

// Yuck; this isn't really designed so much as
// jammed together :/
func NewGapWatcher(q tramtracker.Query, gapsize time.Duration) (Watcher, chan Gap) {
	out := make(chan Gap)
	return Watcher{
		Query:             q,
		Recognizer:        getLastTramId,
		OnEventRecognized: gapIdentifier(gapsize, out),
	}, out
}

func gapIdentifier(gapsize time.Duration, out chan Gap) func(tramtracker.TrackerResponse) {
	return func(tt tramtracker.TrackerResponse) {
		start, end, drought := tt.AnyWaitOver(gapsize)
		if drought {
			out <- Gap{start, end}
		}
	}
}

// Watches the tramtracker API
//
// I want to be notified when a gap is coming up.
// I should only be notified for a gap
// the first time it occurs
type Watcher struct {
	Query             tramtracker.Query
	Recognizer        eventRecognizer
	OnEventRecognized func(tramtracker.TrackerResponse)
	OnErr             func(error)
	Ticker            <-chan time.Time
	alreadySeen       map[interface{}]bool
}

type eventRecognizer func(tramtracker.TrackerResponse) interface{}

func (w *Watcher) Init(interval time.Duration) {
	if w.Recognizer == nil {
		panic(errors.New("No event recognizer"))
	}
	if w.OnEventRecognized == nil {
		panic(errors.New("No OnEventRecognized handler"))
	}
	w.alreadySeen = make(map[interface{}]bool)
	w.Ticker = time.NewTicker(interval).C
}

func (w *Watcher) Start() {
	for _ = range w.Ticker {
		w.update()
	}
}

func (w Watcher) update() {
	next, err := tramtracker.NextTrams(w.Query)
	if err != nil {
		if w.OnErr != nil {
			w.OnErr(err)
		} else {
			panic(err)
		}
		return
	}
	resp := w.Recognizer(*next)
	if resp == nil {
		return
	}
	_, alreadySeen := w.alreadySeen[resp]
	if alreadySeen {
		return
	}
	w.alreadySeen[resp] = true
	w.OnEventRecognized(*next)
}

func getLastTramId(tt tramtracker.TrackerResponse) interface{} {
	return tt.Upcoming[len(tt.Upcoming)-1].VehicleNo
}

func idealWiring() {
	// describe a more useful invocation style here
	// Idea: transform it into a list of wait times.
	// E.G. if there's a tram in 3 and 8 minutes then
	// there is a 3 minute wait for tramid #3023, then
	// a 5 minute wait for tramid #1039.

	// Running an update gives me a new list,
	// but since I can track which vehicles I've
	// warned you about waits for, I can avoid
	// repeating myself.
	// I'll need to 'forget' them when they disappear
	// from the list.

	// watchedQuery := tramtracker.WatchQuery(Query{}, time.Minute)
	// watchedQuery.OnChange = func(change WatchedQueryChange) {
	// 	if change.PredictedTime {
	// 		change.PredictedTime
	// 	}
	// 	t := &Tram{}
	// 	for _, tram := range change.ExpectedTrams {
	// 		// a new tram appears
	// 		if t != nil {
	// 			tram.PredictedTime
	// 		}
	// 		t = &tram
	// 	}
	// }
	// watchedQuery.Start()
}
