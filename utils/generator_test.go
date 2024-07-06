package utils

import (
	"testing"
)

func TestGenerateShortLink(t *testing.T) {
	t.Log(GenerateShortLink("localhost:8080/short-link/v1/user"))
}
