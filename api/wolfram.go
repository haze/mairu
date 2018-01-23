package wolfram

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"haze.sh/mairu/util"
)

const wolframGateway = "http://api.wolframalpha.com/v1/result?appid=%s&i=%s"

// Ask -> String -> String -> (*String, Error)
func Ask(key, query string) (*string, error) {
	resp, err := http.Get(fmt.Sprintf(wolframGateway, url.QueryEscape(key), url.QueryEscape(query)))
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return str.S(string(content)), nil
}
