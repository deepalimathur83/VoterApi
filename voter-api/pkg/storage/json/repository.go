package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"drexel.edu/voter-api/pkg/adding"
	"drexel.edu/voter-api/pkg/changing"
	"drexel.edu/voter-api/pkg/listing"
	"drexel.edu/voter-api/pkg/storage"
)

const (
	CollectionVoter        = "voters"
	CollectionVoterHistory = "voter_history"
)

type DbMapVoter map[int]Voter
type DbMapVoterHistory map[int]VoterHistory

type DBVoter struct {
	dbVoterMap        DbMapVoter
	dbVoterHistoryMap DbMapVoterHistory
	dbVoterFileName   string
	dbHistoryFileName string
}

func NewVoterDB(path string) (*DBVoter, error) {

	fileNameVoter := path + "Voter.json"
	fileNameHistory := path + "History.json"

	exists, err := storage.CheckIfFileOrFolderExistsAndNotEmpty(fileNameVoter)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := initDB(fileNameVoter)
		if err != nil {
			return nil, err
		}
	}

	exists, err = storage.CheckIfFileOrFolderExistsAndNotEmpty(fileNameHistory)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := initDB(fileNameHistory)
		if err != nil {
			return nil, err
		}
	}

	DBVoter := &DBVoter{
		dbVoterMap:        make(map[int]Voter),
		dbVoterFileName:   fileNameVoter,
		dbVoterHistoryMap: make(map[int]VoterHistory),
		dbHistoryFileName: fileNameHistory,
	}

	return DBVoter, nil
}

func (v *DBVoter) AddVoter(voter adding.Voter) (int, error) {
	if err := v.loadDB(); err != nil {
		return 0, errors.New("failed to load the database")
	}

	exists := v.findVoterByEmail(voter.Email)
	if exists != nil {
		return 0, errors.New("the email provided already exists in the database")
	}

	id := len(v.dbVoterMap) + 1

	currentTime := time.Now()

	newVoter := Voter{
		VoterId:      uint(id),
		Name:         voter.Name,
		Email:        voter.Email,
		VoterHistory: nil,
		Created:      currentTime,
		Modified:     currentTime,
	}

	v.dbVoterMap[int(newVoter.VoterId)] = newVoter

	if err := v.saveDB(); err != nil {
		return 0, errors.New("Failed to save to the database.")
	}
	fmt.Println("The voter was successfully registered.")
	v.printVoter(newVoter)

	return id, nil
}

func (v *DBVoter) GetVoter(id int) (listing.Voter, error) {

	var voter listing.Voter

	if err := v.loadDB(); err != nil {
		return voter, errors.New("failed to load the database")
	}

	item, exists := v.dbVoterMap[id]
	if exists {
		voter.VoterId = uint(item.VoterId)
		voter.Name = item.Name
		voter.Email = item.Email
		voter.Created = item.Created
		voter.Modified = item.Modified
		if len(item.VoterHistory) > 0 {
			voter.VoterHistory = v.getAllVoterHistory(item)
		}

		return voter, nil
	}

	return voter, errors.New(fmt.Sprintf("Couldn't get item because the id %d doesn't exist", id))
}

func (v *DBVoter) GetVoters() ([]listing.Voter, error) {

	var voters []listing.Voter

	if err := v.loadDB(); err != nil {
		return voters, errors.New("failed to load the database")
	}

	for _, item := range v.dbVoterMap {
		voter := listing.Voter{
			VoterId:      item.VoterId,
			Name:         item.Name,
			Email:        item.Email,
			Created:      item.Created,
			Modified:     item.Modified,
			VoterHistory: v.getAllVoterHistory(item),
		}
		voters = append(voters, voter)
	}

	return voters, nil
}

