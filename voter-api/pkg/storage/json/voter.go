package json

import (
	"time"
)

type Voter struct {
	VoterId      uint      `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	VoterHistory []uint    `json:"history"`
	Created      time.Time `json:"created"`
	Modified     time.Time `json:"modified"`
}
