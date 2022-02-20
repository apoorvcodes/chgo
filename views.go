package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) loadingView() string {
	return m.loading.spinner.View() + m.loading.text
}

func (m model) coursesView() string {
	var list string

	for i, course := range m.courses {
		cursor := " "
		if m.focussedIdx == 1 && m.activeIdx == i {
			cursor = ">"
		}

		list += fmt.Sprintf("%s %s\n", cursor, course.Title)
	}

	return lipgloss.NewStyle().MarginTop(2).Render(list)
}

func (m model) courseView() string {
	if len(m.courses) == 0 || m.focussedIdx != 1 {
		return ""
	}

	course := m.courses[m.activeIdx]

	view := fmt.Sprintf("  Title:    %s\n", course.Title)
	view += fmt.Sprintf("  Source:   %s\n", course.Source)
	view += fmt.Sprintf("  Language: %s\n", course.Language)

	// for books, duration and lessons fields will by empty
	if len(course.Duration) != 0 {
		view += fmt.Sprintf("  Duration: %s\n", course.Duration)
	}

	if len(course.Lessons) != 0 {
		view += fmt.Sprintf("  Lessons:  %s\n", course.Lessons)
	}

	return lipgloss.NewStyle().MarginTop(1).Render(view)
}
