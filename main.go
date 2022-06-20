package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Ticket struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	EventID int    `json:"event_id"`
	Type    string `json:"type"`
}

type Event struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	StartDateTime string   `json:"start_datetime"`
	EndDateTime   string   `json:"end_datetime"`
	ValidTypes    []string `json:"valid_types"`
}

func CompareEvents(a, b Event) bool {
	if len(a.ValidTypes) == len(b.ValidTypes) {
		for i, t := range a.ValidTypes {
			if b.ValidTypes[i] != t {
				return false
			}
		}
		return a.ID == b.ID && a.Name == b.Name && a.StartDateTime == b.StartDateTime && a.EndDateTime == b.EndDateTime
	}
	return false
}

var tickets = []Ticket{
	{ID: 1, Name: "John Doe", Price: 1399, EventID: 1, Type: "A"},
	{ID: 2, Name: "Jack Smith", Price: 999, EventID: 1, Type: "A"},
	{ID: 3, Name: "Ivan Ivanovich", Price: 1399, EventID: 1, Type: "B"},
}

var events = []Event{
	{ID: 1, Name: "The Show", StartDateTime: "6/4/2022", EndDateTime: "6/6/2022", ValidTypes: []string{"A", "B", "C"}},
}

func getEvent(t Ticket) (Event, bool) {
	for _, e := range events {
		if e.ID == t.EventID {
			return e, true
		}
	}
	return Event{}, false
}

func validTicketEvent(t Ticket) bool {
	for _, e := range events {
		if e.ID == t.EventID {
			for _, ty := range e.ValidTypes {
				if ty == t.Type {
					return true
				}
			}
			return false
		}
	}
	return false
}

/*
func getEventTickets(e Event) []Ticket {
	return make([]Ticket, 0)
}*/

func getTickets(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, tickets)
}

func getTicketsForEvent(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	for _, a := range events {
		if a.ID == id {
			res := make([]Ticket, 0)
			for _, t := range tickets {
				if t.EventID == a.ID {
					res = append(res, t)
				}
			}
			ctx.IndentedJSON(http.StatusOK, res)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "event ID " + ctx.Param("id") + " not found"})
}

func getTicketById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	for _, a := range tickets {
		if a.ID == id {
			ctx.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "ticket ID " + ctx.Param("id") + " not found"})
}

func deleteTicketById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	for i, a := range tickets {
		if a.ID == id {
			tickets[i] = tickets[len(tickets)-1]
			tickets = tickets[:len(tickets)-1]
			ctx.IndentedJSON(http.StatusOK, nil)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "ticket ID " + ctx.Param("id") + " not found"})
}

func updateTicket(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var newTicket Ticket
	if err := ctx.BindJSON(&newTicket); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to parse JSON to struct"})
		return
	}
	newTicket.ID = id
	for i, t := range tickets {
		if t.ID == id {
			e, val := getEvent(newTicket)
			if val {
				for _, ty := range e.ValidTypes {
					if newTicket.Type == ty {
						tickets[i] = newTicket
						ctx.IndentedJSON(http.StatusAccepted, newTicket)
						return
					}
				}
				ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ticket type for event"})

			} else {
				ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "event for ticket  not found"})
			}

			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "ticket ID " + ctx.Param("id") + " not found"})
}

func postTickets(ctx *gin.Context) {
	var newTicket Ticket
	if err := ctx.BindJSON(&newTicket); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to parse JSON to struct"})
		return
	}
	ids := make([]int, len(tickets))
	for i, t := range tickets {
		ids[i] = t.ID
	}
	newTicket.ID = generateNewID(ids)
	if validTicketEvent(newTicket) {
		tickets = append(tickets, newTicket)
		ctx.Header("Location", ctx.RemoteIP()+"/tickets/"+strconv.Itoa(newTicket.ID))
		ctx.IndentedJSON(http.StatusCreated, newTicket)
		return

	} else {
		ctx.IndentedJSON(http.StatusBadRequest,
			gin.H{"message": "Event (" + strconv.Itoa(newTicket.EventID) + ") not found, or bad ticket type (" + newTicket.Type + ") for event"})
	}
}

func getEvents(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, events)
}

func getEventsById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	for _, a := range events {
		if a.ID == id {
			ctx.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "event ID " + ctx.Param("id") + " not found"})
}

func deleteEvent(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	for i, a := range events {
		if a.ID == id {
			j := 0
			for j < len(tickets) {
				if tickets[j].EventID == a.ID {
					tickets[j] = tickets[len(tickets)-1]
					tickets = tickets[:len(tickets)-1]
				}
				j++
			}
			events[i] = events[len(events)-1]
			events = events[:len(events)-1]
			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "event ID and associated tickets " + ctx.Param("id") + " deleted"})
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "event ID " + ctx.Param("id") + " not found"})
}

func generateNewID(ids []int) int {
	id := 0
	sort.Ints(ids)
	for _, a := range ids {
		if a == id {
			id++
		}
	}
	return id
}

func postEvents(ctx *gin.Context) {
	var newEvent Event
	if err := ctx.BindJSON(&newEvent); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to parse JSON to struct"})
		return
	}
	ids := make([]int, len(events))
	for i, e := range events {
		ids[i] = e.ID
	}
	newEvent.ID = generateNewID(ids)

	events = append(events, newEvent)
	ctx.IndentedJSON(http.StatusCreated, newEvent)
	ctx.Header("Location", ctx.RemoteIP()+"/events/"+strconv.Itoa(newEvent.ID))
}

func main() {
	router := gin.Default()
	router.GET("/tickets", getTickets)
	router.GET("/tickets/:id", getTicketById)
	router.DELETE("/tickets/:id", deleteTicketById)
	router.POST("/tickets", postTickets)
	router.PUT("/tickets/:id", updateTicket)

	router.GET("/events", getEvents)
	router.GET("/events/:id", getEventsById)
	router.GET("/events/:id/tickets", getTicketsForEvent)
	router.DELETE("/events/:id", deleteEvent)
	router.POST("/events", postEvents)
	router.Run("localhost:8080")
}
