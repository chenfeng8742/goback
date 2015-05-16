# goback
基于WM_COPYDATA，实现go语言和c++之间的双向通讯，可以很方便在windows平台下用golang开发后台服务，c++负责UI界面

目前支持平台：windows

使用方式: 
```
// MainWnd为UI主窗口名称 
obj := goback.Regist("MainWnd") 
for {
     _, ok := <-obj.BufCh 
    if !ok { 
            break 
         } 
    }
close(obj.BufCh) 
// end
```

存在的问题:WM_COPYDATA通讯是同步阻塞的，只能通过SendMessage方式发送数据，所以该项目只能满足业务量不大的情况下的需求