package message

import (
	"fmt"
	"html/template"
	"strings"
)

func init() {
	d := map[string][]string{
		"[what is a purine]": {
			"What is a purine simple definition",
			"Examples of high-purine foods",
			"Purines vs Pyrimidines structure and function",
			"What is the role of purines in DNA and RNA?",
		},

		"[CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";]": {
			"PostgreSQL generate UUID with uuid-ossp",
			"uuid-ossp vs gen_random_uuid() performance",
			"How to use CREATE EXTENSION in PostgreSQL",
			"PostgreSQL best practices for unique identifiers",
			"Troubleshoot \"uuid-ossp extension not found\" error",
		},

		"[caddy zitadel]": {
			"Caddy reverse proxy configuration for Zitadel",
			"Secure Zitadel with Caddy and mutual TLS",
			"Caddy vs Nginx as a reverse proxy for Zitadel",
			"Zitadel and Caddy automatic HTTPS setup",
			"Example Caddyfile for a Zitadel instance",
		},

		"[good smart tvs]": {
			"Best smart TVs under $500 in {{.Year}}",
			"OLED vs QLED vs Mini-LED TV comparison",
			"Samsung vs LG vs Sony smart TVs {{.Year}}",
			"Top-rated 65-inch smart TVs",
			"Smart TV buying guide {{.Year}}",
		},

		"[go remote debug vscode could not find file]": {
			"Go remote debug VSCode \"could not find file\" fix",
			"Configure dlv path mapping for remote Go debugging",
			"VSCode launch.json for remote Go debug attach",
			"Go debugger cannot find source files in container",
			"Troubleshoot Delve debugger 'can't find package' error",
		},

		"[Best toothpaste]": {
			"Best toothpaste for sensitive teeth {{.Year}}",
			"Top-rated whitening toothpaste that works",
			"Dentist recommended toothpaste brands",
			"Fluoride vs fluoride-free toothpaste pros and cons",
			"What is in toothpaste? Ingredients explained",
		},

		"[What is the capital of Japan?]": {
			"What is the current capital of Japan?",
			"Why did the capital of Japan move?",
			"Former capitals of Japan",
		},

		"[Who are the presidential contenders in 2028?]": {
			"Potential 2028 US presidential candidates",
			"Who is running for president in 2028?",
			"Potential 2028 Republican presidential candidates",
			"Potential 2028 Democratic presidential candidates",
			"2028 presidential election predictions",
		},

		"[angular js tutorial]": {
			"Learn AngularJS for beginners tutorial",
			"Should I learn AngularJS in {{.Year}}?",
			"Maintain a legacy AngularJS application",
		},

		"[driving test tips]": {
			"Tips for passing driving test first time",
			"Common mistakes on driving test to avoid",
			"DMV road test checklist",
			"How to parallel park for driving test",
			"How to stay calm before a driving test",
		},

		"[home espresso machine]": {
			"Best home espresso machines for beginners {{.Year}}",
			"Breville vs De'Longhi espresso machines",
			"Best espresso machine under $1000",
			"How to pull a perfect espresso shot at home",
			"Espresso machine cleaning and maintenance guide",
		},
	}
	for k, v := range d {
		QueryExpandData[k] = v
	}
}

var QueryExpandData map[string][]string = make(map[string][]string)

func TemplateToUserAssistant(tm map[string][]string, data any) ([]UserAssistant, error) {
	us := make([]UserAssistant, len(tm))
	for k, t := range tm {
		var sB strings.Builder
		for _, line := range t {
			sB.WriteString(line)
			sB.WriteString("\n")
		}
		// Remove trailing newline
		result := strings.TrimSpace(sB.String())

		tmpl, err := template.New(k).Parse(result)
		if err != nil {
			return nil, fmt.Errorf("unable to parse template: %v error: %w", k, err)
		}

		sB.Reset()
		if err := tmpl.Execute(&sB, data); err != nil {
			return nil, fmt.Errorf("unable to execute template: %v error: %w", k, err)
		}

		us = append(us, UserAssistant{
			User:      k,
			Assistant: sB.String(),
		})
	}
	return us, nil
}
