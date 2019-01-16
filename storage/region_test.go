package storage

import (
	"testing"
)

func TestEndpoint(t *testing.T) {
	type input struct {
		UseHttps bool
		Host     string
	}
	testInputs := []input{
		{UseHttps: true, Host: "rs.qiniu.com"},
		{UseHttps: false, Host: "rs.qiniu.com"},
		{UseHttps: true, Host: ""},
		{UseHttps: false, Host: ""},
		{UseHttps: true, Host: "https://rs.qiniu.com"},
		{UseHttps: false, Host: "https://rs.qiniu.com"},
		{UseHttps: false, Host: "http://rs.qiniu.com"},
	}
	testWants := []string{"https://rs.qiniu.com", "http://rs.qiniu.com", "", "", "https://rs.qiniu.com",
		"http://rs.qiniu.com", "http://rs.qiniu.com"}

	for ind, testInput := range testInputs {
		testGot := endpoint(testInput.UseHttps, testInput.Host)
		testWant := testWants[ind]
		if testGot != testWant {
			t.Fail()
		}
	}
}
