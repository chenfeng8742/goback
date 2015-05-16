package main

import (
	"goback"
	//"log"
)

// SecKillWnd为UI主窗口名称

func main() {
	obj := goback.Regist("SecKill")
	for {
		_, ok := <-obj.BufCh
		if !ok {
			break
		}
	}
	close(obj.BufCh)
}
