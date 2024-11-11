package pkg

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"time"
)

// Transporter 运输者
type Transporter struct {
	client *SMTP
}

func NewTransporter(name string) *Transporter {
	transporter := &Transporter{}
	transporter.init(name)
	return transporter
}

func (t *Transporter) init(name string) {
	config := NewSMTPConfig(name)
	t.client = NewSMTPClient(config)
}

// Transmit 转发
func (t *Transporter) Transmit(path, name, to, prefix, body string, size uint64) {
	if path == "" && name == "" {
		log.Fatalln("[ERROR] path and name can neither be empty")
	}
	var fullnames []string
	if path != "" && Exists(path) {
		err := filepath.Walk(path, func(innerpath string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}
			if innerpath == path {
				return nil
			}

			if info.IsDir() {
				fmt.Printf("[INFO] skip path %v\n", info.Name())
				return filepath.SkipDir
			}

			if info.Size() > int64(size) {
				return errors.New(fmt.Sprintf("the size of file %v exceeds the range %v", info.Name(), size))
			}
			fullnames = append(fullnames, innerpath)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] walk the path %v failed, %v\n", path, err)
			return
		}
	}

	if name != "" && Exists(name) {
		fullnames = append(fullnames, name)
	}

	for _, fullname := range fullnames {
		_, filename := filepath.Split(fullname)
		err := t.client.Send(to, prefix+filename, filename, fullname)
		if err != nil {
			log.Printf("[ERROR] file %v send into queue failed\n", fullname)
		}
		log.Printf("[INFO] file %v send into queue succeed\n", fullname)
		time.Sleep(1 * time.Second)
	}

	// 等待邮件发送完成
	for {
		select {
		case <-t.client.Done():
			log.Println("[INFO] send all mails")
			return
		}
	}
}
