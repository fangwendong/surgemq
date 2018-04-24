package acl

import (
	"fmt"
	"sync"
)

const (
	TopicAlwaysVerifyType = "topicAlwaysVerify"
	TopicNumAuthType      = "topicNumAuth"
	TopicSetAuthType      = "topicSetAuth"
	userTopicKeyFmt       = "%s:%s"
)

type TopicAlwaysVerify bool

var topicAlwaysVerify TopicAlwaysVerify = true

var _ Authenticator = (*TopicAlwaysVerify)(nil)

func init() {
	Register(TopicAlwaysVerifyType, topicAlwaysVerify)
	Register(TopicNumAuthType, new(topicNumAuth))
	Register(TopicSetAuthType, new(topicSetAuth))
}

func (this TopicAlwaysVerify) CheckPub(clientInfo *ClientInfo, topic string) bool {
	return true

}

func (this TopicAlwaysVerify) CheckSub(clientInfo *ClientInfo, topic string) bool {
	return true

}

func (this TopicAlwaysVerify) ProcessUnSub(clientInfo *ClientInfo, topic string) {
}

type topicNumAuth struct {
	topicTotalNowM sync.Map
	topicUserM     sync.Map
}

var _ Authenticator = (*topicNumAuth)(nil)

func (this *topicNumAuth) CheckPub(clientInfo *ClientInfo, topic string) bool {
	return true
}

func (this *topicNumAuth) CheckSub(clientInfo *ClientInfo, topic string) bool {
	key := fmt.Sprintf(userTopicKeyFmt, clientInfo.UserName, topic)
	if _, ok := this.topicUserM.Load(key); ok {
		return true
	}

	totalLimit, ok := getAuth(clientInfo.UserName, topic).(int)
	if !ok || totalLimit == 0 {
		return false
	}

	totalNow, ok := this.topicTotalNowM.Load(clientInfo.UserName)
	if !ok {
		this.topicTotalNowM.Store(clientInfo.UserName, 1)
		this.topicUserM.Store(key, true)
		return true
	}

	if totalNow.(int) >= totalLimit {
		return false
	}

	this.topicTotalNowM.Store(clientInfo.UserName, totalNow.(int)+1)
	this.topicUserM.Store(key, true)
	return true
}

func (this *topicNumAuth) ProcessUnSub(clientInfo *ClientInfo, topic string) {
	key := fmt.Sprintf(userTopicKeyFmt, clientInfo.UserName, topic)
	if _, ok := this.topicUserM.Load(key); !ok {
		return
	}
	this.topicUserM.Delete(key)
	totalNow, ok := this.topicTotalNowM.Load(clientInfo.UserName)
	if ok {
		this.topicTotalNowM.Store(clientInfo.UserName, totalNow.(int)-1)
	}
}

type topicSetAuth struct {
	topicM sync.Map
}

var _ Authenticator = (*topicSetAuth)(nil)

func (this *topicSetAuth) CheckPub(clientInfo *ClientInfo, topic string) bool {
	return this.CheckSub(clientInfo, topic)
}

func (this *topicSetAuth) CheckSub(clientInfo *ClientInfo, topic string) bool {
	key := fmt.Sprintf(userTopicKeyFmt, clientInfo.UserName, topic)
	if _, ok := this.topicM.Load(key); ok {
		return true
	}

	exists, ok := getAuth(clientInfo.UserName, topic).(bool)
	if !ok {
		return false
	}

	if exists {
		this.topicM.Store(key, true)
	}

	return exists
}

func (this *topicSetAuth) ProcessUnSub(clientInfo *ClientInfo, topic string) {
	return
}
