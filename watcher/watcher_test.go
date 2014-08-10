package tramtracker

import (
	"fmt"
	"github.com/danielheath/tramtracker"
	"net/http"
	"net/http/httptest"
	_ "testing"
	"time"
)

func stubTicker() chan time.Time {
	return make(chan time.Time)
}

const apiResponseWithNewTram = `{"errorMessage":null,"hasError":false,"hasResponse":true,"responseObject":[{"__type":"NextPredictedRoutesCollectionInfo","AirConditioned":false,"Destination":"Port Melbourne","DisplayAC":false,"DisruptionMessage":{"DisplayType":"Text","MessageCount":0,"Messages":[]},"HasDisruption":false,"HasSpecialEvent":true,"HeadBoardRouteNo":"109","InternalRouteNo":109,"IsLowFloorTram":true,"IsTTAvailable":true,"PredictedArrivalDateTime":"\/Date(1407550475481+1000)\/","RouteNo":"109","SpecialEventMessage":"From 27 July Route 12 operates between Stop 24 Victoria Gardens, Burnley St and St Kilda via Collins St","TripID":0,"VehicleNo":3004},{"__type":"NextPredictedRoutesCollectionInfo","AirConditioned":false,"Destination":"Port Melbourne","DisplayAC":false,"DisruptionMessage":{"DisplayType":"Text","MessageCount":0,"Messages":[]},"HasDisruption":false,"HasSpecialEvent":true,"HeadBoardRouteNo":"109","InternalRouteNo":109,"IsLowFloorTram":true,"IsTTAvailable":true,"PredictedArrivalDateTime":"\/Date(1407551124000+1000)\/","RouteNo":"109","SpecialEventMessage":"From 27 July Route 12 operates between Stop 24 Victoria Gardens, Burnley St and St Kilda via Collins St","TripID":0,"VehicleNo":3011},{"__type":"NextPredictedRoutesCollectionInfo","AirConditioned":false,"Destination":"Port Melbourne","DisplayAC":false,"DisruptionMessage":{"DisplayType":"Text","MessageCount":0,"Messages":[]},"HasDisruption":false,"HasSpecialEvent":true,"HeadBoardRouteNo":"109","InternalRouteNo":109,"IsLowFloorTram":true,"IsTTAvailable":true,"PredictedArrivalDateTime":"\/Date(1407551628000+1000)\/","RouteNo":"109","SpecialEventMessage":"From 27 July Route 12 operates between Stop 24 Victoria Gardens, Burnley St and St Kilda via Collins St","TripID":0,"VehicleNo":3035}],"timeRequested":"\/Date(1407550465405+1000)\/","timeResponded":"\/Date(1407550465481+1000)\/","webMethodCalled":"GetNextPredictedRoutesCollection"}`
const defaultApiResponse = `{"errorMessage":null,"hasError":false,"hasResponse":true,"responseObject":[{"__type":"NextPredictedRoutesCollectionInfo","AirConditioned":false,"Destination":"Port Melbourne","DisplayAC":false,"DisruptionMessage":{"DisplayType":"Text","MessageCount":0,"Messages":[]},"HasDisruption":false,"HasSpecialEvent":true,"HeadBoardRouteNo":"109","InternalRouteNo":109,"IsLowFloorTram":true,"IsTTAvailable":true,"PredictedArrivalDateTime":"\/Date(1407550475481+1000)\/","RouteNo":"109","SpecialEventMessage":"From 27 July Route 12 operates between Stop 24 Victoria Gardens, Burnley St and St Kilda via Collins St","TripID":0,"VehicleNo":3004},{"__type":"NextPredictedRoutesCollectionInfo","AirConditioned":false,"Destination":"Port Melbourne","DisplayAC":false,"DisruptionMessage":{"DisplayType":"Text","MessageCount":0,"Messages":[]},"HasDisruption":false,"HasSpecialEvent":true,"HeadBoardRouteNo":"109","InternalRouteNo":109,"IsLowFloorTram":true,"IsTTAvailable":true,"PredictedArrivalDateTime":"\/Date(1407551124000+1000)\/","RouteNo":"109","SpecialEventMessage":"From 27 July Route 12 operates between Stop 24 Victoria Gardens, Burnley St and St Kilda via Collins St","TripID":0,"VehicleNo":3011},{"__type":"NextPredictedRoutesCollectionInfo","AirConditioned":false,"Destination":"Port Melbourne","DisplayAC":false,"DisruptionMessage":{"DisplayType":"Text","MessageCount":0,"Messages":[]},"HasDisruption":false,"HasSpecialEvent":true,"HeadBoardRouteNo":"109","InternalRouteNo":109,"IsLowFloorTram":true,"IsTTAvailable":true,"PredictedArrivalDateTime":"\/Date(1407551628000+1000)\/","RouteNo":"109","SpecialEventMessage":"From 27 July Route 12 operates between Stop 24 Victoria Gardens, Burnley St and St Kilda via Collins St","TripID":0,"VehicleNo":3033}],"timeRequested":"\/Date(1407550465405+1000)\/","timeResponded":"\/Date(1407550465481+1000)\/","webMethodCalled":"GetNextPredictedRoutesCollection"}`

func tramTrackerUrl(responses []string) string {
	i := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, responses[i])
		i += 1
	}))
	return srv.URL
}

func ExampleWatcher_Start() {
	url := tramTrackerUrl([]string{defaultApiResponse, apiResponseWithNewTram})
	w := Watcher{
		// Find trams stop 1234
		Query: tramtracker.Query{StopId: 1234, Url: url},
		// Trigger OnEventRecognized when a new tram
		// appears in the response
		Recognizer: getLastTramId,
		OnEventRecognized: func(tt tramtracker.TrackerResponse) {
			fmt.Println(getLastTramId(tt))
		},
		OnErr: func(e error) { fmt.Println(e) },
	}
	w.Init(time.Minute)
	t := stubTicker()
	w.Ticker = t
	go w.Start()     // Waits for the w.Ticker
	t <- time.Time{} // A minute passes
	t <- time.Time{} // Another minute passes
	t <- time.Time{} // A third minute passes

	// Output:
	// 3033
	// 3035
}
