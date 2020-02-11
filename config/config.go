package config

import (
	"os"
	"strconv"
)

const (
	RedisAddr           = "REDIS_ADDR"            //Address of the Redis instance to scrape
	RedisPwd            = "REDIS_PWD"             //Password of the Redis instance to scrape
	Namespace           = "NAMESPACE"             //Namespace for metrics
	CheckKeys           = "CHECK_KEYS"            //Comma separated list of key-patterns to export value and length/size, searched for with SCAN
	CheckSingleKeys     = "CHECK_SINGLE_KEYS"     //Comma separated list of single keys to export value and length/size
	ScriptPath          = "SCRIPT_PATH"           //Path to Lua Redis script for collecting extra metrics
	ListenAddress       = "LISTEN_ADDRESS"        //Address to listen on for web interface and telemetry.
	MetricPath          = "METRIC_PATH"           //Path under which to expose metrics.
	LogFormat           = "LOG_FORMAT"            //Log format, valid options are txt and json
	ConfigCommand       = "CONFIG_COMMAND"        //What to use for the CONFIG command
	ConnectionTimeout   = "CONNECTION_TIMEOUT"    //Timeout for connection to Redis instance
	TlsClientKeyFile    = "TLS_CLIENT_KEY_FILE"   //Name of the client key file (including full path) if the server requires TLS client authentication
	TlsClientCertFile   = "TLS_CLIENT_CERT_FILE"  //Name of the client certificate file (including full path) if the server requires TLS client authentication
	IsDebug             = "IS_DEBUG"              //Output verbose debug information
	IsTile38            = "IS_TILE38"             //Whether to scrape Tile38 specific metrics
	ExportClientList    = "EXPORT_CLIENT_LIST"    //Whether to scrape Client List specific metrics
	SetClientName       = "SET_CLIENT_NAME"       //Whether to set client name to redis_exporter
	ShowVersion         = "SHOW_VERSION"          //Show version information and exit
	RedisMetricsOnly    = "REDIS_METRICS_ONLY"    //Whether to also export go runtime metrics
	InclSystemMetrics   = "INCL_SYSTEM_METRICS"   //Whether to include system metrics like e.g. redis_total_system_memory_bytes
	SkipTLSVerification = "SKIP_TLS_VERIFICATION" //Whether to skip TLS verification
	PushGatewayAddr     = "PUSH_GATEWAY_ADDR"     //Address of the push gateway
	MetricsJobName      = "METRICS_JOB_NAME"      //Job name of the metrics to be pushed
	AppName             = "APP_NAME"              //Rock App Name
	InstanceID          = "INSTANCE_ID"           //Rock instance ID
	PushIntervalInSec   = "PUSH_INTERVAL_IN_SEC"  //Interval to push metrics to push gateway in seconds
)

const (
	defaultRedisAddr           = "redis://localhost:9851"
	defaultRedisPwd            = ""
	defaultNamespace           = "namespace"
	defaultCheckKeys           = ""
	defaultCheckSingleKeys     = ""
	defaultScriptPath          = ""
	defaultListenAddress       = ":9121"
	defaultMetricPath          = "/metrics"
	defaultLogFormat           = "txt"
	defaultConfigCommand       = "CONFIG"
	defaultConnectionTimeout   = "15s"
	defaultTlsClientKeyFile    = ""
	defaultTlsClientCertFile   = ""
	defaultIsDebug             = false
	defaultIsTile38            = true
	defaultExportClientList    = false
	defaultSetClientName       = false
	defaultShowVersion         = false
	defaultRedisMetricsOnly    = false
	defaultInclSystemMetrics   = true
	defaultSkipTLSVerification = true
	defaultPushGatewayAddr     = "http://localhost:8081"
	defaultMetricsJobName      = ""
	defaultAppName             = ""
	defaultInstanceID          = ""
	defaultPushIntervalInSec   = 30
)

type Config struct {
	RedisAddr           string
	RedisPwd            string
	Namespace           string
	CheckKeys           string
	CheckSingleKeys     string
	ScriptPath          string
	ListenAddress       string
	MetricPath          string
	LogFormat           string
	ConfigCommand       string
	ConnectionTimeout   string
	TlsClientKeyFile    string
	TlsClientCertFile   string
	IsDebug             bool
	IsTile38            bool
	ExportClientList    bool
	SetClientName       bool
	ShowVersion         bool
	RedisMetricsOnly    bool
	InclSystemMetrics   bool
	SkipTLSVerification bool
	PushGatewayAddr     string
	MetricsJobName      string
	AppName             string
	InstanceID          string
	PushIntervalInSec   int
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		RedisAddr:           getEnv(RedisAddr, defaultRedisAddr),
		RedisPwd:            getEnv(RedisPwd, defaultRedisPwd),
		Namespace:           getEnv(Namespace, defaultNamespace),
		CheckKeys:           getEnv(CheckKeys, defaultCheckKeys),
		CheckSingleKeys:     getEnv(CheckSingleKeys, defaultCheckSingleKeys),
		ScriptPath:          getEnv(ScriptPath, defaultScriptPath),
		ListenAddress:       getEnv(ListenAddress, defaultListenAddress),
		MetricPath:          getEnv(MetricPath, defaultMetricPath),
		LogFormat:           getEnv(LogFormat, defaultLogFormat),
		ConfigCommand:       getEnv(ConfigCommand, defaultConfigCommand),
		ConnectionTimeout:   getEnv(ConnectionTimeout, defaultConnectionTimeout),
		TlsClientKeyFile:    getEnv(TlsClientKeyFile, defaultTlsClientKeyFile),
		TlsClientCertFile:   getEnv(TlsClientCertFile, defaultTlsClientCertFile),
		IsDebug:             getEnvAsBool(IsDebug, defaultIsDebug),
		IsTile38:            getEnvAsBool(IsTile38, defaultIsTile38),
		ExportClientList:    getEnvAsBool(ExportClientList, defaultExportClientList),
		SetClientName:       getEnvAsBool(SetClientName, defaultSetClientName),
		ShowVersion:         getEnvAsBool(ShowVersion, defaultShowVersion),
		RedisMetricsOnly:    getEnvAsBool(RedisMetricsOnly, defaultRedisMetricsOnly),
		InclSystemMetrics:   getEnvAsBool(InclSystemMetrics, defaultInclSystemMetrics),
		SkipTLSVerification: getEnvAsBool(SkipTLSVerification, defaultSkipTLSVerification),
		PushGatewayAddr:     getEnv(PushGatewayAddr, defaultPushGatewayAddr),
		MetricsJobName:      getEnv(MetricsJobName, defaultMetricsJobName),
		AppName:             getEnv(AppName, defaultAppName),
		InstanceID:          getEnv(InstanceID, defaultInstanceID),
		PushIntervalInSec:   getEnvAsInt(PushIntervalInSec, defaultPushIntervalInSec),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
