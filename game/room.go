package game

import (
	"sync"
)

type GameRoom struct {
	sync.Mutex

	freeTables [] *GameTable
	playingTables [] *GameTable
}

func (r *GameRoom) Sitdown() {
	r.Lock()
	defer r.Unlock()

	if len(r.freeTables) == 0 {
		r.freeTables = append(r.freeTables, new(GameTable))
	}

	gt := r.freeTables[0]
}