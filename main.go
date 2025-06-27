package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("--- Time to Leave Calculator ---")

	fmt.Print("Enter your arrival time (HH:MM AM/PM, e.g., 9:00 AM): ")
	arrivalTimeStr, _ := reader.ReadString('\n')
	arrivalTimeStr = strings.TrimSpace(arrivalTimeStr)

	arrivalTime, err := time.Parse("3:04 PM", arrivalTimeStr) // Adjust format as needed
	if err != nil {
		fmt.Println("Error parsing arrival time:", err)
		return
	}

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

	// --- Google Calendar Integration (Conceptual - Requires API setup) ---
	// In a real application, you'd call a function here to interact with Google Calendar
	fmt.Println("\n(Soon to add: 'Out for the day' event to your Google Calendar!)")
}
