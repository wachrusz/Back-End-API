package user

import (
	"github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/service/categories"
	"mime/multipart"
)

type Service struct {
	repo       *mydatabase.Database
	categories categories.Categories
}

func NewService(repo *mydatabase.Database, cats categories.Categories) *Service {
	return &Service{
		repo:       repo,
		categories: cats,
	}
}

type Users interface {
	DeleteTokens(email, deviceID string) error
	GetTokenPairsAmount(email string) (int, error)
	Logout(device, userID string) error
	GetProfile(userID string) (*UserProfile, error)
	UpdateUserNameInDB(userID, newName, newSurname string) error
	UploadAvatar(userID string, f multipart.File) (string, error)
	GetAvatar(id string) ([]byte, error)
	UploadIcon(file multipart.File) (string, error)
	GetIcon(id string) ([]byte, error)
	GetIconsFromDataSource() ([]Icon, error)
	GetUserByEmail(email string) (IdentificationData, bool)
	Register(email, password string) error
	SaveSessionToDatabase(userID, deviceID, token string) error
	RemoveSessionFromDatabase(deviceID, userID string) error
	GetUserIDFromUsersDatabase(usernameOrDeviceID string) (string, error)
}
