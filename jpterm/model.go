package jpterm

import (
	"encoding/json"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goccy/go-yaml"
	"github.com/jmespath/go-jmespath"
)

func NewModel(data interface{}, mode int, expressionsValue string) *Model {

	expression := textinput.New()
	expression.Prompt = " JMESPath Expression: "
	expression.Focus()
	if expressionsValue != "" {
		expression.SetValue(expressionsValue)
	}

	content, err := yaml.Marshal(data)
	if err != nil {
		log.Fatal("Broken data")
	}

	input := viewport.New(1, 1)
	input.SetContent(yamlFormatter(content))

	output := viewport.New(1, 1)
	output.SetContent(" ")

	return &Model{
		mode:       Mode(mode),
		expression: expression,
		input:      input,
		data:       data,
		output:     output,
	}

}

type Mode int8

const (
	Result Mode = iota
	Expression
	Quiet
)

type Model struct {
	height int
	width  int
	mode   Mode

	expression textinput.Model
	input      viewport.Model
	data       interface{}
	result     interface{}
	output     viewport.Model
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg { return tea.KeyMsg{} }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			switch m.mode {
			case Result:
				dat, _ := json.Marshal(m.result)
				print(string(dat), m.height)
			case Expression:
				print(m.expression.Value(), m.height)
			case Quiet:
			}
			return m, tea.Quit
		}
		if msg.String() == "tab" {
			m.mode += 1
			if m.mode > 2 {
				m.mode -= 3
			}
		}
		m.expression, cmd = m.expression.Update(msg)
		result, err := jmespath.Search(m.expression.Value(), m.data)
		if err != nil || result == nil {
			expressionStyle.Background(subtle)
			IOStyle.BorderForeground(highlight)
		} else {
			m.result = result
			bytes, _ := yaml.Marshal(result)
			m.output.SetContent(yamlFormatter(bytes))

			expressionStyle.Background(special)
			IOStyle.BorderForeground(special)
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		// msg.Height -= 2
		// msg.Width -= 4
		// w := ((m.width / 2) - 6)
		h := m.height - 7

		expressionStyle.Width(m.width - 19)
		IOStyle.Width((m.width / 2) - 6)

		// m.input.Width = w
		// m.output.Width = w
		m.input.Height = h
		m.output.Height = h
	}

	switch m.mode {
	case Result:
		rModeStyle.Background(special)
		eModeStyle.Background(subtle)
		qModeStyle.Background(subtle)
	case Expression:
		rModeStyle.Background(subtle)
		eModeStyle.Background(special)
		qModeStyle.Background(subtle)
	case Quiet:
		rModeStyle.Background(subtle)
		eModeStyle.Background(subtle)
		qModeStyle.Background(special)
	}
	return m, cmd
}

func (m Model) View() string {
	s := lipgloss.NewStyle().MaxHeight(m.height).MaxWidth(m.width).Padding(1, 2, 1, 2)
	return s.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.JoinHorizontal(lipgloss.Top, expressionStyle.Render(m.expression.View()), rModeStyle.Render(" R "), eModeStyle.Render(" E "), qModeStyle.Render(" Q ")),
			lipgloss.JoinHorizontal(lipgloss.Top, IOStyle.Render(m.input.View()), IOStyle.Render(m.output.View())),
		),
	)
}
