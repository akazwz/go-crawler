package config

type Conf struct {
	URL       string `yaml:"url"`
	Token     string `yaml:"token"`
	Bucket    string `yaml:"bucket"`
	Org       string `yaml:"org"`
	SecretId  string `yaml:"secretId"`
	SecretKey string `yaml:"secretKey"`
}
