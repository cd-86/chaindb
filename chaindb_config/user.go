package chaindb_config

// 该 HTTP 服务端口用于响应终端用户的各类请求.
// 比如, 向该服务发送添加 Tx 的请求, 它会负责
// 转换请求的格式并按需转发到其它每个节点的
// TCPPortMonitor 服务.
//
// 也就是说, 终端用户只需要关注这一个端口, 它
// 提供了所有 API 接口.
const TCPPortUserService = 56784
