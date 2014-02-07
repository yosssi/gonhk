// Copyright 2014 Keiji Yoshida.  All rights reserved.
// 情報提供:ＮＨＫ

// Package gonhk implements a NHK API client.
package gonhk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	baseUrl = "http://api.nhk.or.jp/"
)

// A Client represents a NHK API client.
type Client struct {
	apikey string
}

// An NhkList represents a list data of NHK API.
type NhkList struct {
	List map[string][]*NhkProgram `json:"list"`
}

// An NhkDescriptionList represents a description list data of NHK API.
type NhkDescriptionList struct {
	List map[string][]*NhkDescription `json:"list"`
}

// An NhkNowOnAirList represents a now on air list data of NHK API.
type NhkNowOnAirList struct {
	NowOnAirList map[string]*NhkNowOnAir `json:"nowonair_list"`
}

// An NhkNowOnAir represents a now on air data of NHK API.
type NhkNowOnAir struct {
	Previous  NhkProgram `json:"previous"`
	Present   NhkProgram `json:"present"`
	Following NhkProgram `json:"following"`
}

// An NhkProgram represents a program data of NHK API.
type NhkProgram struct {
	Id        string     `json:"id"`
	EventId   string     `json:"event_id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
	Area      NhkArea    `json:"area"`
	Service   NhkService `json:"service"`
	Title     string     `json:"title"`
	Subtitle  string     `json:"subtitle"`
	Genres    []string   `json:"genres"`
}

// An NhkDescription represents a description data of NHK API.
type NhkDescription struct {
	Id          string     `json:"id"`
	EventId     string     `json:"event_id"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	Area        NhkArea    `json:"area"`
	Service     NhkService `json:"service"`
	Title       string     `json:"title"`
	Subtitle    string     `json:"subtitle"`
	Genres      []string   `json:"genres"`
	ProgramLogo NhkLogo    `json:"program_logo"`
	ProgramUrl  string     `json:"program_url"`
	EpisodeUrl  string     `json:"episode_url"`
	Hashtags    []string   `json:"hashtags"`
	Extras      NhkExtra   `json:"extras"`
}

// An NhkArea represents an area data of NHK API.
type NhkArea struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// An NhkService represents a service data of NHK API.
type NhkService struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	LogoS NhkLogo `json:"logo_s"`
	LogoM NhkLogo `json:"logo_m"`
	LogoL NhkLogo `json:"logo_l"`
}

// An NhkLogo represents a logo data of NHK API.
type NhkLogo struct {
	Url    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

// An NhkExtra represents an extra data of NHK API.
type NhkExtra struct {
	OndemandProgram NhkLink `json:"ondemand_program"`
	OndemandEpisode NhkLink `json:"ondemand_episode"`
}

// An NhkLink represents a link data of NHK API.
type NhkLink struct {
	Url   string `json:"url"`
	Title string `json:"title"`
	Id    string `json:"id"`
}

// An NhkError represents an error data of NHK API.
type NhkError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ProgramList calls NHK Program List API and returns its result.
func (c *Client) ProgramList(version, area, service, date string) (*NhkList, error) {
	url := fmt.Sprintf(baseUrl+"%s/pg/list/%s/%s/%s.json?key=%s", version, area, service, date, c.apikey)
	return getNhkList(url)
}

// ProgramGenre calls NHK Program Genre API and returns its result.
func (c *Client) ProgramGenre(version, area, service, genre, date string) (*NhkList, error) {
	url := fmt.Sprintf(baseUrl+"%s/pg/genre/%s/%s/%s/%s.json?key=%s", version, area, service, genre, date, c.apikey)
	return getNhkList(url)
}

// ProgramInfo calls NHK Program Info API and returns its result.
func (c *Client) ProgramInfo(version, area, service, id string) (*NhkDescriptionList, error) {
	url := fmt.Sprintf(baseUrl+"%s/pg/info/%s/%s/%s.json?key=%s", version, area, service, id, c.apikey)
	return getNhkDescriptionList(url)
}

// NowOnAir calls NHK Now On Air API and returns its result.
func (c *Client) NowOnAir(version, area, service string) (*NhkNowOnAirList, error) {
	url := fmt.Sprintf(baseUrl+"%s/pg/now/%s/%s.json?key=%s", version, area, service, c.apikey)
	return getNhkNowOnAirList(url)
}

// NewClient generates and returns a new client.
func NewClient(apikey string) Client {
	return Client{apikey: apikey}
}

// getNhkList gets and returns NhkList.
func getNhkList(url string) (*NhkList, error) {
	dec, err := decoder(url)
	if err != nil {
		return nil, err
	}
	var result *NhkList
	err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getNhkDescriptionList gets and returns NhkDescriptionList.
func getNhkDescriptionList(url string) (*NhkDescriptionList, error) {
	dec, err := decoder(url)
	if err != nil {
		return nil, err
	}
	var result *NhkDescriptionList
	err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getNhkNowOnAirList gets and returns NhkNowOnAirList.
func getNhkNowOnAirList(url string) (*NhkNowOnAirList, error) {
	dec, err := decoder(url)
	if err != nil {
		return nil, err
	}
	var result *NhkNowOnAirList
	err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// decoder sends an HTTP GET request and returns a JSON decoder which includes the response body data.
func decoder(url string) (*json.Decoder, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	statusCode := res.StatusCode
	if statusCode != http.StatusOK {
		return nil, apiError(dec, statusCode)
	}
	return dec, nil
}

// apiError returns an error which includes an HTTP statsu code, an API error code and an API error message.
func apiError(dec *json.Decoder, statusCode int) error {
	errResult := make(map[string]*NhkError)
	err := dec.Decode(&errResult)
	if err != nil {
		return err
	}
	errMessage := fmt.Sprintf("An error occurred during calling NHK API. [status: %d]", statusCode)
	if nhkError, prs := errResult["error"]; prs {
		errMessage += fmt.Sprintf("[code: %d][message: %s]", nhkError.Code, nhkError.Message)
	}
	return errors.New(errMessage)
}
