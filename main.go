package main

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"

	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Id           uint   `gorm:"autoIncrement,primaryKey"`
	Company_name string `valid:"string"`
	Zip_Code     string `gorm:"size:5" valid:"length(5),numeric"`
	Website      string `valid:"url,optional"`
}

func init() {
	valid.SetFieldsRequiredByDefault(true)
}

func read_csv(file_path string) [][]string {
	// 1. Open the file
	f, err := os.Open(file_path)
	if err != nil {
		log.Panic("Unable to load CSV file, be sure to check the file path.\n", err)
	}
	log.Info("Read the CSV file sucessfully")
	// 2. Parse CSV contents (base data)
	csv_reader := csv.NewReader(f)
	csv_reader.Comma = ';' // Custom separator

	csv_companies, err := csv_reader.ReadAll()
	if err != nil {
		log.Panic("Unable to parse CSV contents, check the file syntax.\n", err)
	}

	// return the file data as a list of lists of strings,
	// removing the first entry as it is a header
	return csv_companies[1:]
}

func format_company_data(company_data [][]string) [][]string {
	for _, company := range company_data {
		// Titlecase all names
		company[0] = strings.Title(company[0])
		if len(company) == 3 {
			company[2] = strings.ToLower(company[2])
		}
	}

	return company_data

}

func setup_database(remove_current bool) *gorm.DB {
	if remove_current {
		log.Info("Removing existing data...")
		os.Remove("test.db")
	}
	// Database connection/setup
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Panic("Failed to connect database!\n", err)
	}
	log.Info("Database connection sucessfull...")

	// Migrate the schema
	db.AutoMigrate(&Company{})

	return db
}

func populate_database(company_data [][]string, db *gorm.DB) {
	log.Info("Populating database with ", len(company_data), " entries...")
	for _, company := range company_data {
		db.Create(&Company{Company_name: company[0], Zip_Code: company[1]})
	}
}

func merge_data(current_db *gorm.DB, additional_data [][]string) {
	log.Info("Attempting to merge ", len(additional_data), " entries into the database...")
	for _, company_data := range additional_data {
		var comp Company
		result := current_db.Where("company_name LIKE ? AND zip_code = ?", "%"+company_data[0]+"%", company_data[1]).First(&comp)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Error("Company not found...skipping...")
			continue
		}
		current_db.Model(&comp).Update("website", company_data[2])
	}
}

/*
System requirements:
  1. Open CSV, parse and save the data on DB;
  2. Acquire new data (CSV) and merge with existing data;
  3. Serve information as a REST API
*/
func main() {
	db := setup_database(true)

	// Data acquisition and filtering
	company_data := read_csv("input_data/q1_catalog.csv")

	company_data = format_company_data(company_data)
	populate_database(company_data, db)

	mergeable_data := read_csv("input_data/q2_clientData.csv")
	mergeable_data = format_company_data(mergeable_data)

	merge_data(db, mergeable_data)
}
