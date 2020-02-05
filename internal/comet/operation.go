package comet

import (
	"context"
	"encoding/json"
	"time"

	model "github.com/Terry-Mao/goim/api/comet/grpc"
	logic "github.com/Terry-Mao/goim/api/logic/grpc"
	"github.com/Terry-Mao/goim/pkg/strings"
	log "github.com/golang/glog"

	msg "github.com/Terry-Mao/goim/internal/logic/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

// Connect connected a connection.
func (s *Server) Connect(c context.Context, p *model.Proto, cookie string) (mid int64, key, rid string, accepts []int32, heartbeat time.Duration, err error) {
	reply, err := s.rpcClient.Connect(c, &logic.ConnectReq{
		Server: s.serverID,
		Cookie: cookie,
		Token:  p.Body,
	})
	if err != nil {
		return
	}
	return reply.Mid, reply.Key, reply.RoomID, reply.Accepts, time.Duration(reply.Heartbeat), nil
}

// Disconnect disconnected a connection.
func (s *Server) Disconnect(c context.Context, mid int64, key string) (err error) {
	_, err = s.rpcClient.Disconnect(context.Background(), &logic.DisconnectReq{
		Server: s.serverID,
		Mid:    mid,
		Key:    key,
	})
	return
}

// Heartbeat heartbeat a connection session.
func (s *Server) Heartbeat(ctx context.Context, mid int64, key string) (err error) {
	_, err = s.rpcClient.Heartbeat(ctx, &logic.HeartbeatReq{
		Server: s.serverID,
		Mid:    mid,
		Key:    key,
	})
	return
}

// RenewOnline renew room online.
func (s *Server) RenewOnline(ctx context.Context, serverID string, rommCount map[string]int32) (allRoom map[string]int32, err error) {
	reply, err := s.rpcClient.RenewOnline(ctx, &logic.OnlineReq{
		Server:    s.serverID,
		RoomCount: rommCount,
	}, grpc.UseCompressor(gzip.Name))
	if err != nil {
		return
	}
	return reply.AllRoomCount, nil
}

// Receive receive a message.
func (s *Server) Receive(ctx context.Context, mid int64, p *model.Proto) (err error) {
	_, err = s.rpcClient.Receive(ctx, &logic.ReceiveReq{Mid: mid, Proto: p})
	return
}

// Operate operate.
func (s *Server) Operate(ctx context.Context, p *model.Proto, ch *Channel, b *Bucket) error {
	log.Info("operate proto", *p)
	//是否发送Ready 信号。
	sendSignal := false
	//是否需要调用Receive转发
	forwarding := false
	switch p.Op {
	case model.OpChangeRoom:
		if err := b.ChangeRoom(string(p.Body), ch); err != nil {
			log.Errorf("b.ChangeRoom(%s) error(%v)", p.Body, err)
		}
		p.Op = model.OpChangeRoomReply
		log.Info("change room: ", b.rooms, b.RoomsCount())
	case model.OpSub:
		if ops, err := strings.SplitInt32s(string(p.Body), ","); err == nil {
			ch.Watch(ops...)
		}
		p.Op = model.OpSubReply
	case model.OpUnsub:
		if ops, err := strings.SplitInt32s(string(p.Body), ","); err == nil {
			ch.UnWatch(ops...)
		}
		p.Op = model.OpUnsubReply
	case model.OpJoinRoom:
		joinRoom(p, b, ch)
		sendSignal = true
		log.Info("create room rooms:", b.rooms)
	case model.OpLeaveRoom:
		leaveRoom(p, b, ch)
		sendSignal = true
	case model.OpAnyoneCall:
		p.Op = model.OpAnyoneCallReply
		p.Body = nil
		sendSignal = true
		forwarding = true
	case model.OpIncomingCallResp:
		p.Op = model.OpIncomingCallRespReply
		p.Body = nil
		sendSignal = true
		forwarding = true
	default:
		// TODO ack ok&failed
		log.Info("default", p.Op)
		if err := s.Receive(ctx, ch.Mid, p); err != nil {
			log.Errorf("s.Report(%d) op:%d error(%v)", ch.Mid, p.Op, err)
		}
		p.Body = nil
	}

	if sendSignal {
		ch.CliProto.SetAdv()
		ch.Signal()
	}
	if forwarding {
		if err := s.Receive(ctx, ch.Mid, p); err != nil {
			log.Errorf("s.Report(%d) op:%d error(%v)", ch.Mid, p.Op, err)
		}
	}
	return nil
}

func joinRoom(p *model.Proto, b *Bucket, ch *Channel) {
	req := msg.JoinRoomReq{}
	err := json.Unmarshal(p.Body, &req)
	if err != nil {
		log.Error("Unmarshal", err)
		return
	}
	if req.RoomID == "" {
		log.Error("create room id empty")
		return
	}
	if err := b.ChangeRoom(req.RoomID, ch); err != nil {
		log.Errorf("create room error(%v) body(%v)", err, p.Body)
		return
	}
	p.Op = model.OpJoinRoomReply
	body := msg.JoinRoomReply{
		ID:             "joinRoomResponse",
		MasterID:       b.Room(req.RoomID).MasterId(),
		OnlineUserList: b.Room(req.RoomID).Users(),
	}
	p.Body = body.ToBytes()

	return
}

func leaveRoom(p *model.Proto, b *Bucket, ch *Channel) {
	if err := b.ChangeRoom("", ch); err != nil {
		log.Errorf("create room error(%v) body(%v)", err, p.Body)
		return
	}
	p.Op = model.OpLeaveRoomReply
	p.Body = nil

	return
}
