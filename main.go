package main

import (
	"crypto/tls"
	"github.com/oliver006/redis_exporter/config"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	log "github.com/sirupsen/logrus"
)

var (
	// BuildVersion, BuildDate, BuildCommitSha are filled in by the build script
	BuildVersion   = "<<< filled in by build >>>"
	BuildDate      = "<<< filled in by build >>>"
	BuildCommitSha = "<<< filled in by build >>>"
)

func getEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if envVal, ok := os.LookupEnv(key); ok {
		envBool, err := strconv.ParseBool(envVal)
		if err == nil {
			return envBool
		}
	}
	return defaultVal
}

func main() {
	conf := config.New()
	switch conf.LogFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}
	log.Printf("Redis Metrics Exporter %s    build date: %s    sha1: %s    Go: %s    GOOS: %s    GOARCH: %s",
		BuildVersion, BuildDate, BuildCommitSha,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
	if conf.IsDebug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Enabling debug output")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if conf.ShowVersion {
		return
	}

	to, err := time.ParseDuration(conf.ConnectionTimeout)
	if err != nil {
		log.Fatalf("Couldn't parse connection timeout duration, err: %s", err)
	}

	var tlsClientCertificates []tls.Certificate
	if (conf.TlsClientKeyFile != "") != (conf.TlsClientCertFile != "") {
		log.Fatal("TLS client key file and cert file should both be present")
	}
	if conf.TlsClientKeyFile != "" && conf.TlsClientCertFile != "" {
		cert, err := tls.LoadX509KeyPair(conf.TlsClientCertFile, conf.TlsClientKeyFile)
		if err != nil {
			log.Fatalf("Couldn't load TLS client key pair, err: %s", err)
		}
		tlsClientCertificates = append(tlsClientCertificates, cert)
	}

	var ls []byte
	if conf.ScriptPath != "" {
		if ls, err = ioutil.ReadFile(conf.ScriptPath); err != nil {
			log.Fatalf("Error loading script file %s    err: %s", conf.ScriptPath, err)
		}
	}

	registry := prometheus.NewRegistry()
	if !conf.RedisMetricsOnly {
		registry = prometheus.DefaultRegisterer.(*prometheus.Registry)
	}

	exp, err := NewRedisExporter(
		conf.RedisAddr,
		ExporterOptions{
			Password:            conf.RedisPwd,
			Namespace:           conf.Namespace,
			ConfigCommandName:   conf.ConfigCommand,
			CheckKeys:           conf.CheckKeys,
			CheckSingleKeys:     conf.CheckSingleKeys,
			LuaScript:           ls,
			InclSystemMetrics:   conf.InclSystemMetrics,
			IsTile38:            conf.IsTile38,
			ExportClientList:    conf.ExportClientList,
			SkipTLSVerification: conf.SkipTLSVerification,
			SetClientName:       conf.SetClientName,
			ClientCertificates:  tlsClientCertificates,
			ConnectionTimeouts:  to,
			MetricsPath:         conf.MetricPath,
			RedisMetricsOnly:    conf.RedisMetricsOnly,
			Registry:            registry,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	go func(registry *prometheus.Registry){
		for {
			err = push.New(conf.PushGatewayAddr, conf.AppName).Gatherer(registry).Push()
			if err != nil {
				log.Infof("Error pushing metrics to push gateway. Error: %+v", err)
			}
			time.Sleep(time.Duration(conf.PushIntervalInSec) * time.Second)
		}
	}(registry)

	log.Fatal(http.ListenAndServe(conf.ListenAddress, exp))
}
