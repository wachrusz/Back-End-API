package user

import (
	"fmt"
)

func (s *Service) Logout(device, userID string) error {
	err := s.RemoveSessionFromDatabase(device, userID)
	if err != nil {
		return fmt.Errorf("error removing session from db: %v", err)
	}

	delete(s.ActiveUsers, userID)
	return nil
}
