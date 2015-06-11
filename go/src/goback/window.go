package goback

import (
	"github.com/chenfeng/goback/win"
	"log"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	moduser32      = syscall.NewLazyDLL("user32.dll")
	procFindWindow = moduser32.NewProc("FindWindowW")
)

func FindWindow(className, capName string) win.HWND {
	param1 := uintptr(unsafe.Pointer(nil))
	param2 := uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(capName)))
	ret, _, _ := procFindWindow.Call(param1, param2)

	return win.HWND(ret)
}

const (
	CAPTIONNAME = "goback_Wnd"
	CLASSNAME   = "goback_Class"
)

type BackWnd struct {
	bStart bool
	hwnd   win.HWND
	cppwnd string
}

type COPYDATASTRUCT struct {
	dwData uintptr
	cbData uint32
	lpData uintptr
}

type COMMDATA struct {
	data string
}

func NewBackWnd(cppwndName string) *BackWnd {

	obj := new(BackWnd)
	obj.bStart = true

	// 注册窗口
	obj.regist(win.CS_HREDRAW | win.CS_VREDRAW | win.CS_OWNDC | win.CS_DBLCLKS)

	// 创建隐藏的后台服务窗口
	hwnd := obj.create(CLASSNAME, CAPTIONNAME, 0, win.WS_OVERLAPPEDWINDOW|win.WS_CLIPCHILDREN|win.WS_CLIPSIBLINGS)
	if hwnd == 0 {
		log.Println("Create BackServerWnd failed")
		return nil
	}

	log.Println("Create BackServerWnd success:", hwnd)
	win.ShowWindow(hwnd, win.SW_HIDE)
	win.UpdateWindow(hwnd)

	// 获取主窗口句柄
	obj.cppwnd = cppwndName
	obj.hwnd = FindWindow("", cppwndName)
	//go obj.WaitMessage()

	return obj
}

func (p *BackWnd) regist(style uint) {
	hInst := win.GetModuleHandle(nil)
	if hInst == 0 {
		panic("GetModuleHandle")
	}

	hIcon := win.LoadIcon(0, (*uint16)(unsafe.Pointer(uintptr(win.IDI_APPLICATION))))
	if hIcon == 0 {
		panic("LoadIcon")
	}

	hCursor := win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
	if hCursor == 0 {
		panic("LoadCursor")
	}

	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = syscall.NewCallback(p.BackWndProc)
	wc.HInstance = hInst
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = win.COLOR_BTNFACE + 1
	wc.LpszClassName = syscall.StringToUTF16Ptr(CLASSNAME)

	if atom := win.RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}
}

func (p *BackWnd) create(class string, caption string, exStyle, style uint32) win.HWND {
	hInst := win.GetModuleHandle(nil)

	var hwnd win.HWND
	hwnd = win.CreateWindowEx(
		exStyle,
		syscall.StringToUTF16Ptr(class),
		syscall.StringToUTF16Ptr(caption),
		style|win.WS_CLIPSIBLINGS,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		0,
		0,
		hInst,
		nil)

	return hwnd
}

func (p *BackWnd) WaitMessage() {
	var msg win.MSG
	for win.GetMessage(&msg, 0, 0, 0) != -1 {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)

		if msg.Message == win.WM_QUIT {
			log.Println("exit app")
			break
		}
	}
}

func (p *BackWnd) SendSyn(data string) {
	if p.hwnd == 0 {
		p.hwnd = FindWindow("", p.cppwnd)
	}
	arrUtf16, _ := syscall.UTF16FromString(data)

	pCopyData := new(COPYDATASTRUCT)
	pCopyData.dwData = 0
	pCopyData.cbData = uint32(len(arrUtf16)*2 + 1)
	pCopyData.lpData = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(data)))
	win.SendMessage(p.hwnd, win.WM_COPYDATA, 0, uintptr(unsafe.Pointer(pCopyData)))
}

func (p *BackWnd) SendASyn(data string) {

}

func (p *BackWnd) SetStartState(state bool) {
	p.bStart = state
}

func uintptrToString(cstr uintptr) string {
	if cstr != 0 {
		us := make([]uint16, 0, 256)
		for p := cstr; ; p += 2 {
			u := *(*uint16)(unsafe.Pointer(p))
			if u == 0 {
				return string(utf16.Decode(us))
			}
			us = append(us, u)
		}
	}
	return ""
}

func (p *BackWnd) BackWndProc(hwnd win.HWND, msg uint, wparam, lparam uintptr) uintptr {
	if msg == win.WM_COPYDATA {
		ldata := (*COPYDATASTRUCT)(unsafe.Pointer(lparam))
		go backsev.Accept(uintptrToString(ldata.lpData))
		return 0
	}
	return win.DefWindowProc(hwnd, uint32(msg), wparam, lparam)
}
