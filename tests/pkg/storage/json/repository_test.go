package tests

import (
	"log"
	"os"
	"testing"
	"time"

	"drexel.edu/VoterApi/pkg/adding"
	"drexel.edu/VoterApi/pkg/changing"
	"drexel.edu/VoterApi/pkg/storage/json"
	"github.com/stretchr/testify/assert"
)

var dbFileVoter string

func init() {
	log.Println("initializing json repository tests...")
	dbFileVoter = "../../../../data/"
}

func TestCanInstantiateVoter(t *testing.T) {

	_, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

}

func TestCanInstantiateVoterHistory(t *testing.T) {

	_, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

}

func TestAddVoter(t *testing.T) {

	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	_, err = db.AddVoter(voter)
	assert.NoError(t, err)
}

func TestErrorOnAddTwice(t *testing.T) {

	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	_, err = db.AddVoter(voter)
	assert.NoError(t, err)

	_, err = db.AddVoter(voter)
	assert.Error(t, err)
}

func TestGetVoter(t *testing.T) {

	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	id, err := db.AddVoter(voter)
	assert.NoError(t, err)

	retrievedVoter, err := db.GetVoter(id)
	assert.NoError(t, err)

	assert.Equal(t, retrievedVoter.VoterId, uint(id))
	assert.Equal(t, retrievedVoter.Name, voter.Name)
	assert.Equal(t, retrievedVoter.Email, voter.Email)
	assert.NotNil(t, retrievedVoter.Created)
	assert.NotNil(t, retrievedVoter.Modified)

}

func TestGetVoterNotFound(t *testing.T) {

	cleanDB()

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	_, err = db.GetVoter(1)
	assert.Error(t, err)

}

func TestGetAllVoters(t *testing.T) {

	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	voter2 := adding.Voter{
		Name:  "john smit",
		Email: "js3762@drexel.edu",
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	_, err = db.AddVoter(voter)
	assert.NoError(t, err)

	_, err = db.AddVoter(voter2)
	assert.NoError(t, err)

	retrievedVoters, err := db.GetVoters()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(retrievedVoters))

}

func TestUpdateVoter(t *testing.T) {

	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	newId, err := db.AddVoter(voter)
	assert.NoError(t, err)

	updatedVoter := changing.Voter{
		VoterId: uint(newId),
		Name:    "Deepali M Mathur",
		Email:   voter.Email,
	}

	err = db.UpdateVoter(updatedVoter)
	assert.NoError(t, err)

	retrievedVoter, err := db.GetVoter(int(updatedVoter.VoterId))
	assert.NoError(t, err)

	assert.Equal(t, retrievedVoter.VoterId, updatedVoter.VoterId)
	assert.Equal(t, retrievedVoter.Name, updatedVoter.Name)
	assert.Equal(t, retrievedVoter.Email, updatedVoter.Email)
}

