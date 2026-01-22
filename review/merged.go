package review

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"time"

	s21client "github.com/arseniisemenow/s21auto-client-go"
	"github.com/arseniisemenow/s21auto-client-go/requests"
)

// GetMergedReviewsWithProjects merges review bookings with project names from notifications
// It fetches both calendar events and user notifications, then matches them by time.
// Uses ConvertInterface to convert bookings from interface{} to typed Booking struct.
// Times are converted to UTC+3 to match the notification format.
func GetMergedReviewsWithProjects(ctx context.Context, client *s21client.Client) ([]MergedReview, error) {
	// Get events using existing API (contains bookings as interface{})
	from := time.Now()
	to := time.Now().AddDate(0, 0, 7) // Next week

	eventsResp, err := client.R().SetContext(ctx).CalendarGetEvents(requests.CalendarGetEvents_Variables{
		From: from,
		To:   to,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar events: %w", err)
	}

	// Get notifications
	notificationsResp, err := client.R().SetContext(ctx).GetUserNotifications(requests.GetUserNotifications_Variables{
		Paging: requests.GetUserNotifications_Variables_Paging{
			Offset: 0,
			Limit:  200, // Get more notifications to find matches
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	events := eventsResp.CalendarEventS21.GetMyCalendarEvents
	notifications := notificationsResp.S21Notification.GetS21Notifications

	// Build a map of time -> project name from notifications
	// Parse times from notifications (they are in UTC+3 format)
	timeToProject := make(map[string]string)

	for _, n := range notifications.Notifications {
		// Look for messages like: "Someone registered for a review of the project <b>DO6_CICD</b> by you on <b>2026.01.22, 19:45</b>"
		message := n.Message
		if n.GroupName == "PROJECTS" {
			projectName, reviewTime := parseProjectFromNotification(message)
			if projectName != "" && reviewTime != "" {
				timeToProject[reviewTime] = projectName
			}
		}
	}

	// Extract bookings from events using ConvertInterface
	var mergedReviews []MergedReview

	for _, event := range events {
		for _, bookingI := range event.Bookings {
			// Convert interface{} to typed Booking struct
			var booking Booking
			if err := ConvertInterface(bookingI, &booking); err != nil {
				// Continue on conversion error - skip this booking
				continue
			}

			// Parse time from event slot
			startTime, err := time.Parse(time.RFC3339, booking.EventSlot.Start)
			if err != nil {
				// Skip if time parsing fails
				continue
			}

			// Convert to UTC+3 to match notification format
			// Todo: This need to be configurable or auto
			localTime := startTime.In(time.FixedZone("UTC+3", 3*3600))
			timeKey := localTime.Format("2006.01.02, 15:04")

			// Extract goal name
			goalName := "N/A"
			if booking.Task != nil && booking.Task.GoalName != "" {
				goalName = booking.Task.GoalName
			}

			// Extract verifier login
			verifier := "N/A"
			if booking.VerifierUser != nil && booking.VerifierUser.Login != "" {
				verifier = booking.VerifierUser.Login
			}

			// Create merged review
			review := MergedReview{
				Time:     timeKey,
				Booking:  booking,
				GoalName: goalName,
				Verifier: verifier,
			}

			// Try to find project name from notifications (notifications use UTC+3)
			if projectName, exists := timeToProject[timeKey]; exists {
				review.ProjectName = projectName
			}

			mergedReviews = append(mergedReviews, review)
		}
	}

	// Sort by time, maybe redundant
	sort.Slice(mergedReviews, func(i, j int) bool {
		timeI, _ := time.Parse("2006.01.02, 15:04", mergedReviews[i].Time)
		timeJ, _ := time.Parse("2006.01.02, 15:04", mergedReviews[j].Time)
		return timeI.Before(timeJ)
	})

	return mergedReviews, nil
}

// parseProjectFromNotification extracts project name and review time from notification message
// Example message: "Someone registered for a review of the project <b>DO6_CICD</b> by you on <b>2026.01.22, 19:45</b>"
// Returns: projectName, reviewTime
func parseProjectFromNotification(message string) (projectName, reviewTime string) {
	// Extract project name from <b> tags
	projectRegex := regexp.MustCompile(`<b>([^<]+)</b>\s+by you on`)
	matches := projectRegex.FindStringSubmatch(message)
	if len(matches) >= 2 {
		projectName = matches[1]
	}

	// Extract time from the second <b> tag
	timeRegex := regexp.MustCompile(`on\s+<b>([^<]+)</b>`)
	timeMatches := timeRegex.FindStringSubmatch(message)
	if len(timeMatches) >= 2 {
		reviewTime = timeMatches[1]
	}

	return projectName, reviewTime
}
