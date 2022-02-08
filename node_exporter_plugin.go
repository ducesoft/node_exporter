package main

import (
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/node_exporter/collector"
	"net/http"
	"os"
	"strconv"

	"os/user"
)

var (
	DftPromLogLevel = &promlog.AllowedLevel{}
	DftPromLogFmt   = &promlog.AllowedFormat{}
)

type NodeExporterConfig struct {
	// such as info
	PromLogLevel string
	// such as logfmt
	PromLogFmt string
	// such as 9100
	ListenPort int
	// such as /metrics
	MetricsPath string
	// such as /
	ExporterPath string
	// such as false
	DisableExporterMetrics bool
	// such as 40
	MaxRequests int
	// such as false
	DisableDefaultCollectors bool
	// such as ""
	ConfigFile string
}

// StartNodeExporter access for quick integration
// NOT FORGET to run with goroutine
func StartNodeExporter(conf *NodeExporterConfig) {
	initConfig(conf)
	logger := promlog.New(&promlog.Config{
		Level:  DftPromLogLevel,
		Format: DftPromLogFmt,
	})
	if conf.DisableDefaultCollectors {
		collector.DisableDefaultCollectors()
	}
	level.Info(logger).Log("msg", "Starting node_exporter", "version", "1.3.1")
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())
	if curUser, err := user.Current(); err == nil && curUser.Uid == "0" {
		level.Warn(logger).Log("msg", "Node Exporter is running as root user. This exporter is designed to run as unpriviledged user, root is not required.")
	}

	http.Handle(conf.MetricsPath, newHandler(true, conf.MaxRequests, logger))
	http.HandleFunc(conf.ExporterPath, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Node Exporter</title></head>
			<body>
			<h1>Node Exporter</h1>
			<p><a href="` + conf.MetricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	listenAddr := ":" + strconv.Itoa(conf.ListenPort)
	level.Info(logger).Log("msg", "Node Exporter Listening on", "address", listenAddr)
	server := &http.Server{Addr: listenAddr}
	if err := web.ListenAndServe(server, conf.ConfigFile, logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}

func initConfig(conf *NodeExporterConfig) {
	if conf.PromLogLevel == "" {
		DftPromLogLevel.Set("info")
	}
	if conf.PromLogFmt == "" {
		DftPromLogFmt.Set("logfmt")
	}
	if conf.ListenPort <= 0 {
		conf.ListenPort = 9100
	}
	if conf.MetricsPath == "" {
		conf.MetricsPath = "/metrics"
	}
	if conf.ExporterPath == "" {
		conf.ExporterPath = "/"
	}
	if conf.MaxRequests <= 0 {
		conf.MaxRequests = 40
	}
}
