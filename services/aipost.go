package services

import (
	"apost/models"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func GeneratePost(topic string) (*models.AiResult, error) {
	var err error
	var client = &http.Client{}
	var data map[string]interface{}

	model := "google/gemma-2-2b-it"
	jsonBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": "Create article about " + topic,
			},
		},
		"max_tokens": 2000,
		//"stream":     true,
	}

	jsonMarshal, err := json.Marshal(jsonBody)

	bodyReader := bytes.NewReader(jsonMarshal)

	request, err := http.NewRequest("POST", os.Getenv("HUGGINGFACE_API_URL")+"/"+model+"/v1/chat/completions", bodyReader)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+os.Getenv("HUGGINGFACE_TOKEN"))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	var result models.AiResult
	res, err := json.Marshal(&data)
	json.Unmarshal(res, &result)

	return &result, nil
}
