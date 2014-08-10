// Package tramtracker fetches tram arrival data from the tramtracker JSON api.
//
// MIT licenced, (c) Daniel Heath

package tramtracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type TrackerResponse struct {
	WebMethodCalled string          `json:"webMethodCalled"`
	HasResponse     bool            `json:"hasResponse"`
	TimeResponded   TramTrackerDate `json:"timeResponded"`
	TimeRequested   TramTrackerDate `json:"timeRequested"`
	HasError        bool            `json:"hasError"`
	ErrorMessage    *string         `json:"errorMessage"`
	Upcoming        []Tram          `json:"responseObject"`
}

type Tram struct {
	IsTTAvailable       bool            `json:"IsTTAvailable"`
	InternalRouteNo     int             `json:"InternalRouteNo"`
	HasSpecialEvent     bool            `json:"HasSpecialEvent"`
	DisplayAC           bool            `json:"DisplayAC"`
	DisruptionMessage   Disruption      `json:"DisruptionMessage"`
	VehicleNo           int             `json:"VehicleNo"`
	Destination         string          `json:"Destination"`
	RouteNo             string          `json:"RouteNo"`
	IsLowFloorTram      bool            `json:"IsLowFloorTram"`
	TripID              int             `json:"TripID"`
	PredictedTime       TramTrackerDate `json:"PredictedArrivalDateTime"`
	SpecialEventMessage string          `json:"SpecialEventMessage"`
	HasDisruption       bool            `json:"HasDisruption"`
	HeadBoardRouteNo    string          `json:"HeadBoardRouteNo"`
	AirConditioned      bool            `json:"AirConditioned"`
	Type                string          `json:"__type"`
}

type Disruption struct {
	DisplayType  string
	MessageCount int
	// Messages [] // Not sure what the format of this array is.
}

type TramTrackerDate struct {
	time.Time
}

func (d *TramTrackerDate) UnmarshalJSON(json []byte) error {
	n, err := strconv.ParseInt(string(json[8:18]), 10, 64)
	if err != nil {
		return err
	}
	d.Time = time.Unix(n, 0)
	return nil
}

type Query struct {
	StopId   int
	RouteNo  int    // Optional
	LowFloor bool   // Optional
	Url      string // Optional; overrides all other query params
}

func NextTrams(q Query) (*TrackerResponse, error) {
	if q.Url == "" {
		q.Url = fmt.Sprintf(defaultUrl, q.StopId, q.RouteNo, q.LowFloor)
	}
	resp, err := http.Get(q.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := TrackerResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

var now func() time.Time

const defaultUrl = "http://tramtracker.com.au/Controllers/GetNextPredictionsForStop.ashx?stopNo=%d&routeNo=%d&isLowFloor=%t"

func init() {
	now = time.Now
}

func (r TrackerResponse) WaitTime(startAfter time.Duration) time.Duration {
	ignoreTramsBefore := now().Add(startAfter)
	for _, tram := range r.Upcoming {
		if tram.PredictedTime.After(ignoreTramsBefore) {
			return tram.PredictedTime.Sub(now()) - startAfter
		}
	}
	return time.Hour * 48
}

func (r TrackerResponse) AnyWaitOver(maxWait time.Duration) (droughtStart time.Duration, droughtEnd time.Duration, longWait bool) {
	for idx, tram := range r.Upcoming {
		// Is there another tram after this one?
		if idx < len(r.Upcoming)-1 {
			followingTram := r.Upcoming[idx+1]
			wait := followingTram.PredictedTime.Sub(tram.PredictedTime.Time)
			if wait > maxWait {
				droughtStart = tram.PredictedTime.Sub(now())
				droughtEnd = followingTram.PredictedTime.Sub(now())
				longWait = true
				return
			}
		}
	}
	return
}
