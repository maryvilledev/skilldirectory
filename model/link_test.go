package model

import "testing"

func TestLink_NewLink(t *testing.T) {
	linkOne := NewLink("Google", "http://www.google.com", "1234", WebpageLinkType)
	linkTwo := Link{
		Name:     "Google",
		URL:      "http://www.google.com",
		SkillID:  "1234",
		LinkType: WebpageLinkType,
	}

	if linkOne != linkTwo {
		t.Errorf("model/link\".NewLink()\" produced incorrect Link.")
	}
}

func TestIsValidLinkType1(t *testing.T) {
	if IsValidLinkType("InvalidLinkType") {
		t.Errorf("func IsValidLinkType() failed to detect invalid LinkType")
	}
}

func TestIsValidLinkType2(t *testing.T) {
	if !IsValidLinkType(WebpageLinkType) {
		t.Errorf("func IsValidLinkType() flagged valid LinkType as invalid.")
	}
}
