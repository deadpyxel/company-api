package main

import (
	"strings"
	"testing"
)

const DB_NAME_TEST = "test.db"

func TestFormattingListWithTwoComponents(t *testing.T) {
	initial_list := [][]string{
		{"teste teste", "12345"},
	}

	result := format_company_data(initial_list)
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

	result := format_company_data(initial_list)
	expected_result := [][]string{
		{"Teste Teste", "12345", "http://www.site.com"},
	}

	if strings.Compare(result[0][0], expected_result[0][0]) != 0 || strings.Compare(result[0][1], expected_result[0][1]) != 0 || strings.Compare(result[0][2], expected_result[0][2]) != 0 {
		t.Log("Error, the result was ", result, " while the expectation was ", expected_result)
		t.Fail()
	}
}
