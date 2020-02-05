package logic

import (
	"context"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"time"

	"github.com/Terry-Mao/goim/api/comet/grpc"
	"github.com/Terry-Mao/goim/internal/logic/model"
	log "github.com/golang/glog"
	"github.com/google/uuid"
)

// Connect connected a conn.
func (l *Logic) Connect(c context.Context, server, cookie string, token []byte) (mid int64, key, roomID string, accepts []int32, hb int64, err error) {
	log.Info("connect token:", token)
	var params struct {
		Mid      int64   `json:"mid"`
		Key      string  `json:"key"`
		RoomID   string  `json:"room_id"`
		Platform string  `json:"platform"`
		Accepts  []int32 `json:"accepts"`
	}
	if err = json.Unmarshal(token, &params); err != nil {
		log.Errorf("json.Unmarshal(%s) error(%v)", token, err)
		return
	}
	mid = params.Mid
	roomID = params.RoomID
	accepts = params.Accepts
	hb = int64(l.c.Node.Heartbeat) * int64(l.c.Node.HeartbeatMax)
	if key = params.Key; key == "" {
		// CUSTOM
		// 这里可以返回应用内用户唯一id，之后可以按key（即用户id）发送消息
		key = uuid.New().String()
	}
	if err = l.dao.AddMapping(c, mid, key, server); err != nil {
		log.Errorf("l.dao.AddMapping(%d,%s,%s) error(%v)", mid, key, server, err)
	}
	log.Infof("conn connected key:%s server:%s mid:%d token:%s", key, server, mid, token)
	return
}

// Disconnect disconnect a conn.
func (l *Logic) Disconnect(c context.Context, mid int64, key, server string) (has bool, err error) {
	if has, err = l.dao.DelMapping(c, mid, key, server); err != nil {
		log.Errorf("l.dao.DelMapping(%d,%s) error(%v)", mid, key, server)
		return
	}
	log.Infof("conn disconnected key:%s server:%s mid:%d", key, server, mid)
	return
}

// Heartbeat heartbeat a conn.
func (l *Logic) Heartbeat(c context.Context, mid int64, key, server string) (err error) {
	has, err := l.dao.ExpireMapping(c, mid, key)
	if err != nil {
		log.Errorf("l.dao.ExpireMapping(%d,%s,%s) error(%v)", mid, key, server, err)
		return
	}
	if !has {
		if err = l.dao.AddMapping(c, mid, key, server); err != nil {
			log.Errorf("l.dao.AddMapping(%d,%s,%s) error(%v)", mid, key, server, err)
			return
		}
	}
	log.Infof("conn heartbeat key:%s server:%s mid:%d", key, server, mid)
	return
}

// RenewOnline renew a server online.
func (l *Logic) RenewOnline(c context.Context, server string, roomCount map[string]int32) (map[string]int32, error) {
	online := &model.Online{
		Server:    server,
		RoomCount: roomCount,
		Updated:   time.Now().Unix(),
	}
	if err := l.dao.AddServerOnline(context.Background(), server, online); err != nil {
		return nil, err
	}
	return l.roomCount, nil
}

// Receive receive a message.
func (l *Logic) Receive(c context.Context, mid int64, proto *grpc.Proto) (err error) {
	//CUSTOM
	// TODO  根据 proto.Op 处理业务消息
	log.Info("receive ", mid, proto)

	switch proto.Op {
	case grpc.OpPushMessage:
		req := model.PushMessageReq{}
		err = json.Unmarshal(proto.Body, &req)
		if err != nil {
			break
		}
		switch req.ConversationType {
		case grpc.ConversationTypePrivate:
			err = l.PushKeys(c, proto.Op, []string{req.TargetId}, req.ToPushToClient().ToBytes())
		case grpc.ConversationTypeChatRoom:
			roomType, roomId, err := model.DecodeRoomKey(req.TargetId)
			if err != nil {
				log.Error("room id format err", req.TargetId)
			}
			err = l.PushRoom(c, proto.Op, roomType, roomId, req.ToPushToClient().ToBytes())
		}
		/*	case grpc.OpBroadcast:
			err = l.PushAll(c, proto.Op, 0, proto.Body)*/
	case grpc.OpAnyoneCall:
		err = linkMikeApply(l, c, proto)
	case grpc.OpIncomingCallResp:
		err = linkMikeResp(l, c, proto)
	default:
		log.Info("unknown op ", mid, proto)
		err = l.PushAll(c, proto.Op, 0, proto.Body)
	}
	if err != nil {
		log.Error("push err" + err.Error())
	}

	log.Infof("receive mid:%d message:%+v", mid, proto)
	return
}

//连麦申请
func linkMikeApply(l *Logic, c context.Context, proto *grpc.Proto) error {
	// 转发给媒体服务器
	err := l.PushKeys(c, proto.Op, []string{model.UserId3}, proto.Body)
	// 发送给主播
	toId := jsoniter.Get(proto.Body, "toId").ToString()
	err = l.PushKeys(c, proto.Op, []string{toId}, proto.Body)

	return err
}

// 主播连麦回复
func linkMikeResp(l *Logic, c context.Context, proto *grpc.Proto) error {
	toId := jsoniter.Get(proto.Body, "toId").ToString()
	//发送给申请连麦的用户
	err := l.PushKeys(c, proto.Op, []string{toId}, proto.Body)
	return err
}
