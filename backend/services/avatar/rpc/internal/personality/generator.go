package personality

import (
	"time"
)

type Personality struct {
	Warmth      int32
	Adventurous int32
	Social      int32
	Creative    int32
	Calm        int32
	Energetic   int32
}

// GeneratePersonality 基于人口统计学数据生成初始人格
func GeneratePersonality(gender int32, birthDate, occupation string, maritalStatus int32) *Personality {
	p := &Personality{
		Warmth:      50,
		Adventurous: 50,
		Social:      50,
		Creative:    50,
		Calm:        50,
		Energetic:   50,
	}

	age := calculateAge(birthDate)

	// 基于年龄调整
	if age < 25 {
		p.Adventurous += 15
		p.Energetic += 10
		p.Social += 10
	} else if age >= 25 && age < 35 {
		p.Adventurous += 5
		p.Energetic += 5
	} else if age >= 35 && age < 50 {
		p.Calm += 10
		p.Warmth += 5
	} else {
		p.Calm += 15
		p.Warmth += 10
		p.Adventurous -= 10
	}

	// 基于职业调整
	switch occupation {
	case "creative", "艺术", "设计", "创意":
		p.Creative += 20
		p.Adventurous += 10
	case "technical", "技术", "工程师", "程序员":
		p.Creative += 10
		p.Calm += 10
	case "business", "商务", "销售", "市场":
		p.Social += 15
		p.Energetic += 10
	case "education", "教育", "教师":
		p.Warmth += 15
		p.Calm += 10
	case "medical", "医疗", "医生", "护士":
		p.Warmth += 15
		p.Calm += 15
	}

	// 基于婚姻状态调整
	switch maritalStatus {
	case 1: // 单身
		p.Adventurous += 5
		p.Social += 5
	case 2: // 恋爱中
		p.Warmth += 10
		p.Social += 5
	case 3: // 已婚
		p.Warmth += 10
		p.Calm += 10
	}

	// 基于性别调整（轻微）
	if gender == 2 { // 女性
		p.Warmth += 5
	} else if gender == 1 { // 男性
		p.Adventurous += 5
	}

	// 确保所有值在 0-100 范围内
	p.normalize()

	return p
}

func (p *Personality) normalize() {
	p.Warmth = clamp(p.Warmth, 0, 100)
	p.Adventurous = clamp(p.Adventurous, 0, 100)
	p.Social = clamp(p.Social, 0, 100)
	p.Creative = clamp(p.Creative, 0, 100)
	p.Calm = clamp(p.Calm, 0, 100)
	p.Energetic = clamp(p.Energetic, 0, 100)
}

func clamp(value, min, max int32) int32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func calculateAge(birthDate string) int {
	birth, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 30 // 默认年龄
	}
	now := time.Now()
	age := now.Year() - birth.Year()
	if now.YearDay() < birth.YearDay() {
		age--
	}
	return age
}
