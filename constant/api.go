package constant

type APICodeType uint32

// 状态数字参考了 http 状态码，意义类似
const (
	// 空响应，没有设置响应处理
	API_EMPTY APICodeType = 0

	// 正常响应
	API_OK APICodeType = 200

	// 正常响应，但是没有找到对应的数据
	API_NO_DATA APICodeType = 204

	// 常规错误
	API_ERROR APICodeType = 400

	// 用户登录令牌无效，含过期
	API_TOKEN_INVALID APICodeType = 401

	// 签名无效
	API_SIGN_INVALID APICodeType = 402

	// AccessToken 无效
	API_ACCESS_TOKEN_INVALID APICodeType = 407

	// 限制请求，一般是用于请求限流
	API_REJECT APICodeType = 429

	// 响应异常
	API_CRASH APICodeType = 500
)
