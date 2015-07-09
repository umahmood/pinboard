package pinboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var baseURL string = "https://api.pinboard.in/v1/%s/?%s"

type Pinboard struct {
	token  string // e.g. username:TOKEN
	authed bool   //  Authenticated with Pinboard service?
}

type Bookmark struct {
	URL   string
	Title string
	Desc  string
	Tags  []string

	// Creation time for this bookmark.
	Created time.Time
	// Replace any existing bookmark with this URL.
	Replace bool
	// Make bookmark public.
	Shared bool
	// Marks the bookmarks as unread.
	ToRead bool

	Hash []byte // 32 character hexadecimal MD5 hash.
	Meta []byte // 32 character hexadecimal MD5 hash.
}

type Tag struct {
	Name  string
	Count int
}

type Post struct {
	Date  time.Time
	Count int
}

type NoteMetadata struct {
	ID      string
	Title   string
	Length  int       // Size of the note in bytes.
	Hash    []byte    // 20 character long sha1 hash of note text.
	Created time.Time // Time note was created. UTC time.
	Updated time.Time // Time note was last updated. UTC time.
}

type Note struct {
	NoteMetadata
	Text string
}

type Popular []string
type Recommended []string

// decodeJSON decodes a JSON structure into map key:string val:string. Decodes
// response data from Auth, LastUpdate, Add, Del, Tags, DelTag, RenTag.
func decodeJSON(jsonBlob []byte) (map[string]string, error) {
	var j map[string]string
	err := json.Unmarshal(jsonBlob, &j)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// decodeJSONIFace decodes a JSON structure into map key:string val:interface{}.
// Decodes response data from Get, Dates, Recent, Notes, NoteID.
func decodeJSONIFace(jsonBlob []byte) (map[string]interface{}, error) {
	var j map[string]interface{}
	err := json.Unmarshal(jsonBlob, &j)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// decodeJSONListIFace decodes a JSON structure into []interface{}. Decodes
// response data from Bookmarks, Suggest.
func decodeJSONListIFace(jsonBlob []byte) ([]interface{}, error) {
	var j []interface{}
	err := json.Unmarshal(jsonBlob, &j)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// parseDateTime parses a date time in the format "2010-02-11 03:46:56" and
// returns a time.Time with location UTC.
func parseDateTime(dt string) time.Time {
	if len(dt) != 19 {
		return time.Time{}
	}
	yr := atoi(dt[0:4])
	mo := atoi(dt[6:7])
	dy := atoi(dt[8:10])
	hr := atoi(dt[11:13])
	mn := atoi(dt[14:16])
	sd := atoi(dt[17:19])
	return time.Date(yr, time.Month(mo), dy, hr, mn, sd, 0, time.UTC)
}

// stringToBool returns true if s == "yes", otherwise false.
func stringToBool(s string) bool {
	if s == "yes" {
		return true
	}
	return false
}

// boolToString return "yes" if b == true, otherwise "no".
func boolToString(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// atoi converts a string to an integer, if there is an error in the conversion
// returns 0.
func atoi(s string) int {
	l, err := strconv.Atoi(s)
	if err != nil {
		l = 0
	}
	return l
}

// do performs a HTTP GET on a URL.
func do(url string) (data []byte, err error) {
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	c := rsp.StatusCode
	if c != http.StatusOK {
		return nil, errors.New("HTTP " + strconv.Itoa(c) + " " + http.StatusText(c))
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	// server sent back no data or just returned '\n'.
	if body == nil || (len(body) == 1 && body[0] == 10) {
		return nil, errors.New("No data returned from server.")
	}
	return body, nil
}

// New returns a new unauthorized instance of Pinboard.
func New() *Pinboard {
	return &Pinboard{}
}

// Token returns the users token in the format username:TOKEN.
func (p Pinboard) Token() string {
	return p.token
}

// makeURL constructs a valid URL which can be used to make a request to the
// Pinboard service.
func (p Pinboard) makeURL(method string, vals url.Values) string {
	if vals == nil {
		vals = url.Values{}
	}
	vals.Set("format", "json")

	a := vals.Get("auth_token")
	if a == "" {
		vals.Set("auth_token", p.token)
	}

	return fmt.Sprintf(baseURL, method, vals.Encode())
}

// performRequest performs a request to the Pinboard service.
func (p Pinboard) performRequest(method string, vals url.Values) ([]byte, error) {
	// Only calls from Auth(...) can pass when p.auth == false, as we are
	// requesting authentication. For all other calls the API must already be
	// authorized by a previous call to Auth(...).
	if !p.authed && method != "user/api_token" {
		return nil, errors.New("API not authorized.")
	}
	url := p.makeURL(method, vals)
	data, err := do(url)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// USER

// Auth validates the provided 'token' with the Pinboard service and returns the
// user's API token, The API token must be in the format: username:TOKEN.
func (p *Pinboard) Auth(token string) (string, error) {
	v := url.Values{}
	v.Set("auth_token", token)
	data, err := p.performRequest("user/api_token", v)
	if err != nil {
		return "", err
	}
	j, err := decodeJSON(data)
	if err != nil {
		return "", err
	}
	p.token = token
	p.authed = true
	return j["result"], nil
}

// POSTS

// LastUpdate returns the last time a bookmark was added, updated or deleted.
// Use this before calling Bookmarks() to see if the data has changed since the
// last fetch.
func (p Pinboard) LastUpdate() (time.Time, error) {
	data, err := p.performRequest("posts/update", nil)
	if err != nil {
		return time.Time{}, err
	}
	j, err := decodeJSON(data)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.Parse(time.RFC3339, j["update_time"])
	if err != nil {
		return time.Time{}, err
	}
	return t, err
}

// Add adds a bookmark.
func (p *Pinboard) Add(b Bookmark) (bool, error) {
	v := url.Values{}
	// mandatory.
	v.Set("url", b.URL)
	v.Set("description", b.Title)
	// optional
	if b.Desc != "" {
		v.Set("extended", b.Desc)
	}
	if b.Tags != nil {
		var g string
		for _, t := range b.Tags {
			g += t + ","
		}
		v.Set("tags", strings.Trim(g, ","))
	}

	var dt string
	if b.Created.IsZero() {
		// no creation time provided.
		dt = time.Now().UTC().Format(time.RFC3339)
	} else {
		dt = b.Created.UTC().Format(time.RFC3339)
	}
	v.Set("dt", dt)
	v.Set("replace", boolToString(b.Replace))
	v.Set("shared", boolToString(b.Shared))
	v.Set("toread", boolToString(b.ToRead))

	data, err := p.performRequest("posts/add", v)
	if err != nil {
		return false, err
	}
	j, err := decodeJSON(data)
	if err != nil {
		return false, nil
	}
	if j["result_code"] != "done" {
		return false, errors.New(j["result_code"])
	}
	return true, nil
}

// Del deletes a bookmark.
func (p Pinboard) Del(URL string) (bool, error) {
	v := url.Values{}
	v.Set("url", URL)
	data, err := p.performRequest("posts/delete", v)
	if err != nil {
		return false, err
	}
	j, err := decodeJSON(data)
	if err != nil {
		return false, err
	}
	if j["result_code"] != "done" {
		return false, errors.New(j["result_code"])
	}
	return true, nil
}

// Get returns one or more posts on a single day matching the arguments. If no
// date or URL is given, date of most recent bookmark will be used. The meta
// flag allows the ability to include a change detection signature for each
// bookmark.
func (p Pinboard) Get(dt time.Time, URL string, tags []string, meta bool) ([]Bookmark, error) {
	v := url.Values{}
	if !dt.IsZero() {
		v.Set("dt", dt.UTC().Format(time.RFC3339))
	}
	if URL != "" {
		v.Set("url", URL)
	}
	if tags != nil {
		var g string
		for _, t := range tags {
			g += t + ","
		}
		v.Set("tag", strings.Trim(g, ","))
	}
	if meta {
		v.Set("meta", "yes")
	} else {
		v.Set("meta", "no")
	}
	data, err := p.performRequest("posts/get", v)
	if err != nil {
		return nil, err
	}
	j, err := decodeJSONIFace(data)
	if err != nil {
		return nil, err
	}
	bmarks := make([]Bookmark, 0)
	posts := j["posts"].([]interface{})
	for i := 0; i < len(posts); i++ {
		p := posts[i].(map[string]interface{})

		t, err := time.Parse(time.RFC3339, p["time"].(string))
		if err != nil {
			t = time.Time{}
		}
		b := Bookmark{
			URL:     p["href"].(string),
			Title:   p["description"].(string),
			Desc:    p["extended"].(string),
			Tags:    strings.Split(p["tags"].(string), " "),
			Created: t,
			Shared:  stringToBool(p["shared"].(string)),
			ToRead:  stringToBool(p["toread"].(string)),
			Hash:    []byte(p["hash"].(string)),
			Meta:    []byte(p["meta"].(string)),
		}
		bmarks = append(bmarks, b)
	}
	return bmarks, nil
}

// Dates returns a list of dates with the number of posts at each date.
func (p Pinboard) Dates(tags []string) ([]Post, error) {
	v := url.Values{}
	if tags != nil {
		var g string
		for _, t := range tags {
			g += t + ","
		}
		v.Set("tag", strings.Trim(g, ","))
	}
	data, err := p.performRequest("posts/dates", v)
	if err != nil {
		return nil, err
	}
	j, err := decodeJSONIFace(data)
	if err != nil {
		return nil, err
	}
	d := j["dates"].(map[string]interface{})
	s := make([]Post, 0)
	for k, v := range d {
		w, err := time.Parse("2006-01-02", k)
		if err != nil {
			return nil, err
		}
		s = append(s, Post{Date: w, Count: atoi(v.(string))})
	}
	return s, nil
}

// Recent returns a list of the user's most recent posts, filtered by tag. count
// indicates the number results to return, default is 15 max is 100.
func (p Pinboard) Recent(tags []string, count int) ([]Bookmark, error) {
	v := url.Values{}
	if tags != nil {
		var g string
		for _, t := range tags {
			g += t + ","
		}
		v.Set("tag", strings.Trim(g, ","))
	}
	v.Set("count", strconv.Itoa(count))
	data, err := p.performRequest("posts/recent", v)
	if err != nil {
		return nil, err
	}
	j, err := decodeJSONIFace(data)
	if err != nil {
		return nil, err
	}
	bmarks := make([]Bookmark, 0)
	posts := j["posts"].([]interface{})
	for i := 0; i < len(posts); i++ {
		p := posts[i].(map[string]interface{})

		t, err := time.Parse(time.RFC3339, p["time"].(string))
		if err != nil {
			t = time.Time{}
		}
		b := Bookmark{
			URL:     p["href"].(string),
			Title:   p["description"].(string),
			Desc:    p["extended"].(string),
			Tags:    strings.Split(p["tags"].(string), " "),
			Created: t,
			Shared:  stringToBool(p["shared"].(string)),
			ToRead:  stringToBool(p["toread"].(string)),
			Hash:    []byte(p["hash"].(string)),
			Meta:    []byte(p["meta"].(string)),
		}
		bmarks = append(bmarks, b)
	}
	return bmarks, nil
}

// Bookmarks returns all bookmarks in the user's account. Provides the ability
// to:
// - 'tags' Filter by tags.
// - 'offset' Set offset value.
// - 'count' Set the number of results to return. Default is all.
// - 'start' Return only bookmarks created after this time.
// - 'end' Return only bookmarks created before this time.
// - 'meta' A meta flag to include a change detection signature for each bookmark.
func (p Pinboard) Bookmarks(tags []string, offset int, count int,
	start time.Time, end time.Time, meta bool) ([]Bookmark, error) {

	v := url.Values{}
	if tags != nil {
		var g string
		for _, t := range tags {
			g += t + ","
		}
		v.Set("tag", strings.Trim(g, ","))
	}
	v.Set("start", strconv.Itoa(offset))
	v.Set("results", strconv.Itoa(count))

	if !start.IsZero() {
		v.Set("fromdt", start.UTC().Format(time.RFC3339))
	}
	if !end.IsZero() {
		v.Set("todt", end.UTC().Format(time.RFC3339))
	}
	if meta {
		v.Set("meta", "yes")
	} else {
		v.Set("meta", "no")
	}
	data, err := p.performRequest("posts/all", v)
	if err != nil {
		return nil, err
	}
	j, err := decodeJSONListIFace(data)
	if err != nil {
		return nil, err
	}
	bmarks := make([]Bookmark, 0)
	for i := 0; i < len(j); i++ {
		p := j[i].(map[string]interface{})

		t, err := time.Parse(time.RFC3339, p["time"].(string))
		if err != nil {
			t = time.Time{}
		}
		b := Bookmark{
			URL:     p["href"].(string),
			Title:   p["description"].(string),
			Desc:    p["extended"].(string),
			Tags:    strings.Split(p["tags"].(string), " "),
			Created: t,
			Shared:  stringToBool(p["shared"].(string)),
			ToRead:  stringToBool(p["toread"].(string)),
			Hash:    []byte(p["hash"].(string)),
			Meta:    []byte(p["meta"].(string)),
		}
		bmarks = append(bmarks, b)
	}
	return bmarks, nil
}

// Suggest returns a list of popular tags and recommended tags for a given URL.
// Popular tags are tags used site-wide for the url; recommended tags are drawn
// from the user's own tags.
func (p Pinboard) Suggest(URL string) (Popular, Recommended, error) {
	v := url.Values{}
	v.Set("url", URL)
	data, err := p.performRequest("posts/suggest", v)
	if err != nil {
		return nil, nil, err
	}
	j, err := decodeJSONListIFace(data)
	if err != nil {
		return nil, nil, err
	}
	pop := make(Popular, 0)
	rec := make(Recommended, 0)
	if len(j) >= 2 {
		a := j[0].(map[string]interface{})
		b := a["popular"].([]interface{})
		for _, t := range b {
			pop = append(pop, t.(string))
		}
		a = j[1].(map[string]interface{})
		b = a["recommended"].([]interface{})
		for _, t := range b {
			rec = append(rec, t.(string))
		}
	}
	return pop, rec, nil
}

// TAGS

// Tags returns a full list of the user's tags along with the number of times
// they were used.
func (p Pinboard) Tags() ([]Tag, error) {
	data, err := p.performRequest("tags/get", nil)
	if err != nil {
		return nil, err
	}
	j, err := decodeJSON(data)
	if err != nil {
		return nil, err
	}
	t := make([]Tag, 0)
	for k, v := range j {
		t = append(t, Tag{Name: k, Count: atoi(v)})
	}
	return t, nil
}

// DelTag delete an existing tag, returns true if the delete operation
// succeeds.
func (p Pinboard) DelTag(tag string) (bool, error) {
	v := url.Values{}
	v.Set("tag", tag)
	data, err := p.performRequest("tags/delete", v)
	if err != nil {
		return false, err
	}
	j, err := decodeJSON(data)
	if err != nil {
		return false, err
	}
	if j["result"] != "done" {
		return false, errors.New(j["result"])
	}
	return true, nil
}

// RenTag rename a tag, or fold it in to an existing tag. Match is not case
// sensitive, returns true if the rename operation succeeds.
func (p Pinboard) RenTag(oldTag, newTag string) (bool, error) {
	v := url.Values{}
	v.Set("old", oldTag)
	v.Set("new", newTag)
	data, err := p.performRequest("tags/rename", v)
	if err != nil {
		return false, err
	}
	j, err := decodeJSON(data)
	if err != nil {
		return false, err
	}
	if j["result"] != "done" {
		return false, errors.New(j["result"])
	}
	return true, nil
}

// NOTES

// Notes returns a list of the user's notes
func (p Pinboard) Notes() ([]NoteMetadata, error) {
	data, err := p.performRequest("notes/list", nil)
	if err != nil {
		return nil, err
	}
	j, err := decodeJSONIFace(data)
	if err != nil {
		return nil, err
	}
	meta := make([]NoteMetadata, 0)
	count := int(j["count"].(float64))
	notes := j["notes"].([]interface{})
	for i := 0; i < count; i++ {
		n := notes[i].(map[string]interface{})
		t := NoteMetadata{ID: n["id"].(string),
			Title:   n["title"].(string),
			Length:  atoi(n["length"].(string)),
			Hash:    []byte(n["hash"].(string)),
			Created: parseDateTime(n["created_at"].(string)),
			Updated: parseDateTime(n["updated_at"].(string)),
		}
		meta = append(meta, t)
	}
	return meta, nil
}

// NoteID given a notes ID, returns an individual user note.
func (p *Pinboard) NoteID(id string) (Note, error) {
	data, err := p.performRequest("notes/"+id, nil)
	if err != nil {
		return Note{}, err
	}
	j, err := decodeJSONIFace(data)
	if err != nil {
		return Note{}, err
	}
	n := Note{}
	n.ID = j["id"].(string)
	n.Title = j["title"].(string)
	n.Length = int(j["length"].(float64))
	n.Hash = []byte(j["hash"].(string))
	n.Created = parseDateTime(j["created_at"].(string))
	n.Updated = parseDateTime(j["updated_at"].(string))
	n.Text = j["text"].(string)
	return n, nil
}
