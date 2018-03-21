使用json作为配置文件
文件内容为json语法
核心由github.com/tidwall/gjson驱动
配置文件放在config/下
如
config/config.json
config/config.prod.json

ENV GIN_MODE=release
prod环境下，会先读取config.json，然后用config.prod.json进行覆盖