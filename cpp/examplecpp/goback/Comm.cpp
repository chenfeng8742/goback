#include "stdafx.h"
#include "Comm.h"
#include <atlstr.h>
#include <windows.h>

// 通讯规则：两个字段 Type, Data，具体内部规则由业务指定
// 1、DUIToBack需要指定流程状态，比如Init,Run,Over,再附带请求数据
// 示例:{"Type":"Init", Data:""}, {"Type":"Run", Data:""}, {"Type":"Over", Data:""}
// 2、BackToDUI按业务数据类型区分，包括更新UI数据，更新日志数据，更新系统信息，再附带更新数据
// 示例:{"Type":"Edit", "Data":""}, {"Type":"Log", Data:""}, {"Type":"Sys", Data:""}

namespace Service
{
	namespace Net
	{
		callback_func_type * Comm::s_pFunc;
		void Comm::Init(callback_func_type *pFunc)
		{
			s_pFunc = pFunc;

			// 启动后台服务
			CString filePath;
			GetModuleFileName(NULL, filePath.GetBuffer(255), MAX_PATH);
			filePath.ReleaseBuffer();//获取带.exe文件名的路径
			filePath = filePath.Left(filePath.ReverseFind('\\'));
			CString szCmdLine = filePath + "\\example.exe";

			PROCESS_INFORMATION piProcInfoGPS;
			STARTUPINFO siStartupInfo;
			SECURITY_ATTRIBUTES saProcess, saThread;
			ZeroMemory(&siStartupInfo, sizeof(siStartupInfo));
			siStartupInfo.cb = sizeof(siStartupInfo);
			saProcess.nLength = sizeof(saProcess);
			saProcess.lpSecurityDescriptor = NULL;
			saProcess.bInheritHandle = true;
			saThread.nLength = sizeof(saThread);
			saThread.lpSecurityDescriptor = NULL;
			saThread.bInheritHandle = true;
			BOOL bRet = ::CreateProcess(NULL, (LPTSTR)(szCmdLine.GetString()), &saProcess, &saThread, false, CREATE_DEFAULT_ERROR_MODE, NULL, NULL, &siStartupInfo, &piProcInfoGPS);

			if (!bRet)
			{
				// 提示无法启动
				CString strPrint;
				if (2 == GetLastError())
				{
					strPrint.Format(_T("系统无法启动后台服务:\n%s"), szCmdLine);
				}
				::MessageBox(NULL, strPrint, _T("提示"), MB_OK);
			}
		}

		bool Comm::BackToDUI(CString strData)
		{
			if (strData.Compare(_T("goback for golang is OK")) == 0)
			{
				DUIToBack(_T("goback for c++ is OK"));
			}
			
			s_pFunc(strData);
			return true;
		}

		bool Comm::DUIToBack(CString strData)
		{
			int count = strData.GetLength() * 2 + 1;
			WCHAR *wstr = new WCHAR[count];
			mbstowcs(wstr, strData, count);

			COPYDATASTRUCT copyData;
			COPYDATASTRUCT *pCopyData = &copyData/*new COPYDATASTRUCT*/;
			pCopyData->dwData = 0;
			pCopyData->cbData = count;
			pCopyData->lpData = (void*)(wstr);

			Sleep(1000);
			HWND hwndCopy = ::FindWindow(NULL, _T("goback_Wnd"));
			if (hwndCopy == NULL)
				return false;

			LRESULT ret = ::GetLastError();
			::SendMessage(hwndCopy, WM_COPYDATA, 0, (LPARAM)(pCopyData));
			delete[]wstr;
			return true;
		}
	}
}
