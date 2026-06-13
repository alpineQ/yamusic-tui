package wavesettings

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/dece2183/yamusic-tui/config"
)

type helpKeyMap struct {
	move   key.Binding
	change key.Binding
	apply  key.Binding
	cancel key.Binding
}

func newHelpMap() *helpKeyMap {
	controls := config.Current.Controls
	return &helpKeyMap{
		move:   key.NewBinding(key.WithKeys("up", "down"), key.WithHelp("↑/↓", "move")),
		change: key.NewBinding(key.WithKeys("left", "right"), key.WithHelp("←/→", "change")),
		apply:  key.NewBinding(controls.Apply.Binding(), controls.Apply.Help("apply")),
		cancel: key.NewBinding(controls.Cancel.Binding(), controls.Cancel.Help("cancel")),
	}
}

func (k *helpKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.move, k.change, k.apply, k.cancel}
}

func (k *helpKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
