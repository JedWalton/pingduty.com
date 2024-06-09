package uptime

import (
	"time"
)

type BespokeConfig struct {
	keysFornotifications KeysForNotifications
	endpoints            map[string]*EndpointState
	pingConfig           PingConfig
}

func pingduty() {
	bespokeConfig := BespokeConfig{
		keysFornotifications: KeysForNotifications{
			discordWebhookURL: "placeholder_discord_webhook_url",
		},
		endpoints: map[string]*EndpointState{
			"pingduty": {
				URL:              "https://www.pingduty.com",
				IsDown:           false,
				LastNotification: time.Now(),
			},
		},
		pingConfig: PingConfig{
			pingFrequency:        1 * time.Second,
			notificationInterval: 300 * time.Second,
		},
	}

	// Start health checks
	HealthCheckScheduler(bespokeConfig)
}
