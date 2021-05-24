package config

type Config struct {
  Server struct {
    Port string `yaml:"port"`
    Host string `yaml:"host"`
  } `yaml:"server"`

  Mysql struct {
    Port string `yaml:"port"`
    Host string `yaml:"host"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Dbname string `yaml:"dbname"`
  } `yaml:"database"`

  Rabbitmq struct {
    Port string `yaml:"port"`
    Host string `yaml:"host"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
  } `yaml:"rabbitmq"`
}
