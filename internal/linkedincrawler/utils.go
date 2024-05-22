package linkedincrawler

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func getJobPostingID(company_id string, position string) string {
	data := company_id + "|" + position

	hasher := sha256.New()

	hasher.Write([]byte(data))

	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}

func containsBlacklistedWord(blacklistedWords []string, str string) bool {
	for _, blacklistedWord := range blacklistedWords {
		if strings.Contains(strings.ToLower(str), strings.ToLower(blacklistedWord)) {
			return true
		}
	}
	return false
}
