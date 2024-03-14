package adding

import "errors"

var ErrDuplicateVoter = errors.New("voter already exists")
var ErrDuplicateVoterHistory = errors.New("voter history already exists")
var ErrBadId = errors.New("all Ids must be greater than 0")

type Service interface {
	RegisterVoter(Voter) (int, error)
	AddVoterHistory(VoterHistory) error
}

type Repository interface {
	AddVoter(Voter) (int, error)
	AddVoterHistory(VoterHistory) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) RegisterVoter(v Voter) (int, error) {

	id, err := s.r.AddVoter(v)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *service) AddVoterHistory(v VoterHistory) error {

	if v.VoterId < 1 {
		return ErrBadId
	}

	if v.PollId < 1 {
		return ErrBadId
	}

	if v.VoteId < 1 {
		return ErrBadId
	}

	err := s.r.AddVoterHistory(v)
	if err != nil {
		return err
	}
	return nil
}
