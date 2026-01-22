# Review Package

This package provides functionality to merge review booking information with project names from notifications.

## Usage

```go
import (
    "context"
    "github.com/arseniisemenow/s21auto-client-go"
    "github.com/arseniisemenow/s21auto-client-go/review"
)

func main() {
    ctx := context.Background()
    client := s21client.New(
        s21client.DefaultAuth(
            os.Getenv("S21_USERNAME"),
            os.Getenv("S21_PASSWORD"),
        ),
    )

    reviews, err := review.GetMergedReviewsWithProjects(ctx, client)
    if err != nil {
        log.Fatal(err)
    }

    for _, r := range reviews {
        fmt.Printf("Time: %s\n", r.Time)
        fmt.Printf("Project: %s\n", r.ProjectName)
        fmt.Printf("Goal: %s\n", r.GoalName)
        fmt.Printf("Booking ID: %s\n", r.BookingID)
        fmt.Printf("Slot ID: %s\n", r.SlotID)
        fmt.Printf("Verifier: %s\n", r.Verifier)
        fmt.Printf("Status: %s\n", r.Status)
        fmt.Printf("Online: %v\n", r.IsOnline)
        if r.VcLinkUrl != "" {
            fmt.Printf("VC Link: %s\n", r.VcLinkUrl)
        }
        fmt.Println()
    }
}
```

## Functions

### GetMergedReviewsWithProjects

Merges review bookings with project names from notifications.

**Parameters:**
- `ctx context.Context`: Context for the request
- `client *s21client.Client`: Authenticated s21 client

**Returns:**
- `[]MergedReview`: List of merged reviews sorted by time
- `error`: Error if API calls fail

**Details:**
- Fetches calendar events as raw JSON for the next 7 days
- Fetches recent notifications (last 200)
- Extracts bookings from raw JSON response (instead of event slots)
- Matches bookings with project names from notifications
- Extracts rich booking information (ID, status, online flag, video link, verifier, goal)
- Converts times to UTC+3 to match notification format
- Sorts results chronologically

## Types

### MergedReview

Represents a review with project name merged from notifications.

**Fields:**
- `Time string`: Review time in format "2006.01.02, 15:04" (UTC+3)
- `ProjectName string`: Name of the project being reviewed
- `GoalName string`: Goal name (N/A if not available)
- `BookingID string`: Booking ID for the review
- `SlotID string`: Slot ID for the review (kept for backward compatibility)
- `Verifier string`: Verifier login
- `Status string`: Booking status
- `IsOnline bool`: Whether the review is online
- `VcLinkUrl string`: Video call link URL

## Implementation Notes

- Time conversion: Notifications use UTC+3, API returns UTC
- Raw JSON: Uses raw JSON response to access bookings array (which has `interface{}` type in generated code)
- Booking IDs: Uses booking ID (e.g., "1690503") instead of event/slot ID
- Matching: Reviews are matched by exact time between bookings and notifications
- Project names: Extracted from notification messages using regex
- Rich data: Bookings contain goal name, verifier info, online flag, and video call link
