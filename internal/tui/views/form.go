package views

import (
	"errors"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/omangm/dwaar/internal/tunnel"
)

type RuleFormModel struct {
	Form *huh.Form

	Name       string
	LocalPort  string
	RemoteHost string
	RemotePort string
	AutoStart  bool
	UseSSH     bool
	SSHHost    string
	SSHPort    string
	SSHUser    string
	SSHKey     string
	SSHPass    string
}

func NewRuleForm(rule *tunnel.ForwardRule) *RuleFormModel {
	m := &RuleFormModel{SSHPort: "22"}

	if rule != nil {
		m.Name = rule.Name
		m.LocalPort = strconv.Itoa(rule.LocalPort)
		m.RemoteHost = rule.RemoteHost
		m.RemotePort = strconv.Itoa(rule.RemotePort)
		m.AutoStart = rule.AutoStart
		if rule.Via != nil {
			m.UseSSH = true
			m.SSHHost = rule.Via.Host
			m.SSHUser = rule.Via.User
			m.SSHKey = rule.Via.IdentityFile
			m.SSHPass = rule.Via.Password
			if rule.Via.Port != 0 {
				m.SSHPort = strconv.Itoa(rule.Via.Port)
			}
		}
	}

	m.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Rule Name").
				Value(&m.Name).
				Validate(required),
			huh.NewInput().
				Title("Local Port").
				Value(&m.LocalPort).
				Validate(validatePort),
			huh.NewInput().
				Title("Remote Host").
				Value(&m.RemoteHost).
				Validate(required),
			huh.NewInput().
				Title("Remote Port").
				Value(&m.RemotePort).
				Validate(validatePort),
			huh.NewConfirm().
				Title("Start automatically on launch?").
				Value(&m.AutoStart),
			huh.NewConfirm().
				Title("Forward over SSH?").
				Value(&m.UseSSH),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("SSH Host").
				Value(&m.SSHHost).
				Validate(m.requiredIfSSH),
			huh.NewInput().
				Title("SSH Port").
				Value(&m.SSHPort).
				Validate(m.portIfSSH),
			huh.NewInput().
				Title("SSH User").
				Value(&m.SSHUser).
				Validate(m.requiredIfSSH),
			huh.NewInput().
				Title("Identity File (optional, e.g. ~/.ssh/id_ed25519)").
				Value(&m.SSHKey),
			huh.NewInput().
				Title("Password (optional)").
				Value(&m.SSHPass).
				EchoMode(huh.EchoModePassword),
		).WithHideFunc(func() bool { return !m.UseSSH }),
	)

	return m
}

func (m *RuleFormModel) ToRule() tunnel.ForwardRule {
	lp, _ := strconv.Atoi(m.LocalPort)
	rp, _ := strconv.Atoi(m.RemotePort)
	rule := tunnel.ForwardRule{
		Name:       m.Name,
		LocalPort:  lp,
		RemoteHost: m.RemoteHost,
		RemotePort: rp,
		Protocol:   tunnel.ProtoTCP,
		AutoStart:  m.AutoStart,
	}
	if m.UseSSH {
		sp, _ := strconv.Atoi(m.SSHPort)
		rule.Via = &tunnel.SSHHost{
			Host:         m.SSHHost,
			Port:         sp,
			User:         m.SSHUser,
			IdentityFile: m.SSHKey,
			Password:     m.SSHPass,
		}
	}
	return rule
}

func required(s string) error {
	if s == "" {
		return errors.New("cannot be empty")
	}
	return nil
}

func (m *RuleFormModel) requiredIfSSH(s string) error {
	if !m.UseSSH {
		return nil
	}
	return required(s)
}

func (m *RuleFormModel) portIfSSH(s string) error {
	if !m.UseSSH {
		return nil
	}
	return validatePort(s)
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
