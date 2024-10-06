//go:build !exclude_swagger
// +build !exclude_swagger

// Package profile provides profile information and it's functionality.
package user

import (
	"fmt"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/service/categories"
)

type UserProfile struct {
	Surname   string `json:"surname"` //*changed
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
	AvatarURL string `json:"avatar_url"`
}

var (
	limitStr  string = "20"
	offsetStr string = "0"
)

func (s *Service) GetProfile(userID string) (*UserProfile, error) {
	surname, name, err := categories.GetUserInfoFromDB(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user information")
	}

	avatarURL, err := GetAvatarInfo(userID)
	if err != nil {
		avatarURL = "null"
	}

	return &UserProfile{
		UserID:    userID,
		Surname:   surname,
		Name:      name,
		AvatarURL: avatarURL,
	}, nil
}

func (s *Service) UpdateUserNameInDB(userID, newName, newSurname string) error {
	_, err := mydb.GlobalDB.Exec("UPDATE users SET name = $1, surname = $3 WHERE id = $2", newName, userID, newSurname)
	return err
}
