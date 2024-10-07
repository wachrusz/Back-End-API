package service

import (
	"github.com/wachrusz/Back-End-API/internal/models"
	"github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/service/categories"
	"github.com/wachrusz/Back-End-API/internal/service/currency"
	"github.com/wachrusz/Back-End-API/internal/service/email"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	"mime/multipart"
	"time"
)

type Services struct {
	Users      Users
	Categories Categories
	Emails     Emails
	//Reports    *report.Service
	Currency CurrencyService
}

type Dependencies struct {
	Repo *mydatabase.Database
}

func NewServices(deps Dependencies) (*Services, error) {
	u := user.NewService(deps.Repo)
	c, err := currency.NewService(deps.Repo)
	if err != nil {
		return nil, err
	}
	return &Services{
		Users:      u,
		Categories: categories.NewService(deps.Repo, c),
		Emails:     email.NewService(deps.Repo, u),
		//Reports:    report.NewService(deps.Repo),
		Currency: c,
	}, nil
}

type Emails interface {
	SendEmail(to, subject, body string) error
	SendConfirmationEmail(email, token string) error
	DecryptToken(encryptedToken string) (string, error)
	SaveConfirmationCode(email, confirmationCode, token string) error
	CheckConfirmationCode(email, token, enteredCode string) email.CheckResult
	DeleteConfirmationCode(email string, code string)
	GetConfirmationCode(email string) (string, error)
	ConfirmEmail(token, code, deviceID string) (*email.TokenDetails, error)
	ConfirmEmailLogin(token, code, deviceID string) (*email.TokenDetails, error)
	ResetPasswordConfirm(token, code string) error
	ResetPassword(email, password string) error
}

type CurrencyService interface {
	ScheduleCurrencyUpdates()
}

type Categories interface {
	GetAnalyticsFromDB(userID, currencyCode, limitStr, offsetStr, startDateStr, endDateStr string) (*categories.Analytics, error)
	GetTrackerFromDB(userID, currencyCode, limitStr, offsetStr string) (*categories.Tracker, error)
	GetUserInfoFromDB(userID string) (string, string, error)
	GetMoreFromDB(userID string) (*categories.More, error)
	GetAppFromDB(userID string) (*models.App, error)
	GetSubscriptionFromDB(userID string) (*models.Subscription, error)
	GetConnectedAccountsFromDB(userID string) ([]models.ConnectedAccount, error)
	GetCategorySettingsFromDB(userID string) (*models.CategorySettings, error)
	GetOperationArchiveFromDB(userID, limit, offset string) ([]models.Operation, error)
}

type Users interface {
	Login(email, password string) (string, error)
	RefreshToken(rt, userID string) (string, string, error)
	GenerateToken(userID string, device_id string, duration time.Duration) (*email.TokenDetails, error)
	DeleteTokens(email, deviceID string) error
	GetTokenPairsAmount(email string) (int, error)
	Logout(device, userID string) error
	GetProfile(userID string) (*user.UserProfile, error)
	UpdateUserNameInDB(userID, newName, newSurname string) error
	PrimaryRegistration(email, password string) (string, error)
	ResetPassword(email string) error
	ChangePasswordForRecover(email, password, resetToken string) error
	Register(email, password string) error
	HashPassword(password string) (string, error)
	InitActiveUsers()
	GetActiveUser(userID string) user.ActiveUser
	IsUserActive(userID string) bool
	AddActiveUser(userID, email, deviceID, token string)
	RemoveActiveUser(userID string)
	SetAccessToken(userID, newAccessToken string)
	SaveSessionToDatabase(email, deviceID, userID, token string) error
	UpdateLastActivity(userID string) error
	CheckSessionInDatabase(email, deviceID string) (bool, error)
	RemoveSessionFromDatabase(deviceID, userID string) error
	IsDeviceIDAlreadyUsed(email, deviceID string) (error, bool)
	GetUserIDFromUsersDatabase(usernameOrDeviceID string) (string, error)
	GetUserIDFromSessionDatabase(usernameOrDeviceID string) (string, error)
	UploadAvatar(userID string, f multipart.File) (string, error)
	GetAvatar(id string) ([]byte, error)
	UploadIcon(file multipart.File) (string, error)
	GetIcon(string) ([]byte, error)
	GetIconsFromDataSource() ([]user.Icon, error)
}
