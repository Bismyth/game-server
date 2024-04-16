export enum IPacketType {
  Error = 'server_error',
  RoomInfo = 'server_room_info',
  RoomUserChange = 'server_room_user_change',
  GameEvent = 'server_game_event',
  GameState = 'server_game_state',
  GameAction = 'server_game_action',
  GameError = 'server_error_game',
  RoomKick = 'server_room_kick',
}

export enum OPacketType {
  LeaveRoom = 'client_room_leave',
  RoomUsers = 'client_room_users',
  RoomChange = 'client_room_change',
  RoomOptions = 'client_room_options',
  RoomKick = 'client_room_kick',
  UserNameChange = 'client_user_name_change',
  GameAction = 'client_game_action',
  GameStart = 'client_game_start',
  GameReady = 'client_game_ready',
}
