package domain_tests

import (
	"testing"

	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
)

func TestSkillLevelParsingAndComparison(t *testing.T) {
	s, err := model.ParseSkillLevel("Junior")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if s != model.SkillJunior {
		t.Fatalf("expected SkillJunior got %v", s)
	}

	a := model.SkillJuniorPlus
	b := model.SkillJunior
	if !a.AtLeast(b) {
		t.Fatalf("%v should be >= %v", a, b)
	}

	if model.SkillIntern.AtLeast(model.SkillJunior) {
		t.Fatalf("intern should not be >= junior")
	}
}
