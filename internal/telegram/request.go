package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Create a struct that mimics the webhook response body
// https://core.telegram.org/bots/api#update

type requestBody struct {
	Message RequestMessage `json:"message"`
}

type RequestMessage struct {
	ID      int64   `json:"message_id"`
	From    User    `json:"from"`
	Chat    Chat    `json:"chat"`
	Text    string  `json:"text"`
	Date    int64   `json:"date"`
	Photo   []Photo `json:"photo"`
	Caption string  `json:"caption"`
}

type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Username string `json:"username"`
}

type Photo struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

func (r RequestMessage) TextOrCaption() string {
	if len(r.Text) > 0 {
		return r.Text
	}

	if len(r.Caption) > 0 {
		return r.Caption
	}

	return ""
}

// Returns the username of a user or if he has none, the firstname and lastname
func (u User) UsernameOrName() string {
	if len(u.Username) > 0 {
		return "@" + u.Username
	}

	return strings.Trim(fmt.Sprintf("%s %s", u.FirstName, u.LastName), " ")
}

// Parse a request body and returns the message
func ParseRequest(_ http.ResponseWriter, req *http.Request) (*RequestMessage, error) {
	body := &requestBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		return nil, fmt.Errorf("could not decode request body: %s", err.Error())
	}

	return &body.Message, nil
}
