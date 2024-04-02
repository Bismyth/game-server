export enum IPacketType {
  Chat = 'server_chat',
  Error = 'server_error',
  LobbyChange = 'server_lobby_change',
  UserInit = 'server_user_init',
  UserChange = 'server_user_change',
  LobbyJoin = 'server_lobby_join',
  GameEvent = 'server_game_event',
  GameState = 'server_game_state',
  GameAction = 'server_game_action',
}

export enum OPacketType {
  Chat = 'client_chat',
  CreateLobby = 'client_lobby_create',
  JoinLobby = 'client_lobby_join',
  LeaveLobby = 'client_lobby_leave',
  LobbyUsers = 'client_lobby_users',
  UserChange = 'client_user_change',
  GameAction = 'client_game_action',
  GameNew = 'client_game_new',
  GameReady = 'client_game_ready',
}
