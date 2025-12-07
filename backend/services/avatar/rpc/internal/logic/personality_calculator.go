package logic

// PersonalityTypeInfo 人格类型信息
type PersonalityTypeInfo struct {
	Code        string
	Name        string
	Description string
}

// PersonalityTypes 人格类型定义
var PersonalityTypes = map[string]PersonalityTypeInfo{
	"comforter": {
		Code:        "comforter",
		Name:        "温柔治愈师",
		Description: "温暖体贴，善于倾听和安慰他人",
	},
	"explorer": {
		Code:        "explorer",
		Name:        "调皮冒险家",
		Description: "充满好奇，勇于探索未知世界",
	},
	"mentor": {
		Code:        "mentor",
		Name:        "理性导师",
		Description: "逻辑清晰，善于引导和解决问题",
	},
	"poet": {
		Code:        "poet",
		Name:        "文学诗人",
		Description: "富有诗意，热爱创作和艺术",
	},
	"tsundere": {
		Code:        "tsundere",
		Name:        "傲娇少女",
		Description: "外冷内热，充满个性魅力",
	},
	"companion": {
		Code:        "companion",
		Name:        "深情伴侣",
		Description: "情感丰富，专注于陪伴与关怀",
	},
	"creator": {
		Code:        "creator",
		Name:        "创意大师",
		Description: "想象力丰富，擅长创新思维",
	},
	"socializer": {
		Code:        "socializer",
		Name:        "社交引导者",
		Description: "活力四射，善于社交和沟通",
	},
}

// CalculatePersonalityType 根据6维人格计算人格类型
func CalculatePersonalityType(warmth, adventurous, social, creative, calm, energetic int32) string {
	// 定义各个维度的权重阈值
	const (
		HIGH   = 60 // 高于60认为是高
		MEDIUM = 40 // 40-60认为是中等
		LOW    = 40 // 低于40认为是低
	)

	// 计算各维度的特征分数
	scores := make(map[string]float64)

	// 温柔治愈师：warmth 高 + calm 高
	scores["comforter"] = 0
	if warmth >= HIGH {
		scores["comforter"] += float64(warmth) * 1.5
	}
	if calm >= HIGH {
		scores["comforter"] += float64(calm) * 1.2
	}
	if social >= MEDIUM {
		scores["comforter"] += float64(social) * 0.5
	}

	// 调皮冒险家：adventurous 高 + energetic 高
	scores["explorer"] = 0
	if adventurous >= HIGH {
		scores["explorer"] += float64(adventurous) * 1.5
	}
	if energetic >= HIGH {
		scores["explorer"] += float64(energetic) * 1.3
	}
	if creative >= MEDIUM {
		scores["explorer"] += float64(creative) * 0.5
	}

	// 理性导师：calm 高 + social 中等，adventurous 不太高
	scores["mentor"] = 0
	if calm >= HIGH {
		scores["mentor"] += float64(calm) * 1.5
	}
	if social >= MEDIUM && social < HIGH {
		scores["mentor"] += float64(social) * 1.0
	}
	if adventurous < MEDIUM {
		scores["mentor"] += 30.0
	}
	if creative >= MEDIUM {
		scores["mentor"] += float64(creative) * 0.5
	}

	// 文学诗人：creative 高 + calm 高
	scores["poet"] = 0
	if creative >= HIGH {
		scores["poet"] += float64(creative) * 1.5
	}
	if calm >= HIGH {
		scores["poet"] += float64(calm) * 1.2
	}
	if social < MEDIUM {
		scores["poet"] += 20.0
	}

	// 傲娇少女：social 低 + creative 高 + energetic 高
	scores["tsundere"] = 0
	if social < MEDIUM {
		scores["tsundere"] += 40.0
	}
	if creative >= HIGH {
		scores["tsundere"] += float64(creative) * 1.2
	}
	if energetic >= HIGH {
		scores["tsundere"] += float64(energetic) * 1.2
	}

	// 深情伴侣：warmth 高 + social 高
	scores["companion"] = 0
	if warmth >= HIGH {
		scores["companion"] += float64(warmth) * 1.5
	}
	if social >= HIGH {
		scores["companion"] += float64(social) * 1.5
	}
	if calm >= MEDIUM {
		scores["companion"] += float64(calm) * 0.5
	}

	// 创意大师：creative 高 + adventurous 中高
	scores["creator"] = 0
	if creative >= HIGH {
		scores["creator"] += float64(creative) * 1.5
	}
	if adventurous >= MEDIUM {
		scores["creator"] += float64(adventurous) * 1.2
	}
	if energetic >= MEDIUM {
		scores["creator"] += float64(energetic) * 0.5
	}

	// 社交引导者：social 高 + energetic 高
	scores["socializer"] = 0
	if social >= HIGH {
		scores["socializer"] += float64(social) * 1.5
	}
	if energetic >= HIGH {
		scores["socializer"] += float64(energetic) * 1.5
	}
	if warmth >= MEDIUM {
		scores["socializer"] += float64(warmth) * 0.5
	}

	// 找出得分最高的人格类型
	maxScore := 0.0
	maxType := "comforter" // 默认为温柔治愈师

	for typeCode, score := range scores {
		if score > maxScore {
			maxScore = score
			maxType = typeCode
		}
	}

	return maxType
}

// GetPersonalityTypeName 获取人格类型的中文名称
func GetPersonalityTypeName(typeCode string) string {
	if info, ok := PersonalityTypes[typeCode]; ok {
		return info.Name
	}
	return "未知类型"
}

// GetPersonalityTypeDescription 获取人格类型的描述
func GetPersonalityTypeDescription(typeCode string) string {
	if info, ok := PersonalityTypes[typeCode]; ok {
		return info.Description
	}
	return ""
}
