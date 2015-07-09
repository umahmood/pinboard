package pinboard

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// startTestServer for testing purposes all requests to the Pinboard service are
// routed to here. The test server sends back JSON responses based on the
// incoming URL.
func startTestServer() *httptest.Server {
	s := make(chan bool)
	var ts *httptest.Server
	go func() {
		ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
			r *http.Request) {
			p := r.URL.Path
			switch p {
			case "/user/api_token/":
				fmt.Fprint(w, `{"result":"0123456789"}`)
			case "/posts/update/":
				fmt.Fprint(w, `{"update_time":"2015-07-02T17:03:45Z"}`)
			case "/posts/add/", "/posts/delete/":
				fmt.Fprint(w, `{"result_code": "done"}`)
			case "/posts/get/", "/posts/recent/":
				data := `{"date":"2015-07-02T07:56:40Z",
                "user":"mango",
                "posts":[
                {"href":"http://aaa.com/",
                "description":"AAA",
                "extended":"AAA",
                "meta":"0feee4bcd1ee2724ef8b266c8baaa29c",
                "hash":"d67b75105b042e87b54342de46aca979",
                "time":"2015-07-02T07:56:40Z",
                "shared":"no",
                "toread":"yes",
                "tags":"aazzaa bbzzbb"},
                {"href":"https://bbb.org/",
                "description":"BBB",
                "extended":"BBB",
                "meta":"aca07c3d2676549f454129ffc47bdd3d",
                "hash":"0356c24cf4856410e04a82289e57d79a",
                "time":"2015-07-02T17:03:45Z",
                "shared":"no",
                "toread":"yes",
                "tags":"cczzcc"}]
                }`
				fmt.Fprint(w, data)
			case "/posts/dates/":
				data := `{"user":"mango",
                "tag":"news",
                "dates":{"2015-07-03":"4","2015-07-02":"2","2015-07-01":"1"}}`
				fmt.Fprint(w, data)
			case "/posts/all/":
				data := `[{ "href":"https://foo.com",
                "description":"foo desc",
                "extended":"foo extended",
                "meta":"0ed700a46ec65d5eaea9f95618dd66ad",
                "hash":"4969cc1ba7205f7cd19376ae6930dfcc",
                "time":"2015-07-01T12:37:35Z",
                "shared":"no",
                "toread":"no",
                "tags":"foo_tag"
                },
                {
                "href":"https://bar.com",
                "description":"bar desc",
                "extended":"bar extended",
                "meta":"eb96f3fa1d3f334ee544bc5cbc03c7fe",
                "hash":"c0c2fd08d877b6ea9b4faf3657eeadbd",
                "time":"2015-06-29T08:51:22Z",
                "shared":"no",
                "toread":"no",
                "tags":"bar_tag1 bar_tag2"}]`
				fmt.Fprint(w, data)
			case "/posts/suggest/":
				data := `[{"popular":["security"]},
                {"recommended":["privacy", 
                "security", 
                "internet", 
                "politics",
                "eff",
                "freedom",
                "rights",
                "Technology",
                "Copyright",
                "opensource"]}]`
				fmt.Fprint(w, data)
			case "/tags/get/":
				fmt.Fprint(w, `{"foo":"27","bar":"2","zip":"1","zap":"37"}`)
			case "/tags/delete/", "/tags/rename/":
				fmt.Fprint(w, `{"result":"done"}`)
			case "/notes/list/":
				data := `{"count":2,
                "notes":[{"id":"364bd4c30a2b9654d0e1",
                "hash":"b7c4cd9c55bb946c0216",
                "title":"shopping list",
                "length":"26",
                "created_at":"2015-04-20 13:51:58",
                "updated_at":"2015-04-20 13:51:58"},
                {"id":"364bd9876a2b96543453",
                "hash":"b7c4876y55bb946c0216",
                "title":"coding notes",
                "length":"42",
                "created_at":"2015-04-20 13:51:58",
                "updated_at":"2015-04-20 13:51:58"}]}`
				fmt.Fprint(w, data)
			case "/notes/1234/":
				data := `{"id":"364bd4c30a2b9654d0e1",
                "hash":"b7c4cd9c55bb946c0216",
                "title":"shopping list",
                "length":26,
                "text":"some note text.",
                "created_at":"2015-04-20 13:51:58",
                "updated_at":"2015-04-20 13:51:58"}`
				fmt.Fprint(w, data)
			}
		}))
		baseURL = ts.URL + "/%s/?%s"
		s <- true
	}()
	_ = <-s
	return ts
}

