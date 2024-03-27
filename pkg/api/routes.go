package api

var router = map[IPacketType]func(i HandlerInput) error{
	pt_IChat:        handleChat,
	pt_ICreateLobby: createLobby,
	pt_IJoinLobby:   joinLobby,
	pt_ILeaveLobby:  leaveLobby,
	pt_IUserChange:  handleUserChange,
}
