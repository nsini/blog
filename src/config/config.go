package config

import "github.com/Unknwon/goconfig"

const (
	SectionServer = "server"
	ImageFilePath = "image_file_path"
	ImageDomain   = "image_domain"
)

type Config struct {
	*goconfig.ConfigFile
}

func NewConfig(path string) (*Config, error) {
	// 处理配置文件

	cfg, err := goconfig.LoadConfigFile(path)
	if err != nil {
		return nil, err
	}
	return &Config{cfg}, nil
}

func (c *Config) GetString(section, key string) string {
	val, _ := c.GetValue(section, key)
	return val
}
