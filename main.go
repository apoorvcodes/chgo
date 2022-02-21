package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	config        Config
	baseURL       = "https://coursehunter.net"
	lessonsWidth  = 30
	lessonsHeight = 23
)

type Loading struct {
	text    string
	status  bool
	spinner spinner.Model
}

type model struct {
	courses         []Course
	lessons         []Lesson
	searchInput     textinput.Model
	lessonsList     list.Model
	loading         Loading
	err             error
	focussedIdx     int
	activeCourseIdx int
	activeLessonIdx int
	screenIdx       int
}

type coursesMsg []Course
type lessonsMsg []Lesson
type errMsg error

func fetchCourses(url string) tea.Cmd {
	return func() tea.Msg {
		courses, err := searchCourses(url)
		if err != nil {
			return errMsg(err)
		}

		return coursesMsg(courses)
	}
}

func fetchLessons(url string) tea.Cmd {
	return func() tea.Msg {
		lessons, err := getLessons(url)
		if err != nil {
			return errMsg(err)
		}

		return lessonsMsg(lessons)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, m.loading.spinner.Tick)
}

func (m model) View() string {
	appStyle := lipgloss.NewStyle().Padding(1)

	if m.loading.status {
		return appStyle.Render(m.loadingView())
	}

	if m.screenIdx == 1 {
		return appStyle.Render(m.lessonsList.View())
	}

	list := m.coursesView()
	return appStyle.Render(m.searchInput.View() + list + m.courseView())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
  
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "k", "up":
			if m.activeCourseIdx == 0 {
				m.activeCourseIdx = len(m.courses) - 1
			} else {
				m.activeCourseIdx--
			}
		case "j", "down":
			if m.activeCourseIdx == len(m.courses)-1 {
				m.activeCourseIdx = 0
			} else {
				m.activeCourseIdx++
			}
		case "P":
			m.screenIdx = 0
		case "enter":
			currIdx := m.activeCourseIdx

      if m.screenIdx == 0 {
        m.activeCourseIdx = 0
        m.loading.status = true
      }

			if m.focussedIdx == 0 && m.screenIdx == 0 {
				m.loading.text = "searching for course: " + m.searchInput.Value()
				return m, fetchCourses(m.searchInput.Value())
			}

			if m.focussedIdx == 1 && m.screenIdx == 0 {
				course := m.courses[currIdx]
				m.loading.text = "fetching lessons for: " + course.Title
				return m, fetchLessons(course.URL)
			}

      if m.screenIdx == 1 {
        i := m.lessonsList.Index()
        url := m.lessons[i].File
        playLesson(url)
      }
		case "tab":
			if m.screenIdx != 0 {
				return m, nil
			}

			next := m.focussedIdx + 1

			if m.focussedIdx == 0 {
				m.searchInput.Blur()
			}

			if m.focussedIdx == 1 {
				m.searchInput.Focus()
				next = 0
			}

			m.focussedIdx = next
		}

	case coursesMsg:
		m.loading.status = false
		m.courses = msg

	case lessonsMsg:
		m.loading.status = false
		m.lessons = msg

		var items []list.Item
		for _, lesson := range msg {
			items = append(items, lessonItem(lesson.Title))
		}

		m.lessonsList.SetItems(items)
		m.screenIdx = 1

	case errMsg:
		m.loading.status = false
		return m, tea.Quit

  case spinner.TickMsg:
    m.loading.spinner, cmd = m.loading.spinner.Update(msg)
    return m, cmd
	}

  if m.screenIdx == 0 {
    m.searchInput, cmd = m.searchInput.Update(msg)
  }

	if m.screenIdx == 1 {
		m.lessonsList, cmd = m.lessonsList.Update(msg)
	}

	return m, cmd
}

func NewModel() tea.Model {
	i := textinput.New()
	i.Placeholder = "search for courses"
	i.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(accentColor)

	l := list.New([]list.Item{}, lessonItemDelegate{}, lessonsWidth, lessonsHeight)
	l.Title = "Lessons"
	l.Styles.Title = headerStyle
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	return model{
		searchInput: i,
		lessonsList: l,
		loading: Loading{
			status:  false,
			text:    "",
			spinner: s,
		},
	}
}

func startTUI() {
	if err := setConfig(); err != nil {
		fmt.Println(`failed to login`)
		return
	}

	if err := isTokenExpired(); err != nil {
		fmt.Println(`token expired`)
		return
	}

	if err := tea.NewProgram(NewModel()).Start(); err != nil {
		fmt.Printf("uh oh: %s", err)
		os.Exit(1)
	}
}

func main() {
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	email := loginCmd.String("u", "", "email")
	password := loginCmd.String("p", "", "password")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "login":
			loginCmd.Parse(os.Args[2:])

			if len(*email) == 0 || len(*password) == 0 {
				fmt.Println("missing credentials")
				return
			}

			config, err := login(*email, *password)
			if err != nil {
				fmt.Printf("failed to login: %s\n", err)
				return
			}

			err = createConfig(config)
			if err != nil {
				fmt.Printf("failed to create config: %s\n", err)
			}

			fmt.Println("Logged in")
		default:
			fmt.Println("expected subcommand 'login'")
		}
	} else {
		startTUI()
	}
}
