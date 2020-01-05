package message

// 定义常量类型
const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMesType"
)

// 定义用户状态
const (
	UserOffline = iota
	UserOnline
	UserBusy
)

// 消息结构体
type Message struct {
	Type string `json:"type"` // 消息类型
	Data string `json:"data"` // 消息内容
}

// 登陆所需信息结构体
type LoginMes struct {
	UserId   int    `json:"user_id"`   // 用户id
	UserPwd  string `json:"user_pwd"`  // 用户密码
	UserName string `json:"user_name"` // 用户名
}

// 登陆成功返回信息结构体
type LoginResMes struct {
	Code    int    `json:"code"`     // 状态码
	Error   string `json:"error"`    // 返回错误信息
	UserIds []int  `json:"user_ids"` // 增加字段，保存用户id的切片
}

type RegisterMes struct {
	User User `json:"user"` // user 结构体
}

type RegisterResMes struct {
	Code  int    `json:"code"`  // 状态码
	Error string `json:"error"` // 返回错误信息
}

// 为了配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"user_id"` // 用户id
	Status int `json:"status"`  // 用户的状态
}

// 发送信息结构体
type SmsMes struct {
	Content string `json:"content"` // 内容
	User           // 匿名User结构体
}
