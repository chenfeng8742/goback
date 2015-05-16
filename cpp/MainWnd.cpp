#include "stdafx.h"
#include "MainFrame.h"

namespace UI
{
	LRESULT CMainFrame::HandleMessage(UINT uMsg, WPARAM wParam, LPARAM lParam)
	{
		LRESULT lRes = 0;
		BOOL bHandled = TRUE;
		switch (uMsg)
		{
		case WM_CREATE:        lRes = OnCreate(uMsg, wParam, lParam, bHandled); break;
		case WM_NCCALCSIZE:    lRes = OnNcCalcSize(uMsg, wParam, lParam, bHandled); break;
		case WM_SIZE:          lRes = OnSize(uMsg, wParam, lParam, bHandled); break;
		case WM_NCACTIVATE:    lRes = OnNcActivate(uMsg, wParam, lParam, bHandled); break;
		case WM_NCHITTEST:     lRes = OnNcHitTest(uMsg, wParam, lParam, bHandled); break;
		case WM_GETMINMAXINFO: lRes = OnGetMinMaxInfo(uMsg, wParam, lParam, bHandled); break;
		case WM_SYSCOMMAND:		lRes = OnSysCommand(uMsg, wParam, lParam, bHandled); break;
		case WM_CLOSE:         lRes = OnClose(uMsg, wParam, lParam, bHandled); break;
		default:               bHandled = FALSE;
		}

		if (uMsg == WM_COPYDATA)
		{
			COPYDATASTRUCT* pCDS = (COPYDATASTRUCT*)lParam;
			int count = (int)(pCDS->cbData);

			char *pData = new char[count];
			memset(pData, 0, count);
			wcstombs(pData, (wchar_t*)(pCDS->lpData), count);

			Service::Net::Comm::BackToDUI(pData);
			delete[]pData;
		}

		if (bHandled)
			return lRes;
		if (m_pm.MessageHandler(uMsg, wParam, lParam, lRes))
			return lRes;
		return CWindowWnd::HandleMessage(uMsg, wParam, lParam);
	}
}