func (v *DBVoter) UpdateVoter(item changing.Voter) error {

	if err := v.loadDB(); err != nil {
		return errors.New("failed to load the database")
	}

	if previousVoter, exists := v.dbVoterMap[int(item.VoterId)]; exists {

		if item.Email != previousVoter.Email {
			exists := v.findVoterByEmail(item.Email)
			if exists != nil {
				return errors.New("the email you attempted to update is already associated with another voter.")
			}
		}

		updatedVoter := Voter{
			VoterId:      item.VoterId,
			Name:         item.Name,
			Email:        item.Email,
			VoterHistory: previousVoter.VoterHistory,
			Created:      previousVoter.Created,
			Modified:     time.Now(),
		}

		v.dbVoterMap[int(updatedVoter.VoterId)] = updatedVoter

		if err := v.saveDB(); err != nil {
			return errors.New("Failed to save to the database.")
		}

		return nil
	}

	return errors.New(fmt.Sprintf("Couldn't update item because the id %d doesn't exist", item.VoterId))
}

func (v *DBVoter) DeleteVoter(id uint) error {

	if err := v.loadDB(); err != nil {
		return errors.New("failed to load the database")
	}

	if voter, exists := v.dbVoterMap[int(id)]; exists {

		for _, historyId := range voter.VoterHistory {
			err := v.DeleteVoterHistory(id, historyId)
			if err != nil {
				return err
			}
		}

		delete(v.dbVoterMap, int(id))

		if err := v.saveDB(); err != nil {
			return errors.New("Failed to save to the database.")
		}

		return nil
	}

	return errors.New(fmt.Sprintf("Couldn't find voter with the vote id %d ", id))
}

func (v *DBVoter) AddVoterHistory(history adding.VoterHistory) error {
	if err := v.loadDB(); err != nil {
		return errors.New("failed to load the database")
	}

	if _, exists := v.dbVoterHistoryMap[int(history.PollId)]; exists {
		return errors.New(fmt.Sprintf("Couldn't add the poll because the id %d already exists", history.PollId))
	}

	if voter, exists := v.dbVoterMap[int(history.VoterId)]; exists {

		currentTime := time.Now()

		newPoll := VoterHistory{
			PollId:   history.PollId,
			VoteId:   history.VoteId,
			VoteDate: history.VoteDate,
			Created:  currentTime,
			Modified: currentTime,
		}

		v.dbVoterHistoryMap[int(history.PollId)] = newPoll

		voter.VoterHistory = append(voter.VoterHistory, history.PollId)

		v.dbVoterMap[int(voter.VoterId)] = voter

		if err := v.saveDB(); err != nil {
			return errors.New("Failed to save to the database.")
		}

		fmt.Println("History Was Successfully added to the voter")
		v.printHistory(newPoll)

		return nil
	}

	return errors.New(fmt.Sprintf("Couldn't add the poll because the voter id %d doesn't exist", history.VoterId))
}

func (v *DBVoter) GetVoterHistory(voterId uint, voteId uint) (listing.VoterHistory, error) {
	var history listing.VoterHistory

	if err := v.loadDB(); err != nil {
		return history, errors.New("failed to load the database")
	}

	if voter, exists := v.dbVoterMap[int(voterId)]; exists {

		for _, storedVoteIds := range voter.VoterHistory {
			if storedVoteIds == voteId {
				if item, exists := v.dbVoterHistoryMap[int(voteId)]; exists {

					history.VoteId = item.VoteId
					history.PollId = item.PollId
					history.VoteDate = item.VoteDate
					history.Created = item.Created
					history.Modified = item.Modified

					return history, nil
				}
			}
		}
	}

	return history, errors.New(fmt.Sprintf("Couldn't find history with id %d history for the vote id %d ", voteId, voterId))
}

func (v *DBVoter) UpdateVoterHistory(item changing.VoterHistory) error {
	if err := v.loadDB(); err != nil {
		return errors.New("failed to load the database")
	}

	if voter, exists := v.dbVoterMap[int(item.VoterId)]; exists {

		for _, voteId := range voter.VoterHistory {
			if voteId == item.VoteId {
				if previousHistory, exists := v.dbVoterHistoryMap[int(item.VoteId)]; exists {
					newHistory := VoterHistory{
						PollId:   item.PollId,
						VoteId:   previousHistory.VoteId,
						VoteDate: item.VoteDate,
						Created:  previousHistory.Created,
						Modified: time.Now(),
					}

					v.dbVoterHistoryMap[int(item.VoterId)] = newHistory

					if err := v.saveDB(); err != nil {
						return errors.New("Failed to save to the database.")
					}

					return nil
				}
			}
		}
		return errors.New(fmt.Sprintf("Couldn't update history because the vote id %d isn't associated with the voter", item.VoteId))
	}
	return errors.New(fmt.Sprintf("Couldn't update history because the voter %d doesn't exist", item.VoterId))
}

