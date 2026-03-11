package util

import (
	"fmt"
	"regexp"
)

func CheckGoogleDriveLink(link string) bool {
	re := regexp.MustCompile(`https?://(www\.)?drive\.google\.com/file/d/[^/]+/view\?usp=sharing`)
	return re.MatchString(link)
}

func ExtractDriveFileID(link string) (string, error) {
	re := regexp.MustCompile(`/d/([^/]+)`)
	match := re.FindStringSubmatch(link)

	if len(match) < 2 {
		return "", fmt.Errorf("invalid google drive link")
	}

	return match[1], nil
}

func ConvertToDownloadURL(link string) (string, error) {
	id, err := ExtractDriveFileID(link)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"https://drive.google.com/uc?export=download&id=%s",
		id,
	)

	return url, nil
}
