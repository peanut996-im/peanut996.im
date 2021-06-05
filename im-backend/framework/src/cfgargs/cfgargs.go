package cfgargs

import (
	"flag"
	"fmt"
	"framework/tool"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	defaultConfigPath = "./etc/config-local.yaml"
)

var (
	lastSrvCfg *SrvConfig
)

//SrvConfig records for all conf
type SrvConfig struct {
	Build
	Mongo    `yaml:"mongo" json:"Mongo,omitempty"`
	Redis    `yaml:"redis" json:"Redis,omitempty"`
	Log      `yaml:"log" json:"Log,omitempty"`
	HTTP     `yaml:"http" json:"HTTP,omitempty"`
	SocketIO `yaml:"socket.io" json:"SocketIO,omitempty"`
	Logic    `yaml:"logic" json:"Logic,omitempty"`
	Gate     `yaml:"gate" json:"Gate,omitempty"`
	AppKey   string `yaml:"appkey" json:"AppKey"`
}

type Build struct {
	BuildTime    string
	BuildUser    string
	BuildVersion string
	BuildMachine string
}

type SocketIO struct {
	Port int  `yaml:"port"`
	Cors bool `yaml:"cors"`
}

type HTTP struct {
	Cors    bool   `yaml:"cors"`
	Port    string `yaml:"port"`
	Release bool   `yaml:"release"`
	Sign    bool   `yaml:"sign"`
}

type Log struct {
	Level   string `yaml:"level"`
	Console bool   `yaml:"console"`
	Path    string `yaml:"path"`
	Sync    bool   `yaml:"sync"`
}

//Mongo conf for mongoDB
type Mongo struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       string `yaml:"db"`
	Password string `yaml:"password"`
	Panic    bool   `yaml:"panic"`
}

//Redis configure for Redis
type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Panic    bool   `yaml:"panic"`
}

//Gate configure for Gate Node.
type Gate struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
	// 初始化Mode模式的Broker失败是否panic
	Panic bool `yaml:"panic"`
}

//Logic configure for Logic Server
type Logic struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
	// 初始化Mode模式的Broker失败是否panic
	Panic bool `yaml:"panic"`
}

func InitSrvCfg(build *Build, flagParse func()) (*SrvConfig, error) {
	srvCfg := newSrvConfig()
	if nil != build {
		srvCfg.Build = *build
	}

	yamlPath := ""

	flag.StringVar(&yamlPath, "c", defaultConfigPath, "App configuration file. Relative to the path of repository.")
	flag.StringVar(&yamlPath, "config", defaultConfigPath, "App configuration file. Relative to the path of repository.")
	if nil != flagParse {
		flagParse()
	}
	flag.Parse()

	if err := srvCfg.loadLocalYaml(yamlPath); err != nil {
		return nil, err
	}

	lastSrvCfg = srvCfg
	return srvCfg, nil
	// yamlPath := flag.String("c", "", "Yaml config file path.(Relative to the path of the executable file)")
}

func newSrvConfig() *SrvConfig {
	return &SrvConfig{}
}

func (s *SrvConfig) loadLocalYaml(path string) error {
	yamlFile := path
	data, err := ioutil.ReadFile(yamlFile)
	if nil != err {
		fmt.Printf("load local yaml err:%v path: %v\n", err, yamlFile)
		return err
	}

	err = yaml.Unmarshal([]byte(data), s)
	if nil != err {
		fmt.Printf("unmarshal local yaml err:%v path: %v\n", err, yamlFile)
		return err
	}

	if "" == s.Log.Level {
		s.Log.Level = "INFO"
	}

	s.Log.Level = strings.ToUpper(s.Log.Level)
	switch s.Log.Level {
	case "DEBUG":
	case "INFO":
	case "WARN":
	case "ERROR":
	case "FATAL":
	default:
		s.Log.Level = "INFO"
	}

	return nil
}

// GetRedisAddr printf the redis addr from srvconfig
func GetRedisAddr(config *SrvConfig) string {
	return fmt.Sprintf("%v:%v", config.Redis.Host, config.Redis.Port)
}

//GetLastSrvConfig 返回全局
func GetLastSrvConfig() *SrvConfig {
	return lastSrvCfg
}

func (s *SrvConfig) Print() {
	fmt.Println("BuildInfo:")
	json, err := tool.PrettyPrint(s)
	if err == nil {
		fmt.Println(json + "\n")
	} else {
		fmt.Printf("BuildTime: %v\n", s.BuildTime)
		fmt.Printf("BuildUser: %v\n", s.BuildUser)
		fmt.Printf("BuildVersion: %v\n", s.BuildVersion)
		fmt.Printf("BuildMachine: %v\n", s.BuildMachine)
		fmt.Printf("Log: %+v\n", s.Log)
		fmt.Printf("HTTP: %+v\n", s.HTTP)
		fmt.Printf("Mongo: %+v\n", s.Mongo)
		fmt.Printf("Redis: %+v\n", s.Redis)
		fmt.Println("")
	}

}
