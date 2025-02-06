package goals

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	repo "github.com/wachrusz/Back-End-API/internal/repository"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
)

type Goals interface {
	Create(goal *models.Goal) (int64, error)
	Update(goal *models.Goal) error
	Delete(id int64, userID int64) error
	ListByUserID(userID int64) ([]models.Goal, error)
	Details(id int64, userID int64) (*models.GoalDetails, error)
	NewTransaction(transaction *models.GoalTransaction, userID int64) (*models.GoalDetails, error)
}

type Service struct {
	goalRepo        repo.GoalRepo
	transactionRepo repo.GoalTransactionRepo
}

func NewService(gr repo.GoalRepo, tr repo.GoalTransactionRepo) *Service {
	return &Service{goalRepo: gr, transactionRepo: tr}
}

func (s *Service) NewTransaction(transaction *models.GoalTransaction, userID int64) (*models.GoalDetails, error) {
	_, err := s.transactionRepo.Create(transaction, userID)
	if err != nil {
		return nil, err
	}

	return s.Details(transaction.GoalID, userID)
}

func (s *Service) Details(goalID, userID int64) (*models.GoalDetails, error) {
	details, err := s.goalRepo.Details(goalID, userID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no goal (id %d) found: %w", goalID, myerrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}

	return details, nil
}

func (s *Service) Create(goal *models.Goal) (int64, error) {
	return s.goalRepo.Create(goal)
}

func (s *Service) Update(goal *models.Goal) error {
	return s.goalRepo.Update(goal)
}

func (s *Service) Delete(id int64, userID int64) error {
	return s.goalRepo.Delete(id, userID)
}

func (s *Service) ListByUserID(userID int64) ([]models.Goal, error) {
	return s.goalRepo.ListByUserID(userID)
}
