package acl

import (
	"errors"

	"go.uber.org/zap"
)

type GetAuthFunc func(userName, topic string) interface{}

var getAuth GetAuthFunc //callback get authInfo of user

type ClientInfo struct {
	UserName string
}

type Authenticator interface {
	CheckPub(clientInfo *ClientInfo, topic string) bool
	CheckSub(clientInfo *ClientInfo, topic string) bool
	ProcessUnSub(clientInfo *ClientInfo, topic string)
}

var providers = make(map[string]Authenticator)

type TopicAclManger struct {
	p Authenticator
}

func (this *TopicAclManger) CheckPub(clientInfo *ClientInfo, topic string) bool {
	logger.Debug("pub", zap.String("user_name", clientInfo.UserName))
	return this.p.CheckPub(clientInfo, topic)
}

func (this *TopicAclManger) CheckSub(clientInfo *ClientInfo, topic string) bool {
	logger.Debug("sub", zap.String("user_name", clientInfo.UserName))
	return this.p.CheckSub(clientInfo, topic)
}

func (this *TopicAclManger) ProcessUnSub(clientInfo *ClientInfo, topic string) {
	logger.Debug("unSub", zap.String("user_name", clientInfo.UserName))
	this.p.ProcessUnSub(clientInfo, topic)
	return
}

func NewTopicAclManger(providerName string, f GetAuthFunc) (*TopicAclManger, error) {
	getAuth = f
	v, ok := providers[providerName]
	if !ok {
		return nil, errors.New("providers not exist this name:" + providerName)
	}

	return &TopicAclManger{v}, nil
}

func Register(name string, provider Authenticator) {
	if provider == nil {
		panic("auth: Register provide is nil")
	}

	if _, dup := providers[name]; dup {
		panic("auth: Register called twice for provider " + name)
	}

	providers[name] = provider
}

func UnRegister(name string) {
	delete(providers, name)
}
