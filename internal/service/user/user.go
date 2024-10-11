package user

import (
	"github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/service/categories"
	"mime/multipart"
	"sync"
)

type Service struct {
	repo        *mydatabase.Database
	ActiveUsers map[string]ActiveUser
	categories  categories.Categories
	mutex       sync.Mutex
	activeMu    sync.Mutex
}

func NewService(repo *mydatabase.Database) *Service {
	return &Service{
		repo:        repo,
		ActiveUsers: make(map[string]ActiveUser),
		mutex:       sync.Mutex{},
		activeMu:    sync.Mutex{},
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
	InitActiveUsers()
	GetActiveUser(userID string) ActiveUser
	IsUserActive(userID string) bool
	AddActiveUser(user ActiveUser)
	RemoveActiveUser(userID string)
	SaveSessionToDatabase(email, deviceID, userID, token string) error
	UpdateLastActivity(userID string) error
	//CheckSessionInDatabase(email, deviceID string) (bool, error)
	RemoveSessionFromDatabase(deviceID, userID string) error
	//IsDeviceIDAlreadyUsed(email, deviceID string) (error, bool)
	GetUserIDFromUsersDatabase(usernameOrDeviceID string) (string, error)
	GetUserIDFromSessionDatabase(usernameOrDeviceID string) (string, error)
	SetAccessToken(userID, newAccessToken string)
}
