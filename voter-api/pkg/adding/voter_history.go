package adding

import "time"

type VoterHistory struct {
	VoterId  uint      `json:"voter_id"`
	PollId   uint      `json:"poll_id"`
	VoteId   uint      `json:"vote_id"`
	VoteDate time.Time `json:"vote_date"`
}
