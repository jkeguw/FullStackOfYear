package upload

type UploadConfig struct {
	SavePath   string   `yaml:"savePath"`
	AllowTypes []string `yaml:"allowTypes"`
	MaxSize    int64    `yaml:"maxSize"` // MB
	UrlPrefix  string   `yaml:"urlPrefix"`
}

var DefaultConfig = UploadConfig{
	SavePath:   "./uploads",
	AllowTypes: []string{".jpg", ".jpeg", ".png"},
	MaxSize:    5, // 5MB
	UrlPrefix:  "/static/uploads",
}
