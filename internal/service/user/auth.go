package user

import (
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
)

func (s *Service) DeleteTokens(email, deviceID string) error {
	if email != "" {
		err := s.deleteForEmail(email)
		if err != nil {
			return fmt.Errorf("%w: %v", myerrors.ErrDeletingTokens, err)
		}
		userID, err := GetUserIDFromUsersDatabase(email)
		if err != nil {
			return fmt.Errorf("%w: %v", myerrors.ErrDeletingTokens, err)
		}
		RemoveActiveUser(userID)
	}
	if deviceID != "" {
		err := s.deleteForDeviceID(deviceID)
		if err != nil {
			return fmt.Errorf("%w: %v", myerrors.ErrDeletingTokens, err)
		}
		userID, err := GetUserIDFromSessionDatabase(deviceID)
		if err != nil {
			return fmt.Errorf("%w: %v", myerrors.ErrDeletingTokens, err)
		}
		RemoveActiveUser(userID)
	}

	return nil
}

func (s *Service) GetTokenPairsAmount(email string) (int, error) {
	var amount int
	err := s.repo.QueryRow("SELECT COUNT(*) FROM sessions WHERE email = $1", email).Scan(&amount)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func (s *Service) deleteForEmail(email string) error {
	_, err := s.repo.Exec("DELETE FROM sessions WHERE email = $1", email)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) deleteForDeviceID(deviceID string) error {
	_, err := s.repo.Exec("DELETE FROM sessions WHERE device_id = $1", deviceID)
	if err != nil {
		return err
	}
	return nil
}
