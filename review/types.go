package review

// MergedReview represents a review with project name merged from notifications
type MergedReview struct {
	Time        string // Review time in format "2006.01.02, 15:04"
	ProjectName string // Name of the project being reviewed
	GoalName    string // Goal name (N/A if not available)
	SlotID      string // Slot ID for the review
	Verifier    string // Verifier login
	Status      string // Booking status
}
