package sendbit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

// Creates a new instance of sendbit.Client for
// concreted SendGrid Account
func NewClient(username, password string) *Client {
	return &Client{
		Auth: &Auth{
			Username: username,
			Password: password,
		},
	}
}

func (client *Client) authenticate(data url.Values) (url.Values, error) {
	if client.Auth == nil ||
		client.Auth.Username == "" ||
		client.Auth.Password == "" {
		return data, fmt.Errorf("The client credentails are missing or invalid.")
	}

	if data == nil {
		data = make(url.Values)
	}

	data.Add("api_user", client.Auth.Username)
	data.Add("api_key", client.Auth.Password)
	return data, nil
}

func (client *Client) post(path string, data url.Values) (io.Reader, error) {
	values, err := client.authenticate(data)
	if err != nil {
		return nil, err
	}

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	host := fmt.Sprintf("https://api.sendgrid.com/api/%s", path)
	response, err := http.PostForm(host, values)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var message Response
	if err := json.Unmarshal(body, &message); err == nil {
		if message.Error != "" {
			return nil, fmt.Errorf("%s", message.Error)
		}

		if message.Message != "success" {
			return nil, fmt.Errorf("Non success message: %s", message.Message)
		}
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("The response status code is %d", response.StatusCode)
	}

	return bytes.NewReader(body), nil
}

func (client *Client) errorf(method string, err error) error {
	return fmt.Errorf("sendbit: client.%s error: %s", method, err)
}
