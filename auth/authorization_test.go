package auth

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

var at *Credentials

func init() {
	at = New("ak", "sk")
}

func TestNew(t *testing.T) {
	if at.AccessKey != "ak" || string(at.SecretKey) != "sk" {
		t.Fail()
	}
}

func TestAuthSign(t *testing.T) {
	testStrs := []struct {
		Data   string
		Signed string
	}{
		{Data: "hello", Signed: "ak:NDN8cM0rwosxhHJ6QAcI7ialr0g="},
		{Data: "world", Signed: "ak:wZ-sw41ayh3PFDmQA-D3o7eBJIY="},
		{Data: "-test", Signed: "ak:oJ59sZasiWSqL1o7ugZs5OInEK4="},
		{Data: "ba#a-", Signed: "ak:tqHL8V2BbNI0dVDXsvueZp_2QnI="},
	}

	for _, b := range testStrs {
		got := at.Sign([]byte(b.Data))
		if got != b.Signed {
			t.Errorf("Sign %q, want=%q, got=%q\n", b.Data, b.Signed, got)
		}
	}

}

func TestAuthSignWithData(t *testing.T) {
	testStrs := []struct {
		Data   string
		Signed string
	}{
		{Data: "hello", Signed: "ak:2pn0qs-2kfEsQFuHI2pAYlo0hpc=:aGVsbG8="},
		{Data: "world", Signed: "ak:vzqcP6VeODVu_youBJnyr_nefT4=:d29ybGQ="},
		{Data: "-test", Signed: "ak:uV60zWZgj-Jbrg9VHc06Nok64Bw=:LXRlc3Q="},
		{Data: "ba#a-", Signed: "ak:RLvTUx_kizrrbpSrinkdxC4jCy8=:YmEjYS0="},
	}
	for _, b := range testStrs {
		got := at.SignWithData([]byte(b.Data))
		if got != b.Signed {
			t.Errorf("SignWithData %q, want=%q, got=%q\n", b.Data, b.Signed, got)
		}
	}
}

func TestCollectData(t *testing.T) {
	inputs := []ReqParams{
		{Method: "", Url: "", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "GET", Url: "", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "POST", Url: "", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}, Body: strings.NewReader(`name=test&language=go`)},
		{Method: "", Url: "http://upload.qiniup.com?v=2", Headers: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}, Body: strings.NewReader(`name=test&language=go`)},
		{Method: "", Url: "http://upload.qiniup.com/find/sdk?v=2", Headers: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}, Body: strings.NewReader(`name=test&language=go`)},
	}
	wants := []string{"\n", "\n", "\n", "\n", "\n", "\n", "\nname=test&language=go", "?v=2\nname=test&language=go", "/find/sdk?v=2\nname=test&language=go"}
	reqs, gErr := genRequests(inputs)
	if gErr != nil {
		t.Errorf("generate requests: %v\n", gErr)
	}

	for ind, req := range reqs {
		data, err := collectData(req)
		if err != nil {
			t.Error("collectData: ", err)
		}
		if string(data) != wants[ind] {
			t.Errorf("collectData, want = %q, got = %q\n", wants[ind], data)
		}
	}

}

func TestCollectDataV2(t *testing.T) {
	inputs := []ReqParams{
		{Method: "", Url: "", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "GET", Url: "", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "POST", Url: "", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}, Body: strings.NewReader(`name=test&language=go`)},
		{Method: "", Url: "http://upload.qiniup.com?v=2", Headers: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}, Body: strings.NewReader(`name=test&language=go`)},
		{Method: "", Url: "http://upload.qiniup.com/find/sdk?v=2", Headers: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}, Body: strings.NewReader(`name=test&language=go`)},
	}

	wants := []string{
		"GET \nHost: \nContent-Type: application/json\n\n{\"name\": \"test\"}",
		"GET \nHost: \n\n",
		"GET \nHost: \nContent-Type: application/json\n\n{\"name\": \"test\"}",
		"GET \nHost: \n\n",
		"POST \nHost: \nContent-Type: application/json\n\n{\"name\": \"test\"}",
		"GET \nHost: upload.qiniup.com\n\n",
		"GET \nHost: upload.qiniup.com\nContent-Type: application/json\n\n{\"name\": \"test\"}",
		"GET \nHost: upload.qiniup.com\nContent-Type: application/x-www-form-urlencoded\n\nname=test&language=go",
		"GET ?v=2\nHost: upload.qiniup.com\nContent-Type: application/x-www-form-urlencoded\n\nname=test&language=go",
		"GET /find/sdk?v=2\nHost: upload.qiniup.com\nContent-Type: application/x-www-form-urlencoded\n\nname=test&language=go",
	}
	reqs, gErr := genRequests(inputs)
	if gErr != nil {
		t.Errorf("generate requests: %v\n", gErr)
	}

	for ind, req := range reqs {
		data, err := collectDataV2(req)
		if err != nil {
			t.Error("collectDataV2: ", err)
		}
		if string(data) != wants[ind] {
			t.Errorf("collectDataV2, want = %q, got = %q\n", wants[ind], data)
		}
	}
}

