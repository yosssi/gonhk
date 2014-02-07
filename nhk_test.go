// Copyright 2014 Keiji Yoshida.  All rights reserved.
// 情報提供:ＮＨＫ

// Package gonhk implements a NHK API client.
package gonhk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	apikey    string
	c         Client
	today     string
	programId string
)

// init executes initial processes.
func init() {
	apikey = os.Getenv("NHK_API_KEY")
	if apikey == "" {
		fmt.Println("You have to set the environment variable `NHK_API_KEY` to execute the tests.")
		os.Exit(1)
	}
	c = NewClient(apikey)
	today = time.Now().Format("2006-01-02")
}

// TestProgramList tests ProgramList.
func TestProgramList(t *testing.T) {
	// Normal end is expected.
	result, err := c.ProgramList("v1", "130", "g1", today)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("\n==== ProgramList Result List ==================================")
	for k, v := range result.List {
		for i, p := range v {
			fmt.Printf("\n%s %d %+v\n", k, i, p)
			if k == "g1" && i == 0 {
				programId = p.Id
			}
		}
	}
	fmt.Println("\n==============================================================\n")

	// Abnormal end is expected.
	result, err = c.ProgramList("v1", "130", "g1", "2001-01-01")
	if err == nil {
		t.Error("An expected error did not occur")
	}
	fmt.Println(err)
}

// TestProgramGenre tests ProgramGenre.
func TestProgramGenre(t *testing.T) {
	// Normal end is expected.
	result, err := c.ProgramGenre("v1", "130", "g1", "0000", today)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("\n==== ProgramGenre Result List ================================")
	for k, v := range result.List {
		for i, p := range v {
			fmt.Printf("\n%s %d %+v\n", k, i, p)
		}
	}
	fmt.Println("\n==============================================================\n")

	// Abnormal end is expected.
	result, err = c.ProgramGenre("v1", "130", "g1", "0000", "2001-01-01")
	if err == nil {
		t.Error("An expected error did not occur")
	}
	fmt.Println(err)
}

// TestProgramInfo tests ProgramInfo.
func TestProgramInfo(t *testing.T) {
	// Normal end is expected.
	result, err := c.ProgramInfo("v1", "130", "g1", programId)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("\n==== ProgramInfo Result List =================================")
	for k, v := range result.List {
		for i, p := range v {
			fmt.Printf("\n%s %d %+v\n", k, i, p)
		}
	}
	fmt.Println("\n==============================================================\n")

	// Abnormal end is expected.
	result, err = c.ProgramInfo("v1", "130", "g1", "0000000000000")
	if err == nil {
		t.Error("An expected error did not occur")
	}
	fmt.Println(err)
}

// TestNowOnAir tests NowOnAir.
func TestNowOnAir(t *testing.T) {
	// Normal end is expected.
	result, err := c.NowOnAir("v1", "130", "g1")
	if err != nil {
		t.Error(err)
	}

	fmt.Println("\n==== NowOnAir Result List ====================================")
	for k, v := range result.NowOnAirList {
		fmt.Printf("\n%s %+v\n", k, v)
	}
	fmt.Println("\n==============================================================\n")

	// Abnormal end is expected.
	result, err = c.NowOnAir("v1", "130", "00")
	if err == nil {
		t.Error("An expected error did not occur")
	}
	fmt.Println(err)
}

// TestNewClient tests NewClient.
func TestNewClient(t *testing.T) {
	testApikey := "testApikey"
	expected := Client{apikey: testApikey}
	actual := NewClient(testApikey)
	if expected.apikey != actual.apikey {
		t.Error(fmt.Sprintf("Client is invalid. [expected: %+v][actual: %+v]"), expected, actual)
	}
}

// TestGetNhkList tests getNhkList
func TestGetNhkList(t *testing.T) {
	// Normal end is expected.
	url := fmt.Sprintf("http://api.nhk.or.jp/v1/pg/list/130/g1/%s.json?key=%s", today, c.apikey)
	_, err := getNhkList(url)
	if err != nil {
		t.Error(err)
	}

	// Abnormal end is expected.
	url = ""
	_, err = getNhkList(url)
	if err == nil {
		t.Error("An expected error did not occur.")
	}
	fmt.Printf("%#v\n", err)

	// Abnormal end is expected.
	url = fmt.Sprintf("http://api.nhk.or.jp/v1/pg/list/130/g1/%s.json?key=%s", "2001-01-01", c.apikey)
	_, err = getNhkList(url)
	if err == nil {
		t.Error("An expected error did not occur.")
	}
	fmt.Printf("%#v\n", err)
}

