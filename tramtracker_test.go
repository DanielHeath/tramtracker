package tramtracker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	_ "testing"
	"time"
)

func init() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, cannedTramQuery)
	}))
	TramTrackerUrl = srv.URL + "/%d_%d_%t"
	now = func() time.Time {
		return time.Unix(1407550175, 0)
	}
}

func ExampleTramTime() {
	trams, _ := NextTrams(Query{StopId: 1234})
	for _, tram := range trams.Upcoming {
		minutesAway := int(tram.PredictedTime.Sub(now()).Minutes())
		arrival := tram.PredictedTime.Format("3:04pm")
		fmt.Printf("%s in %d minutes (%s)\n", tram.RouteNo, minutesAway, arrival)
	}
	// Output:
	// 109 in 5 minutes (12:14pm)
	// 109 in 15 minutes (12:25pm)
	// 109 in 24 minutes (12:33pm)
}

func ExampleWaitTime() {
	trams, _ := NextTrams(Query{StopId: 1234})

	// If I run and take 4 minutes to get to the stop
	fmt.Println(trams.WaitTime(time.Minute * 4))

	// If I walk and take 6 minutes to get to the stop
	fmt.Println(trams.WaitTime(time.Minute * 6))

	// Output:
	// 1m0s
	// 9m49s
}

func ExampleAlert() {
	trams, _ := NextTrams(Query{StopId: 1234})

	droughtStart, droughtEnd, longWait := trams.AnyWaitOver(time.Minute * 8)
	if longWait {
		fmt.Printf("Long wait: there's a tram in %s, then nothing for %s\n", droughtStart, droughtEnd-droughtStart)
	}

	droughtStart, droughtEnd, longWait = trams.AnyWaitOver(time.Minute * 12)
	if !longWait {
		fmt.Printf("No wait over 12 minutes\n")
	}

	// Output:
	// Long wait: there's a tram in 5m0s, then nothing for 10m49s
	// No wait over 12 minutes
}

const cannedTramQuery = `
{"errorMessage":null,"hasError":false,"hasResponse":true,"responseObject":[{"__type":"NextPredictedRoutesCollectionInfo","AirConditioned":false,"Destination":"Port Melbourne","DisplayAC":false,"DisruptionMessage":{"DisplayType":"Text","MessageCount":0,"Messages":[]},"HasDisruption":false,"HasSpecialEvent":true,"HeadBoardRouteNo":"109","InternalRouteNo":109,"IsLowFloorTram":true,"IsTTAvailable":true,"PredictedArrivalDateTime":"\/Date(1407550475481+1000)\/","RouteNo":"109","SpecialEventMessage":"From 27 July Route 12 operates between Stop 24 Victoria Gardens, Burnley St and St Kilda via Collins St","TripID":0,"VehicleNo":3004},{"__type":"NextPredictedRoutesCollectionInfo","AirConditioned":false,"Destination":"Port Melbourne","DisplayAC":false,"DisruptionMessage":{"DisplayType":"Text","MessageCount":0,"Messages":[]},"HasDisruption":false,"HasSpecialEvent":true,"HeadBoardRouteNo":"109","InternalRouteNo":109,"IsLowFloorTram":true,"IsTTAvailable":true,"PredictedArrivalDateTime":"\/Date(1407551124000+1000)\/","RouteNo":"109","SpecialEventMessage":"From 27 July Route 12 operates between Stop 24 Victoria Gardens, Burnley St and St Kilda via Collins St","TripID":0,"VehicleNo":3011},{"__type":"NextPredictedRoutesCollectionInfo","AirConditioned":false,"Destination":"Port Melbourne","DisplayAC":false,"DisruptionMessage":{"DisplayType":"Text","MessageCount":0,"Messages":[]},"HasDisruption":false,"HasSpecialEvent":true,"HeadBoardRouteNo":"109","InternalRouteNo":109,"IsLowFloorTram":true,"IsTTAvailable":true,"PredictedArrivalDateTime":"\/Date(1407551628000+1000)\/","RouteNo":"109","SpecialEventMessage":"From 27 July Route 12 operates between Stop 24 Victoria Gardens, Burnley St and St Kilda via Collins St","TripID":0,"VehicleNo":3033}],"timeRequested":"\/Date(1407550465405+1000)\/","timeResponded":"\/Date(1407550465481+1000)\/","webMethodCalled":"GetNextPredictedRoutesCollection"}
`