// compareBookmarks helper method for tests which compares each field and
// returns an error if there is a difference.
func compareBookmarks(a, b Bookmark) error {
	f := "%s: got %v want %v"
	if a.URL != b.URL {
		return fmt.Errorf(f, "url", a.URL, b.URL)
	}
	if a.Title != b.Title {
		return fmt.Errorf(f, "title", a.Title, b.Title)
	}
	if a.Desc != b.Desc {
		return fmt.Errorf(f, "desc", a.Desc, b.Desc)
	}
	if len(a.Tags) != len(b.Tags) {
		return fmt.Errorf(f, "tags", len(a.Tags), len(b.Tags))
	}
	if a.Created != b.Created {
		return fmt.Errorf(f, "created", a.Created, b.Created)
	}
	if string(a.Meta) != string(b.Meta) {
		return fmt.Errorf(f, "meta", a.Meta, b.Meta)
	}
	if string(a.Hash) != string(b.Hash) {
		return fmt.Errorf(f, "hash", a.Hash, b.Hash)
	}
	if a.Shared != b.Shared {
		return fmt.Errorf(f, "shared", a.Shared, b.Shared)
	}
	if a.ToRead != b.ToRead {
		return fmt.Errorf(f, "to read", a.ToRead, b.ToRead)
	}
	return nil
}

func TestAuthSuccess(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	in := "mango:0123456789"
	want := "0123456789"

	p := New()
	got, err := p.Auth(in)

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	if got != want {
		t.Errorf("auth: got %s want %s", got, want)
	}
}

func TestAuthFail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer ts.Close()
	baseURL = ts.URL + "/%s/?%s"
	in := "bad:token"
	wantError := "HTTP 401 Unauthorized"
	p := New()
	got, err := p.Auth(in)
	if err.Error() != wantError {
		t.Errorf("error: got %v want %s", err, wantError)
	}
	if got != "" {
		t.Errorf("auth: got %v want \"\" (empty string)", got)
	}
	if p.token != "" {
		t.Errorf("auth: got %s want \"\" (empty string)", got)
	}
	if p.authed != false {
		t.Errorf("auth: got %t want false", p.authed)
	}
}

func TestLastUpdate(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	want := time.Date(2015, 7, 2, 17, 3, 45, 0, time.UTC)
	p := New()
	_, err := p.Auth("mango:0123456789")
	if err != nil {
		t.Errorf("error: got %s want nil", err)
	}
	got, err := p.LastUpdate()
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	if got != want {
		t.Errorf("time: got %v want %v", got, want)
	}
}

func TestAdd(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()

	_, err := p.Auth("mango:0123456789")

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	tg := []string{"zza", "zzb", "zzc"}

	b := Bookmark{
		URL:     "http://www.food.com/",
		Title:   "Foody",
		Desc:    "some food website",
		Tags:    tg,
		Created: time.Now(),
		Replace: true,
		Shared:  true,
		ToRead:  true,
	}

	want := true

	got, err := p.Add(b)

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	if got != want {
		t.Errorf("add: got %v want %v", got, want)
	}
}

func TestDel(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()
	_, err := p.Auth("mango:0123456789")
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}
	want := true
	got, err := p.Del("aazzbbzz")
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}
	if got != want {
		t.Errorf("error: got %v want %v", got, want)
	}
}

