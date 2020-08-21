package channel

import (
	"rummy/base"
)

type Paw struct {
}

func (p *Paw) configure() error {
	return nil
}

func (p *Paw) Init() error {
	return nil
}

func (p *Paw) Login(loginParam map[string] string) (player base.Player, err error) {
	return player, nil
}