package message

import (
	"fmt"
	"strconv"
	"strings"
)

func systemQueryUserString(data MessageData) string {
	var s strings.Builder
	if data.Age != 0 {
		s.WriteString(fmt.Sprintf("Age \"%d\"", data.Age))
	}
	if data.Age != 0 && data.Gender != "" {
		s.WriteString(" and ")
	}
	if data.Gender != "" {
		s.WriteString(fmt.Sprintf("Gender \"%s\"", data.Gender))
	}
	return s.String()
}
func SystemQueryExpand(min, max int, data MessageData) string {
	var s strings.Builder
	s.WriteString(
		"You are a Smart Search Assistant. Your task is to generate ",
	)
	s.WriteString(fmt.Sprintf("%d-%d", min, max))
	s.WriteString("diverse and highly relevant search queries based on a user's initial query and profile. The goal is to anticipate the user's next search actions.\n")
	s.WriteString("# RULES\n")
	if data.Age != 0 || data.Gender != "" {
		s.WriteString("- Where relevant (e.g., shopping, health, lifestyle), tailor suggestions to the User's")
		s.WriteString(systemQueryUserString(data))
		s.WriteString(". In other cases, ignore the profile.\n")
	}
	s.WriteString("- Use the year \"")
	s.WriteString(strconv.Itoa(data.Year))
	s.WriteString("\" for any suggestions about reviews, best products, or trends.\n")
	s.WriteString("- Output ONLY a newline-separated list of the suggestions.\n")
	s.WriteString("- Do not include numbers, bullet points, headers, or any introductory text.")
	return s.String()
}