func TestGet(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()
	_, err := p.Auth("mango:0123456789")
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	want := make([]Bookmark, 0)

	tg := []string{"aazzaa", "bbzzbb"}

	b := Bookmark{
		URL:     "http://aaa.com/",
		Title:   "AAA",
		Desc:    "AAA",
		Tags:    tg,
		Created: time.Date(2015, 7, 2, 7, 56, 40, 0, time.UTC),
		Meta:    []byte("0feee4bcd1ee2724ef8b266c8baaa29c"),
		Hash:    []byte("d67b75105b042e87b54342de46aca979"),
		Shared:  false,
		ToRead:  true,
	}

	want = append(want, b)

	tg = []string{"aazzaa"}

	b = Bookmark{
		URL:     "https://bbb.org/",
		Title:   "BBB",
		Desc:    "BBB",
		Tags:    tg,
		Created: time.Date(2015, 7, 2, 17, 3, 45, 0, time.UTC),
		Meta:    []byte("aca07c3d2676549f454129ffc47bdd3d"),
		Hash:    []byte("0356c24cf4856410e04a82289e57d79a"),
		Shared:  false,
		ToRead:  true,
	}

	want = append(want, b)

	got, err := p.Get(time.Time{}, "", nil, false)

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	if len(got) != 2 {
		t.Errorf("len: got %v want 2", got)
	}

	err = compareBookmarks(got[0], want[0])
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = compareBookmarks(got[1], want[1])
	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func TestDates(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()
	_, err := p.Auth("mango:0123456789")
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	want := make([]Post, 0)
	want = append(want, Post{Date: time.Date(2015, 7, 3, 0, 0, 0, 0, time.UTC), Count: 4})
	want = append(want, Post{Date: time.Date(2015, 7, 2, 0, 0, 0, 0, time.UTC), Count: 2})
	want = append(want, Post{Date: time.Date(2015, 7, 1, 0, 0, 0, 0, time.UTC), Count: 1})

	tg := []string{"news"}

	got, err := p.Dates(tg)

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	checkPost := func(o Post) bool {
		for _, s := range want {
			if s.Date == o.Date && s.Count == o.Count {
				return true
			}
		}
		return false
	}

	for _, s := range got {
		if !checkPost(s) {
			t.Errorf("dates: got %v want %v", got, want)
		}
	}
}

func TestRecent(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()
	_, err := p.Auth("mango:0123456789")
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	want := make([]Bookmark, 0)

	tg := []string{"aazzaa", "bbzzbb"}

	b := Bookmark{
		URL:     "http://aaa.com/",
		Title:   "AAA",
		Desc:    "AAA",
		Tags:    tg,
		Created: time.Date(2015, 7, 2, 7, 56, 40, 0, time.UTC),
		Meta:    []byte("0feee4bcd1ee2724ef8b266c8baaa29c"),
		Hash:    []byte("d67b75105b042e87b54342de46aca979"),
		Shared:  false,
		ToRead:  true,
	}

	want = append(want, b)
	tg = []string{"aazzaa"}

	b = Bookmark{
		URL:     "https://bbb.org/",
		Title:   "BBB",
		Desc:    "BBB",
		Tags:    tg,
		Created: time.Date(2015, 7, 2, 17, 3, 45, 0, time.UTC),
		Meta:    []byte("aca07c3d2676549f454129ffc47bdd3d"),
		Hash:    []byte("0356c24cf4856410e04a82289e57d79a"),
		Shared:  false,
		ToRead:  true,
	}
	want = append(want, b)

	got, err := p.Recent(nil, 2)

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	if len(got) != 2 {
		t.Errorf("len: got %v want 2", got)
	}

	err = compareBookmarks(got[0], want[0])
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = compareBookmarks(got[1], want[1])
	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func TestBookmarks(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()
	_, err := p.Auth("mango:0123456789")
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	want := make([]Bookmark, 0)
	tg := []string{"foo_tag"}

	b := Bookmark{
		URL:     "https://foo.com",
		Title:   "foo desc",
		Desc:    "foo extended",
		Tags:    tg,
		Created: time.Date(2015, 7, 1, 12, 37, 35, 0, time.UTC),
		Meta:    []byte("0ed700a46ec65d5eaea9f95618dd66ad"),
		Hash:    []byte("4969cc1ba7205f7cd19376ae6930dfcc"),
		Shared:  false,
		ToRead:  false,
	}

	want = append(want, b)
	tg = []string{"bar_tag1", "bar_tag2"}

	b = Bookmark{
		URL:     "https://bar.com",
		Title:   "bar desc",
		Desc:    "bar extended",
		Tags:    tg,
		Created: time.Date(2015, 6, 29, 8, 51, 22, 0, time.UTC),
		Meta:    []byte("eb96f3fa1d3f334ee544bc5cbc03c7fe"),
		Hash:    []byte("c0c2fd08d877b6ea9b4faf3657eeadbd"),
		Shared:  false,
		ToRead:  false,
	}
	want = append(want, b)

	got, err := p.Bookmarks(tg, 0, 2, time.Time{}, time.Time{}, false)

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	if len(got) != 2 {
		t.Errorf("len: got %v want 2", got)
	}

	err = compareBookmarks(got[0], want[0])
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = compareBookmarks(got[1], want[1])
	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func TestSuggest(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()
	_, err := p.Auth("mango:0123456789")
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	wantPop := Popular{"security"}

	wantRec := Recommended{"privacy", "security", "internet", "politics", "eff",
		"freedom", "rights", "Technology", "Copyright", "opensource"}

	pop, rec, err := p.Suggest("http://www.eef.org")

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}
	if len(pop) != len(wantPop) {
		t.Errorf("popular len: got %v want %v", len(pop), len(wantPop))
	}
	if len(rec) != len(wantRec) {
		t.Errorf("rec len: got %v want %v", len(rec), len(wantRec))
	}
	for i := 0; i < len(pop); i++ {
		if pop[i] != wantPop[i] {
			t.Errorf("popular: got %v want %v", pop[i], wantPop[i])
		}
	}
	for i := 0; i < len(rec); i++ {
		if rec[i] != wantRec[i] {
			t.Errorf("rec: got %v want %v", rec[i], wantRec[i])
		}
	}
}

func TestTags(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()
	_, err := p.Auth("mango:0123456789")
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}
	got, err := p.Tags()
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}
	if len(got) != 4 {
		t.Errorf("len tags: got %v want 4", got)
	}
}

func TestDelTag(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()
	_, err := p.Auth("mango:0123456789")

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	want := true

	got, err := p.DelTag("zap")

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	if got != want {
		t.Errorf("del tag: got %v want %v", got, want)
	}
}

func TestRenTag(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()
	_, err := p.Auth("mango:0123456789")

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	want := true

	got, err := p.RenTag("foo", "bar")

	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}

	if got != want {
		t.Errorf("ren tag: got %v want %v", got, want)
	}
}

func TestNotes(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()
	_, err := p.Auth("mango:0123456789")
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}
	got, err := p.Notes()
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}
	if len(got) != 2 {
		t.Errorf("len: got %v want 2", got)
	}
}

func TestNoteID(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	p := New()
	_, err := p.Auth("mango:0123456789")
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}
	got, err := p.NoteID("1234")
	if err != nil {
		t.Errorf("error: got %v want nil", err)
	}
	if got.Text != "some note text." {
		t.Errorf("id: got %v want \"some note text.\"", got.Text)
	}
}
