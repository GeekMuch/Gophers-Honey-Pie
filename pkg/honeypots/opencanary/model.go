package opencanary

/* Canaryconf structs start */
type canaryConf struct {
	DeviceNodeId                 string       `json:"device.node_id"`
	IpIgnoreList                 []string     `json:"ip.ignorelist"`
	GitEnabled                   bool         `json:"git.enabled"`
	GitPort                      uint16       `json:"git.port"`
	FtpEnabled                   bool         `json:"ftp.enabled"`
	FtpPort                      uint16       `json:"ftp.port"`
	FtpBanner                    string       `json:"ftp.banner"`
	HttpBanner                   string       `json:"http.banner"`
	HttpEnabled                  bool         `json:"http.enabled"`
	HttpPort                     uint16       `json:"http.port"`
	HttpSkin                     string       `json:"http.skin"`
	HttpProxyEnabled             bool         `json:"httpproxy.enabled"`
	HttpProxyPort                uint16       `json:"httpproxy.port"`
	HttpProxySkin                string       `json:"httpproxy.skin"`
	Logger                       logger       `json:"logger"`
	PortScanEnabled              bool         `json:"portscan.enabled"`
	PortScanIgnoreLocalhost      bool         `json:"portscan.ignore_localhost"`
	PortScanLogfile              string       `json:"portscan.logfile"`
	PortScanSynrate              uint8        `json:"portscan.synrate"`
	PortScanNmapOsRate           uint8        `json:"portscan.nmaposrate"`
	PortScanLoRate               uint8        `json:"portscan.lorate"`
	SmbAuditFile                 string       `json:"smb.auditfile"`
	SmbEnabled                   bool         `json:"smb.enabled"`
	MysqlEnabled                 bool         `json:"mysql.enabled"`
	MysqlPort                    uint16       `json:"mysql.port"`
	MysqlBanner                  string       `json:"mysql.banner"`
	SshEnabled                   bool         `json:"ssh.enabled"`
	SshPort                      uint16       `json:"ssh.port"`
	SshVersion                   string       `json:"ssh.version"`
	RedisEnabled                 bool         `json:"redis.enabled"`
	RedisPort                    uint16       `json:"redis.port"`
	RdpEnabled                   bool         `json:"rdp.enabled"`
	RdpPort                      uint16       `json:"rdp.port"`
	SipEnabled                   bool         `json:"sip.enabled"`
	SipPort                      uint16       `json:"sip.port"`
	SnmpEnabled                  bool         `json:"snmp.enabled"`
	SnmpPort                     uint16       `json:"snmp.port"`
	NtpEnabled                   bool         `json:"ntp.enabled"`
	NtpPort                      uint16       `json:"ntp.port"`
	TftpEnabled                  bool         `json:"tftp.enabled"`
	TftpPort                     uint16       `json:"tftp.port"`
	TcpBannerMaxnum              uint         `json:"tcpbanner.maxnum"`
	TcpBannerEnabled             bool         `json:"tcpbanner.enabled"`
	TcpBanner1Enabled            bool         `json:"tcpbanner_1.enabled"`
	TcpBanner1Port               uint16       `json:"tcpbanner_1.port"`
	TcpBanner1DataReceivedBanner string       `json:"tcpbanner_1.datareceivedbanner"`
	TcpBanner1InitBanner         string       `json:"tcpbanner_1.initbanner"`
	TcpBanner1AlertstringEnabled bool         `json:"tcpbanner_1.alertstring.enabled"`
	TcpBanner1Alertstring        string       `json:"tcpbanner_1.alertstring"`
	TcpBanner1KeepAliveEnabled   bool         `json:"tcpbanner_1.keep_alive.enabled"`
	TcpBanner1KeepAliveSecret    string       `json:"tcpbanner_1.keep_alive_secret"`
	TcpBanner1KeepAliveProbes    uint         `json:"tcpbanner_1.keep_alive_probes"`
	TcpBanner1KeepAliveInterval  uint         `json:"tcpbanner_1.keep_alive_interval"`
	TcpBanner1KeepAliveIdle      uint         `json:"tcpbanner_1.keep_alive_idle"`
	TelnetEnabled                bool         `json:"telnet.enabled"`
	TelnetPort                   uint16       `json:"telnet.port"`
	TelnetBanner                 string       `json:"telnet.banner"`
	TelnetHoneycreds             []honeycreds `json:"telnet.honeycreds"`
	MssqlEnabled                 bool         `json:"mssql.enabled"`
	MssqlVersion                 string       `json:"mssql.version"`
	MssqlPort                    uint16       `json:"mssql.port"`
	VncEnabled                   bool         `json:"vnc.enabled"`
	VncPort                      uint16       `json:"vnc.port"`
}

type honeycreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type logger struct {
	Class  string `json:"class"`
	Kwargs kwargs `json:"kwargs"`
}

type kwargs struct {
	Formatters formatters `json:"formatters"`
	Handlers   handlers   `json:"handlers"`
}
type formatters struct {
	Plain     plain     `json:"plain"`
	SyslogRfc syslogRfc `json:"syslog_rfc"`
}
type plain struct {
	Format string `json:"format"`
}
type syslogRfc struct {
	Format string `json:"format"`
}
type handlers struct {
	Console console `json:"console"`
	File    file    `json:"file"`
}

type console struct {
	Class  string `json:"class"`
	Stream string `json:"stream"`
}

type file struct {
	Class    string `json:"class"`
	FileName string `json:"filename"`
}

/* CanaryConf structs end */
