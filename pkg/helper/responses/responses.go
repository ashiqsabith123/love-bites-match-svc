package responses

import "github.com/ashiqsabith123/love-bytes-proto/match/pb"

type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Content string `json:"content"`
}

type Result struct {
	Result []Match `json:"result"`
}

type Match struct {
	UserID     uint   `json:"user_id"`
	Name       string `json:"name"`
	MatchScore float32    `json:"matchscore"`
	Age        int    `json:"age"`
	Place      string `json:"place"`
	Photos     []*pb.Images
}
