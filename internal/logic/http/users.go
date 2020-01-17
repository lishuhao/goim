package http

import (
	"fmt"
	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

//获取GoIM auth时需要的token
func (s *Server) getGoIMToken(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	if mid, ok := model.Uid2Mid[userId]; ok {
		token := fmt.Sprintf(`{"mid":%d,"key":"%s","room_id":"","platform":"web","accepts":[%s]}`, mid, userId, watchOp())
		result(ctx, token, OK)
	} else {
		result(ctx, "", RequestErr)
	}
}

//当发送**广播**和**私信**时，op字段如果不在accepts里边则收不到消息
//
func watchOp() string {
	ops := make([]string, 0)
	ops = append(ops, opsInterval(2001, 2005)...)
	ops = append(ops, opsInterval(10001, 10015)...)
	return strings.Join(ops, ",")
}

//operation 区间
//start <= end
func opsInterval(start, end int) []string {
	if start > end {
		return nil
	}
	ops := make([]string, 0, end-start+1)
	for i := start; i <= end; i++ {
		ops = append(ops, strconv.Itoa(i))
	}
	return ops
}
