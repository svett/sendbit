package sendbit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Response struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Represent an SendBit credentials for SendGrid API
type Auth struct {
	// Username - An username of your SendGrid Account
	Username string
	// Password - A password of your SendGrid Account
	Password string
}

// A sendbit client for SendGrid REST API
// You should instanciate it in the following manner
//
// client := sendbit.NewClient("your_username", "your_password")
type Client struct {
	Auth *Auth
}

// Creates a new client from Environment variables
// SENDGRID_USER
// SENDGRID_PASS
func NewClientFromEnv() (*Client, error) {
	user := os.Getenv("SENDGRID_USER")
	pass := os.Getenv("SENDGRID_PASS")
	return NewClient(user, pass)
}

// Creates a new instance of sendbit.Client for
// concreted SendGrid Account
func NewClient(username, password string) (*Client, error) {
	if username == "" {
		return nil, errors.New("sendbit: Username argument cannot be empty.")
	}
	if password == "" {
		return nil, errors.New("sendbit: Password argument cannot be empty.")
	}
	return &Client{
		Auth: &Auth{
			Username: username,
			Password: password,
		},
	}, nil
}

func (client *Client) authenticate(data url.Values) error {
	if client.Auth == nil ||
		client.Auth.Username == "" ||
		client.Auth.Password == "" {
		return fmt.Errorf("The client credentails are missing or invalid.")
	}

	data.Add("api_user", client.Auth.Username)
	data.Add("api_key", client.Auth.Password)
	return nil
}

func (client *Client) post(path string, data url.Values) (io.Reader, error) {
	if data == nil {
		data = url.Values{}
	}
	if err := client.authenticate(data); err != nil {
		return nil, err
	}

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	host := fmt.Sprintf("https://api.sendgrid.com/api/%s", path)

	request, err := http.NewRequest("POST", host, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "sendbit/0.0.1;go")

	httpClient := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   5 * time.Second,
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var message Response
	if err := json.Unmarshal(body, &message); err == nil && message.Error != "" {
		return nil, errors.New(message.Error)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("The response status code is %d", response.StatusCode)
	}

	return bytes.NewReader(body), nil
}

func (client *Client) errorf(method string, err error) error {
	return fmt.Errorf("sendbit: client.%s error: %s", method, err)
}
