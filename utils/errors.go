package utils

import (
	"restman/components/config"

	"github.com/charmbracelet/lipgloss"
)

func RenderErrors(errors []string) string {
	errorMessages := ""
	if len(errors) > 0 {
		errorsToRender := make([]string, len(errors))
		for i, e := range errors {
			errorsToRender[i] = config.ErrorStyle.Render("* " + e)
		}

		errorMessages = lipgloss.JoinVertical(lipgloss.Left, errorsToRender...) + "\n"
	}
	return errorMessages
}
