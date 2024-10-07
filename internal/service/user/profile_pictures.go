package user

import (
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/pkg/encryption"
	"github.com/wachrusz/Back-End-API/secret"
	"math/rand"
	"mime/multipart"
	"strconv"
	"time"

	"io/ioutil"
)

type Icon struct {
	id        string `json:"id"`
	url       string `json:"url"`
	serviceID string `json:"service_id"`
}

func (s *Service) UploadAvatar(userID string, f multipart.File) (string, error) {
	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("%w: failed to read avatar: %v", myerrors.ErrInternal, err)
	}

	encryptedID, err := encryption.EncryptID(userID)
	if err != nil {
		return "", fmt.Errorf("%w: failed to encrypt avatar: %v", myerrors.ErrInternal, err)
	}

	err = s.saveAvatarInfo(userID, fileBytes, encryptedID)
	if err != nil {
		return "", fmt.Errorf("%w: failed to save avatar: %v", myerrors.ErrInternal, err)
	}

	return encryptedID, nil
}

// ! ЛИКВИДИРОВАТЬ
func (s *Service) GetAvatar(id string) ([]byte, error) {
	encryptedID, err := encryption.DecryptID(id)
	if err != nil {
		return nil, fmt.Errorf("%w: decryption failed: %v", myerrors.ErrInternal, err)
	}

	var bytes []byte
	err = s.repo.QueryRow("SELECT image_data FROM profile_images WHERE profile_id = $1", encryptedID).Scan(&bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: error getting avatar: %v", myerrors.ErrInternal, err)
	}

	return bytes, nil
}

func (s *Service) UploadIcon(file multipart.File) (string, error) {
	rand.Seed(time.Now().UnixNano())
	userID_i := rand.Intn(20000000)
	userID := strconv.Itoa(userID_i)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("%w: failed to read icon: %v", myerrors.ErrInternal, err)
	}

	encryptedID, err := encryption.EncryptID(userID)
	if err != nil {
		return "", fmt.Errorf("%w: failed to encrypt icon: %v", myerrors.ErrInternal, err)
	}

	err = s.saveIconInfo(userID, fileBytes, encryptedID)
	if err != nil {
		return "", fmt.Errorf("%w: failed to save icon info: %v", myerrors.ErrInternal, err)
	}

	return encryptedID, nil
}

func (s *Service) GetIcon(id string) ([]byte, error) {
	encryptedID, err := encryption.DecryptID(id)
	if err != nil {
		return nil, fmt.Errorf("%w: decryption failed: %v", myerrors.ErrInternal, err)
	}

	var bytes []byte
	err = s.repo.QueryRow("SELECT image_data FROM service_images WHERE service_id = $1", encryptedID).Scan(&bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: error getting avatar: %v", myerrors.ErrInternal, err)
	}

	return bytes, nil
}

func (s *Service) GetIconsFromDataSource() ([]Icon, error) {
	query := "SELECT id, url, service_id FROM service_images"
	rows, err := s.repo.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var icons []Icon

	for rows.Next() {
		var icon Icon

		err = rows.Scan(&icon.id, &icon.url, &icon.serviceID)
		if err != nil {
			return nil, err
		}

		icons = append(icons, icon)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return icons, nil
}

func (s *Service) saveAvatarInfo(userID string, imageBytes []byte, encryptedID string) error {
	url, err := encryption.EncryptID("https://" + secret.Secret.BaseURL + "/v1/profile/image/get/" + encryptedID)
	if err != nil {
		return err
	}
	_, err = s.repo.Exec("INSERT INTO service_images (profile_id, image_data, url) VALUES ($1, $2, $3) ON CONFLICT (profile_id) DO UPDATE SET image_data = $2", userID, imageBytes, url)
	if err != nil {
		return err
	}
	return err
}

func (s *Service) saveIconInfo(userID string, imageBytes []byte, encryptedID string) error {
	url, err := encryption.EncryptID("https://" + secret.Secret.BaseURL + "/v1/api/emojis/get/" + encryptedID)
	if err != nil {
		return err
	}
	_, err = s.repo.Exec("INSERT INTO service_images (service_id, image_data, url) VALUES ($1, $2, $3) ON CONFLICT (service_id) DO UPDATE SET image_data = $2", userID, imageBytes, url)
	if err != nil {
		return err
	}
	return err
}

func (s *Service) getAvatarInfo(userID string) (string, error) {
	var avatarURL string
	err := s.repo.QueryRow("SELECT url FROM profile_images WHERE profile_id = $1", userID).Scan(&avatarURL)
	if err != nil {
		return "null", err
	}
	decryptedURL, err := encryption.DecryptID(avatarURL)
	if err != nil {
		return "null", err
	}
	return decryptedURL, err
}
