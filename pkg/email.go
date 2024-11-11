package pkg

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/go-mail/mail"
	"log"
	"os"
	"time"
)

type SMTP struct {
	Config *SMTPConfig
	ch     chan *mail.Message
	chOpen bool
	done   chan bool
}

type SMTPConfig struct {
	Name       string `toml:"name"`       // 发送者名
	Address    string `toml:"address"`    // 发送者地址
	ReplyTo    string `toml:"replyTo"`    // 回复地址
	Host       string `toml:"host"`       // 服务器主机名
	Port       int    `toml:"port"`       // 服务器端口
	User       string `toml:"user"`       // 用户名
	Password   string `toml:"password"`   // 密码
	Encryption bool   `toml:"encryption"` // 是否启用加密
	Keepalive  int    `toml:"keepalive"`  // SMTP连接保留时长
}

func NewSMTPClient(config *SMTPConfig) *SMTP {
	client := &SMTP{
		Config: config,
		ch:     make(chan *mail.Message, 30),
		chOpen: false,
		done:   make(chan bool),
	}
	go client.eventLoop()
	return client
}

func (client *SMTP) eventLoop() {
	var s mail.SendCloser
	var err error
	open := false
	client.chOpen = true
	for {
		d := mail.NewDialer(client.Config.Host, client.Config.Port, client.Config.User, client.Config.Password)
		d.Timeout = time.Duration(client.Config.Keepalive+5) * time.Second
		d.SSL = false
		if client.Config.Encryption {
			d.SSL = true
		}
		d.StartTLSPolicy = mail.OpportunisticStartTLS

		select {
		case msg, ok := <-client.ch:
			if !ok {
				log.Println("[ERROR] mail queue closed")
				client.chOpen = false
				return
			}
			if s, err = d.Dial(); err != nil {
				panic(err)
			}
			if err := mail.Send(s, msg); err != nil {
				log.Printf("[ERROR] send file failed, %s\n", err)
			} else {
				log.Println("[INFO] send file succeeded")
			}
		case <-time.After(time.Duration(client.Config.Keepalive) * time.Second):
			if open {
				if err := s.Close(); err != nil {
					log.Printf("[ERROR] cannot close SMTP connection, %v", err)
				}
				open = false
			}
			close(client.done)
		}
	}
}

// Send 发送邮件
func (client *SMTP) Send(to, title, body, filename string) error {
	if !client.chOpen {
		return errors.New("mail queue is still not open")
	}
	m := mail.NewMessage()
	m.SetAddressHeader("From", client.Config.Address, client.Config.Name)
	m.SetAddressHeader("Reply-To", client.Config.ReplyTo, client.Config.Name)
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)
	m.Attach(filename)
	client.ch <- m
	return nil
}

func (client *SMTP) Close() {
	if client.ch != nil {
		close(client.ch)
	}
}

func (client *SMTP) Done() chan bool {
	return client.done
}

func NewSMTPConfig(name string) *SMTPConfig {
	config := &SMTPConfig{}
	config.init(name)
	return config
}

func (config *SMTPConfig) init(name string) {
	content, err := os.ReadFile(name)
	if err != nil {
		log.Fatalf("[ERROR] cannot load config file, %v\n", err)
	}
	err = toml.Unmarshal(content, config)
	if err != nil {
		log.Fatalf("[ERROR] unmarshal config file failed, %v\n", err)
	}
}
