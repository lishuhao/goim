[ - ] comet 重启后cometServer count 为0，需等待一段时间(几分钟)  
[ - ] 聊天室管理  
[ x ] rest 转 socket
[ x ] 连续创建聊天室会导致comet panic  
[ - ] panic recover  
[ x ] auth token 里的accepts:[1001,1002,1003]是什么意思？
    - ~~https://github.com/Terry-Mao/goim/issues/296#issuecomment-491197075~~
    - 当发送**广播**和**私信**时，op字段如果不在accepts里边则收不到消息  
[ x ] 创建聊天室后，立即获取聊天室列表，获取不到，需要等一会儿左右才能获取到列表  
    因为logic接口获取的聊天室列表是从comet 同步过来的（RenewOnline），
    10秒同步一次同步到logic，logic再存储到redis