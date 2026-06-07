package onboarding

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// RunGuided runs interactive onboarding and returns the collected answers.
func RunGuided() (*Answers, error) {
	answers := &Answers{}
	if err := survey.Ask(Questions, answers); err != nil {
		return nil, fmt.Errorf("onboarding: guided prompt failed: %w", err)
	}
	return answers, nil
}
