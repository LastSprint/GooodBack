package entries

type NewFeedback struct {
	// Message is feedback text
	Message string `json:"message"`
	// Target is the point of the feedback
	Target string `json:"target"`
}
