package grpc

const (
	OpPushKeys = int32(2001)
	OpPushKeysReply = int32(2002)

	OpPushRoom = int32(2003)
	OpPushRoomReply = int32(2004)

	OpBroadcast = int32(2005)
	OpBroadcastReply = int32(2006)

	//加入房间
	OpJoinRoom      = int32(10001)
	OpJoinRoomReply = int32(10002)

	//离开房间
	OpLeaveRoom = int32(10003)

	//申请连麦
	OpAnyoneCall = int32(10005)

	//主播回复连麦申请
	OpIncomingCallResp = int32(10007)
)
