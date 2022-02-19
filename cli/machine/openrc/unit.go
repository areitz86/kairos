package openrc

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/c3os-io/c3os/cli/utils"
)

type ServiceUnit struct {
	content string
	name    string
	rootdir string
}

type ServiceOpts func(*ServiceUnit) error

func WithRoot(n string) ServiceOpts {
	return func(su *ServiceUnit) error {
		su.rootdir = n
		return nil
	}
}

func WithName(n string) ServiceOpts {
	return func(su *ServiceUnit) error {
		su.name = n
		return nil
	}
}

func WithUnitContent(n string) ServiceOpts {
	return func(su *ServiceUnit) error {
		su.content = n
		return nil
	}
}

func NewService(opts ...ServiceOpts) (ServiceUnit, error) {
	s := &ServiceUnit{}
	for _, o := range opts {
		if err := o(s); err != nil {
			return *s, err
		}
	}
	return *s, nil
}

func (s ServiceUnit) WriteUnit() error {
	uname := s.name

	if err := ioutil.WriteFile(filepath.Join(s.rootdir, fmt.Sprintf("/etc/init.d/%s", uname)), []byte(s.content), 0755); err != nil {
		return err
	}

	return nil
}

func (s ServiceUnit) OverrideCmd(cmd string) error {

	svcDir := filepath.Join(s.rootdir, fmt.Sprintf("/etc/init.d/%s", s.name))

	d, err := ioutil.ReadFile(svcDir)
	if err != nil {
		return err
	}

	ss := strings.ReplaceAll(string(d), "command_args=\"agent \\", fmt.Sprintf("command_args=\"%s \\", cmd))
	ss = strings.ReplaceAll(ss, "command_args=\"server \\", fmt.Sprintf("command_args=\"%s \\", cmd))

	return ioutil.WriteFile(svcDir, []byte(ss), 0600)
}

func (s ServiceUnit) Start() error {
	_, err := utils.SH(fmt.Sprintf("/etc/init.d/%s start", s.name))
	return err
}

func (s ServiceUnit) Restart() error {
	_, err := utils.SH(fmt.Sprintf("/etc/init.d/%s restart", s.name))
	return err
}

func (s ServiceUnit) Enable() error {
	_, err := utils.SH(fmt.Sprintf("ln -sf /etc/init.d/%s /etc/runlevels/default/%s", s.name, s.name))
	return err
}

func (s ServiceUnit) StartBlocking() error {
	return s.Start()
}

func (s ServiceUnit) SetEnvFile(es string) error {
	svcDir := filepath.Join(s.rootdir, fmt.Sprintf("/etc/init.d/%s", s.name))

	d, err := ioutil.ReadFile(svcDir)
	if err != nil {
		return err
	}

	ss := string(d) + "\nsource " + es + "\n"

	return ioutil.WriteFile(svcDir, []byte(ss), 0600)
}
