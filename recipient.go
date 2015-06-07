package sendbit

import (
	"encoding/json"
	"errors"
	"net/url"
)

// An email recipient subscribed to particular recipient list
type Recipient struct {
	// This is a recipient's name
	Name string `json:"name"`
	// This is a recipient's email
	Email string `json:"email"`
}

// Add an email recipient to a list
func (client *Client) AddRecipient(list string, recipient *Recipient) error {
	errorf := func(err error) error {
		return client.errorf("AddRecipient", err)
	}

	if list == "" {
		return errorf(errors.New("The list is empty."))
	}

	if recipient == nil || recipient.Email == "" {
		return errorf(errors.New("The recipeint is nil or has invalid email."))
	}

	body, err := json.Marshal(&recipient)
	if err != nil {
		return errorf(err)
	}

	data := url.Values{}
	data.Add("list", list)
	data.Add("data", string(body))

	response, err := client.post("/newsletter/lists/email/add.json", data)

	if err != nil {
		return errorf(err)
	}

	var stats struct {
		AffectedRows int `json:"inserted"`
	}

	if err := json.NewDecoder(response).Decode(&stats); err != nil {
		return errorf(err)
	}

	if stats.AffectedRows == 0 {
		return errorf(errors.New("The recipient is not added."))
	}

	return nil
}

// Remove one or more emails from a Recipient List.
func (client *Client) DeleteRecipient(list, email string) error {
	errorf := func(err error) error {
		return client.errorf("DeleteRecipient", err)
	}

	if list == "" {
		return errorf(errors.New("The list is empty."))
	}

	if email == "" {
		return errorf(errors.New("The recipeint email is empty."))
	}

	data := url.Values{}
	data.Add("list", list)
	data.Add("email[]", email)

	response, err := client.post("/newsletter/lists/email/delete.json", data)
	if err != nil {
		return errorf(err)
	}

	var stats struct {
		AffectedRows int `json:"removed"`
	}

	if err := json.NewDecoder(response).Decode(&stats); err != nil {
		return errorf(err)
	}

	if stats.AffectedRows == 0 {
		return errorf(errors.New("The recipient is not removed."))
	}

	return nil
}

// Get the email and associated fields for a Recipient List.
func (client *Client) Recipient(list, email string) (*Recipient, error) {
	errorf := func(err error) error {
		return client.errorf("Recipient", err)
	}

	if list == "" {
		return nil, errorf(errors.New("The list is empty."))
	}

	if email == "" {
		return nil, errorf(errors.New("The email is empty."))
	}

	data := url.Values{}
	data.Add("list", list)
	data.Add("email", email)
	response, err := client.post("/newsletter/lists/email/get.json", data)
	if err != nil {
		return nil, errorf(err)
	}

	var recipients []Recipient
	err = json.NewDecoder(response).Decode(&recipients)
	if err != nil {
		return nil, errorf(err)
	}

	if len(recipients) == 0 {
		return nil, nil
	}

	return &recipients[0], nil
}

// Get the email addresses and associated fields for a Recipient List.
func (client *Client) Recipients(list string) ([]Recipient, error) {
	errorf := func(err error) error {
		return client.errorf("Recipients", err)
	}

	if list == "" {
		return nil, errorf(errors.New("The list is empty."))
	}

	data := url.Values{}
	data.Add("list", list)
	response, err := client.post("/newsletter/lists/email/get.json", data)
	if err != nil {
		return nil, errorf(err)
	}

	var recipients []Recipient
	err = json.NewDecoder(response).Decode(&recipients)
	if err != nil {
		return nil, errorf(err)
	}

	return recipients, nil
}

// Retrieve the number of entries on a list.
func (client *Client) RecipientCount(list string) (uint64, error) {
	errorf := func(err error) error {
		return client.errorf("RecipientCount", err)
	}

	if list == "" {
		return 0, errorf(errors.New("The list is empty."))
	}

	data := url.Values{}
	data.Add("list", list)
	response, err := client.post("/newsletter/lists/email/count.json", data)
	if err != nil {
		return 0, errorf(err)
	}

	var stats struct {
		Count uint64 `json:"count"`
	}

	err = json.NewDecoder(response).Decode(&stats)
	if err != nil {
		return 0, errorf(err)
	}

	return stats.Count, nil
}
