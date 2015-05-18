#pragma once

// 通讯规则：两个字段 Type, Data，具体内部规则由业务指定
// 1、DUIToBack需要指定流程状态，比如Init,Run,Over,再附带请求数据
// 示例:{"Type":"Init", Data:""}, {"Type":"Run", Data:""}, {"Type":"Over", Data:""}
// 2、BackToDUI按业务数据类型区分，包括更新UI数据，更新日志数据，更新系统信息，再附带更新数据
// 示例:{"Type":"Edit", "Data":""}, {"Type":"Log", Data:""}, {"Type":"Sys", Data:""}

namespace Service
{
	namespace Net
	{
		typedef void callback_func_type(CString data);
		class Comm
		{
		public:
			static bool BackToDUI(CString strData);
			static bool DUIToBack(CString strData);
			static void Init(callback_func_type *pFunc);

			static callback_func_type *s_pFunc;
		};
	}
}

