export enum IPacketType {
  Chat = 'server_chat',
  Error = 'server_error',
  LobbyChange = 'server_lobby_change',
  UserInit = 'server_user_init',
  UserChange = 'server_user_change',
}

export enum OPacketType {
  Chat = 'client_chat',
  CreateLobby = 'client_create_lobby',
  JoinLobby = 'client_join_lobby',
  LeaveLobby = 'client_leave_lobby',
  UserChange = 'client_user_change',
}
