package process

import (
	"encoding/json"
	"fmt"
	"sggStudyGo/chatroom/common/message"
)

func outputGroupMes(mes *message.Message) { // 这个地方mes一定是SmsMes
	// 反序列化
	var smsMes message.SmsMes

	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &smsMes) err: ", err)
		return
	}

	// 显示信息
	info := fmt.Sprintf("用户id:\t%d 对大家说：\t%s", smsMes.UserId, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}
