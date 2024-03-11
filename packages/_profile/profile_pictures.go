package profile

import (
	auth "main/packages/_auth"
	encryption "main/packages/_encryption"
	jsonresponse "main/packages/_json_response"
	mydb "main/packages/_mydatabase"
	"main/secret"

	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("Error getting userID: "), http.StatusUnauthorized)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error retrieving the file: "+err.Error()), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error reading the file: "+err.Error()), http.StatusInternalServerError)
		return
	}

	encryptedID, err := encryption.EncryptID(userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error encrypting ID: "+err.Error()), http.StatusInternalServerError)
		return
	}

	err = saveAvatarInfo(userID, fileBytes, encryptedID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error saving avatar info: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfuly uploaded a file",
		"status_code": http.StatusCreated,
		"avatar_url":  "https://" + secret.Secret.BaseURL + "/v1/profile/image/get/" + encryptedID,
	}
	json.NewEncoder(w).Encode(response)
}

// ! ЛИКВИДИРОВАТЬ
func GetAvatarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" || len(id) < 20 {
		jsonresponse.SendErrorResponse(w, errors.New("Something went wrong: "), http.StatusBadRequest)
		return
	}

	encryptedID, err := encryption.DecryptID(id)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error decrypting ID: "+err.Error()), http.StatusInternalServerError)
		return
	}

	image, err := GetAvatarBytes(encryptedID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error getting avatar: "+err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(image)
}

func saveAvatarInfo(userID string, imageBytes []byte, encryptedID string) error {
	url, err := encryption.EncryptID("https://" + secret.Secret.BaseURL + "/v1/profile/image/get/" + encryptedID)
	if err != nil {
		return err
	}
	_, err = mydb.GlobalDB.Exec("INSERT INTO profile_images (profile_id, image_data, url) VALUES ($1, $2, $3) ON CONFLICT (profile_id) DO UPDATE SET image_data = $2", userID, imageBytes, url)
	if err != nil {
		return err
	}
	return err
}

func GetAvatarBytes(userID string) ([]byte, error) {
	var bytes []byte
	err := mydb.GlobalDB.QueryRow("SELECT image_data FROM profile_images WHERE profile_id = $1", userID).Scan(&bytes)
	return bytes, err
}

func GetAvatarInfo(userID string) (string, error) {
	var avatarURL string
	err := mydb.GlobalDB.QueryRow("SELECT url FROM profile_images WHERE profile_id = $1", userID).Scan(&avatarURL)
	if err != nil {
		return "null", err
	}
	decryptedURL, err := encryption.DecryptID(avatarURL)
	if err != nil {
		return "null", err
	}
	return decryptedURL, err
}
