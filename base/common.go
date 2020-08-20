package base

type Player struct {
	Clover string `json:"clover"`
	Life string `json:"life"`
	Coin string `json:"coin"`
	Ticket string `json:"ticket"`
	Level string `json:"level"`
	Check_point string `json:"check_point"`
	PlayerId string `json:"playerId"`
	PlayerName string `json:"playerName"`
	Current_game_server_id string `json:"current_game_server_id"`
	Current_room_id string `json:"current_room_id"`
	Current_game_id string `json:"current_game_id"`
	Current_field_id string `json:"current_field_id"`
	Is_admin string `json:"is_admin"`
	PetId string `json:"petId"`
	Current_gate_id string `json:"current_gate_id"`
	Current_match_status string `json:"current_match_status"`
	Play_num string `json:"play_num"`
	Win_num string `json:"win_num"`
	Item_ddz string `json:"item_ddz"`
	Played_player string `json:"played_player"`
}