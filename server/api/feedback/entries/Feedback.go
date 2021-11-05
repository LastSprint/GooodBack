package entries

import "time"

type Feedback struct {
	NewFeedback
	CreationDate time.Time `json:"creation_date"`
}

type FeedbackSortable []Feedback

func (f FeedbackSortable) Len() int {
	return len(f)
}

func (f FeedbackSortable) Less(i, j int) bool {
	return f[i].CreationDate.After(f[j].CreationDate)
}

func (f FeedbackSortable) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
