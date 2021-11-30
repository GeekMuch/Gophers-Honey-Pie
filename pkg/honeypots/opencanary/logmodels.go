package opencanary

import "encoding/json"

// Logs severity levels.
var (
	CRITICAL      = 0
	SCAN          = 1
	INFORMATIONAL = 2
)

// Type severity levels.
var (
	CriticalTypes      = []int{2000, 3000, 3001, 4000, 4002, 5000, 6001, 7001, 8001, 9001, 9002}
	ScanTypes          = []int{5001, 5002, 5003, 5004, 5005}
	InformationalTypes = []int{1000, 1001, 1002, 1003, 1004, 1005, 1006, 10001, 11001, 12001, 13001, 14001, 15001,
		17001, 18001, 18002, 18003, 18004, 18005, 99000, 99001, 99002, 99003, 99004, 99005, 99006, 99007, 99008, 99009}
)

// OpenCanaryLogTypes from https://github.com/thinkst/opencanary/blob/master/opencanary/logger.py
var OpenCanaryLogTypes = map[int]string{
	1000:  "LOG_BASE_BOOT",
	1001:  "LOG_BASE_MSG",
	1002:  "LOG_BASE_DEBUG",
	1003:  "LOG_BASE_ERROR",
	1004:  "LOG_BASE_PING",
	1005:  "LOG_BASE_CONFIG_SAVE",
	1006:  "LOG_BASE_EXAMPLE",
	2000:  "LOG_FTP_LOGIN_ATTEMPT",
	3000:  "LOG_HTTP_GET",
	3001:  "LOG_HTTP_POST_LOGIN_ATTEMPT",
	4000:  "LOG_SSH_NEW_CONNECTION",
	4001:  "LOG_SSH_REMOTE_VERSION_SENT",
	4002:  "LOG_SSH_LOGIN_ATTEMPT",
	5000:  "LOG_SMB_FILE_OPEN",
	5001:  "LOG_PORT_SYN",
	5002:  "LOG_PORT_NMAPOS",
	5003:  "LOG_PORT_NMAPNULL",
	5004:  "LOG_PORT_NMAPXMAS",
	5005:  "LOG_PORT_NMAPFIN",
	6001:  "LOG_TELNET_LOGIN_ATTEMPT",
	7001:  "LOG_HTTPPROXY_LOGIN_ATTEMPT",
	8001:  "LOG_MYSQL_LOGIN_ATTEMPT",
	9001:  "LOG_MSSQL_LOGIN_SQLAUTH",
	9002:  "LOG_MSSQL_LOGIN_WINAUTH",
	10001: "LOG_TFTP",
	11001: "LOG_NTP_MONLIST",
	12001: "LOG_VNC",
	13001: "LOG_SNMP_CMD",
	14001: "LOG_RDP",
	15001: "LOG_SIP_REQUEST",
	16001: "LOG_GIT_CLONE_REQUEST",
	17001: "LOG_REDIS_COMMAND",
	18001: "LOG_TCP_BANNER_CONNECTION_MADE",
	18002: "LOG_TCP_BANNER_KEEP_ALIVE_CONNECTION_MADE",
	18003: "LOG_TCP_BANNER_KEEP_ALIVE_SECRET_RECEIVED",
	18004: "LOG_TCP_BANNER_KEEP_ALIVE_DATA_RECEIVED",
	18005: "LOG_TCP_BANNER_DATA_RECEIVED",
	99000: "LOG_USER_0",
	99001: "LOG_USER_1",
	99002: "LOG_USER_2",
	99003: "LOG_USER_3",
	99004: "LOG_USER_4",
	99005: "LOG_USER_5",
	99006: "LOG_USER_6",
	99007: "LOG_USER_7",
	99008: "LOG_USER_8",
	99009: "LOG_USER_9",
}

type OpenCanaryLogData struct {
	PASSWORD string               `json:"PASSWORD"`
	USERNAME string               `json:"USERNAME"`
	PROTO    string               `json:"PROTO"`
	Msg      OpenCanaryLogDataMsg `json:"msg,omitempty"`
}

type OpenCanaryLogDataMsg struct {
	LogData string `json:"logdata"`
}

type OpenCanaryLog struct {
	DstHost string            `json:"dst_host"`
	DstPort json.Number       `json:"dst_port"`
	SrcHost string            `json:"src_host"`
	SrcPort json.Number       `json:"src_port"`
	UTCTime string            `json:"utc_time"`
	Logdata OpenCanaryLogData `json:"logdata"`
	LogType int               `json:"log_type"`
	NodeID  string            `json:"node_id"`
}
