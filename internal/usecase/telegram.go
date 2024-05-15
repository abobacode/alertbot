package usecase

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

	"alertbot/config"
	"alertbot/internal/models"
)

type Client struct {
	Host   string
	Path   string
	Client http.Client
}

func (c *Client) Updates(offset, limit int) ([]models.Update, error) {
	query := url.Values{}

	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	data, err := c.SendRequest("getUpdates", query)
	if err != nil {
		return nil, err
	}

	var resp models.UpdatesResponse

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	query := url.Values{}

	query.Add("chat_id", strconv.Itoa(chatID))
	query.Add("text", text)

	_, err := c.SendRequest("sendMessage", query)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendPhoto(chatID int, text, path string, buttons [][]string) error {
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

	if err := SendFile(params, "photo", filepath.Base(path), file); err != nil {
		return err
	}

	return nil
}

func SendFile(params map[string]string, paramName, fileName string, file *os.File) error {
	cfg, err := config.New("./config.yaml")
	if err != nil {
		return err
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, fileName)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal(err)
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	_ = writer.Close()

	urlSendPhoto := fmt.Sprintf("https://api.telegram.org/bot%s/%s", cfg.AlertBot.Token, "sendPhoto")
	req, err := http.NewRequest("POST", urlSendPhoto, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Ошибка HTTP запроса: %s", resp.Status)
	}

	return nil
}

func (c *Client) SendRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.Host,
		Path:   path.Join(c.Path, method),
	}

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	request.URL.RawQuery = query.Encode()

	response, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func New(host, token string) *Client {
	return &Client{
		Host:   host,
		Path:   "bot" + token,
		Client: http.Client{},
	}
}
