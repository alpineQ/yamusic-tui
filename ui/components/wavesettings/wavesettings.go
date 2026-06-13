package wavesettings

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dece2183/yamusic-tui/config"
	"github.com/dece2183/yamusic-tui/ui/model"
	"github.com/dece2183/yamusic-tui/ui/style"
)

type Control uint

const (
	APPLY Control = iota
	CANCEL
)

type option struct {
	value string
	label string
}

type axis struct {
	title   string
	options []option
	index   int
}

var (
	moodOptions = []option{
		{"all", "Любое"},
		{"calm", "Спокойное"},
		{"sad", "Грустное"},
		{"fun", "Весёлое"},
		{"active", "Бодрое"},
	}
	diversityOptions = []option{
		{"default", "Любое"},
		{"favorite", "Любимое"},
		{"popular", "Популярное"},
		{"discover", "Незнакомое"},
	}
	languageOptions = []option{
		{"any", "Любой"},
		{"russian", "Русский"},
		{"not-russian", "Иностранный"},
		{"without-words", "Без слов"},
	}
)

type Model struct {
	help     help.Model
	helpKeys *helpKeyMap
	width    int
	cursor   int
	axes     []axis
}

func New() *Model {
	return &Model{
		help:     help.New(),
		helpKeys: newHelpMap(),
		axes: []axis{
			{title: "Под настроение", options: moodOptions},
			{title: "По характеру", options: diversityOptions},
			{title: "По языку", options: languageOptions},
		},
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(message tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := message.(type) {
	case tea.KeyMsg:
		controls := config.Current.Controls
		keypress := msg.String()

		switch {
		case controls.Apply.Contains(keypress):
			cmds = append(cmds, model.Cmd(APPLY))
		case controls.Cancel.Contains(keypress):
			cmds = append(cmds, model.Cmd(CANCEL))
		case controls.CursorUp.Contains(keypress):
			if m.cursor > 0 {
				m.cursor--
			}
		case controls.CursorDown.Contains(keypress):
			if m.cursor < len(m.axes)-1 {
				m.cursor++
			}
		case keypress == "left" || keypress == "h":
			a := &m.axes[m.cursor]
			a.index = (a.index - 1 + len(a.options)) % len(a.options)
		case keypress == "right" || keypress == "l":
			a := &m.axes[m.cursor]
			a.index = (a.index + 1) % len(a.options)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	titleStyle := lipgloss.NewStyle().Foreground(style.InactiveTextColor).Width(16)
	valueStyle := lipgloss.NewStyle().Foreground(style.NormalTextColor)

	rows := make([]string, 0, len(m.axes))
	for i, a := range m.axes {
		label := a.options[a.index].label
		title := titleStyle.Render(a.title)
		var value string
		if i == m.cursor {
			value = style.AccentTextStyle.Render("‹ " + label + " ›")
		} else {
			value = valueStyle.Render("  " + label + "  ")
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, title, value))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		style.DialogTitleStyle.Render("Настройки Моей волны"),
		style.DialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...)),
		style.DialogHelpStyle.Render(m.help.View(m.helpKeys)),
	)
}

func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m *Model) SetValues(moodEnergy, diversity, language string) {
	set := func(a *axis, v string) {
		for i := range a.options {
			if a.options[i].value == v {
				a.index = i
				return
			}
		}
		a.index = 0
	}
	set(&m.axes[0], moodEnergy)
	set(&m.axes[1], diversity)
	set(&m.axes[2], language)
	m.cursor = 0
}

func (m *Model) MoodEnergy() string {
	return m.axes[0].options[m.axes[0].index].value
}

func (m *Model) Diversity() string {
	return m.axes[1].options[m.axes[1].index].value
}

func (m *Model) Language() string {
	return m.axes[2].options[m.axes[2].index].value
}
