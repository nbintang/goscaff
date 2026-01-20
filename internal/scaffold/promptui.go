package scaffold

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

type Option struct {
	Label string
	Value string
}

type Wizard struct{}

func NewWizard() *Wizard { return &Wizard{} }

func (w *Wizard) SelectOption(ctx context.Context, label string, items []Option, def string) (string, error) {
	start := 0
	for i, it := range items {
		if it.Value == def {
			start = i
			break
		}
	}

	labels := make([]string, len(items))
	for i, it := range items {
		labels[i] = it.Label
	}

	p := promptui.Select{
		Label:     label,
		Items:     labels,
		CursorPos: start,
		Size:      len(items),
		Stdout:    os.Stderr,
	}

	idx, _, err := p.Run()
	if err != nil {
		return "", err
	}
	return items[idx].Value, nil
}

func (w *Wizard) Input(ctx context.Context, label, def string) (string, error) {
	fmt.Printf("%s [%s]: ", label, def)

	reader := bufio.NewReader(os.Stdin)

	done := make(chan struct{})
	var line string
	var err error

	go func() {
		line, err = reader.ReadString('\n')
		close(done)
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-done:
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			return def, nil
		}
		return line, nil
	}
}
