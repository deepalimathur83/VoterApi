package listing

import "errors"

var ErrDuplicateVoter = errors.New("voter already exists")
var ErrDuplicateVoterHistory = errors.New("voter history already exists")

type Service interface {
	GetVoter(id int) (Voter, error)
	GetVoterHistory(voterId uint, voteId uint) (VoterHistory, error)
	GetAllVoters() ([]Voter, error)
}

type Repository interface {
	GetVoter(id int) (Voter, error)
	GetVoterHistory(voterId uint, voteId uint) (VoterHistory, error)
	GetVoters() ([]Voter, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetVoter(id int) (Voter, error) {
	var voter Voter
	voter, err := s.r.GetVoter(id)
	if err != nil {
		return voter, err
	}
	return voter, nil
}

func (s *service) GetVoterHistory(voterId uint, voteId uint) (VoterHistory, error) {

	var history VoterHistory
	history, err := s.r.GetVoterHistory(voterId, voteId)
	if err != nil {
		return history, err
	}
	return history, nil
}

func (s *service) GetAllVoters() ([]Voter, error) {

	var voterList []Voter

	voterList, err := s.r.GetVoters()
	if err != nil {
		return voterList, err
	}
	return voterList, nil
}
