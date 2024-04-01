export enum IPacketType {
  Chat = 'server_chat',
  Error = 'server_error',
  LobbyChange = 'server_lobby_change',
  UserInit = 'server_user_init',
  UserChange = 'server_user_change',
  LobbyJoin = 'server_lobby_join',
}

export enum OPacketType {
  Chat = 'client_chat',
  CreateLobby = 'client_lobby_create',
  JoinLobby = 'client_lobby_join',
  LeaveLobby = 'client_lobby_leave',
  LobbyUsers = 'client_lobby_users',
  UserChange = 'client_user_change',
}
