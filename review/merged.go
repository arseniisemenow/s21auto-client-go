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
// Times are converted to UTC+3 to match the notification format.
func GetMergedReviewsWithProjects(ctx context.Context, client *s21client.Client) ([]MergedReview, error) {
	// Get events (contains both available and booked slots)
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

	// Extract booked slots from events (where user is the verifier)
	var bookedSlots []struct {
		ID    string
		Start time.Time
		End   time.Time
	}
	for _, event := range events {
		for _, slot := range event.EventSlots {
			if slot.Type != "FREE_TIME" {
				bookedSlots = append(bookedSlots, struct {
					ID    string
					Start time.Time
					End   time.Time
				}{
					ID:    slot.ID,
					Start: slot.Start,
					End:   slot.End,
				})
			}
		}
	}

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

	// Merge data
	var mergedReviews []MergedReview
	for _, slot := range bookedSlots {
		var review MergedReview
		// Convert to UTC+3 to match notification format
		localTime := slot.Start.In(time.FixedZone("UTC+3", 3*3600))
		review.Time = localTime.Format("2006.01.02, 15:04")
		review.SlotID = slot.ID

		// Try to find project name from notifications (notifications use UTC+3)
		timeKey := localTime.Format("2006.01.02, 15:04")
		if projectName, exists := timeToProject[timeKey]; exists {
			review.ProjectName = projectName
		}

		// For booked slots where user is verifier, we don't have goal name/verifier info in the event
		// This info might be in the notifications or requires another API call
		review.GoalName = "N/A"
		review.Verifier = "You (verifier)"
		review.Status = "BOOKED"

		if review.Time != "" {
			mergedReviews = append(mergedReviews, review)
		}
	}

	// Sort by time
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
