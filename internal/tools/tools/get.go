package tools

import (
	"fmt"
	"nixie/internal/art"
	"os"
)

func GetTxt(targetPATH string) (string, error) {
	bytes, err := os.ReadFile(targetPATH)
	if err != nil {
		return "", err
	}
	text := string(bytes)
	return text, nil
}

func GetArt(kind string, name string) string {

	path := fmt.Sprintf("assets/%s/%s.txt", kind, name)
	bytes, err := art.GetFS().ReadFile(path)
	if err != nil {
		return ""
	}
	txt := string(bytes)
	return txt
}
