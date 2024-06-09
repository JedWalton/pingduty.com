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

type KeysForNotifications struct {
	discordWebhookURL string
	anotherWebhookURL string
}

type PingConfig struct {
	pingFrequency        time.Duration
	notificationInterval time.Duration
}

func pingduty() {
	// Read the Discord webhook URL from an environment variable
	keysForNotifications := KeysForNotifications{
		discordWebhookURL: "placeholder_discord_webhook_url",
	}

	// Configuration
	endpoints := map[string]*EndpointState{
		"pingduty": {
			URL:              "http://pingduty.com",
			IsDown:           false,
			LastNotification: time.Now(),
		},
	}

	pingConfig := PingConfig{
		pingFrequency:        1 * time.Second,
		notificationInterval: 2 * time.Second,
	}

	// Start health checks
	HealthCheckScheduler(endpoints, pingConfig, keysForNotifications)
}

func HealthCheckScheduler(endpoints map[string]*EndpointState,
	pingConfig PingConfig, keysForNotifications KeysForNotifications) {
	ticker := time.NewTicker(pingConfig.pingFrequency)
	defer ticker.Stop()

	for range ticker.C {
		for name, state := range endpoints {
			go checkEndpointHealth(name, state, pingConfig.notificationInterval,
				keysForNotifications)
		}
	}
}

func checkEndpointHealth(name string, state *EndpointState,
	notifyInterval time.Duration, keysForNotifications KeysForNotifications) {
	resp, err := http.Get(state.URL)
	currentTime := time.Now()

	if err != nil || resp.StatusCode != http.StatusOK ||
		resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		if !state.IsDown {
			state.IsDown = true
			// Record the time when the service first went down
			state.DownSince = currentTime
			state.LastNotification = currentTime
			// Initially, down time is 0
			notifyChannel(name, "DOWN", 0, keysForNotifications)
		} else if currentTime.Sub(state.LastNotification) >= notifyInterval {
			downDuration := currentTime.Sub(state.DownSince)
			state.LastNotification = currentTime
			notifyChannel(name, "STILL DOWN", downDuration, keysForNotifications)
		}
	} else {
		if state.IsDown {
			downDuration := currentTime.Sub(state.DownSince)
			state.IsDown = false
			state.LastNotification = currentTime
			notifyChannel(name, "UP", downDuration, keysForNotifications)
		}
	}
	if resp != nil {
		resp.Body.Close()
	}
}

func notifyChannel(serviceName, status string, duration time.Duration,
	keysForNotifications KeysForNotifications) {
	message := fmt.Sprintf("Service %s is %s", serviceName, status)
	if duration > 0 {
		// Truncate to remove microseconds
		message += fmt.Sprintf(" for %s", duration.Truncate(time.Second))
	}

	// Implement notification logic here to use discordWebhookURL
	log.Println("Notification sent:", message) // Placeholder for actual notification implementation

	if keysForNotifications.discordWebhookURL != "" {
		// Placeholder for actual notification implementation
		notifyDiscord(message, keysForNotifications.discordWebhookURL)
	}
	if keysForNotifications.anotherWebhookURL != "" {
		// Placeholder for actual notification implementation
		notifyAnotherWebhook(message, keysForNotifications.anotherWebhookURL)
	}

}

func notifyAnotherWebhook(message, s string) {
	panic("unimplemented")
}

func notifyDiscord(message string, discordWebhookURL string) {
	// Placeholder for actual notification implementation
}

func main() {
	pingduty()
}
