package api

var router = map[IPacketType]func(i HandlerInput) error{
	pt_IChat:           handleChat,
	pt_ICreateLobby:    createLobby,
	pt_IJoinLobby:      joinLobby,
	pt_ILeaveLobby:     leaveLobby,
	pt_IUserNameChange: handleUserNameChange,
	pt_ILobbyUsers:     lobbyUsers,
	pt_IGameNew:        gameNew,
	pt_IGameReady:      gameReady,
	pt_IGameAction:     gameAction,
}