func genRequest(method, url string, headers http.Header, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	req.Header = headers
	return
}

type ReqParams struct {
	Method  string
	Url     string
	Headers http.Header
	Body    io.Reader
}

func genRequests(params []ReqParams) (reqs []*http.Request, err error) {
	for _, reqParam := range params {
		req, rErr := genRequest(reqParam.Method, reqParam.Url, reqParam.Headers, reqParam.Body)
		if rErr != nil {
			err = rErr
			return
		}
		reqs = append(reqs, req)
	}
	return
}

func TestAuthSignRequest(t *testing.T) {
	inputs := []ReqParams{
		{Method: "", Url: "", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "GET", Url: "", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "POST", Url: "", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}, Body: strings.NewReader(`name=test&language=go`)},
	}
	wants := []string{
		"ak:qfWnqF1E_vfzjZnReCVkcSMl29M=",
		"ak:qfWnqF1E_vfzjZnReCVkcSMl29M=",
		"ak:qfWnqF1E_vfzjZnReCVkcSMl29M=",
		"ak:qfWnqF1E_vfzjZnReCVkcSMl29M=",
		"ak:qfWnqF1E_vfzjZnReCVkcSMl29M=",
		"ak:qfWnqF1E_vfzjZnReCVkcSMl29M=",
		"ak:h8gBb1Adb2Jgoys1N8sRVAnNvpw=",
	}
	reqs, gErr := genRequests(inputs)
	if gErr != nil {
		t.Errorf("generate requests: %v\n", gErr)
	}
	for ind, req := range reqs {
		got, sErr := at.SignRequest(req)
		if sErr != nil {
			t.Errorf("SignRequest: %v\n", sErr)
		}
		if got != wants[ind] {
			t.Errorf("SignRequest, want = %q, got = %q\n", wants[ind], got)
		}
	}
}

func TestAuthSignRequestV2(t *testing.T) {
	inputs := []ReqParams{
		{Method: "", Url: "", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "GET", Url: "", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "POST", Url: "", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: nil, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: http.Header{"Content-Type": []string{"application/json"}}, Body: strings.NewReader(`{"name": "test"}`)},
		{Method: "", Url: "http://upload.qiniup.com", Headers: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}, Body: strings.NewReader(`name=test&language=go`)},
	}
	wants := []string{
		"ak:XNay-AIghhXfytRKsKNj0DQqV2E=",
		"ak:K1DI0goT05yhGizDFE5FiPJxAj4=",
		"ak:XNay-AIghhXfytRKsKNj0DQqV2E=",
		"ak:0ujEjW_vLRZxebsveBgqa3JyQ-w=",
		"ak:Eadl-_gKUNECGo3mcikiTBoNfqI=",
		"ak:Pkuq20x3HNWJlHDbRLW1kDYmXF0=",
		"ak:rZjOJKtlePVSegqoSO8p6Gpsr64=",
	}
	reqs, gErr := genRequests(inputs)
	if gErr != nil {
		t.Errorf("generate requests: %v\n", gErr)
	}
	for ind, req := range reqs {
		got, sErr := at.SignRequestV2(req)
		if sErr != nil {
			t.Errorf("SignRequest: %v\n", sErr)
		}
		if got != wants[ind] {
			t.Errorf("SignRequest, want = %q, got = %q\n", wants[ind], got)
		}
	}
}
