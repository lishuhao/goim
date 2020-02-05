package grpc

const (
	OpPushMessage      = int32(2001)
	OpPushMessageReply = int32(2002)

	OpPushRoom      = int32(2003)
	OpPushRoomReply = int32(2004)

	OpBroadcast      = int32(2005)
	OpBroadcastReply = int32(2006)

	//销毁房间
	OpDestroyRoom = int32(2007)

	//创建房间
	OpJoinRoom      = int32(10001)
	OpJoinRoomReply = int32(10002)

	//离开房间
	OpLeaveRoom      = int32(10003)
	OpLeaveRoomReply = int32(10004)

	//申请连麦
	OpAnyoneCall      = int32(10005)
	OpAnyoneCallReply = int32(10006)

	//主播回复连麦申请
	OpIncomingCallResp      = int32(10007)
	OpIncomingCallRespReply = int32(10008)
)

//会话类型
const (
	ConversationTypePrivate  = 1
	ConversationTypeChatRoom = 2
	ConversationTypeSystem   = 3
)
