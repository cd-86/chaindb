# 测试

我默认你把测试机器的防火墙完全关闭了.

服务发现需要打开 局域网组播, 这样 mDNS 协议才能够运行:
- MS-Windows: 打开注册表, 新建 `计算机\HKEY_LOCAL_MACHINE\SOFTWARE\Policies\Microsoft\Windows NT\DNSClient` (DOWRD), 名称为 `EnableMulticast`, 值为 `1`.
- Linux in VirtualBox on MS-Windows: 在 MS-Windows 上的控制面板中把 VirtualBox Host-Only Ethernet Adapter 禁用.
