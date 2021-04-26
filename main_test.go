package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const DB_NAME_TEST = "test.db"

func TestFormattingListWithTwoComponents(t *testing.T) {
	initial_list := [][]string{
		{"teste teste", "12345"},
	}

	result := formatCompanyData(initial_list)
	expected_result := [][]string{
		{"Teste Teste", "12345"},
	}

	if strings.Compare(result[0][0], expected_result[0][0]) != 0 || strings.Compare(result[0][1], expected_result[0][1]) != 0 {
		t.Log("Error, the result was ", result, " while the expectation was ", expected_result)
		t.Fail()
	}
}

func TestFormattingListWithThreeComponents(t *testing.T) {
	initial_list := [][]string{
		{"teste teste", "12345", "http://WWW.SIte.com"},
	}

	result := formatCompanyData(initial_list)
	expected_result := [][]string{
		{"Teste Teste", "12345", "http://www.site.com"},
	}

	if strings.Compare(result[0][0], expected_result[0][0]) != 0 || strings.Compare(result[0][1], expected_result[0][1]) != 0 || strings.Compare(result[0][2], expected_result[0][2]) != 0 {
		t.Log("Error, the result was ", result, " while the expectation was ", expected_result)
		t.Fail()
	}
}

func TestSearchingWithoutBodyGivesError(t *testing.T) {
	req, err := http.NewRequest("GET", "/companies/search", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(searchCompany)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("HTTP handler returned wrong status code: got %v but expected %v", status, http.StatusBadRequest)
	}
}

func TestSearchingReturns404IfNoResults(t *testing.T) {
	jsonBody := []byte(`{"name":"ZZZZZZ", "zip_code": "99999"}`)
	req, err := http.NewRequest("GET", "/companies/search", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(searchCompany)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("HTTP handler returned wrong status code: got %v but expected %v", status, http.StatusBadRequest)
	}
}

func TestSearchingReturnsDesiredCompanyIfFoundWithStatus200(t *testing.T) {
	comp := Company{Company_name: "Test Company", Zip_Code: "99999", Id: 999}
	db := createConnection()
	db.Create(&comp)
	jsonBody := []byte(`{"name":"Test", "zip_code": "99999"}`)
	req, err := http.NewRequest("GET", "/companies/search", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(searchCompany)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HTTP handler returned wrong status code: got %v but expected %v", status, http.StatusOK)
	}

	expected := `{"id":999,"name":"Test Company","zip":"99999","website":""}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
