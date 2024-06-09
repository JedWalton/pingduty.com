package uptime

import (
	"net/http"
	"time"
)

func Run() {
	pingduty()
}

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
