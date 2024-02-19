package profile

import (
	auth "main/packages/_auth"
	jsonresponse "main/packages/_json_response"
	mydb "main/packages/_mydatabase"

	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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
		response := map[string]interface{}{
			"message":     "Error retrieving the file",
			"status_code": http.StatusBadRequest,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		response := map[string]interface{}{
			"message":     "Error reading the file",
			"status_code": http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	err = saveAvatarInfo(userID, fileBytes)
	if err != nil {
		response := map[string]interface{}{
			"message":     "Error saving avatar info",
			"status_code": http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfuly uploaded a file",
		"status_code": http.StatusCreated,
	}
	json.NewEncoder(w).Encode(response)
}

func GetAvatarHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("Error getting userID"), http.StatusUnauthorized)
		return
	}

	imageBytes, err := getAvatarInfo(userID)
	if err != nil {
		response := map[string]interface{}{
			"message":     "Error getting avatar info",
			"image_bytes": []byte{},
			"status_code": http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	response := map[string]interface{}{
		"message":     "Successfuly got an image",
		"image_bytes": imageBytes,
		"status_code": http.StatusCreated,
	}
	json.NewEncoder(w).Encode(response)
}

func saveAvatarInfo(userID string, imageBytes []byte) error {
	_, err := mydb.GlobalDB.Exec("INSERT INTO profile_images (profile_id, image_data) VALUES ($1, $2)", userID, imageBytes)
	return err
}

func getAvatarInfo(userID string) ([]byte, error) {
	var imageBytes []byte
	err := mydb.GlobalDB.QueryRow("SELECT image_data FROM profile_images WHERE profile_id = $1", userID).Scan(&imageBytes)
	return imageBytes, err
}
