//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package report

import (
	"fmt"
	"main/internal/auth"
	"net/http"
	"os"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type ReportData struct {
	Date        string
	Description string
	Amount      string
}

func generateExcelReport(data []ReportData) (*excelize.File, error) {
	file := excelize.NewFile()

	index, err := file.NewSheet("Sheet1")
	if err != nil {
		return file, err
	}

	headers := map[string]string{
		"A1": "Дата",
		"B1": "Описание",
		"C1": "Сумма",
	}

	for cell, value := range headers {
		file.SetCellValue("Sheet1", cell, value)
	}

	for i, entry := range data {
		row := i + 2
		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), entry.Date)
		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), entry.Description)
		file.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), entry.Amount)
	}

	file.SetActiveSheet(index)

	return file, nil
}

func generatePDFReport(data []ReportData) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 12)

	for _, entry := range data {
		pdf.Cell(40, 10, fmt.Sprintf("Дата: %s", entry.Date))
		pdf.Ln(10)
		pdf.Cell(40, 10, fmt.Sprintf("Описание: %s", entry.Description))
		pdf.Ln(10)
		pdf.Cell(40, 10, fmt.Sprintf("Сумма: %s", entry.Amount))
		pdf.Ln(10)
	}

	return pdf.OutputFileAndClose("report.pdf")
}

// @Summary Exports report
// @Description Get a financial report .
// @Tags App
// @Param expense body models.ConnectedAccount true "ConnectedAccount object"
// @Success 201 {string} string "Connected account created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error adding connected account"
// @Security JWT
// @Router /app/report [get]
func ExportHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data, err := fetchDataFromDatabase(userID)
	if err != nil {
		http.Error(w, "Failed to fetch data from database", http.StatusInternalServerError)
		return
	}

	excelReport, err := generateExcelReport(data)
	if err != nil {
		http.Error(w, "Failed to generate Excel report", http.StatusInternalServerError)
		return
	}

	excelFilename := "report.xlsx"
	if err := excelReport.SaveAs(excelFilename); err != nil {
		http.Error(w, "Failed to save Excel report", http.StatusInternalServerError)
		return
	}

	if err := generatePDFReport(data); err != nil {
		http.Error(w, "Failed to generate PDF report", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, excelFilename)
	http.ServeFile(w, r, "report.pdf")

	defer os.Remove(excelFilename)
	defer os.Remove("report.pdf")
}
