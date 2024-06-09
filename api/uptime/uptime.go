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

func HealthCheckScheduler(bespokeConfig BespokeConfig) {
	ticker := time.NewTicker(bespokeConfig.pingConfig.pingFrequency)
	defer ticker.Stop()

	for range ticker.C {
		for name, state := range bespokeConfig.endpoints {
			go checkEndpointHealth(name, state,
				bespokeConfig.pingConfig.notificationInterval,
				bespokeConfig.keysFornotifications)
		}
	}
}

func checkEndpointHealth(name string, state *EndpointState,
	notifyInterval time.Duration, keysForNotifications KeysForNotifications) {
	resp, err := http.Get(state.URL)
	currentTime := time.Now()

	if err != nil || resp.StatusCode != http.StatusOK ||
		resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if !state.IsDown {
			state.IsDown = true
			state.DownSince = currentTime
			state.LastNotification = currentTime
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
