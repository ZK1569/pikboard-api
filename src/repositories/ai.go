package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	util "github.com/zk1569/pikboard-api/src/utils"
)

type ChatGpt struct {
}

var singleGPTInterface IAInterface

func GetChatGPTInstance() IAInterface {
	if singleGPTInterface == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleGPTInterface == nil {
			singleGPTInterface = &ChatGpt{}
		}
	}
	return singleGPTInterface
}

func (self *ChatGpt) ImageToFem(img string) (string, error) {
	payload := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]interface{}{
			{
				"role":    "developer",
				"content": "You will receive an image of a chess board, you must return the fem that represents the possitions pieces.",
			},
			{
				"role":    "developer",
				"content": "Your answer must only contain a FEM, nothing more. ",
			},
			{
				"role":    "assistant",
				"content": "brknnrqb/ppp1pppp/3p4/8/3PQN2/1N6/PPPKPPPP/BR3R1B b - - 0 1",
			},
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": img,
						},
					},
				},
			},
		},
		"max_tokens": 300,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("erreur lors du marshal du payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création de la requête: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+util.EnvVariable.Ai.OpenAiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'exécution de la requête: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("réponse non-200: %d, corps: %s", resp.StatusCode, string(bodyBytes))
	}

	var responseData struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", fmt.Errorf("erreur lors du décodage de la réponse: %w", err)
	}

	if len(responseData.Choices) == 0 {
		return "", fmt.Errorf("aucun choix dans la réponse")
	}

	return responseData.Choices[0].Message.Content, nil
}
