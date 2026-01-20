package scaffold

import (
	"context"
	"errors"

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
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	idx, _, err := p.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			return "", context.Canceled
		}
		return "", err
	}

	return items[idx].Value, nil
}

func (w *Wizard) Input(ctx context.Context, label, def string) (string, error) {
	p := promptui.Prompt{
		Label:   label,
		Default: def,
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	v, err := p.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			return "", context.Canceled
		}
		return "", err
	}
	return v, nil
}
