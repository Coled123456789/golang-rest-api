package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestEvent(t *testing.T) {
	r := SetUpRouter()
	r.POST("/events", postEvents)
	r.GET("/events/:id", getEventsById)
	event := event_t{
		ID:            0,
		Name:          "New Event",
		StartDateTime: time.Now().UTC().String(),
		EndDateTime:   time.Now().Add(time.Hour * 3).UTC().String(),
		ValidTypes:    []string{"A", "B", "C"},
	}
	jsonValue, _ := json.Marshal(event)
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusCreated != w.Code {
		t.Fatalf("Failed with response code %d\n", w.Code)
	}

	var resp *http.Response = w.Result()
	body, _ := io.ReadAll(resp.Body)
	var data event_t
	json.Unmarshal(body, &data)
	var id int = data.ID

	req, _ = http.NewRequest("GET", "/events/"+strconv.Itoa(id), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	body, _ = io.ReadAll(w.Result().Body)
	var newEvent event_t
	json.Unmarshal(body, &newEvent)
	if !CompareEvents(event, newEvent) {
		t.Log("Retrieved Event: ")
		t.Log(newEvent)
		t.Log("Added Event: ")
		t.Log(event)
		t.Fatalf("Retrieved event does not match added event")
	}

}
