package profile

import (
	"github.com/wachrusz/Back-End-API/internal/auth"
	"github.com/wachrusz/Back-End-API/pkg/encryption"
	mydb "github.com/wachrusz/Back-End-API/pkg/mydatabase"
	"github.com/wachrusz/Back-End-API/secret"
	"math/rand"
	"strconv"
	"time"

	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type List struct {
	icons       []Icon `json:"icons"`
	message     string `json:"message"`
	status_code string `json:"code"`
}
type Icon struct {
	id        string `json:"id"`
	url       string `json:"url"`
	serviceID string `json:"service_id"`
}

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
	w.WriteHeader(response["status_code"].(int))
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

func UploadIconHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	rand.Seed(time.Now().UnixNano())
	userID_i := rand.Intn(20000000)
	userID := strconv.Itoa(userID_i)

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

	err = saveIconInfo(userID, fileBytes, encryptedID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error saving avatar info: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfuly uploaded a file",
		"status_code": http.StatusCreated,
		"avatar_url":  "https://" + secret.Secret.BaseURL + "/v1/api/emojis/get/" + encryptedID,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func GetIconHandler(w http.ResponseWriter, r *http.Request) {
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

	image, err := GetIconBytes(encryptedID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error getting avatar: "+err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(image)
}

func GetIconsURLs(w http.ResponseWriter, r *http.Request) {
	icons, err := getIconsFromDataSource()
	if err != nil {
		list := List{
			icons:       nil,
			message:     err.Error(),
			status_code: "500",
		}
		response := map[string]interface{}{
			"response": list,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	list := List{
		icons:       icons,
		message:     "Successfully got icons",
		status_code: "200",
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"response": list,
	}
	json.NewEncoder(w).Encode(response)
}

func getIconsFromDataSource() ([]Icon, error) {
	query := "SELECT id, url, service_id FROM service_images"
	rows, err := mydb.GlobalDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	icons := []Icon{}

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

func saveAvatarInfo(userID string, imageBytes []byte, encryptedID string) error {
	url, err := encryption.EncryptID("https://" + secret.Secret.BaseURL + "/v1/profile/image/get/" + encryptedID)
	if err != nil {
		return err
	}
	_, err = mydb.GlobalDB.Exec("INSERT INTO service_images (profile_id, image_data, url) VALUES ($1, $2, $3) ON CONFLICT (profile_id) DO UPDATE SET image_data = $2", userID, imageBytes, url)
	if err != nil {
		return err
	}
	return err
}

func saveIconInfo(userID string, imageBytes []byte, encryptedID string) error {
	url, err := encryption.EncryptID("https://" + secret.Secret.BaseURL + "/v1/api/emojis/get/" + encryptedID)
	if err != nil {
		return err
	}
	_, err = mydb.GlobalDB.Exec("INSERT INTO service_images (service_id, image_data, url) VALUES ($1, $2, $3) ON CONFLICT (service_id) DO UPDATE SET image_data = $2", userID, imageBytes, url)
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

func GetIconBytes(userID string) ([]byte, error) {
	var bytes []byte
	err := mydb.GlobalDB.QueryRow("SELECT image_data FROM service_images WHERE service_id = $1", userID).Scan(&bytes)
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
