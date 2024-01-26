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

func (M *Utils) MakeMatchesByPrefrences(userData *authPb.UserRepsonse, usersData []*authPb.UserRepsonse, person1Prefrences []domain.UserPreferences, person2sPrefrences []domain.UserPreferences) (responses.Result, error) {

	var content1 string
	var content2 string

	content1 += fmt.Sprintf("user_id %v - name: %v, height: %v, marital_status: %v, faith: %v, mother_tongue: %v, smoke_status: %v, alcohol_status: %v, settle_status: %v, hobbies: %v, tea_person: %v, love_language: %v, date_of_birth: %v", userData.UserID, userData.Fullname, person1Prefrences[0].Height, person1Prefrences[0].MaritalStatus, person1Prefrences[0].Faith, person1Prefrences[0].MotherTounge, person1Prefrences[0].SmokeStatus, person1Prefrences[0].AlcoholStatus, person1Prefrences[0].SettleStatus, person1Prefrences[0].Hobbies, person1Prefrences[0].TeaPerson, person1Prefrences[0].LoveLanguage, userData.Dob)

	for i, v := range usersData {

		c := fmt.Sprintf("user_id %v - name: %v, height: %v, marital_status: %v, faith: %v, mother_tongue: %v, smoke_status: %v, alcohol_status: %v, settle_status: %v, hobbies: %v, tea_person: %v, love_language: %v, date_of_birth: %v, place: %v", v.UserID, v.Fullname, person2sPrefrences[i].Height, person2sPrefrences[i].MaritalStatus, person2sPrefrences[i].Faith, person2sPrefrences[i].MotherTounge, person2sPrefrences[i].SmokeStatus, person2sPrefrences[i].AlcoholStatus, person2sPrefrences[i].SettleStatus, person2sPrefrences[i].Hobbies, person2sPrefrences[i].TeaPerson, person2sPrefrences[i].LoveLanguage, v.Dob, v.Location)

		content2 += c + " "

	}

	// match, err := M.OpenAi(content1, content2)

	// if err != nil {
	// 	return responses.Result{}, err
	// }

	//return match, nil
	return responses.Result{
		Result: []responses.Match{
			{
				UserID:     11,
				Name:       "Ashiq Sabith",
				MatchScore: 100,
				Age:        22,
				Place:      "Kottayam, Kerala",
			}, {
				UserID:     5,
				Name:       "Hisham Cs",
				MatchScore: 70,
				Age:        33,
				Place:      "Trissur, Kerala",
			}, {
				UserID:     3,
				Name:       "Abin V",
				MatchScore: 60,
				Age:        23,
				Place:      "Alappuzha, Kerala",
			},
		},
	}, nil
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
		Temperature: 0.5,
		Messages: []Message{
			{Role: "system", Content: `You're a matchmaker. I'll share details about a person as person1 and provide a list of potential matches with their name, preferences, and date of birth  as person2. Your task is to identify the top 5 perfect matches for person1, ordered by match score. Keep it concise. Format the response as 
			{
				result:[
					{
						"user_id": id of the user,
						"name": "name of the user",
						"marchscore": percenteage of match
						"age": age of the user. calculate it from the date of birth
						"palce":"place of the user"
					}
				] 
			}.
			`},
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