// TestGetNhkDescriptionList tests getNhkDescriptionList
func TestGetNhkDescriptionList(t *testing.T) {
	// Normal end is expected.
	url := fmt.Sprintf("http://api.nhk.or.jp/v1/pg/info/130/g1/%s.json?key=%s", programId, c.apikey)
	_, err := getNhkDescriptionList(url)
	if err != nil {
		t.Error(err)
	}

	// Abnormal end is expected.
	url = ""
	_, err = getNhkDescriptionList(url)
	if err == nil {
		t.Error("An expected error did not occur.")
	}
	fmt.Printf("%#v\n", err)

	// Abnormal end is expected.
	url = fmt.Sprintf("http://api.nhk.or.jp/v1/pg/info/130/g1/%s.json?key=%s", "0000000000000", c.apikey)
	_, err = getNhkDescriptionList(url)
	if err == nil {
		t.Error("An expected error did not occur.")
	}
	fmt.Printf("%#v\n", err)
}

// TestGetNhkNowOnAirList tests getNhkNowOnAirList
func TestGetNhkNowOnAirList(t *testing.T) {
	// Normal end is expected.
	url := fmt.Sprintf("http://api.nhk.or.jp/v1/pg/now/130/g1.json?key=%s", c.apikey)
	_, err := getNhkNowOnAirList(url)
	if err != nil {
		t.Error(err)
	}

	// Abnormal end is expected.
	url = ""
	_, err = getNhkNowOnAirList(url)
	if err == nil {
		t.Error("An expected error did not occur.")
	}
	fmt.Printf("%#v\n", err)

	// Abnormal end is expected.
	url = fmt.Sprintf("http://api.nhk.or.jp/v1/pg/now/130/xxx.json?key=%s", c.apikey)
	_, err = getNhkNowOnAirList(url)
	if err == nil {
		t.Error("An expected error did not occur.")
	}
	fmt.Printf("%#v\n", err)
}

// TestDecoder tests decoder.
func TestDecoder(t *testing.T) {
	// Normal end is expected.
	url := fmt.Sprintf("http://api.nhk.or.jp/v1/pg/list/130/g1/%s.json?key=%s", today, c.apikey)
	_, err := decoder(url)
	if err != nil {
		t.Error(err)
	}

	// Abnormal end is expected.
	url = fmt.Sprintf("http://api.nhk.or.jp/v1/pg/list/130/g1/%s.json?key=%s", "2001-01-01", c.apikey)
	_, err = decoder(url)
	if err == nil {
		t.Error("An expected error did not occur.")
	}

	// Abnormal end is expected.
	url = ""
	_, err = decoder(url)
	if err == nil {
		t.Error("An expected error did not occur.")
	}
}

// TestApiError tests apiError.
func TestApiError(t *testing.T) {
	// Case1. Not returned NHK API error.
	url := fmt.Sprintf("http://api.nhk.or.jp/v1/pg/list/130/g1/%s.json?key=%s", today, c.apikey)
	res, _ := http.Get(url)
	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)
	dec := json.NewDecoder(bytes.NewReader(b))
	err := apiError(dec, 404)
	expected := "An error occurred during calling NHK API. [status: 404]"
	actual := err.Error()
	if expected != actual {
		t.Error(fmt.Sprintf("An error is invalid. [expected: %s][actual: %s]", expected, actual))
	}

	// Case2. Returned NHK API error.
	url = fmt.Sprintf("http://api.nhk.or.jp/v1/pg/list/130/g1/%s.json?key=%s", "2001-01-01", c.apikey)
	res, _ = http.Get(url)
	defer res.Body.Close()
	b, _ = ioutil.ReadAll(res.Body)
	dec = json.NewDecoder(bytes.NewReader(b))
	err = apiError(dec, 404)
	expected = "An error occurred during calling NHK API. [status: 404][code: 1][message: Invalid parameters]"
	actual = err.Error()
	if expected != actual {
		t.Error(fmt.Sprintf("An error is invalid. [expected: %s][actual: %s]", expected, actual))
	}
}
