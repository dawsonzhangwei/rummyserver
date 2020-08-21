package channel

import (
	"rummy/base"
)

type ThirdChannel interface {
	Login(map[string] string) (base.Player, error)
}