package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bodsch/container-service-discovery/container"
	"github.com/bodsch/container-service-discovery/detect"
	"github.com/bodsch/container-service-discovery/discover"
	"github.com/bodsch/container-service-discovery/full_list"
	"github.com/bodsch/container-service-discovery/health"
	"github.com/bodsch/container-service-discovery/utils"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
	"time"
)

var Version string = "0.10.0"
var Description string = "Prometheus Service Discovery for docker"

// Args command-line parameters
type Args struct {
	ConfigPath string
}

type ConfigRestAPI struct {
	Address string `yaml:"address" env-default: "127.0.0.1"`
	Port    int    `yaml:"port" env-default: "8088"`
}

type ConfigDockerHosts []struct {
	Host         string            `yaml:"host"`
	Username     interface{}       `yaml:"username,omitempty"`
	Password     interface{}       `yaml:"password,omitempty"`
	Services     []string          `yaml:"services,omitempty"`
	MetricsPorts map[string]string `yaml:"metrics_ports,omitempty"`
}

type Config struct {
	LogFile     string            `yaml:"log_file"`
	RestAPI     ConfigRestAPI     `yaml:"rest_api"`
	DockerHosts ConfigDockerHosts `yaml:"docker_hosts"`
}

var DockerHost = "unix:///run/docker.sock"
var DockerMetrics = make(map[string]string)

var debug bool

const (
	// YYYY-MM-DD: 2022-03-23
	YYYYMMDD = "2006-01-02"
	// 24h hh:mm:ss: 14:23:20
	HHMMSS24h = "15:04:05"
)

func ping() map[string]string {
	s := make(map[string]string)
	p, _ := container.Ping(DockerHost)
	json.Unmarshal([]byte(p), &s)
	return s
}

func healthCheck() map[string]string {
	s := make(map[string]string)
	h, _ := health.Check(DockerHost)
	json.Unmarshal([]byte(h), &s)
	return s
}

func fullList() map[string]interface{} {
	s := make(map[string]interface{})
	h, _ := full_list.FullList(DockerHost, debug)
	json.Unmarshal([]byte(h), &s)
	return s
}

func detectContainer() map[string]interface{} {
	s := make(map[string]interface{})
	l, _ := detect.Detect(DockerHost, debug)
	json.Unmarshal([]byte(l), &s)
	return s
}

func serviceDiscover() []discover.ServiceDiscover {
	l, _ := discover.Discover(DockerHost, debug)
	return l
}

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["proto"] = params.Request.Proto
			log["start_time"] = params.TimeStamp.Format("02-01-2006 15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()
			log["user_agent"] = params.Request.UserAgent()

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}

func main() {

	var listenAddress string = "127.0.0.1"
	var listenPort int = 8088

	var cfg Config
	args := ProcessArgs(&cfg)

	// read configuration from the file and environment variables
	err := cleanenv.ReadConfig(args.ConfigPath, &cfg)
	if err != nil {
		log.Printf("Unable to read config file: %v", err)
		log.Printf("use default values")

	} else {

		logWriter, err := syslog.New(syslog.LOG_SYSLOG, "docker-sd")
		if err != nil {
			log.Fatalln("Unable to set logfile:", err.Error())
		}
		log.SetFlags(log.Lmsgprefix)
		log.SetPrefix(time.Now().UTC().Format(YYYYMMDD+" "+HHMMSS24h) + ": ")
		log.SetOutput(logWriter)

		if debug {
			fmt.Printf("[DEBUG] cfg: %v\n", cfg)
		}
		listenAddress = cfg.RestAPI.Address
		listenPort = cfg.RestAPI.Port

		dockerHosts := cfg.DockerHosts

		if len(dockerHosts) > 0 {

			dockerConfiguration := dockerHosts[0]
			DockerHost = dockerConfiguration.Host

			if len(DockerHost) > 0 {

				DockerMetrics = dockerConfiguration.MetricsPorts

				if len(DockerMetrics) > 0 {
					// read path definition of the configuration
					if debug {
						fmt.Printf("[DEBUG] read path definition of the configuration\n")
					}
					for k, v := range DockerMetrics {
						// convert string into unint16
						value, _ := strconv.ParseUint(k, 0, 16)
						// strconv.ParseUint returns only uint64!
						port := uint16(value)

						utils.KnownMetricsPorts[port] = v
					}
				} else {
					// use default metrics definition
					if debug {
						fmt.Printf("[DEBUG] use default metrics definition\n")
					}
					utils.KnownMetricsPorts = map[uint16]string{
						8080: "/metrics", // cadvisor
						9216: "/metrics", // mongodb
						8090: "/metrics", // mgob
						9090: "/metrics", // prometheus
						9100: "/metrics", // node_exporter
					}
				}
			}
		}
	}

	/*
	 * default values
	 */
	if listenAddress == "" {
		listenAddress = "127.0.0.1"
	}
	if listenPort == 0 {
		listenPort = 8088
	}

	listener := fmt.Sprintf("%s:%d", listenAddress, listenPort)

	if debug {
		fmt.Printf("[DEBUG] listener: %s\n", listener)
		fmt.Printf("[DEBUG] known ports: %v\n", utils.KnownMetricsPorts)
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(jsonLoggerMiddleware())
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {

		multiline := "\n" + Description + " - Version " + Version + "\n\n" +
			"Routes:\n\n" +
			"  - http://localhost:$PORT/\n" +
			"  - http://localhost:$PORT/sd\n" +
			"  - http://localhost:$PORT/discover\n" +
			"  - http://localhost:$PORT/full_list\n" +
			"  - http://localhost:$PORT/health\n" +
			"  - http://localhost:$PORT/version\n" +
			"  - http://localhost:$PORT/detect\n"

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(multiline))
	})

	router.GET("/ping", func(c *gin.Context) {
		ping := ping()
		c.JSON(http.StatusOK, ping)
	})

	router.GET("/health", func(c *gin.Context) {
		health := healthCheck()
		c.JSON(http.StatusOK, health)
	})

	router.GET("/version", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{"version": Version})
	})

	router.GET("/full_list", func(c *gin.Context) {
		list := fullList()
		c.JSON(http.StatusOK, list)
	})

	router.GET("/detect", func(c *gin.Context) {
		list := detectContainer()
		c.JSON(http.StatusOK, list)
	})

	router.GET("/discover", func(c *gin.Context) {
		list := serviceDiscover()
		c.JSON(http.StatusOK, list)
	})

	router.GET("/sd", func(c *gin.Context) {
		list := serviceDiscover()
		c.JSON(http.StatusOK, list)
	})

	router.Run(listener)
}

// ProcessArgs processes and handles CLI arguments
func ProcessArgs(cfg interface{}) Args {
	var a Args
	var configHelp bool

	flag.BoolVar(&configHelp, "help", false, "This help")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.StringVar(&a.ConfigPath, "config", "docker-sd.yml", "configuration file")

	flag.Parse()

	if configHelp {
		multiline := "\n" + Description + " - Version " + Version + "\n"
		fmt.Println(multiline)
		flag.PrintDefaults()
		os.Exit(0)
	}

	return a
}