func TestErrorIfUpdateToConflictingEmail(t *testing.T) {

	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	voterWithConflictingEmail := adding.Voter{
		Name:  "John Donald",
		Email: "jd2862@drexel.edu",
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	newId, err := db.AddVoter(voter)
	assert.NoError(t, err)

	_, err = db.AddVoter(voterWithConflictingEmail)
	assert.NoError(t, err)

	updatedVoter := changing.Voter{
		VoterId: uint(newId),
		Name:    "Deepali M Malhotra",
		Email:   voterWithConflictingEmail.Email,
	}

	err = db.UpdateVoter(updatedVoter)
	assert.Error(t, err)
}

func TestAddVoterHistory(t *testing.T) {
	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	voterHistory := adding.VoterHistory{
		VoterId:  1,
		PollId:   uint(1),
		VoteId:   uint(1),
		VoteDate: time.Now(),
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	_, err = db.AddVoter(voter)
	assert.NoError(t, err)

	err = db.AddVoterHistory(voterHistory)
	assert.NoError(t, err)
}

func TestErrorIfCreateSameHistoryTwice(t *testing.T) {
	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	voterHistory := adding.VoterHistory{
		VoterId:  1,
		PollId:   uint(1),
		VoteId:   uint(1),
		VoteDate: time.Now(),
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	_, err = db.AddVoter(voter)
	assert.NoError(t, err)

	err = db.AddVoterHistory(voterHistory)
	assert.NoError(t, err)

	err = db.AddVoterHistory(voterHistory)
	assert.Error(t, err)
}

func TestGetVoterHistory(t *testing.T) {
	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	voterHistory := adding.VoterHistory{
		VoterId:  1,
		PollId:   uint(1),
		VoteId:   uint(1),
		VoteDate: time.Now(),
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	voterId, err := db.AddVoter(voter)
	assert.NoError(t, err)

	err = db.AddVoterHistory(voterHistory)
	assert.NoError(t, err)

	poll, err := db.GetVoterHistory(uint(voterId), voterHistory.VoterId)
	assert.NoError(t, err)
	assert.Equal(t, voterHistory.VoteId, poll.VoteId)
	assert.Equal(t, voterHistory.PollId, poll.PollId)
	assert.True(t, voterHistory.VoteDate.Equal(poll.VoteDate))
}

func TestAllVoterHistoryPresent(t *testing.T) {
	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	voterHistory := adding.VoterHistory{
		VoterId:  1,
		PollId:   uint(1),
		VoteId:   uint(1),
		VoteDate: time.Now(),
	}

	voterHistory2 := adding.VoterHistory{
		VoterId:  1,
		PollId:   uint(2),
		VoteId:   uint(2),
		VoteDate: time.Now(),
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	voterId, err := db.AddVoter(voter)
	assert.NoError(t, err)

	err = db.AddVoterHistory(voterHistory)
	assert.NoError(t, err)

	err = db.AddVoterHistory(voterHistory2)
	assert.NoError(t, err)

	result, err := db.GetVoter(voterId)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.VoterHistory))
}

func TestUpdateHistory(t *testing.T) {
	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	voterHistory := adding.VoterHistory{
		VoterId:  1,
		PollId:   uint(1),
		VoteId:   uint(1),
		VoteDate: time.Now(),
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	voterId, err := db.AddVoter(voter)
	assert.NoError(t, err)

	err = db.AddVoterHistory(voterHistory)
	assert.NoError(t, err)

	voterHistoryEdit := changing.VoterHistory{
		VoterId:  1,
		PollId:   uint(34),
		VoteId:   uint(1),
		VoteDate: time.Now(),
	}

	err = db.UpdateVoterHistory(voterHistoryEdit)
	assert.NoError(t, err)

	actualStoredHistory, err := db.GetVoterHistory(uint(voterId), voterHistory.VoteId)
	assert.NoError(t, err)
	assert.Equal(t, voterHistoryEdit.PollId, actualStoredHistory.PollId)

}

func TestDeleteHistory(t *testing.T) {
	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	voterHistory := adding.VoterHistory{
		VoterId:  1,
		PollId:   uint(1),
		VoteId:   uint(1),
		VoteDate: time.Now(),
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	voterId, err := db.AddVoter(voter)
	assert.NoError(t, err)

	err = db.AddVoterHistory(voterHistory)
	assert.NoError(t, err)

	err = db.DeleteVoterHistory(uint(voterId), 1)
	assert.NoError(t, err)

	result, err := db.GetVoter(voterId)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(result.VoterHistory))

}

func TestDeleteVoter(t *testing.T) {
	cleanDB()

	voter := adding.Voter{
		Name:  "deepali mathur",
		Email: "dm3729@drexel.edu",
	}

	voterHistory := adding.VoterHistory{
		VoterId:  1,
		PollId:   uint(1),
		VoteId:   uint(1),
		VoteDate: time.Now(),
	}

	db, err := json.NewVoterDB(dbFileVoter)
	assert.NoError(t, err)

	voterId, err := db.AddVoter(voter)
	assert.NoError(t, err)

	err = db.AddVoterHistory(voterHistory)
	assert.NoError(t, err)

	err = db.DeleteVoter(uint(voterId))
	assert.NoError(t, err)

	result, err := db.GetVoter(voterId)
	assert.Error(t, err)
	assert.Equal(t, 0, len(result.VoterHistory))

	_, err = db.GetVoterHistory(uint(voterId), voterHistory.VoteId)
	assert.Error(t, err)

}

func cleanDB() {
	os.RemoveAll(dbFileVoter)
}
