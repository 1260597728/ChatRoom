package model

import (
	"encoding/json"
	"fmt"
	"sggStudyGo/chatroom/common/message"

	"github.com/garyburd/redigo/redis"
)

var MyUserDao *UserDao

// 完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

// 创建UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	return &UserDao{
		pool: pool,
	}
}

// 根据用户id查找用户
func (this *UserDao) GetUserById(conn redis.Conn, id int) (user *User, err error) {
	// 通过给定id去redis查询这个用户
	res, err := redis.String(conn.Do("Hget", "users", id))
	if err != nil {
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	user = &User{}

	// 存在数据
	// json反序列化
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("redis data json.Unmarshal err: ", err)
		return
	}

	return
}

// 登陆校验
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	// 从redis连接池中获取一个连接
	conn := this.pool.Get()
	defer conn.Close()

	// 通过userId获取一个User
	user, err = this.GetUserById(conn, userId)
	if err != nil {
		return
	}

	// 获取到一个User，比对下密码
	if userPwd != user.UserPwd {
		err = ERROR_USER_PWD
		return
	}

	return
}

// 注册
func (this *UserDao) Register(user *message.User) (err error) {
	// 从redis连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()

	// 通过userid获取用户
	_, err = this.GetUserById(conn, user.UserId)
	if err == nil {
		// 表示用户已经存在
		err = ERROR_USER_EXISTS
		return
	}

	// 说明用户不存在
	// 注册用户进入redis
	// 将user结构体进行序列化
	userByte, err := json.Marshal(user)
	if err != nil {
		return
	}
	_, err = conn.Do("Hset", "users", user.UserId, string(userByte))
	if err != nil {
		fmt.Println("注册用户失败：", err)
		return
	}

	return
}