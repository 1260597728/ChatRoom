package process

import "fmt"

var userMgr *UserMgr

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 完成对userMgr的初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 增加在线用户
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

// 删除在线用户
func (this *UserMgr) DelOnlineUser(userid int) {
	delete(this.onlineUsers, userid)
}

// 返回当前所有在线用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}


// 根据userid找对应的值
func (this *UserMgr) GetOnlineUser(userid int) (up *UserProcess, err error) {
	up, ok := this.onlineUsers[userid]
	if !ok {
		// 表示当前要查找的用户不在线
		err = fmt.Errorf("%d 用户不在线", userid)
		return
	}

	return
}