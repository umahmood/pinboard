package pinboard

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// DECODE JSON

func TestDecodeJSONValidDecode(t *testing.T) {
	in := []byte(`{"_id": "55885ac73a66c0a52f732f52",
    "index": "42",
    "guid": "e3356187-5496-4409-a0ff-353a499a30ac",
    "balance": "$3,881.62",
    "picture": "http://placehold.it/32x32",
    "eyeColor": "brown"}`)

	wantData := make(map[string]string)
	wantData["_id"] = "55885ac73a66c0a52f732f52"
	wantData["index"] = "42"
	wantData["guid"] = "e3356187-5496-4409-a0ff-353a499a30ac"
	wantData["balance"] = "$3,881.62"
	wantData["picture"] = "http://placehold.it/32x32"
	wantData["eyeColor"] = "brown"

	var wantError error = nil

	gotData, gotError := decodeJSON(in)

	if gotError != wantError {
		t.Errorf("error: got %v want %v", gotError, wantError)
	}

	wantLen := len(wantData)
	gotLen := len(gotData)

	if wantLen != gotLen {
		t.Errorf("data len: got %d want %d", gotLen, wantLen)
	}

	for k, v := range wantData {
		if gotData[k] != v {
			t.Errorf("data: got %v want %v", gotData, wantData)
		}
	}
}

func TestDecodeJSONInvalidDecode(t *testing.T) {
	// the value 42 is an integer not a string.
	in := []byte(`{"_id": "55885ac73a66c0a52f732f52",
    "index": 42, 
    "guid": "e3356187-5496-4409-a0ff-353a499a30ac",
    "balance": "$3,881.62",
    "picture": "http://placehold.it/32x32",
    "eyeColor": "brown"}`)

	wantError := "json: cannot unmarshal number into Go value of type string"

	gotData, gotError := decodeJSON(in)

	if gotError.Error() != wantError {
		t.Errorf("error: got %s want %s", gotError, wantError)
	}

	if gotData != nil {
		t.Errorf("data: got %v want nil", gotData)
	}
}

func TestDecodeJSONIFaceValidDecode(t *testing.T) {
	in := []byte(`{"count":1,
		"notes":[{"id":"19239102",
				"hash":"b378rhef8herh4f",
				"title":"some note",
				"length":"26",
				"created_at":"2015-04-20 13:51:58",
				"updated_at":"2015-04-20 13:51:58"}]}`)

	gotData, gotError := decodeJSONIFace(in)

	if gotError != nil {
		t.Errorf("error: got %s want nil", gotError)
	}

	if gotData == nil {
		t.Errorf("data: got %v want not nil", gotData)
	}
}

func TestDecodeJSONListIFaceValidDecode(t *testing.T) {
	in := []byte(`[{"popular":["news"]},
		{"recommended":["#socent","internet",
		"iPlayer","news","podcast","socent"]}]`)

	gotData, gotError := decodeJSONListIFace(in)

	if gotError != nil {
		t.Errorf("error: got %v want nil", gotError)
	}

	if gotData == nil {
		t.Errorf("data: got %v want not nil", gotData)
	}
}

func TestDecodeJSONInvalidJSON(t *testing.T) {
	in := []byte(`hello`)

	wantError := "invalid character 'h' looking for beginning of value"

	gotData1, gotError := decodeJSON(in)

	if gotError.Error() != wantError {
		t.Errorf("json error: got %s want %s", gotError, wantError)
	}

	if gotData1 != nil {
		t.Errorf("data: got %v want nil", gotData1)
	}

	gotData2, gotError := decodeJSONIFace(in)

	if gotError.Error() != wantError {
		t.Errorf("json iface error: got %s want %s", gotError, wantError)
	}

	if gotData2 != nil {
		t.Errorf("data: got %v want nil", gotData2)
	}

	gotData3, gotError := decodeJSONListIFace(in)

	if gotError.Error() != wantError {
		t.Errorf("json list iface error: got %s want %s", gotError, wantError)
	}

	if gotData3 != nil {
		t.Errorf("data: got %v want nil", gotData3)
	}
}

func TestDecodeJSONWithEmptyJSON(t *testing.T) {
	in := []byte(``)

	wantError := "unexpected end of JSON input"

	gotData1, gotError := decodeJSON(in)

	if gotError.Error() != wantError {
		t.Errorf("json error: got %s want %s", gotError, wantError)
	}

	if gotData1 != nil {
		t.Errorf("data: got %v want nil", gotData1)
	}

	gotData2, gotError := decodeJSONIFace(in)

	if gotError.Error() != wantError {
		t.Errorf("json iface error: got %s want %s", gotError, wantError)
	}

	if gotData2 != nil {
		t.Errorf("data: got %v want nil", gotData2)
	}

	gotData3, gotError := decodeJSONListIFace(in)

	if gotError.Error() != wantError {
		t.Errorf("json list iface error: got %s want %s", gotError, wantError)
	}

	if gotData3 != nil {
		t.Errorf("data: got %v want nil", gotData3)
	}
}

