package departures

import (
	"github.com/patrickbr/gtfsparser"
	"github.com/patrickbr/gtfsparser/gtfs"
	"sort"
	"time"
)

var pragueLocation, _ = time.LoadLocation("Europe/Prague")

// StopDepartures allows live tracking of line departures
// from a given public transport stop described in a
// GTFS feed.
type StopDepartures struct {
	Stop        *gtfs.Stop
	Departures  Departures
	lastService gtfs.Date
}

// NewStopDepartures constructs a new Departures struct using
// necessary data from passed feed
func NewStopDepartures(id string, feed *gtfsparser.Feed) (d StopDepartures, err error) {
	d = StopDepartures{
		Stop:        feed.Stops[id],
		Departures:  make([]Departure, 0),
		lastService: gtfs.Date{0, 0, 0},
	}
	for _, t := range feed.Trips {
		includes := false
		at := gtfs.Time{}
		for _, stop := range t.StopTimes {
			if stop.Stop.Id == id {
				includes = true
				at = stop.Departure_time
				break
			}
		}
		if includes {
			d.Departures = append(d.Departures, Departure{at, t})
			if dateAfter(t.Service.End_date, d.lastService) {
				d.lastService = t.Service.End_date
			}
		}
	}
	sort.Sort(d.Departures)
	return
}

func (d StopDepartures) After(t time.Time, limit int) (result Departures) {
	result = make([]Departure, 0, limit)
	timeInPrague := t.In(pragueLocation)
	currentDate := gtfs.GetGtfsDateFromTime(timeInPrague)
	startTime := gtfs.Time{
		Hour:   int8(timeInPrague.Hour()),
		Minute: int8(timeInPrague.Minute()),
		Second: int8(timeInPrague.Second()),
	}
	i := 0
	for d.Departures[i].At.Minus(startTime) < 0 && i < len(d.Departures) {
		i++
	}
	if i == len(d.Departures) {
		i = 0
		currentDate = incrementDate(currentDate)
	}
	for !dateAfter(currentDate, d.lastService) && len(result) < limit {
		for i < len(d.Departures) && len(result) < limit {
			if d.Departures[i].Trip.Service.IsActiveOn(currentDate) {
				result = append(result, d.Departures[i])
			}
			i++
		}
		currentDate = incrementDate(currentDate)
		i = 0
	}
	return
}

type Departure struct {
	At   gtfs.Time
	Trip *gtfs.Trip
}

type Departures []Departure

func (d Departures) Len() int {
	return len(d)
}

func (d Departures) Less(i, j int) bool {
	atI := d[i].At
	atJ := d[j].At
	return atI.Hour < atJ.Hour ||
		(atI.Hour == atJ.Hour && atI.Minute < atJ.Minute) ||
		(atI.Hour == atJ.Hour && atI.Minute == atJ.Minute && atI.Second < atJ.Second)
}

func (d Departures) Swap(i, j int) {
	temp := d[i]
	d[i] = d[j]
	d[j] = temp
}

func dateAfter(x, y gtfs.Date) bool {
	return y.Year < x.Year ||
		(y.Year == x.Year && y.Month < x.Month) ||
		(y.Year == x.Year && y.Month == x.Month && y.Day < x.Day)
}

func incrementDate(d gtfs.Date) gtfs.Date {
	return gtfs.GetGtfsDateFromTime(d.GetTime().Add(24 * time.Hour))
}
