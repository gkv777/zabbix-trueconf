package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	authPath      = "oauth2/v1/token"
	authGrantType = "client_credentials"

	apiPath    = "api"
	apiVersion = "v3.2"
	apiUser    = "users"
	apiConf    = "conferences"
)

// Client ...
type Client struct {
	TrueConfURL    string
	TrueConfClient string
	TrueConfSecret string
	TrueConfToken  Token
	UserAgent      string
	Debug          bool

	httpClient *http.Client
}

type usersResponse struct {
	Users []User `json:"users"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type clientResponse struct {
	Online int `json:"online"`
	Busy   int `json:"busy"`
	Active int `json:"active"`
	All    int `json:"all"`
}

// NewClient creates new Truconf API client
func NewClient(cfg *Config) (*Client, error) {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.TrueConfTLSInsec},
	}

	client := &Client{
		TrueConfURL:    cfg.TrueConfURL,
		TrueConfClient: cfg.TrueConfClient,
		TrueConfSecret: cfg.TrueConfSecret,
		TrueConfToken:  Token{},
		UserAgent:      cfg.UserAgent,
		Debug:          cfg.LogDebug,
		httpClient: &http.Client{
			Timeout:   5 * time.Minute,
			Transport: transport,
		},
	}

	if err := client.auth(); err != nil {
		return nil, err
	}

	return client, nil
}

// getUserList ...
func (c *Client) getUserList() (*[]User, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s/%s", c.TrueConfURL, apiPath, apiVersion, apiUser), nil)
	if err != nil {
		return nil, err
	}

	var users usersResponse
	if err := c.SendRequest(req, true, &users); err != nil {
		return nil, err
	}

	//log.Println(users)

	return &users.Users, nil
}

// auth ...
func (c *Client) auth() error {
	data := url.Values{}
	data.Add("grant_type", authGrantType)
	data.Add("client_id", c.TrueConfClient)
	data.Add("client_secret", c.TrueConfSecret)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.TrueConfURL, authPath), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//log.Println(req.Body)

	token := Token{}
	if err = c.SendRequest(req, false, &token); err != nil {
		return err
	}

	//log.Println(res)
	c.TrueConfToken = token

	if c.Debug {
		log.Println("API Client Auth OK. ", c.TrueConfToken.AsString())
	}

	return nil
}

// SendRequest ...
func (c *Client) SendRequest(req *http.Request, auth bool, v interface{}) error {
	req.Header.Set("Accept", "application/json; charset=utf-8")
	// if request needs authorization, add query param
	if auth {
		q := req.URL.Query()
		q.Add("access_token", c.TrueConfToken.AccessToken)
		req.URL.RawQuery = q.Encode()
	}

	if c.Debug {
		log.Println("Request URL: ", req.URL.String())
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	//log.Println(res.Body)

	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New("Response error: " + errRes.Message)
		}

		return fmt.Errorf("Response unknown error, status code: %d", res.StatusCode)
	}

	//fullResponse := successResponse{
	//	Data: v,
	//}
	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

// GetTrueConfInfo ...
func (c *Client) GetTrueConfInfo() {
	users, err := c.getUserList()
	if err != nil {
		log.Println("[GetTrueConfInfo|getUserListError]", err)
		fmt.Println("error")
	}
	getUsersInfo(*users, c.Debug)
}

// getUsersInfo ...
func getUsersInfo(us []User, debug bool){
	var i = clientResponse{}
	for _, u := range us {
		if u.IsActive == 1 {
			i.Active++
		}
		if u.Status == 2 {
			i.Busy++
			i.Online++
		}
		if u.Status == 1 {
			i.Online++
		}
	}
	i.All = len(us)
	if debug {
		log.Printf("GetUsersinfo: %d|%d|%d|%d\n", i.Online, i.Busy, i.Active, i.All)
	}
	e, err := json.Marshal(i)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))
}