package uptime

import (
	"time"
)

func pingduty() {
	// Read the Discord webhook URL from an environment variable
	keysForNotifications := KeysForNotifications{
		discordWebhookURL: "placeholder_discord_webhook_url",
	}

	// Configuration
	endpoints := map[string]*EndpointState{
		"pingduty": {
			URL:              "https://www.pingduty.com",
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
