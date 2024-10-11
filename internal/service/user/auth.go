package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/vk"
	"net/http"
)

func (s *Service) DeleteTokens(email, deviceID string) error {
	if email != "" {
		err := s.deleteForEmail(email)
		if err != nil {
			return fmt.Errorf("%w: %v", myerrors.ErrDeletingTokens, err)
		}
		userID, err := s.GetUserIDFromUsersDatabase(email)
		if err != nil {
			return fmt.Errorf("%w: %v", myerrors.ErrDeletingTokens, err)
		}
		s.RemoveActiveUser(userID)
	}
	if deviceID != "" {
		err := s.deleteForDeviceID(deviceID)
		if err != nil {
			return fmt.Errorf("%w: %v", myerrors.ErrDeletingTokens, err)
		}
		userID, err := s.GetUserIDFromSessionDatabase(deviceID)
		if err != nil {
			return fmt.Errorf("%w: %v", myerrors.ErrDeletingTokens, err)
		}
		s.RemoveActiveUser(userID)
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

// ========= VK and GOOGLE =========

var (
	googleConfig = &oauth2.Config{
		ClientID:     "YOUR_GOOGLE_CLIENT_ID",
		ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}

	vkConfig = &oauth2.Config{
		ClientID:     "YOUR_VK_CLIENT_ID",
		ClientSecret: "YOUR_VK_CLIENT_SECRET",
		RedirectURL:  "http://localhost:3000/auth/vk/callback",
		Scopes:       []string{"email", "photos"},
		Endpoint:     vk.Endpoint,
	}

	oauthStateString = "random"
)

func (s *Service) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Service) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("googleConfig.Exchange() failed:", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := googleConfig.Client(context.Background(), token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		fmt.Println("client.Get() failed:", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer userInfo.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(userInfo.Body).Decode(&data); err != nil {
		fmt.Println("json.NewDecoder().Decode() failed:", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Println("Google User Info:", data)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *Service) HandleVKLogin(w http.ResponseWriter, r *http.Request) {
	url := vkConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Service) HandleVKCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := vkConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("vkConfig.Exchange() failed:", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := vkConfig.Client(context.Background(), token)
	userInfo, err := client.Get("https://api.vk.com/method/users.get?v=5.103")
	if err != nil {
		fmt.Println("client.Get() failed:", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer userInfo.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(userInfo.Body).Decode(&data); err != nil {
		fmt.Println("json.NewDecoder().Decode() failed:", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Println("VK User Info:", data)
	// Обработка полученных данных, например, регистрация пользователя в вашей системе.

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
