package utils

import "time"

var (
	// jwt过期时间设定
	JwtExpiresTime = time.Minute * 60
	// 验证码过期时间设定
	SessionMAXAge = 60
)
