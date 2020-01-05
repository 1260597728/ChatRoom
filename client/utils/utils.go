package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"sggStudyGo/chatroom/common/message"
)

type Transfer struct {
	Conn net.Conn   // 连接句柄
	Buf  [8096]byte // 传输时需要用到的缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("读取客户端发送的数据……")
	// conn.Read 在conn没有被关闭的情况下，才会阻塞
	// 如果客户端关闭了，conn就不会阻塞
	// 先接受客户端传过来的数据长度，数据长度是切片类型
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}

	// 将数据长度切片类型 buf[:4] 转换成 uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])
	fmt.Printf("%d\n", pkgLen)

	// 根据pkgLen读取消息内容
	// 从客户端接收传过来的内容数据，根据前面接收的内容长度获取
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		return
	}

	// json 反序列化
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err: ", err)
		return
	}

	return
}

// 用户将数据发送
func (this *Transfer) WritePkg(data []byte) (err error) {
	// 将数据发送给对方
	// 先发送数据的长度
	pkgLen := uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	// 发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) err: ", err)
		return
	}

	// 发送 data 本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(data) fail: ", err)
		return
	}

	return
}
