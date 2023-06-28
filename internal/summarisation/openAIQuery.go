package summarisation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const completionModel = "text-davinci-003"
const completionThread = 4
const summarisationModel = "gpt-4-32k" // also work if we fetch less job descriptions and use gpt-3.5-turbo/gpt-3.5-turbo-16k

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
		"model": "%s",
		"temperature": 0.2,
		"max_tokens": 501,
		"top_p": 0.4,
		"frequency_penalty": 0,
		"presence_penalty": 0,
		"prompt": "%s"}`, completionModel, escapedTxt)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apikey)

	return req, nil

}

func getAPIKey() (string, error) {
	if key, exist := os.LookupEnv("OPEN_AI_API"); exist {
		return key, nil
	} else {
		return "", fmt.Errorf("OPEN AI API KEY not exist")
	}
}

func GetTechnicalQualitiesForAJobPrompt(txt []byte) []byte {
	additionalText := "\nKey technical job qualities required for this job include:\n\n"
	txt = append(txt, []byte(additionalText)...)
	return txt
}

type CompletionResponse struct {
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

	var res CompletionResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		fmt.Println(res.Choices)
		return "", nil
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

func (api *OpenAIClient) GetTechnicalQualitiesForJobs(folderPath string, outputFile string) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	sb := strings.Builder{}
	buffer := make([]string, completionThread)
	errs := make([]error, completionThread)
	wg.Add(completionThread)
	for start := 0; start < completionThread; start++ {
		startIdx := start
		go func(int) {
			res := strings.Builder{}
			for idx := startIdx; idx < len(files); idx += completionThread {
				dt, err := ioutil.ReadFile(path.Join(folderPath, files[idx].Name()))
				if err != nil {
					errs[startIdx] = err
					wg.Done()
					return
				}
				resp, err := api.GetTechnicalQualitiesForAJob(dt)
				if err != nil {
					errs[startIdx] = err
					wg.Done()
					return
				}
				res.WriteString(resp)
				res.WriteString("\n")
			}
			buffer[startIdx] = res.String()
			wg.Done()
		}(startIdx)
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	for _, res := range buffer {
		sb.WriteString(res)
	}
	return ioutil.WriteFile(outputFile, []byte(sb.String()), 0644)

}

// above: extract key job qualities from each post, below: count and summarise

type CountJobQualityOccurrenceRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	MaxTokens        int       `json:"max_tokens"`
	Temperature      float64   `json:"temperature"`
	TopP             float64   `json:"top_p"`
	FrequencyPenalty float64   `json:"frequency_penalty"`
	PresencePenalty  float64   `json:"presence_penalty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (api *OpenAIClient) GetChatCompletion(txt string) (*http.Request, error) {
	apikey, err := getAPIKey()
	if err != nil {
		return nil, err
	}

	sysMessageContent := `You are an HR who is responsible for finding the top 10 technical requirements in short phrases that occurred most often and their corresponding counts and rank based on that. Similar technical qualities can be combined to  a single one with higher level of summarisation. Return the results in a csv format with header "Category Summarisation", "Count". Note only return top 10 results.`
	userMessageContent := txt

	messages := []Message{
		{
			Role:    "system",
			Content: sysMessageContent,
		},
		{
			Role:    "user",
			Content: userMessageContent,
		},
	}

	reqBody := CountJobQualityOccurrenceRequest{
		Model:            summarisationModel,
		Messages:         messages,
		MaxTokens:        500,
		Temperature:      0.0,
		TopP:             1.0,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
	}

	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonReqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apikey)

	return req, nil
}

type CountJobQualityOccurrenceResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func (api *OpenAIClient) GetTechnicalQualitiesCountSummary(txt []byte) (string, error) {
	req, err := api.GetChatCompletion(string(txt))
	if err != nil {
		return "", err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res CountJobQualityOccurrenceResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		fmt.Println(res.Choices)
		return "", nil
	}

	return res.Choices[0].Message.Content, nil
}
