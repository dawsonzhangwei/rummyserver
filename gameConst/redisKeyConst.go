package gameConst

const (
	RedisKey_Player_Prefix string = "player:player_info:"

	REDIS_GAME_CACHE string = "gamecache"
	REDIS_AGC_CACHE string = "agccache"

	/**---------- agc cache start ----------*/

    // 玩家信息缓存  hash类型  拼接玩家id
    PLAYER_PLAYER_INFO string = "player:player_info";

    // 游戏配置缓存  string类型  拼接游戏id
	CONFIG_GAME_CONFIG string = "config:game_config";
	

	/**---------- game cache start ----------*/

    /**
     * game cache g_token
     * hash类型
     * 拼接token
     */
	 GAME_CACHE_G_TOKEN string = "game:cache:g:token:"

	 /**
	  * 房间缓存
	  * hash类型
	  * 拼接游戏ID和房间ID
	  */
	 GAME_CACHE_ROOM_CACHE string = "game:cache:room:cache:"
 
	 /**
	  * 玩家缓存
	  * hash类型
	  * 拼接玩家ID
	  */
	 GAME_CACHE_PLAYER_CACHE string = "game:cache:player:cache:"
)