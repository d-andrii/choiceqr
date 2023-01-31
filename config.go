package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
)

type Client struct {
	config Config
	http   http.Client
}

type Credentials struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type Config struct {
	Domain      string      `json:"domain"`
	Credentials Credentials `json:"credentials"`
}

type OrderItem struct {
	// "cash" | "online"
	PayBy string `json:"payBy"`
	// "closed" | "cancelled"
	Status string `json:"status"`
	// format: 4500 => 45.00
	Total int `json:"total"`
}

type OrdersResponse struct {
	Count int         `json:"count"`
	Items []OrderItem `json:"items"`
	Page  int         `json:"page"`
}

func NewClient() (*Client, error) {
	f, err := os.ReadFile("./config.json")
	if err != nil {
		if err == os.ErrNotExist {
			return nil, errors.New("config.json file not found")
		} else {
			return nil, err
		}
	}

	c := &Client{}
	if err = json.Unmarshal(f, &c.config); err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	c.http = http.Client{
		Jar: jar,
	}

	return c, nil
}

func (c *Client) getUrl(endpoint string) string {
	return "https://" + c.config.Domain + endpoint
}

func (c *Client) patchReq(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", c.getUrl(""))
	req.Header.Set("Referer", c.getUrl("/admin"))
}

func (c *Client) Auth() error {
	body, err := json.Marshal(c.config.Credentials)
	if err != nil {
		return err
	}

	url := c.getUrl("/api/auth/local")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	c.patchReq(req)

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Printf("URL: %s, response: %s", url, rBody)

	return nil
}

func (c *Client) GetChecksForDay(d string) ([]OrderItem, error) {
	var items []OrderItem

	for i := 1; ; i += 1 {
		url := c.getUrl("/api/payments/orders") + fmt.Sprintf("?from=%s&to=%s&page=%d&perPage=100", d, d, i)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		c.patchReq(req)

		res, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		log.Printf("URL: %s, response: %s", url, body)

		var orders OrdersResponse
		if err := json.Unmarshal(body, &orders); err != nil {
			return nil, err
		}

		items = append(items, orders.Items...)

		if len(items) == orders.Count {
			break
		}
	}

	return items, nil
}
