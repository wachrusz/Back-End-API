package v1

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"github.com/wachrusz/Back-End-API/secret"
	"net/http"
)

func (h *MyHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	userID, ok := utility.GetUserIDFromContext(r.Context())
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

	encryptedID, err := h.s.Users.UploadAvatar(userID, file)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
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

func (h *MyHandler) GetAvatarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" || len(id) < 20 {
		jsonresponse.SendErrorResponse(w, errors.New("Something went wrong: "), http.StatusBadRequest)
		return
	}

	image, err := h.s.Users.GetAvatar(id)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(image)
}

func (h *MyHandler) UploadIconHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("image")
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error retrieving the file: "+err.Error()), http.StatusBadRequest)
		return
	}
	defer file.Close()

	encryptedID, err := h.s.Users.UploadIcon(file)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
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

func (h *MyHandler) GetIconHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" || len(id) < 20 {
		jsonresponse.SendErrorResponse(w, errors.New("Something went wrong: "), http.StatusBadRequest)
		return
	}

	image, err := h.s.Users.GetIcon(id)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(image)
}

type List struct {
	Icons      []user.Icon `json:"icons"`
	Message    string      `json:"message"`
	StatusCode string      `json:"code"`
}

func (h *MyHandler) GetIconsURLsHandler(w http.ResponseWriter, r *http.Request) {
	icons, err := h.s.Users.GetIconsFromDataSource()
	if err != nil {
		list := List{
			Icons:      nil,
			Message:    err.Error(),
			StatusCode: "500",
		}
		response := map[string]interface{}{
			"response": list,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	list := List{
		Icons:      icons,
		Message:    "Successfully got icons",
		StatusCode: "200",
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"response": list,
	}
	json.NewEncoder(w).Encode(response)
}
