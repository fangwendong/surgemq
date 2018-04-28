package acl

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type topicSetAuth struct {
	topicM sync.Map
	f      GetAuthFunc
}

var _ Authenticator = (*topicSetAuth)(nil)

func (this *topicSetAuth) CheckPub(clientInfo *ClientInfo, topic string) bool {
	return this.CheckSub(clientInfo, topic)
}

func (this *topicSetAuth) CheckSub(clientInfo *ClientInfo, topic string) (success bool) {
	defer func() {
		Logger.Debug("[sub]", zap.String("userId", clientInfo.UserId), zap.Bool("success", success))
	}()

	token := clientInfo.Token
	key := fmt.Sprintf(userTopicKeyFmt, token, topic)
	if _, ok := this.topicM.Load(key); ok {
		success = true
		return
	}

	var ok bool
	success, ok = this.f(token, topic).(bool)
	if !ok {
		success = false
		return
	}

	if success {
		this.topicM.Store(key, true)
	}

	return
}

func (this *topicSetAuth) ProcessUnSub(clientInfo *ClientInfo, topic string) {
	Logger.Debug("[unSub]", zap.String("userId", clientInfo.UserId))
	return
}

func (this *topicSetAuth) SetAuthFunc(f GetAuthFunc) {
	this.f = f
}
