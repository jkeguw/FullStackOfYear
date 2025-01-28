package config

type UploadConfig struct {
	SavePath   string   `yaml:"savePath"`
	AllowTypes []string `yaml:"allowTypes"`
	MaxSize    int64    `yaml:"maxSize"`
	UrlPrefix  string   `yaml:"urlPrefix"`
}
