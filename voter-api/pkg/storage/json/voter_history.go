package json

import (
	"time"
)

type VoterHistory struct {
	PollId   uint      `json:"id"`
	VoteId   uint      `json:"vote_id"`
	VoteDate time.Time `json:"vote_date"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}
