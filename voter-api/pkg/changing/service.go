package changing

import "errors"

var ErrDuplicateVoter = errors.New("voter already exists")
var ErrDuplicateVoterHistory = errors.New("voter history already exists")

type Service interface {
	UpdateVoterDemographics(Voter) error
	ReviseVoterHistory(VoterHistory) error
}

type Repository interface {
	UpdateVoter(Voter) error
	UpdateVoterHistory(VoterHistory) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) UpdateVoterDemographics(v Voter) error {
	err := s.r.UpdateVoter(v)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ReviseVoterHistory(v VoterHistory) error {

	err := s.r.UpdateVoterHistory(v)
	if err != nil {
		return err
	}
	return nil
}
