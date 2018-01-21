package main

//define dsp service with method:
// /v1  default api as request
// /v1/expose/monitor
// /v1/click/monitor
import (
	"net/http"
	"os/signal"
	"syscall"
	"os"
	"time"
	"dsp_demo/service"
	"github.com/vrischmann/go-metrics-influxdb"
	"github.com/rcrowley/go-metrics"
	"dsp_demo/metric"
	"github.com/astaxie/beego/logs"
)

func main() {
	logs.SetLogger(logs.AdapterConsole, `{"level":6}`)
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"./logs/hm_server.log","level":6}`)
	logs.EnableFuncCallDepth(true)
	logs.Async(1e3)

	var dspSrv service.DSPService
	dspSrv.Service = service.NewService()
	dspSrv.RouterRegister = func() {
		http.HandleFunc("/v1/ads/",  metric.Decorate(dspSrv.AdsHandler))
		http.HandleFunc("/v1/expose/monitor", metric.Decorate(dspSrv.ExposeHandler))
		http.HandleFunc("/v1/click/monitor", metric.Decorate(dspSrv.ClickHandler))
		http.HandleFunc("/v1", metric.Decorate(dspSrv.AdsHandler))
		http.HandleFunc("/v1/admin/conf", metric.Decorate(dspSrv.ConfHandler))
	}
	go dspSrv.Start()

	r := metrics.NewRegistry()
	metrics.RegisterDebugGCStats(r)
	metrics.RegisterRuntimeMemStats(r)

	go metrics.CaptureDebugGCStats(r, time.Second*5)
	go metrics.CaptureRuntimeMemStats(r, time.Second*5)

	//register metric
	go influxdb.InfluxDB(
		r, // metrics registry
		time.Second * 5,        // interval
		"xxx", // the InfluxDB url
		"xxx",                  // your InfluxDB database
		"xxx",                // your InfluxDB user
		"xxx",            // your InfluxDB password
	)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	logs.Info(<-ch)

	// Stop the service gracefully.
	dspSrv.Stop()
}

