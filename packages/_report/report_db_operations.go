package report

import (
	mydb "main/packages/_mydatabase"
)

func fetchDataFromDatabase(userID string) ([]ReportData, error) {

	rows, err := mydb.GlobalDB.Query("SELECT date, description, amount FROM transactions WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []ReportData
	for rows.Next() {
		var date, description, amount string
		if err := rows.Scan(&date, &description, &amount); err != nil {
			return nil, err
		}
		data = append(data, ReportData{Date: date, Description: description, Amount: amount})
	}

	return data, nil
}
