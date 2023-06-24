package extractor

import (
	"regexp"
)

func extractJobDescription(text *string, startStr string, endStr string) string {
	regexPattern := regexp.MustCompile(`(?s)` + startStr + "(.*?)" + endStr)
	matches := regexPattern.FindAllStringSubmatch(*text, -1)
	if len(matches) > 0 {
		return matches[0][1]
	}
	return ""
}

func ExtractIndeedJobDescription(text *string) string {
	startStr := "What\nWhere\nFind Jobs\n"
	endStr := "Report job\nHiring Lab\nCareer Advice\n"
	return extractJobDescription(text, startStr, endStr)

}

func ExtractSeekJobDescription(text *string) string {
	startStr := "Company reviews\nSign in\nEmployer site"
	endStr := "Job seekers\nJob search\n"
	return extractJobDescription(text, startStr, endStr)
}

func ExtractLinkedinJobDescription(text *string) string {
	startStr := "Dismiss\nJoin now\nSign in"
	endStr := "Referrals increase your chances of interviewing at "
	return extractJobDescription(text, startStr, endStr)
}
