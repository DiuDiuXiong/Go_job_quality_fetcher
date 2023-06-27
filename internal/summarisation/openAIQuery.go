package summarisation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type OpenAIClient struct {
	httpClient *http.Client
}

func NewOpenAIClient() *OpenAIClient {
	return &OpenAIClient{httpClient: &http.Client{}}
}

func CompletionRequest(txt []byte) (*http.Request, error) {
	apikey, err := getAPIKey()
	if err != nil {
		return nil, err
	}

	escapedTxt := strconv.Quote(string(txt))
	escapedTxt = escapedTxt[1 : len(escapedTxt)-1]

	jsonStr := fmt.Sprintf(`{
		"model": "text-davinci-003",
		"temperature": 0.2,
		"max_tokens": 501,
		"top_p": 0.4,
		"frequency_penalty": 0,
		"presence_penalty": 0,
		"prompt": "%s"}`, escapedTxt)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apikey)

	return req, nil

}

func GetTechnicalQualitiesForAJobPrompt(txt []byte) []byte {
	additionalText := "\nKey technical job qualities required for this job include:\n\n"
	txt = append(txt, []byte(additionalText)...)
	return txt
}

type Response struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func (api *OpenAIClient) GetTechnicalQualitiesForAJob(txt []byte) (string, error) {
	req, err := CompletionRequest(GetTechnicalQualitiesForAJobPrompt(txt))
	if err != nil {
		return "", err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	return splitAndTrim(res.Choices[0].Text), nil
}

func splitAndTrim(input string) string {
	lines := regexp.MustCompile(`\n+`).Split(input, -1)
	for i, line := range lines {
		trimmedLine := regexp.MustCompile(`^[•\s]*(.*)$`).ReplaceAllString(line, "$1")
		lines[i] = trimmedLine
	}
	return strings.Join(lines, "\n")
}

func getAPIKey() (string, error) {
	if key, exist := os.LookupEnv("OPEN_AI_API"); exist {
		return key, nil
	} else {
		return "", fmt.Errorf("OPEN AI API KEY not exist")
	}
}