func (v *DBVoter) DeleteVoterHistory(voterId uint, voteId uint) error {

	if err := v.loadDB(); err != nil {
		return errors.New("failed to load the database")
	}

	if voter, exists := v.dbVoterMap[int(voteId)]; exists {

		var newArray []uint

		for _, tempVoteId := range voter.VoterHistory {
			if tempVoteId != voterId {
				newArray = append(newArray, tempVoteId)
			}
		}

		voter.VoterHistory = newArray
		v.dbVoterMap[int(voteId)] = voter
		if _, exists := v.dbVoterHistoryMap[int(voteId)]; exists {

			delete(v.dbVoterHistoryMap, int(voteId))

			if err := v.saveDB(); err != nil {
				return errors.New("Failed to save to the database")
			}

			return nil
		}
	}

	return errors.New(fmt.Sprintf("Couldn't find history with the vote id %d ", voterId))
}

func initDB(dbFileName string) error {
	f, err := storage.CreateFile(dbFileName)
	if err != nil {
		return err
	}

	// Given we are working with a json array as our DB structure
	// we should initialize the file with an empty array, which
	// in json is represented as "[]
	_, err = f.Write([]byte("[]"))
	if err != nil {
		return err
	}

	f.Close()

	return nil
}

func (v *DBVoter) saveDB() error {
	//1. Convert our map into a slice
	//2. Marshal the slice into json
	//3. Write the json to our file

	//1. Convert our map into a slice
	var voterList []Voter
	for _, item := range v.dbVoterMap {
		voterList = append(voterList, item)
	}

	var history []VoterHistory
	for _, item := range v.dbVoterHistoryMap {
		history = append(history, item)
	}

	//2. Marshal the slice into json, lets pretty print it, but
	//   this is not required
	dataVoter, err := json.MarshalIndent(voterList, "", "  ")
	if err != nil {
		return err
	}

	dataHistory, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	//3. Write the json to our file
	err = os.WriteFile(v.dbVoterFileName, dataVoter, 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(v.dbHistoryFileName, dataHistory, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (v *DBVoter) loadDB() error {
	voterData, err := os.ReadFile(v.dbVoterFileName)
	if err != nil {
		return err
	}

	historyData, err := os.ReadFile(v.dbHistoryFileName)
	if err != nil {
		return err
	}

	//Now let's unmarshal the data into our map
	var voterList []Voter
	err = json.Unmarshal(voterData, &voterList)
	if err != nil {
		return err
	}

	var history []VoterHistory
	err = json.Unmarshal(historyData, &history)
	if err != nil {
		return err
	}

	//Now let's iterate over our slice and add each item to our map
	for _, item := range voterList {
		v.dbVoterMap[int(item.VoterId)] = item
	}

	for _, item := range history {
		v.dbVoterHistoryMap[int(item.PollId)] = item
	}

	return nil
}

func (v *DBVoter) findVoterByEmail(email string) *Voter {
	for _, voter := range v.dbVoterMap {
		if voter.Email == email {
			return &voter
		}
	}
	return nil
}

func (v *DBVoter) getAllVoterHistory(voter Voter) []listing.VoterHistory {

	var voterHistory []listing.VoterHistory

	for _, voteId := range voter.VoterHistory {

		history := v.dbVoterHistoryMap[int(voteId)]

		item := listing.VoterHistory{
			PollId:   history.PollId,
			VoteId:   history.VoteId,
			VoteDate: history.VoteDate,
			Created:  history.Created,
			Modified: history.Modified,
		}
		voterHistory = append(voterHistory, item)
	}

	return voterHistory
}

func (v *DBVoter) printVoter(item Voter) {
	jsonBytes, _ := json.MarshalIndent(item, "", "  ")
	fmt.Println(string(jsonBytes))
}

func (v *DBVoter) printHistory(item VoterHistory) {
	jsonBytes, _ := json.MarshalIndent(item, "", "  ")
	fmt.Println(string(jsonBytes))
}
