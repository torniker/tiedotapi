package tiedotapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Model interface for using tiedot
type Model interface {
	CollactionName() string
	Migrate() []string
	SetCreatedAt(time.Time)
}

// Query tiedot query object
type Query struct {
	Eq string   `json:"eq"`
	In []string `json:"in"`
}

var ctx *TD

// TD is a struct for tiedot database
type TD struct {
	URL  string
	Port int
}

// NewTD return TD context
func NewTD() TD {
	if ctx == nil {
		port, err := strconv.Atoi(getEnv("TIEDOT_PORT", "5830"))
		if err != nil {
			port = 5830
		}
		ctx = &TD{
			URL:  getEnv("TIEDOT_URL", "http://localhost"),
			Port: port,
		}
	}
	return *ctx
}

func getEnv(name, def string) string {
	env := os.Getenv(name)
	if env == "" {
		env = def
	}
	return env
}

func (td TD) String() string {
	return td.URL + ":" + strconv.Itoa(td.Port) + "/"
}

// GetPage get a page of collection
func (td TD) GetPage(obj Model, page, total int) (*http.Response, error) {
	req, err := http.NewRequest("GET", td.String()+"getpage", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("col", obj.CollactionName())
	q.Add("page", strconv.Itoa(page))
	q.Add("total", strconv.Itoa(total))
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	return client.Do(req)
}

// Query tiedot collection
func (td TD) Query(obj Model, query Query) (*http.Response, error) {
	req, err := http.NewRequest("GET", td.String()+"query", nil)
	if err != nil {
		return nil, err
	}
	newObj := new(bytes.Buffer)
	err = json.NewEncoder(newObj).Encode(query)
	if err != nil {
		return nil, err
	}
	// fmt.Println(newObj.String())
	q := req.URL.Query()
	q.Add("col", obj.CollactionName())
	q.Add("q", newObj.String())
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	return client.Do(req)
}

// Insert inserts given object
func (td TD) Insert(obj Model) error {
	newObj := new(bytes.Buffer)
	obj.SetCreatedAt(time.Now())
	err := json.NewEncoder(newObj).Encode(obj)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", td.String()+"insert", newObj)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("col", obj.CollactionName())
	// q.Add("doc", newObj.String())
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	bodyStr := string(bodyBytes)
	if res.StatusCode != 201 {
		return errors.New(bodyStr)
	}
	// id := bodyStr
	// obj.SetID(id)
	// td.Update(obj, id)
	return nil
}

// Update updates given object
func (td TD) Update(obj Model, id string) error {
	newObj := new(bytes.Buffer)
	err := json.NewEncoder(newObj).Encode(obj)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", td.String()+"update", nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("col", obj.CollactionName())
	q.Add("id", id)
	q.Add("doc", newObj.String())
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	// defer res.Body.Close()
	// bodyBytes, _ := ioutil.ReadAll(res.Body)
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)
	return nil
}

// Delete deletes given document with sepcified ID
func (td TD) Delete(obj Model, id string) error {
	newObj := new(bytes.Buffer)
	err := json.NewEncoder(newObj).Encode(obj)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", td.String()+"delete", nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("col", obj.CollactionName())
	q.Add("id", id)
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	// defer res.Body.Close()
	// bodyBytes, _ := ioutil.ReadAll(res.Body)
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)
	return nil
}

// Migrate is a function to get everything up to date
func (TD TD) Migrate() error {
	// TODO: implement checking migrations for all collections
	return nil
}
