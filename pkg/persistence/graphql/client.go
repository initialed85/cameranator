package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Error struct {
	Message string `json:"message"`
}

type ManyResponseBody struct {
	Data   map[string][]interface{} `json:"data"`
	Errors []Error                  `json:"errors"`
}

type SingleResponseBody struct {
	Data   map[string]interface{} `json:"data"`
	Errors []Error                `json:"errors"`
}

type Client struct {
	url        string
	httpClient *http.Client
}

func NewClient(
	url string,
	timeout time.Duration,
) *Client {
	c := Client{
		url: url,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}

	return &c
}

func (c *Client) Query(
	query string,
) (map[string][]interface{}, error) {
	requestBody := map[string]string{
		"query": strings.TrimSpace(query),
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return map[string][]interface{}{}, err
	}

	response, err := c.httpClient.Post(
		c.url,
		"application/json",
		bytes.NewBuffer(requestBodyJSON),
	)

	defer func() {
		if response == nil || response.Body == nil {
			return
		}
		_ = response.Body.Close()
	}()

	if err != nil {
		return map[string][]interface{}{}, err
	}

	responseBodyJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return map[string][]interface{}{}, err
	}

	manyResponseBody := ManyResponseBody{}
	singleResponseBody := SingleResponseBody{}

	many := true

	err = json.Unmarshal(responseBodyJSON, &manyResponseBody)
	if err != nil {
		err = json.Unmarshal(responseBodyJSON, &singleResponseBody)
		if err != nil {
			return map[string][]interface{}{}, err
		}
		many = false
	}

	var errors []Error

	if many {
		errors = manyResponseBody.Errors
	} else {
		errors = singleResponseBody.Errors
	}

	if len(errors) > 0 {
		return map[string][]interface{}{}, fmt.Errorf(
			"got %#+v attempting query %v",
			errors,
			query,
		)
	}

	data := make(map[string][]interface{})

	if many {
		data = manyResponseBody.Data
	} else {
		for k, v := range singleResponseBody.Data {
			data[k] = make([]interface{}, 0)
			data[k] = append(data[k], v)
		}
	}

	return data, nil
}

func (c *Client) Extract(
	data map[string][]interface{},
	key string,
	result interface{},
) error {
	// try to infer key if none given
	if key == "" {
		keys := make([]string, 0)
		for k := range data {
			keys = append(keys, k)
		}

		if len(keys) == 0 {
			return fmt.Errorf("no possible keys to infer from (is query sane?)")
		}

		if len(keys) > 1 {
			return fmt.Errorf("cannot infer key- more than one possibility (please specify key)")
		}

		key = keys[0]
	}

	dataJSON, err := json.MarshalIndent(data[key], "", "    ")
	if err != nil {
		return err
	}

	return json.Unmarshal(dataJSON, result)
}

func (c *Client) QueryAndExtract(
	query string,
	key string,
	result interface{},
) error {
	data, err := c.Query(query)
	if err != nil {
		return err
	}

	return c.Extract(
		data,
		key,
		result,
	)
}

func (c *Client) QueryAndExtractMultiple(
	query string,
	keys []string,
	results ...interface{},
) error {
	if len(keys) != len(results) {
		return fmt.Errorf("keys and results must be of the same length")
	}

	data, err := c.Query(query)
	if err != nil {
		return err
	}

	for i, key := range keys {
		err = c.Extract(
			data,
			key,
			&results[i],
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Mutate(
	mutation string,
) (map[string][]interface{}, error) {
	requestBody := map[string]interface{}{
		"query":     strings.TrimSpace(mutation),
		"variables": nil,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return map[string][]interface{}{}, err
	}

	response, err := c.httpClient.Post(
		c.url,
		"application/json",
		bytes.NewBuffer(requestBodyJSON),
	)

	defer func() {
		if response == nil || response.Body == nil {
			return
		}
		_ = response.Body.Close()
	}()

	if err != nil {
		return map[string][]interface{}{}, err
	}

	responseBodyJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return map[string][]interface{}{}, err
	}

	manyResponseBody := ManyResponseBody{}
	singleResponseBody := SingleResponseBody{}

	many := true

	err = json.Unmarshal(responseBodyJSON, &manyResponseBody)
	if err != nil {
		err = json.Unmarshal(responseBodyJSON, &singleResponseBody)
		if err != nil {
			return map[string][]interface{}{}, err
		}
		many = false
	}

	var errors []Error

	if many {
		errors = manyResponseBody.Errors
	} else {
		errors = singleResponseBody.Errors
	}

	if len(errors) > 0 {
		return map[string][]interface{}{}, fmt.Errorf(
			"got %#+v attempting query %v",
			errors,
			mutation,
		)
	}

	data := make(map[string][]interface{})

	if many {
		data = manyResponseBody.Data
	} else {
		for k, v := range singleResponseBody.Data {
			data[k] = make([]interface{}, 0)
			data[k] = append(data[k], v)
		}
	}

	return data, nil
}
