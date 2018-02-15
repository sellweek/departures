package main

import (
	"fmt"
	"github.com/patrickbr/gtfsparser"
	"github.com/patrickbr/gtfsparser/gtfs"
	"github.com/sellweek/departures/departures"
	"time"
)

func main() {
	feed := gtfsparser.NewFeed()
	err := feed.Parse("/home/selvek/Dev/go/src/github.com/sellweek/departures/pid")
	if err != nil {
		panic(err)
	}
	stop, err := departures.NewStopDepartures("U699Z2", feed)
	if err != nil {
		panic(err)
	}
	result := stop.After(time.Now(), 75)
	for _, dep := range result {
		departureTime := dep.At.GetLocationTime(gtfs.GetGtfsDateFromTime(time.Now()), dep.Trip.Route.Agency)
		fmt.Printf("%s -> %s %dm\n", dep.Trip.Route.Short_name, dep.Trip.Headsign, int(time.Now().Sub(departureTime).Minutes()))
	}
}
