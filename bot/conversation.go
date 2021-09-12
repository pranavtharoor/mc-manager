package bot

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/pranavtharoor/mc-manager/config"
	"gopkg.in/square/go-jose.v2/json"
)

type InferenceQuery struct {
	Inputs     InferenceInputs     `json:"inputs"`
	Parameters InferenceParameters `json:"parameters"`
}

type InferenceInputs struct {
	PastUserInputs     []string `json:"past_user_inputs"`
	GeneratedResponses []string `json:"generated_responses"`
	Text               string   `json:"text"`
}

type InferenceParameters struct {
	TopK              float32 `json:"top_k"`
	TopP              float32 `json:"top_p"`
	Temperature       float32 `json:"temperature"`
	RepetitionPenalty float32 `json:"repetition_penalty"`
}

type InferenceResult struct {
	GeneratedText string `json:"generated_text"`
}

func newInferenceQuery(config config.ConversationConfiguration, text string, pastUserInputs []string, generatedResponses []string) *InferenceQuery {
	inferenceInputs := InferenceInputs{PastUserInputs: pastUserInputs, GeneratedResponses: generatedResponses, Text: text}
	parameters := InferenceParameters{TopK: config.TopK, TopP: config.TopP, Temperature: config.Temperature, RepetitionPenalty: config.RepetitionPenalty}
	return &InferenceQuery{Inputs: inferenceInputs, Parameters: parameters}
}

func conversation(config config.ConversationConfiguration, text string, pastUserInputs []string, generatedResponses []string, retryNumber int) string {
	inferenceQuery := newInferenceQuery(config, text, pastUserInputs, generatedResponses)
	queryJSON, err := json.Marshal(inferenceQuery)
	if err != nil {
		return err.Error()
	}
	req, err := http.NewRequest("POST", config.URL, bytes.NewBuffer(queryJSON))
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	inferenceResult := &InferenceResult{}

	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(&inferenceResult)
		if err != nil {
			return err.Error()
		}
		return inferenceResult.GeneratedText
	} else {
		var out map[string]interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &out)
		return "I don't feel like speaking right now. Try talking to me in a minute"
	}
}
