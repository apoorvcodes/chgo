package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	/* colors */
	accentColor = lipgloss.Color("#78E2A0")
	foregroundColor  = lipgloss.Color("#FFFFFF")
	headerColor   = lipgloss.Color("#5840a7")

	/* styles */
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(accentColor)
	headerStyle       = lipgloss.NewStyle().Background(headerColor).Foreground(foregroundColor).Padding(0, 1)
)

/* === lessonsView */
type lessonItem string

func (i lessonItem) FilterValue() string { return "" }

type lessonItemDelegate struct{}

func (d lessonItemDelegate) Height() int { return 1 }

func (d lessonItemDelegate) Spacing() int { return 0 }

func (d lessonItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d lessonItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(lessonItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

/* === lessonsView */

func (m model) loadingView() string {
	return m.loading.spinner.View() + m.loading.text
}

func (m model) coursesView() string {
	if len(m.courses) == 0 {
		return ""
	}

	list := headerStyle.Render("Search results")
	list += "\n\n"

	for i, course := range m.courses {
		cursor := " "
		if m.focussedIdx == 1 && m.activeCourseIdx == i {
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

	course := m.courses[m.activeCourseIdx]

	view := headerStyle.Render("Selected course")
	view += "\n\n"
	view += fmt.Sprintf("  Title:    %s\n", course.Title)
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
