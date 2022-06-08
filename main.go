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

type event_t struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	StartDateTime string   `json:"start_datetime"`
	EndDateTime   string   `json:"end_datetime"`
	ValidTypes    []string `json:"valid_types"`
}

func CompareEvents(a, b event_t) bool {
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

var events = []event_t{
	{ID: 1, Name: "The Show", StartDateTime: "6/4/2022", EndDateTime: "6/6/2022", ValidTypes: []string{"A", "B", "C"}},
}

func getTickets(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, tickets)
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
	for _, e := range events {
		if e.ID == newTicket.EventID {
			for _, t := range e.ValidTypes {
				if newTicket.Type == t {
					tickets = append(tickets, newTicket)
					ctx.IndentedJSON(http.StatusCreated, newTicket)
					return
				}
			}
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ticket type (" + newTicket.Type + ") for event (" + strconv.Itoa(e.ID) + ")"})
			return
		}
	}
	ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Event (" + strconv.Itoa(newTicket.ID) + ") not found"})
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
	var newEvent event_t
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
}

func main() {
	router := gin.Default()
	router.GET("/tickets", getTickets)
	router.GET("/tickets/:id", getTicketById)
	router.DELETE("/tickets/:id", deleteTicketById)
	router.POST("/tickets", postTickets)

	router.GET("/events", getEvents)
	router.GET("/events/:id", getEventsById)
	router.POST("/events", postEvents)
	router.Run("localhost:8080")
}
