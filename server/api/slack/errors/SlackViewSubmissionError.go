package errors

import "fmt"

type SlackViewSubmissionError struct {
	Key   string
	Value string
}

func (s *SlackViewSubmissionError) Error() string {
	return fmt.Sprintf("ViewSubmissionError Key: %s; Value: %s", s.Key, s.Value)
}
