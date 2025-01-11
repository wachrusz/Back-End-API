package user

import (
	"database/sql"
	"errors"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
)

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type IdentificationData struct {
	Email          string
	HashedPassword string
}

func (s *Service) GetUserByEmail(email string) (IdentificationData, bool) {
	var userData IdentificationData
	var id int

	row := s.repo.QueryRow("SELECT id, email, hashed_password FROM users WHERE email = $1", email)
	err := row.Scan(&id, &userData.Email, &userData.HashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return userData, false
	} else if err != nil {
		return userData, false
	}
	return userData, true
}

func (s *Service) Register(email, password string) error {
	if _, exists := s.GetUserByEmail(email); exists {
		return errors.New("Already exists")
	}

	if email == "" || password == "" {
		return errors.New("Blank fields are not allowed")
	}

	hashedPassword, err := utility.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.repo.Exec("INSERT INTO users (email, hashed_password) VALUES ($1, $2)", email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
