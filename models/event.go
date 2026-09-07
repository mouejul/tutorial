package models

import "time"

type Event struct {
	Id          int
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Location    string `binding:"required"`
	DatTime     time.Time
	UserId      int
}

// variabel penyimpanan data
var events []Event = []Event{}

// fungsi unutk simpan Event
func (e Event) Save(){
	events = append(events,e)
}

// fungsi menampilkan semua event
func GetAllEvents() []Event {
	return events
}