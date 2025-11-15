package model

import (
	"fmt"
	"strings"
)

type SkillLevel int

const (
	SkillUnknown SkillLevel = iota
	SkillIntern
	SkillJunior
	SkillJuniorPlus
	SkillMiddle
	SkillMiddlePlus
	SkillSenior
)

func (s SkillLevel) String() string {
	switch s {
	case SkillIntern:
		return "Intern"
	case SkillJunior:
		return "Junior"
	case SkillJuniorPlus:
		return "Junior+"
	case SkillMiddle:
		return "Middle"
	case SkillMiddlePlus:
		return "Middle+"
	case SkillSenior:
		return "Senior"
	default:
		return "Unknown"
	}
}

// parses textual representations, case-insensitive
func ParseSkillLevel(in string) (SkillLevel, error) {
	switch strings.ToLower(strings.TrimSpace(in)) {
	case "intern", "стажер", "trainee":
		return SkillIntern, nil
	case "junior", "джун", "jun":
		return SkillJunior, nil
	case "junior+", "джун+", "juniorplus", "jun+":
		return SkillJuniorPlus, nil
	case "middle", "миддл", "mid":
		return SkillMiddle, nil
	case "middle+", "миддл+", "mid+":
		return SkillMiddlePlus, nil
	case "senior", "сеньор", "sen":
		return SkillSenior, nil
	default:
		return SkillUnknown, fmt.Errorf("unknown skill level: %s", in)
	}
}

// returns true if s >= other
func (s SkillLevel) AtLeast(other SkillLevel) bool {
	return s >= other
}
