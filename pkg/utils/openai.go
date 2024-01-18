package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	authPb "github.com/ashiqsabith123/love-bytes-proto/auth/pb"
	"github.com/ashiqsabith123/match-svc/pkg/domain"
	"github.com/ashiqsabith123/match-svc/pkg/helper/responses"
)

func (M *Utils) MakeMatchesByPrefrences(userData *authPb.UserRepsonse, usersData []*authPb.UserRepsonse, usersPrefrences []domain.UserPreferences) (responses.Result, error) {

	var content1 string
	var content2 string

	content1 += fmt.Sprintf("user_id %v - name: %v, height: %v, marital_status: %v, faith: %v, mother_tongue: %v, smoke_status: %v, alcohol_status: %v, settle_status: %v, hobbies: %v, tea_person: %v, love_language: %v, date_of_birth: %v", userData.UserID, userData.Fullname, usersPrefrences[0].Height, usersPrefrences[0].MaritalStatus, usersPrefrences[0].Faith, usersPrefrences[0].MotherTounge, usersPrefrences[0].SmokeStatus, usersPrefrences[0].AlcoholStatus, usersPrefrences[0].SettleStatus, usersPrefrences[0].Hobbies, usersPrefrences[0].TeaPerson, usersPrefrences[0].LoveLanguage, userData.Dob)

	for i, v := range usersData {

		if i < 1 {
			continue
		}

		c := fmt.Sprintf("user_id %v - name: %v, height: %v, marital_status: %v, faith: %v, mother_tongue: %v, smoke_status: %v, alcohol_status: %v, settle_status: %v, hobbies: %v, tea_person: %v, love_language: %v, date_of_birth: %v, place: %v", v.UserID, v.Fullname, usersPrefrences[i].Height, usersPrefrences[i].MaritalStatus, usersPrefrences[i].Faith, usersPrefrences[i].MotherTounge, usersPrefrences[i].SmokeStatus, usersPrefrences[i].AlcoholStatus, usersPrefrences[i].SettleStatus, usersPrefrences[i].Hobbies, usersPrefrences[i].TeaPerson, usersPrefrences[i].LoveLanguage, v.Dob, v.Location)

		content2 += c + " "

	}

	match, err := M.OpenAi(content1, content2)

	if err != nil {
		return responses.Result{}, err
	}

	return match, nil
}

func (O *Utils) OpenAi(content1, content2 string) (responses.Result, error) {

	fmt.Println(content1)
	fmt.Println("")
	fmt.Println(content2)

	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type OpenAiPayload struct {
		Model       string    `json:"model"`
		Messages    []Message `json:"messages"`
		Temperature float64   `json:"temperature"`
	}

	// Create a struct to represent the "messages" field in the JSON payload

	// Construct the payload with dynamic content1 and content2
	payload := OpenAiPayload{
		Model:       "gpt-3.5-turbo",
		Temperature: 0.7,
		Messages: []Message{
			{Role: "system", Content: "You're a matchmaker. I'll share details about a person as person1 and provide a list of potential matches with their name, dob, and preferences as person2 in JSON format. Your task is to identify the top 10 perfect matches for person1, ordered by match score. Keep it concise. Format the response as {result: user_id, name, matchscore, age, place}."},
			{Role: "user", Content: "person1 - " + content1},
			{Role: "user", Content: "person2 list of persons - " + content2},
		},
	}

	// Convert the struct to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {

		return responses.Result{}, err
	}

	url := "https://api.openai.com/v1/chat/completions"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return responses.Result{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+O.config.OpenAi.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return responses.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return responses.Result{}, errors.New("Openai api err:" + resp.Status)
	}

	response, err := io.ReadAll(resp.Body)

	if err != nil {
		return responses.Result{}, err
	}

	fmt.Println(string(response))

	var matchResponse responses.ChatCompletionResponse
	err = json.Unmarshal(response, &matchResponse)

	if err != nil {
		return responses.Result{}, err
	}

	var match responses.Result

	err = json.Unmarshal([]byte(matchResponse.Choices[0].Message.Content), &match)

	if err != nil {
		return responses.Result{}, err
	}

	return match, nil
}
