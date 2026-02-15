package config

var Config *Conf

// 项目全局配置文件

type Conf struct {
	// 经典的数据库连接池配置（最大空闲、最大连接数等）
	Mysql MysqlConfig
	Redis RedisConfig
	// 包含了 ChainId（链 ID）、NetUrl（节点 RPC 地址）以及各种智能合约地址（如 PlgrAddress、PledgePoolToken）。
	// 这说明你的后端会通过 go-ethereum 库与 BSC（币安智能链）等兼容链交互
	TestNet      TestNetConfig
	MainNet      MainNetConfig
	Token        TokenConfig
	Email        EmailConfig
	DefaultAdmin DefaultAdminConfig
	// (阈值控制): PledgePoolTokenThresholdBnb 用于设定某种代币相对于 BNB 的质押阈值，这通常涉及价格预测机（Oracle）的逻辑
	Threshold ThresholdConfig
	// 对应你之前展示的 JWT 工具和 Email 发送工具的参数。
	Jwt JwtConfig
	// (运行环境): 定义了 API 端口、版本以及关键的超时时间（如 TaskDuration 任务持续时间、WssTimeoutDuration WebSocket 超时）
	Env EnvConfig
}

type EnvConfig struct {
	Port               string `toml:"port"`
	Version            string `toml:"version"`
	Protocol           string `toml:"protocol"`
	DomainName         string `toml:"domain_name"`
	TaskDuration       int64  `toml:"task_duration"`
	WssTimeoutDuration int64  `toml:"wss_timeout_duration"`
	TaskExtendDuration int64  `toml:"task_extend_duration"`
}

type ThresholdConfig struct {
	PledgePoolTokenThresholdBnb string `toml:"pledge_pool_token_threshold_bnb"`
}

type EmailConfig struct {
	Username string   `toml:"username"`
	Pwd      string   `toml:"pwd"`
	Host     string   `toml:"host"`
	Port     string   `toml:"port"`
	From     string   `toml:"from"`
	Subject  string   `toml:"subject"`
	To       []string `toml:"to"`
	Cc       []string `toml:"cc"`
}

type DefaultAdminConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type JwtConfig struct {
	SecretKey  string `toml:"secret_key"`
	ExpireTime int    `toml:"expire_time"` // duration, s
}

type TokenConfig struct {
	LogoUrl string `toml:"logo_url"`
}

type MysqlConfig struct {
	Address      string `toml:"address"`
	Port         string `toml:"port"`
	DbName       string `toml:"db_name"`
	UserName     string `toml:"user_name"`
	Password     string `toml:"password"`
	MaxOpenConns int    `toml:"max_open_conns"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxLifeTime  int    `toml:"max_life_time"`
}

type TestNetConfig struct {
	ChainId              string `toml:"chain_id"`
	NetUrl               string `toml:"net_url"`
	PlgrAddress          string `toml:"plgr_address"`
	PledgePoolToken      string `toml:"pledge_pool_token"`
	BscPledgeOracleToken string `toml:"bsc_pledge_oracle_token"`
}

type MainNetConfig struct {
	ChainId              string `toml:"chain_id"`
	NetUrl               string `toml:"net_url"`
	PlgrAddress          string `toml:"plgr_address"`
	PledgePoolToken      string `toml:"pledge_pool_token"`
	BscPledgeOracleToken string `toml:"bsc_pledge_oracle_token"`
}

type RedisConfig struct {
	// Redis服务器地址，如 "127.0.0.1" 或 "redis-server"
	Address string `toml:"address"`

	// Redis服务端口，通常是 "6379"
	Port string `toml:"port"`

	// 使用的数据库索引，默认为 0（Redis默认有16个db）
	Db int `toml:"db"`

	// Redis访问密码，如果没有设置则留空
	Password string `toml:"password"`

	// 连接池中最大空闲连接数
	// 即使没有请求，池里也会保留这么多连接随时待命
	MaxIdle int `toml:"max_idle"`

	// 连接池在同一时间能够分配的最大连接数
	// 设为 0 表示不限制，通常根据服务器并发处理能力设定
	MaxActive int `toml:"max_active"`

	// 空闲连接的超时时间（秒）
	// 超过这个时间的空闲连接会被关闭并从池中移除，防止连接因长时间不用被服务端断开
	IdleTimeout int `toml:"idle_timeout"`
}
