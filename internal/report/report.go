//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package report

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/xuri/excelize/v2"
)

type Service struct {
	repo *mydb.Database
}

func NewService(db *mydb.Database) *Service {
	return &Service{db}
}

type ReportData struct {
	Date        string
	Description string
	Amount      string
}

func (s *Service) generateExcelReport(data []ReportData) (*excelize.File, error) {
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

func (s *Service) generatePDFReport(data []ReportData) error {
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

func (s *Service) ExportHandler(userID string) error {
	data, err := s.fetchDataFromDatabase(userID)
	if err != nil {
		return myerrors.ErrInternal
	}

	excelReport, err := s.generateExcelReport(data)
	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	excelFilename := "report.xlsx"
	if err := excelReport.SaveAs(excelFilename); err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	if err := s.generatePDFReport(data); err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}
	return nil
}
