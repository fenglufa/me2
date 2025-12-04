package logic

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/me2/diary/rpc/diary"
	"github.com/me2/diary/rpc/internal/model"
)

// parseDiaryResponse 解析 AI 返回的日记内容
func parseDiaryResponse(text string) (title, content, mood string, tags []string) {
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "标题：") || strings.HasPrefix(line, "标题:") {
			title = strings.TrimPrefix(strings.TrimPrefix(line, "标题："), "标题:")
			title = strings.TrimSpace(title)
		} else if strings.HasPrefix(line, "心情：") || strings.HasPrefix(line, "心情:") {
			mood = strings.TrimPrefix(strings.TrimPrefix(line, "心情："), "心情:")
			mood = strings.TrimSpace(mood)
		} else if strings.HasPrefix(line, "标签：") || strings.HasPrefix(line, "标签:") {
			tagStr := strings.TrimPrefix(strings.TrimPrefix(line, "标签："), "标签:")
			tagStr = strings.TrimSpace(tagStr)
			tags = strings.Split(tagStr, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		} else if strings.HasPrefix(line, "内容：") || strings.HasPrefix(line, "内容:") {
			content = strings.TrimPrefix(strings.TrimPrefix(line, "内容："), "内容:")
			content = strings.TrimSpace(content)
		} else if content != "" {
			content += "\n" + line
		}
	}

	if title == "" {
		title = "今天的日记"
	}
	if mood == "" {
		mood = "平静"
	}
	if len(tags) == 0 {
		tags = []string{"日常"}
	}
	if content == "" {
		content = text
	}

	return title, content, mood, tags
}

// parseEmotionScore 解析情绪分数
func parseEmotionScore(text string) int32 {
	re := regexp.MustCompile(`[-]?\d+`)
	match := re.FindString(text)
	if match == "" {
		return 0
	}

	var score int32
	fmt.Sscanf(match, "%d", &score)

	if score < -100 {
		score = -100
	} else if score > 100 {
		score = 100
	}

	return score
}

// convertToDiaryInfo 转换 model 到 proto
func convertToDiaryInfo(d *model.Diaries) *diary.DiaryInfo {
	tags := []string{}
	if d.Tags != "" {
		tags = strings.Split(d.Tags, ",")
	}

	replyContent := ""
	if d.ReplyContent.Valid {
		replyContent = d.ReplyContent.String
	}

	emotionScore := int32(d.EmotionScore)

	return &diary.DiaryInfo{
		Id:           d.Id,
		AvatarId:     d.AvatarId,
		Type:         d.Type,
		Date:         d.Date.Format("2006-01-02"),
		Title:        d.Title,
		Content:      d.Content,
		Mood:         d.Mood,
		Tags:         tags,
		ReplyContent: replyContent,
		EmotionScore: emotionScore,
		IsImportant:  d.IsImportant == 1,
		CreatedAt:    d.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
