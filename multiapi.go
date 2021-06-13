package multiapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiUrl = "https://api.itayki.com"

func doRequest(url string) ([]byte, error) {
	var res *http.Response
	var err error
	res, err = http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getUrl(path string, params map[string]string) string {
	toReturn := apiUrl + path + "?"

	if params != nil {
		for key, val := range params {
			toReturn += key + "=" + url.QueryEscape(val) + "&"
		}
	}

	fmt.Println(toReturn)

	return toReturn
}

func unmarshal(data []byte) map[string]interface{} {
	toReturn := make(map[string]interface{})
	json.Unmarshal(data, &toReturn)
	return toReturn
}

func GetExecLangs() (string, error) {
	body, err := doRequest(getUrl("/execlangs", nil))
	if err != nil {
		return "", err
	}

	result := unmarshal(body)
	toReturn := result["langs"].(string)
	return toReturn, nil
}

func ExecCode(lang string, code string) (string, error) {
	body, err := doRequest(
		getUrl("/exec",
			map[string]string{
				"lang": lang,
				"code": code,
			},
		),
	)
	if err != nil {
		return "", err
	}

	result := unmarshal(body)
	if _, ok := result["Errors"]; ok {
		return fmt.Sprintf(
			"Language: %s\n\nCode: %s\n\nResults: %s\n\nErrors: %s",
			result["Language"].(string),
			result["Code"].(string),
			result["Results"].(string),
			result["Errors"].(string),
		), nil
	} else if _, ok := result["Stats"]; ok {
		return fmt.Sprintf(
			"Language: %s\n\nCode: %s\n\nResults: %s\n\nStats: %s",
			result["Language"].(string),
			result["Code"].(string),
			result["Results"].(string),
			result["Stats"].(string),
		), nil
	}

	return result["langs"].(string), nil
}

func Ocr(url string) (string, error) {
	body, err := doRequest(getUrl("/ocr", map[string]string{"url": url}))
	if err != nil {
		return "", err
	}

	result := unmarshal(body)
	if _, ok := result["ocr"]; ok {
		return result["ocr"].(string), nil
	}

	return "Error: " + result["error"].(string), nil
}

func Translate(text string, from string, to string) (string, error) {
	body, err := doRequest(
		getUrl(
			"/tr",
			map[string]string{
				"text":     text,
				"fromlang": from,
				"lang":     to,
			},
		),
	)
	if err != nil {
		return "", err
	}

	result := unmarshal(body)
	if _, ok := result["error"]; ok {
		return result["error"].(string), nil
	}

	return fmt.Sprintf(
		"Text: %s\n\nFrom language: %s\n\nTo language: %s",
		result["text"].(string),
		result["from_language"].(string),
		result["to_language"].(string),
	), nil
}

func Urban(query string) ([]interface{}, error) {
	body, err := doRequest(getUrl("/ud", map[string]string{"query": query}))
	if err != nil {
		return nil, err
	}

	return unmarshal(body)["results"].([]interface{}), nil
}

func Webshot(url string, width string, height string) ([]byte, error) {
	if width == "" || height == "" {
		width = "1280"
		height = "720"
	}

	body, err := doRequest(
		getUrl(
			"/print",
			map[string]string{
				"url":    url,
				"width":  width,
				"height": height,
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func RandomNumber(min string, max string) (int, error) {
	body, err := doRequest(
		getUrl(
			"/random",
			map[string]string{
				"min": min,
				"max": max,
			},
		),
	)
	if err != nil {
		return 0, err
	}

	return unmarshal(body)["number"].(int), nil
}

func PyPiSearch(packageName string) (map[string]interface{}, error) {
	body, err := doRequest(
		getUrl(
			"/pypi",
			map[string]string{
				"package": packageName,
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return unmarshal(body), nil
}

func Paste(content string, title string, author string) (map[string]interface{}, error) {
	body, err := doRequest(
		getUrl(
			"/paste",
			map[string]string{
				"content": content,
				"title":   title,
				"author":  author,
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return unmarshal(body), nil
}

func GetPaste(paste string) (map[string]interface{}, error) {
	body, err := doRequest(
		getUrl(
			"/get_paste",
			map[string]string{
				"paste": paste,
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return unmarshal(body), nil
}
