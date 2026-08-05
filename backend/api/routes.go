package api

var router = map[IPacketType]func(i HandlerInput) error{
	pt_IRoomChange:  roomChange,
	pt_IRoomUsers:   roomUsers,
	pt_IRoomOptions: roomOptions,
	pt_IRoomLeave:   roomLeave,
	pt_IRoomKick:    roomKick,
	// pt_ICreateLobby:    createLobby,
	// pt_IJoinLobby:      joinLobby,
	// pt_ILeaveLobby:     leaveLobby,
	// pt_IUserNameChange: handleUserNameChange,
	// pt_ILobbyUsers:     lobbyUsers,
	// pt_ILobbyInfo:      lobbyInfo,
	// pt_ILobbyChange:    lobbyChange,
	// pt_ILobbyOptions:   lobbyOptions,
	pt_IGameStart:  gameStart,
	pt_IGameReady:  gameReady,
	pt_IGameAction: gameAction,
}
