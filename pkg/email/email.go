package email

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendCode(targetEmail string) (string, error) {
	// 简单设置 log 参数
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	em.From = "LDJ <2271115019@qq.com>"

	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	em.To = []string{targetEmail}

	// 设置主题
	em.Subject = "验证码"

	// 简单设置文件发送的内容，暂时设置成纯文本
	em.Text = []byte("验证码是：" + code)

	//设置服务器相关的配置
	err := em.Send("smtp.qq.com:465", smtp.PlainAuth("", "2271115019@qq.com", "eamrjnyyuckadjhi", "smtp.qq.com"))
	if err != nil {
		// log.Printf("发送失败，详细错误: %+v", err)
		return "", err
	}
	return code, nil
}
