package review

import (
	"encoding/json"
	"fmt"
)

// MergedReview represents a review with project name merged from notifications
type MergedReview struct {
	Time        string  // Review time in format "2006.01.02, 15:04"
	ProjectName string  // Name of the project being reviewed
	GoalName    string  // Goal name (N/A if not available)
	Booking     Booking // Booking ID for the review
	SlotID      string  // Slot ID for the review (kept for backward compatibility)
	Verifier    string  // Verifier login
	Status      string  // Booking status
	IsOnline    bool    // Whether the review is online
	VcLinkUrl   string  // Video call link URL
}

// Booking represents a calendar booking from the raw API response
// Only includes necessary fields for the review package
type Booking struct {
	ID            string        `json:"id"`            // Booking ID (e.g., "1690503")
	EventSlotID   string        `json:"eventSlotId"`   // Associated event slot ID
	Task          *Task         `json:"task"`          // Task information (goal name)
	EventSlot     *EventSlot    `json:"eventSlot"`     // Event slot information (time)
	VerifierUser  *VerifierUser `json:"verifierUser"`  // Verifier information
	BookingStatus string        `json:"bookingStatus"` // Current booking status
	IsOnline      bool          `json:"isOnline"`      // Whether review is online
	VcLinkUrl     string        `json:"vcLinkUrl"`     // Video call link URL
}

// Task represents task information within a booking
type Task struct {
	GoalName string `json:"goalName"` // Name of the goal/project
}

// EventSlot represents event slot information within a booking
type EventSlot struct {
	ID    string `json:"id"`    // Slot ID
	Start string `json:"start"` // Start time in RFC3339 format
	End   string `json:"end"`   // End time in RFC3339 format
}

// VerifierUser represents the verifier user information
type VerifierUser struct {
	ID    string `json:"id"`    // User ID
	Login string `json:"login"` // User login
}

// ConvertInterface converts an interface{} value to a target struct using JSON marshaling/unmarshaling
// This is useful when dealing with interface{} fields from generated API code
// Returns error if conversion fails, allowing caller to continue with other items
func ConvertInterface[T any](i interface{}, target *T) error {
	jsonData, err := json.Marshal(i)
	if err != nil {
		return fmt.Errorf("failed to marshal interface{}: %w", err)
	}

	if err := json.Unmarshal(jsonData, target); err != nil {
		return fmt.Errorf("failed to unmarshal to target type: %w", err)
	}

	return nil
}
