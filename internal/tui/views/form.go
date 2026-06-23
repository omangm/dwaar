package views

import (
	"errors"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/omangm/dwaar/internal/tunnel"
)

type RuleFormModel struct {
	Form       *huh.Form
	Name       string
	LocalPort  string
	RemoteHost string
	RemotePort string
}

func NewRuleForm(rule *tunnel.ForwardRule) *RuleFormModel {
	m := &RuleFormModel{}

	if rule != nil {
		m.Name = rule.Name
		m.LocalPort = strconv.Itoa(rule.LocalPort)
		m.RemoteHost = rule.RemoteHost
		m.RemotePort = strconv.Itoa(rule.RemotePort)
	}

	m.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Rule Name").
				Value(&m.Name).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Local Port").
				Value(&m.LocalPort).
				Validate(validatePort),
			huh.NewInput().
				Title("Remote Host").
				Value(&m.RemoteHost).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Remote Port").
				Value(&m.RemotePort).
				Validate(validatePort),
		),
	)

	return m
}

func (m *RuleFormModel) ToRule() tunnel.ForwardRule {
	lp, _ := strconv.Atoi(m.LocalPort)
	rp, _ := strconv.Atoi(m.RemotePort)
	return tunnel.ForwardRule{
		Name:       m.Name,
		LocalPort:  lp,
		RemoteHost: m.RemoteHost,
		RemotePort: rp,
		Protocol:   tunnel.ProtoTCP,
	}
}

func validatePort(s string) error {
	p, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if p < 1 || p > 65535 {
		return errors.New("port out of range")
	}
	return nil
}
