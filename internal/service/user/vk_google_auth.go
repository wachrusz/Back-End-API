package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/vk"
)

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
