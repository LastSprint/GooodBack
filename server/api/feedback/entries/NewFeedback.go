package entries

type NewFeedback struct {
	// Message is feedback text
	Message string `json:"message"`
	// Target is the point of the feedback
	Target string `json:"target"`
	// Type represent feedback score (or feedback reaction)
	Type int `json:"type"`

	Id string `json:"id,omitempty"`
}
