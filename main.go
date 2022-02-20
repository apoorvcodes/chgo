package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	config  Config
	baseURL = "https://coursehunter.net"
)

type Loading struct {
	text    string
	status  bool
	spinner spinner.Model
}

type model struct {
	activeIdx   int
	courses     []Course
	searchInput textinput.Model
	focussedIdx int
	loading     Loading
	err         error
}

type coursesMsg []Course
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

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, m.loading.spinner.Tick)
}

func (m model) View() string {
	appStyle := lipgloss.NewStyle().Padding(1)

	if m.loading.status {
		return appStyle.Render(m.loadingView())
	}

	list := m.coursesView()
	return appStyle.Render(m.searchInput.View() + list + m.courseView())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "k", "up":
			if m.activeIdx == 0 {
				m.activeIdx = len(m.courses) - 1
			} else {
				m.activeIdx--
			}
		case "j", "down":
			if m.activeIdx == len(m.courses)-1 {
				m.activeIdx = 0
			} else {
				m.activeIdx++
			}
		case "enter":
			if m.focussedIdx == 0 {
				m.loading.text = "searching for course: " + m.searchInput.Value()
				m.loading.status = true
				m.activeIdx = 0
				return m, fetchCourses(m.searchInput.Value())
			}
		case "tab":
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

	case errMsg:
		m.loading.status = false
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	m.loading.spinner, cmd = m.loading.spinner.Update(msg)

	return m, cmd
}

func NewModel() tea.Model {
	i := textinput.New()
	i.Placeholder = "search for courses"
	i.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		searchInput: i,
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
