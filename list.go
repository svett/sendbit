package sendbit

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"
)

// Represents a Recipient List
type List struct {
	// The list identificator
	ID uint64 `json:"id"`
	// The list name
	Name string `json:"list"`
}

// Creates a new recipient list
func (client *Client) CreateList(name string) error {
	errorf := func(err error) error {
		return client.errorf("CreateList", err)
	}

	if name == "" {
		return errorf(errors.New("The list name cannot be empty."))
	}

	data := url.Values{}
	data.Add("list", name)

	_, err := client.post("/newsletter/lists/add.json", data)
	if err != nil {
		return errorf(err)
	}

	return nil
}

// Remove a Recipient List from your account.
func (client *Client) DeleteList(name string) error {
	errorf := func(err error) error {
		return client.errorf("DeleteList", err)
	}

	if name == "" {
		return errorf(errors.New("The list name cannot be empty."))
	}

	data := url.Values{}
	data.Add("list", name)
	_, err := client.post("/newsletter/lists/delete.json", data)
	if err != nil {
		return errorf(err)
	}

	return nil
}

// List a recipient list
func (client *Client) List(name string) (*List, error) {
	errorf := func(err error) error {
		return client.errorf("List", err)
	}

	if name == "" {
		return nil, errorf(errors.New("The list name cannot be empty."))
	}

	data := url.Values{}
	data.Add("list", name)

	response, err := client.post("/newsletter/lists/get.json", data)
	if err != nil {
		return nil, errorf(err)
	}

	var lists []List
	if err := json.NewDecoder(response).Decode(&lists); err != nil && err != io.EOF {
		return nil, errorf(err)
	}

	if len(lists) == 1 {
		return &lists[0], nil
	}

	return nil, nil
}

// List all Recipient Lists on your account, or check if a particular List exists.
func (client *Client) Lists(names ...string) ([]List, error) {
	errorf := func(err error) error {
		return client.errorf("Lists", err)
	}

	response, err := client.post("/newsletter/lists/get.json", nil)
	if err != nil {
		return nil, errorf(err)
	}

	var lists []List
	if err := json.NewDecoder(response).Decode(&lists); err != nil && err != io.EOF {
		return nil, errorf(err)
	}

	return lists, nil
}
