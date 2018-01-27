package wolfram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	str "haze.sh/mairu/strutil"
)

const simpleWolframGateway = "http://api.wolframalpha.com/v1/result?appid=%s&i=%s"
const advancedWolframGateway = "https://api.wolframalpha.com/v2/query?input=%s&format=plaintext&output=JSON&appid=%s"

// type definitions

// ResultWrapper ...
type resultWrapper struct {
	Result Result `json:"queryresult"`
}

// Pod ...
type Pod struct {
	Title          string `json:"title"`
	Scanner        string `json:"scanner"`
	ID             string `json:"id"`
	Position       int    `json:"position"`
	Error          bool   `json:"error"`
	NumbeOfSubPods int    `json:"numsubpods"`
	SubPods        []Pod  `json:"subpods"`
	Plaintext      string `json:"plaintext"`
}

// Result ...
type Result struct {
	Success      bool    `json:"success"`
	Error        bool    `json:"error"`
	NumberOfPods int     `json:"numpods"`
	DataTypes    string  `json:"datatypes"`
	Timing       float32 `json:"timing"`
	ParseTiming  float32 `json:"parsetiming"`
	Pods         []Pod   `json:"pods"`
}

// AskAdvanced -> String -> String -> (*String, Error)
func AskAdvanced(key, query string) (*Result, error) {
	resp, err := http.Get(fmt.Sprintf(advancedWolframGateway, url.QueryEscape(query), url.QueryEscape(key)))
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var respWr resultWrapper
	err = json.Unmarshal(content, &respWr)
	if err != nil {
		return nil, err
	}
	return &respWr.Result, nil
}

// FindPod ...
func FindPod(pods []Pod, name string) *Pod {
	for _, v := range pods {
		fmt.Printf("Looping over %s...\n%+v\n\n", v.Title, v)
		if strings.ToLower(v.ID) == strings.ToLower(name) {
			return &v
		}
	}
	return nil
}

// FindPod ...
func (r Result) FindPod(name string) *Pod {
	return FindPod(r.Pods, name)
}

// GetInputPod ...
func (r Result) getInputPod() *Pod {
	return r.FindPod("Input")
}

func (r Result) getResultPod() *Pod {
	return r.FindPod("Result")
}

// GetInterpretation ...
func (r Result) GetInterpretation() *string {
	pod := r.getInputPod()
	if pod != nil {
		return str.S(pod.SubPods[0].Plaintext)
	}
	return nil
}

// GetResult ...
func (r Result) GetResult() *string {
	pod := r.getResultPod()
	if pod != nil {
		return str.S(pod.SubPods[0].Plaintext)
	}
	return nil
}

// AskSimple -> String -> String -> (*String, Error)
func AskSimple(key, query string) (*string, error) {
	resp, err := http.Get(fmt.Sprintf(simpleWolframGateway, url.QueryEscape(key), url.QueryEscape(query)))
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
