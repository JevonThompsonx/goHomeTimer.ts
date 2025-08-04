package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// CalculateDepartureTime calculates the departure time based on arrival time and break duration.
func (a *App) CalculateDepartureTime(arrivalTimeStr string, breakDurationStr string) string {
	parsedTime, err := time.Parse("3:04 PM", arrivalTimeStr)
	if err != nil {
		return fmt.Sprintf("Error parsing arrival time: %v", err)
	}

	now := time.Now()
	arrivalTime := time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())

	var breakDuration time.Duration
	switch breakDurationStr {
	case "0.5":
		breakDuration = 30 * time.Minute
	case "1":
		breakDuration = 1 * time.Hour
	default:
		return "Invalid break duration"
	}

	workdayDuration := 8 * time.Hour
	totalTimeNeeded := workdayDuration + breakDuration
	departureTime := arrivalTime.Add(totalTimeNeeded)

	return departureTime.Format("3:04 PM")
}

// CreateCalendarEvent creates a Google Calendar event.
func (a *App) CreateCalendarEvent(departureTimeStr string) string {
	departureTime, err := time.Parse("3:04 PM", departureTimeStr)
	if err != nil {
		return fmt.Sprintf("Error parsing departure time: %v", err)
	}

	now := time.Now()
	departureDateTime := time.Date(now.Year(), now.Month(), now.Day(), departureTime.Hour(), departureTime.Minute(), 0, 0, now.Location())

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return fmt.Sprintf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope)
	if err != nil {
		return fmt.Sprintf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		return fmt.Sprintf("Unable to retrieve Calendar client: %v", err)
	}

	event := &calendar.Event{
		Summary: "Out for the day",
		Start: &calendar.EventDateTime{
			DateTime: departureDateTime.Format(time.RFC3339),
			TimeZone: "America/New_York",
		},
		End: &calendar.EventDateTime{
			DateTime: time.Date(departureDateTime.Year(), departureDateTime.Month(), departureDateTime.Day(), 19, 0, 0, 0, departureDateTime.Location()).Format(time.RFC3339),
			TimeZone: "America/New_York",
		},
	}

	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		return fmt.Sprintf("Unable to create event. %v", err)
	}
	return fmt.Sprintf("Event created: %s", event.HtmlLink)
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}