package service

import (
	"encoding/json"
	"errors"
	"net/http"

	jsonresponse "main/packages/_json_response"
	mydb "main/packages/_mydatabase"
)

func DeleteTokensHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	deviceID := r.URL.Query().Get("deviceID")
	if (email == "" && deviceID == "") || (email != "" && deviceID != "") {
		jsonresponse.SendErrorResponse(w, errors.New("Blank fields and two methods are not allowed"), http.StatusBadRequest)
		return
	}
	if email != "" {
		err := deleteForEmail(email)
		if err != nil {
			jsonresponse.SendErrorResponse(w, errors.New("Error while deleting tokens: "+err.Error()), http.StatusInternalServerError)
			return
		}
		userID, err := GetUserIDFromUsersDatabase(email)
		if err != nil {
			jsonresponse.SendErrorResponse(w, errors.New("Error while deleting tokens: "+err.Error()), http.StatusInternalServerError)
			return
		}
		RemoveActiveUser(userID)
	}
	if deviceID != "" {
		err := deleteForDeviceID(deviceID)
		if err != nil {
			jsonresponse.SendErrorResponse(w, errors.New("Error while deleting tokens: "+err.Error()), http.StatusInternalServerError)
			return
		}
		userID, err := GetUserIDFromSessionDatabase(deviceID)
		if err != nil {
			jsonresponse.SendErrorResponse(w, errors.New("Error while deleting tokens: "+err.Error()), http.StatusInternalServerError)
			return
		}
		RemoveActiveUser(userID)
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":     "Successfuly deleted tokens",
		"status_code": http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

func GetTokenPairsAmmountHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		jsonresponse.SendErrorResponse(w, errors.New("Blank fields are not allowed"), http.StatusBadRequest)
		return
	}
	ammount, err := getAmmount(email)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error while counting sessions: "+err.Error()), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message":     "Successfuly got ammount",
		"ammount":     ammount,
		"status_code": http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

func getAmmount(email string) (int, error) {
	var ammount int
	err := mydb.GlobalDB.QueryRow("SELECT COUNT(*) FROM sessions WHERE email = $1", email).Scan(&ammount)
	if err != nil {
		return 0, err
	}
	return ammount, nil
}

func deleteForEmail(email string) error {
	_, err := mydb.GlobalDB.Exec("DELETE FROM sessions WHERE email = $1", email)
	if err != nil {
		return err
	}
	return nil
}

func deleteForDeviceID(deviceID string) error {
	_, err := mydb.GlobalDB.Exec("DELETE FROM sessions WHERE device_id = $1", deviceID)
	if err != nil {
		return err
	}
	return nil
}
