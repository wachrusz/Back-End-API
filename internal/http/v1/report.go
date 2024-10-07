package v1

// @Summary Exports report
// @Description Get a financial report .
// @Tags App
// @Param expense body models.ConnectedAccount true "ConnectedAccount object"
// @Success 201 {string} string "Connected account created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error adding connected account"
// @Security JWT
// @Router /app/report [get]
//func (h *MyHandler) ExportHandler(w http.ResponseWriter, r *http.Request) {
//	userID, ok := utility.GetUserIDFromContext(r.Context())
//	if !ok {
//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
//		return
//	}
//
//	err := h.s.Reports.ExportHandler(userID)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	excelFilename := "report.xlsx"
//	http.ServeFile(w, r, excelFilename)
//	http.ServeFile(w, r, "report.pdf")
//
//	defer os.Remove(excelFilename)
//	defer os.Remove("report.pdf")
//}
