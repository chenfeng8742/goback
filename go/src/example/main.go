package main

import (
	"goback"
	//"log"
)

// MainWnd为UI主窗口名称

func main() {
	obj := goback.Regist("MainWnd")
	for {
		_, ok := <-obj.BufCh
		if !ok {
			break
		}
	}
	close(obj.BufCh)
}