package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// EndpointState holds the state of an endpoint
type EndpointState struct {
	URL              string
	IsDown           bool
	LastNotification time.Time
	DownSince        time.Time
}

func pingduty() {
	// Read the Discord webhook URL from an environment variable
	discordWebhookURL := "placeholder_discord_webhook_url"

	// Configuration
	endpoints := map[string]*EndpointState{
		"pingduty": {
			URL:              "http://pingduty.com",
			IsDown:           false,
			LastNotification: time.Now(),
		},
	}
	pingFrequency := 1 * time.Second        // Frequency of health checks
	notificationInterval := 5 * time.Minute // Interval for re-notification

	// Start health checks
	HealthCheckScheduler(endpoints, pingFrequency, notificationInterval, discordWebhookURL)
}

func HealthCheckScheduler(endpoints map[string]*EndpointState, frequency, notifyInterval time.Duration, discordWebhookURL string) {
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for range ticker.C {
		for name, state := range endpoints {
			go checkEndpointHealth(name, state, notifyInterval, discordWebhookURL)
		}
	}
}

func checkEndpointHealth(name string, state *EndpointState, notifyInterval time.Duration, discordWebhookURL string) {
	resp, err := http.Get(state.URL)
	currentTime := time.Now()

	if err != nil || resp.StatusCode != http.StatusOK {
		if !state.IsDown {
			state.IsDown = true
			state.DownSince = currentTime // Record the time when the service first went down
			state.LastNotification = currentTime
			notifyChannel(name, "DOWN", 0, discordWebhookURL) // Initially, down time is 0
		} else if currentTime.Sub(state.LastNotification) >= notifyInterval {
			downDuration := currentTime.Sub(state.DownSince)
			state.LastNotification = currentTime
			notifyChannel(name, "STILL DOWN", downDuration, discordWebhookURL)
		}
	} else {
		if state.IsDown {
			downDuration := currentTime.Sub(state.DownSince)
			state.IsDown = false
			state.LastNotification = currentTime
			notifyChannel(name, "UP", downDuration, discordWebhookURL)
		}
	}
	if resp != nil {
		resp.Body.Close()
	}
}

func notifyChannel(serviceName, status string, duration time.Duration, discordWebhookURL string) {
	message := fmt.Sprintf("Service %s is %s", serviceName, status)
	if duration > 0 {
		message += fmt.Sprintf(" for %s", duration.Truncate(time.Second)) // Truncate to remove microseconds
	}

	// Implement notification logic here to use discordWebhookURL
	log.Println("Notification sent:", message) // Placeholder for actual notification implementation
}

func main() {
	pingduty()
}
