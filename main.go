package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Company struct {
	Id           uint   `gorm:"autoIncrement,primaryKey" json:"id"`
	Company_name string `valid:"string" json:"name"`
	Zip_Code     string `gorm:"size:5" valid:"length(5),numeric" json:"zip"`
	Website      string `valid:"url,optional" json:"website"`
}

type QueryData struct {
	Name     string `json:"name"`
	Zip_Code string `json:"zip_code"`
}

type APIResponse struct {
	Message string `json:"message"`
}

type DB struct {
	CONNECTION_STRING string
}

var db = DB{CONNECTION_STRING: "prod.db"}

func init() {
	valid.SetFieldsRequiredByDefault(true)
}

// Network functions
func createServer(port int) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/import", importData).Methods("POST")
	router.HandleFunc("/companies/search", searchCompany).Methods("GET")
	router.Use(requestLoggingMiddleware)
	log.Info("Attempting to serve API on ", port)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), router))
}

func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(time.Now(), " - Received ", r.Method, " request on ", r.RequestURI)
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func importData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	file, _, err := req.FormFile("file")
	// either the file is unavailable for in the request body or the field name was wrong
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(APIResponse{Message: "Error while reading file from request."})
		return
	}
	defer file.Close()

	f_name := "uploads/" + uuid.NewString()
	f, err := os.Create(f_name)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(APIResponse{Message: "Error while attempting to read file contents."})
		return
	}
	io.Copy(f, file)
	f.Close()

	new_data := readCsv(f_name)
	new_data = formatCompanyData(new_data)
	db_ref := db.createConnection()
	merge_data(db_ref, new_data)

	json.NewEncoder(res).Encode(APIResponse{Message: "Operation finished Sucessfully"})
	os.Remove(f_name)
}

func searchCompany(res http.ResponseWriter, req *http.Request) {

	var reqBody QueryData
	if req.Body == nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(APIResponse{Message: "Body not found in the request"})
		return
	}
	json.NewDecoder(req.Body).Decode(&reqBody)

	company := queryDB(reqBody)
	res.Header().Set("Content-Type", "application/json")
	if company.Id != 0 {
		json.NewEncoder(res).Encode(company)
		return
	}
	res.WriteHeader(http.StatusNotFound)
	json.NewEncoder(res).Encode(APIResponse{Message: "Company not found."})

}

func queryDB(q QueryData) Company {
	var company Company
	dbConn := db.createConnection()
	result := dbConn.Where("company_name LIKE ? AND zip_code = ?", "%"+q.Name+"%", q.Zip_Code).First(&company)
	if result.Error != nil {
		log.Error("Could not find company...")
	}
	return company
}

// File operations
func readCsv(filePath string) [][]string {
	// 1. Open the file
	f, err := os.Open(filePath)
	if err != nil {
		log.Panic("Unable to load CSV file, be sure to check the file path.\n", err)
	}
	log.Info("Read the CSV file sucessfully")
	// 2. Parse CSV contents (base data)
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';' // Custom separator

	csvCompanies, err := csvReader.ReadAll()
	if err != nil {
		log.Panic("Unable to parse CSV contents, check the file syntax.\n", err)
	}

	// return the file data as a list of lists of strings,
	// removing the first entry as it is a header
	return csvCompanies[1:]
}

func formatCompanyData(company_data [][]string) [][]string {
	// TODO: return error if the data has more than 3 components
	for _, company := range company_data {
		// Titlecase all names
		company[0] = strings.Title(company[0])
		if len(company) == 3 {
			company[2] = strings.ToLower(company[2])
		}
	}

	return company_data
}

func (db *DB) createConnection() *gorm.DB {
	dbConn, err := gorm.Open(sqlite.Open(db.CONNECTION_STRING), &gorm.Config{})
	if err != nil {
		log.Panic("Failed to connect to the database!\n", err)
	}
	return dbConn
}

func setupDatabase(remove_current bool) *gorm.DB {
	if remove_current {
		log.Info("Removing existing data...")
		os.Remove(db.CONNECTION_STRING)
	}
	// Database connection/setup
	dbConn := db.createConnection()
	log.Info("Database connection sucessfull...")

	// Migrate the schema
	dbConn.AutoMigrate(&Company{})

	return dbConn
}

func populateDatabase(companyData [][]string, db *gorm.DB) {
	log.Info("Populating database with ", len(companyData), " entries...")
	for _, company := range companyData {
		db.Create(&Company{Company_name: company[0], Zip_Code: company[1]})
	}
}

func merge_data(currentDB *gorm.DB, additionalData [][]string) {
	log.Info("Attempting to merge ", len(additionalData), " entries into the database...")
	for _, company_data := range additionalData {
		var comp Company
		result := currentDB.Where("company_name LIKE ? AND zip_code = ?", "%"+company_data[0]+"%", company_data[1]).First(&comp)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Error("Company not found...skipping...")
			continue
		}
		currentDB.Model(&comp).Update("website", company_data[2])
	}
}

/*
System requirements:
  1. Open CSV, parse and save the data on DB;
  2. Acquire new data (CSV) and merge with existing data;
  3. Serve information as a REST API
*/
func main() {
	db := setupDatabase(true)

	// Initial data acquisition
	company_data := readCsv("input_data/q1_catalog.csv")
	company_data = formatCompanyData(company_data)
	// Initial data population
	populateDatabase(company_data, db)

	createServer(8000)
}
