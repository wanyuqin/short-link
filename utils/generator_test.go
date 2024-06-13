package utils

import (
	"fmt"
	"testing"
)

func TestGenerateShortLink(t *testing.T) {
	fmt.Print(GenerateShortLink("localhost:8080/short-link/v1/user"))
}
