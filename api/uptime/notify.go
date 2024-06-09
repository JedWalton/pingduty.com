package uptime

import (
	"fmt"
	"log"
	"time"
)

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
