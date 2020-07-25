package exchange

import (
	"errors"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pinmonl-go"
	"github.com/pinmonl/pinmonl/store"
)

type Manager struct {
	addr    string
	configs *store.Configs
	uclient *pinmonl.Client
	mclient *pinmonl.Client
}

func NewManager(configs *store.Configs, serverAddr string) (*Manager, error) {
	m := Manager{
		configs: configs,
		addr:    serverAddr,
	}
	m.setUserClient(configs.GetUserToken())
	m.setMachineClient(configs.GetMachineToken())
	return &m, nil
}

func (m *Manager) setUserClient(token string) {
	m.uclient = newPMClient(m.addr, token)
}

func (m *Manager) setMachineClient(token string) {
	m.mclient = newPMClient(m.addr, token)
}

func (m *Manager) HasUser() bool {
	if m.configs.GetUserToken() == "" {
		return false
	}
	return true
}

func (m *Manager) HasMachine() bool {
	if m.configs.GetMachineToken() == "" {
		return false
	}
	if time.Now().After(m.configs.GetMachineExpireAt()) {
		return false
	}
	return true
}

func (m *Manager) LinkUser(user *model.User) error {
	m.configs.SetUserLinkedUserID(user.ID)
	return m.configs.Save()
}

func (m *Manager) Signup(login, password, name string) error {
	user := &pinmonl.User{Login: login, Password: password, Name: name}
	token, err := m.uclient.Signup(user)
	if err != nil {
		return err
	}
	m.configs.SetUserToken(token.Token)
	m.configs.SetUserLogin(login)
	m.configs.SetUserExpireAt(token.ExpireAt)
	m.setUserClient(token.Token)
	return m.configs.Save()
}

func (m *Manager) Login(login, password string) error {
	user := &pinmonl.User{Login: login, Password: password}
	token, err := m.uclient.Login(user)
	if err != nil {
		return err
	}
	m.configs.SetUserToken(token.Token)
	m.configs.SetUserLogin(login)
	m.configs.SetUserExpireAt(token.ExpireAt)
	m.setUserClient(token.Token)
	return m.configs.Save()
}

func (m *Manager) LoginMe(password string) error {
	login := m.configs.GetUserLogin()
	if login == "" {
		return errors.New("no login name is saved")
	}
	return m.Login(login, password)
}

func (m *Manager) Alive() error {
	token, err := m.uclient.Alive()
	if err != nil {
		return err
	}
	m.configs.SetUserToken(token.Token)
	m.configs.SetUserExpireAt(token.ExpireAt)
	m.setMachineClient(token.Token)
	return m.configs.Save()
}

func (m *Manager) MachineSignup() error {
	token, err := m.mclient.MachineSignup()
	if err != nil {
		return err
	}
	m.configs.SetMachineToken(token.Token)
	m.configs.SetMachineExpireAt(token.ExpireAt)
	m.setMachineClient(token.Token)
	return m.configs.Save()
}

func (m *Manager) MachineAlive() error {
	token, err := m.mclient.Alive()
	if err != nil {
		return err
	}
	m.configs.SetMachineToken(token.Token)
	m.configs.SetMachineExpireAt(token.ExpireAt)
	m.setMachineClient(token.Token)
	return m.configs.Save()
}

func (m *Manager) UserClient() *pinmonl.Client {
	return m.uclient
}

func (m *Manager) MachineClient() *pinmonl.Client {
	return m.mclient
}
