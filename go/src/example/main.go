package main

import (
	"github.com/chenfeng/goback"
	//"log"
)

// MainWnd为UI主窗口名称
// 修复c++端访问goback只能在初始化时访问的bug，需要在隐藏窗口内部等待消息，
// 而不是后台服务example.exe循环等待从隐藏窗口推送过来的管道消息

func main() {
	obj := goback.Regist("MainWnd")
	go func() {
		for {
			_, ok := <-obj.BufCh
			if !ok {
				break
			}
		}
		close(obj.BufCh)
	}()
	goback.Wait()
}