// PARSEDATETIME

func TestParseDateTimeValid(t *testing.T) {
	in := "2010-02-11 03:46:56"

	want := time.Date(2010, time.February, 11, 3, 46, 56, 0, time.UTC)

	got := parseDateTime(in)

	if got != want {
		t.Errorf("parse date time: got %v want %v", got, want)
	}
}

func TestParseDateTimeInvalid(t *testing.T) {
	in := "2010-02-11"

	want := time.Time{}

	got := parseDateTime(in)

	if got.IsZero() != want.IsZero() {
		t.Errorf("parse date time: got %v want %v", got, want)
	}
}

// STRINGTOBOOL

func TestStringToBool(t *testing.T) {
	got := stringToBool("yes")

	if got != true {
		t.Errorf("string to bool: got %v want true", got)
	}

	got = stringToBool("no")

	if got != false {
		t.Errorf("string to bool: got %v want false", got)
	}
}

// BOOLTOSTRING

func TestBoolToString(t *testing.T) {
	got := boolToString(true)

	if got != "yes" {
		t.Errorf("bool to string: got %v want \"yes\"", got)
	}

	got = boolToString(false)

	if got != "no" {
		t.Errorf("bool to string: got %v want \"no\"", got)
	}
}

// ATOI

func TestAtoi(t *testing.T) {
	got := atoi("42")

	if got != 42 {
		t.Errorf("atoi: got %d want 42", got)
	}

	got = atoi("blah")

	if got != 0 {
		t.Errorf("atoi: got %d want 0", got)
	}
}

// DO

func TestDoForResponseData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		fmt.Fprint(w, "{hello:world}")
	}))
	defer ts.Close()

	wantData := []byte(`{hello:world}`)

	var wantError error = nil

	gotData, gotError := do(ts.URL)

	if gotError != wantError {
		t.Errorf("error: got %v want %v", gotError, wantError)
	}

	wantLen := len(wantData)
	gotLen := len(gotData)

	if gotLen != wantLen {
		t.Errorf("data len: got %d want %d", gotLen, wantLen)
	}

	for i := 0; i < wantLen; i++ {
		if gotData[i] != wantData[i] {
			t.Errorf("data: want %v got %v", gotData, wantData)
		}
	}
}

func TestDoWithBadHost(t *testing.T) {
	in := "http://no-such-site-46859755.com"

	wantError := "Get http://no-such-site-46859755.com: dial tcp: lookup " +
		"no-such-site-46859755.com: no such host"

	gotData, gotError := do(in)

	if gotError.Error() != wantError {
		t.Errorf("error: got %s want %s", gotError, wantError)
	}

	if gotData != nil {
		t.Errorf("data: got %v want nil", gotData)
	}
}

func TestDoForServerDown(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	wantError := "HTTP 500 Internal Server Error"

	gotData, gotError := do(ts.URL)

	if gotError.Error() != wantError {
		t.Errorf("error: got %s want %s", gotError, wantError)
	}

	if gotData != nil {
		t.Errorf("data: got %v want nil", gotData)
	}
}

func TestDoFailNilData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		http.Error(w, "", http.StatusOK)
	}))
	defer ts.Close()

	wantError := "No data returned from server."

	gotData, gotError := do(ts.URL)

	if gotData != nil {
		t.Errorf("data: got %v want nil", gotData)
	}

	if gotError.Error() != wantError {
		t.Errorf("error: got %v want %v", gotError, wantError)
	}
}

// TOKEN

func TestToken(t *testing.T) {
	p := New()
	p.token = "mango:49349324"

	want := "mango:49349324"

	got := p.Token()

	if got != want {
		t.Errorf("token: got %s want %s", got, want)
	}
}

// MAKEURL

func TestMakeURLWithAuthToken(t *testing.T) {
	inMethod := "posts/all"
	inVals := url.Values{}
	inVals.Set("auth_token", "hello:1234")

	p := New()

	want := "https://api.pinboard.in/v1/posts/all/?auth_token=hello%3A1234&format=json"

	got := p.makeURL(inMethod, inVals)

	if got != want {
		t.Errorf("make url: got %s want %s", got, want)
	}

}

func TestMakeURLWithNilVals(t *testing.T) {
	inMethod := "posts/all"
	p := New()
	p.token = "mango:1234"

	want := "https://api.pinboard.in/v1/posts/all/?auth_token=mango%3A1234&format=json"

	got := p.makeURL(inMethod, nil)

	if got != want {
		t.Errorf("make url: got %s want %s", got, want)
	}
}
