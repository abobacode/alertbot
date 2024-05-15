package alertbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type client struct {
	Host   string
	Path   string
	Client http.Client

	token string
}

func (c *client) updates(offset, limit int) ([]Update, error) {
	query := url.Values{}

	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	var resp UpdatesResponse
	if err := c.sendRequest("getUpdates", query, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

func (c *client) sendMessage(chatID int, text string) error {
	query := url.Values{}

	query.Add("chat_id", strconv.Itoa(chatID))
	query.Add("text", text)

	var val interface{}

	if err := c.sendRequest("sendMessage", query, &val); err != nil {
		return err
	}

	return nil
}

func (c *client) sendPhoto(chatID int, text, path string, buttons [][]string) error {
	keyboard := map[string]interface{}{
		"keyboard":        buttons,
		"resize_keyboard": true,
	}

	keyboardJSON, err := json.Marshal(keyboard)
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	params := map[string]string{
		"chat_id":      strconv.Itoa(chatID),
		"caption":      text,
		"reply_markup": string(keyboardJSON),
	}

	if err = c.sendFile(params, "photo", filepath.Base(path), file); err != nil {
		return err
	}

	return nil
}

func (c *client) sendFile(params map[string]string, paramName, fileName string, file *os.File) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	_ = writer.Close()

	urlSendPhoto := fmt.Sprintf("https://api.telegram.org/bot%s/%s", c.token, "sendPhoto")
	req, err := http.NewRequest(http.MethodPost, urlSendPhoto, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func (c *client) sendRequest(method string, query url.Values, value interface{}) error {
	u := url.URL{
		Scheme: "https",
		Host:   c.Host,
		Path:   path.Join(c.Path, method),
	}

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	request.URL.RawQuery = query.Encode()

	response, err := c.Client.Do(request)
	if err != nil {
		return err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

	if err = json.NewDecoder(response.Body).Decode(value); err != nil {
		return err
	}

	return nil
}

func newClient(host, token string) *client {
	return &client{
		Host:   host,
		Path:   "bot" + token,
		Client: http.Client{},
	}
}
