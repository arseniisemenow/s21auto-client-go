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
        fmt.Printf("Slot ID: %s\n", r.SlotID)
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
- Fetches calendar events for the next 7 days
- Fetches recent notifications (last 200)
- Matches booked review slots with project names from notifications
- Converts times to UTC+3 to match notification format
- Sorts results chronologically

## Types

### MergedReview

Represents a review with project name merged from notifications.

**Fields:**
- `Time string`: Review time in format "2006.01.02, 15:04" (UTC+3)
- `ProjectName string`: Name of the project being reviewed
- `GoalName string`: Goal name (N/A if not available)
- `SlotID string`: Slot ID for the review
- `Verifier string`: Verifier login
- `Status string`: Booking status

## Implementation Notes

- Time conversion: Notifications use UTC+3, API returns UTC
- Matching: Reviews are matched by time (nearest 15-minute window)
- Booked slots: Only includes slots where user is the verifier
- Project names: Extracted from notification messages using regex
