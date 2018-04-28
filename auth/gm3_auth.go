package auth

var _ Authenticator = (*gm3Authenticator)(nil)

type gm3Authenticator struct {
	authFunc AuthFunc
}

type AuthFunc func(token string) *ClientInfo

func (this *gm3Authenticator) Authenticate(token string) (bool, *ClientInfo) {
	//todo 兼容旧版本，没带token的不需要验证
	if len(token) == 0 {
		return true, nil
	}

	//如果没有设置verifyTokenFunc认为验证通过
	if this.authFunc == nil {
		return true, nil
	}

	//通过token获取信息
	return true, this.authFunc(token)
}

func (this *gm3Authenticator) SetVerifyFunc(f AuthFunc) {
	this.authFunc = f
}

func init() {
	Register("gm3Auth", new(gm3Authenticator))
}
