package vishnu

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

const configName = "vishnu"

type config struct {
	Port    int
	Bridge  string
	Type    string
	Timeout int
}

// Config system, Must LoadConfig() once before use
var Config = &config{
	Port:    1234,
	Bridge:  "localhost:1234",
	Type:    "bridge",
	Timeout: 10,
}

func (c *config) loadFromConfigFile() {
	configFile, _ := homedir.Expand(fmt.Sprintf("~/.%s.toml", configName))
	if _, err := toml.DecodeFile(configFile, Config); err != nil {
		log.Printf("skip load config file %s: %s", configFile, err)
		return
	}
	log.Printf("load config file %s", configFile)
}

// 配置加载顺序
// 1. 默认配置
// 2. ~/.${configName}.toml
// 3. 命令行参数
// 后面的配置会覆盖之前的
func (c *config) LoadConfig() *config {
	c.loadFromConfigFile()

	flag.IntVar(&Config.Port, "port", Config.Port, "server port")
	flag.StringVar(&Config.Bridge, "bridge", Config.Bridge, "bridge addr")
	flag.StringVar(&Config.Type, "type", Config.Type, "service type: {client|bridge|server}")
	flag.Parse()
	return Config
}

func (c *config) ConfigStr() string {
	configList := []string{}
	st := reflect.TypeOf(*Config)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		value := reflect.ValueOf(*Config).FieldByName(field.Name)
		configList = append(configList, fmt.Sprintf("%s: %v", field.Name, value))
	}
	return strings.Join(configList, "\n")
}
