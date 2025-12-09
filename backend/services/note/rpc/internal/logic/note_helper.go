package logic

import (
	"encoding/json"

	"github.com/me2/note/rpc/note"
)

// parseTypesJson 解析 types JSON 字符串
func parseTypesJson(typesJson string) []string {
	var types []string
	if err := json.Unmarshal([]byte(typesJson), &types); err != nil {
		return []string{}
	}
	return types
}

// parseEmotionDataJson 解析 emotion_data JSON 字符串
func parseEmotionDataJson(emotionDataJson string) *note.EmotionData {
	var data struct {
		Primary string  `json:"primary"`
		Score   float64 `json:"score"`
	}
	if err := json.Unmarshal([]byte(emotionDataJson), &data); err != nil {
		return &note.EmotionData{
			Primary: "neutral",
			Score:   0.5,
		}
	}
	return &note.EmotionData{
		Primary: data.Primary,
		Score:   data.Score,
	}
}
