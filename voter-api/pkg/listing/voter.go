package listing

import (
	"time"
)

type Voter struct {
	VoterId      uint           `json:"id"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	VoterHistory []VoterHistory `json:"history"`
	Created      time.Time      `json:"created"`
	Modified     time.Time      `json:"modified"`
}
