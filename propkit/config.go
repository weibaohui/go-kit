package propkit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/tidwall/gjson"
	"github.com/weibaohui/go-kit/cronkit"
)

type config struct {
	path     string
	fileName string
	//内容，每次更换路径都应该清空
	content string
}

var c *config

var once sync.Once

// 是否开发模式
// 由于配合GIN框架
func IsDevMod() bool {
	return os.Getenv("GIN_MODE") != "release"
}

// 初始化config
func Init() *config {
	once.Do(func() {
		c = &config{}
		if IsDevMod() {
			c.path = "./src/config/"
		} else {
			c.path = "./config/"
		}
		c.read()
	})
	return c
}
func (c *config) check() {
	if c.content == "" {
		c.read()
	}
}
func (c *config) read() {
	if len(c.fileName) == 0 {
		envFileName := "config.prod.json"
		if IsDevMod() {
			envFileName = "config.dev.json"
		}

		//环境关联配置文件
		envConfigPath := fmt.Sprintf("%s%s", c.path, envFileName)
		//默认配置文件
		baseConfigPath := fmt.Sprintf("%s%s", c.path, "config.json")

		//没有设置文件名，取默认的config，且按环境变量进行覆盖
		srcBase, err := ioutil.ReadFile(baseConfigPath)
		if err != nil {
			//没有指明具体配置文件，且默认配置文件不存在
			fmt.Println(err.Error())
			fmt.Println("config.json 不存在")
			return
		}

		if _, err := os.Stat(envConfigPath); err != nil {
			//环境配置不存在，那么加载默认配置文件
			str, err := ioutil.ReadFile(baseConfigPath)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("初始化失败", baseConfigPath)
			}
			c.content = string(str)
			return
		}

		//两个配置，都存在，先加载config.json ,再用config.dev.json 环境配置覆盖
		srcEnv, _ := ioutil.ReadFile(fmt.Sprintf("%s%s", c.path, envFileName))
		var m1, m2 map[string]interface{}
		json.Unmarshal(srcBase, &m1)
		json.Unmarshal(srcEnv, &m2)

		merged := mergeMap(m1, m2)
		finalStr, _ := json.Marshal(merged)
		c.content = string(finalStr)
		return
	}

	//设置了文件名，那么不从默认配置中读取，改为从指定的文件中读取
	fullPath := fmt.Sprintf("%s%s", c.path, c.fileName)
	str, err := ioutil.ReadFile(fullPath)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("初始化失败，内容为空", c.path, c.fileName)
	}
	c.content = string(str)
}

// 使用指定的配置文件
// 默认是用init中的文件
func (c *config) Use(fileName string) *config {
	c.fileName = fileName
	c.read()
	return c
}

func (c *config) Get(key string) string {
	return c.GetString(key)
}
func (c *config) GetString(key string) string {
	return gjson.Get(c.content, key).String()
}

func (c *config) GetBool(key string) bool {
	return gjson.Get(c.content, key).Bool()
}
func (c *config) Exists(key string) bool {
	return gjson.Get(c.content, key).Exists()
}
func (c *config) GetFloat64(key string) float64 {
	return gjson.Get(c.content, key).Float()
}
func (c *config) GetInt64(key string) int64 {
	return gjson.Get(c.content, key).Int()
}

func (c *config) GetInt(key string) int {
	return int(gjson.Get(c.content, key).Int())
}
func (c *config) GetTime(key string) time.Time {
	return gjson.Get(c.content, key).Time()
}

func (c *config) GetUint64(key string) uint64 {
	return gjson.Get(c.content, key).Uint()
}

func (c *config) IsObject(key string) bool {
	return gjson.Get(c.content, key).IsObject()
}
