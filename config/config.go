package config

type Config interface {
	Get(key string) string
}

type config struct {
}

func NewConfig(path string) Config {
	// 处理配置文件
	return &config{}
}

func (c *config) Get(key string) string {

	switch key {
	case "image-domain":
		return "http://source.lattecake.com/images"

	}

	return ""
}
