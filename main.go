package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	reader := bufio.NewReader(os.Stdin)
	authCode, _ := reader.ReadString('\n')
	authCode = strings.TrimSpace(authCode)

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("--- Time to Leave Calculator ---")

	fmt.Print("Enter your arrival time (HH:MM AM/PM, e.g., 9:00 AM): ")
	arrivalTimeStr, _ := reader.ReadString('\n')
	arrivalTimeStr = strings.TrimSpace(arrivalTimeStr)

	parsedTime, err := time.Parse("3:04 PM", arrivalTimeStr)
	if err != nil {
		fmt.Println("Error parsing arrival time:", err)
		return
	}

	now := time.Now()
	arrivalTime := time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())

	fmt.Print("Enter your break duration (0.5 for 30 mins, 1 for 1 hour): ")
	breakDurationStr, _ := reader.ReadString('\n')
	breakDurationStr = strings.TrimSpace(breakDurationStr)

	var breakDuration time.Duration
	switch breakDurationStr {
	case "0.5":
		breakDuration = 30 * time.Minute
	case "1":
		breakDuration = 1 * time.Hour
	default:
		fmt.Println("Invalid break duration. Using 1 hour.")
		breakDuration = 1 * time.Hour
	}

	// Assuming an 8-hour workday
	workdayDuration := 8 * time.Hour
	totalTimeNeeded := workdayDuration + breakDuration

	departureTime := arrivalTime.Add(totalTimeNeeded)

	fmt.Printf("\nBased on your input:\n")
	fmt.Printf("  Arrival Time: %s\n", arrivalTime.Format("3:04 PM"))
	fmt.Printf("  Break Duration: %s\n", breakDuration)
	fmt.Printf("  Recommended Departure Time: %s\n", departureTime.Format("3:04 PM"))

	// --- Google Calendar Integration ---
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	// Create the event
	event := &calendar.Event{
		Summary: "Out for the day",
		Start: &calendar.EventDateTime{
			DateTime: departureTime.Format(time.RFC3339),
			TimeZone: "America/New_York", // Change to your timezone
		},
		End: &calendar.EventDateTime{
			DateTime: time.Date(departureTime.Year(), departureTime.Month(), departureTime.Day(), 19, 0, 0, 0, departureTime.Location()).Format(time.RFC3339),
			TimeZone: "America/New_York", // Change to your timezone
		},
	}

	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
}
