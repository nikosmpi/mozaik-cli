package mozconfig

type Config struct {
	MySQLPath string            `yaml:"mysql_path"`
	Stagings  map[string]Remote `yaml:"stagings"`
}

type Remote struct {
	DBHost     string `json:"DBHost"`
	DBUser     string `json:"DBUser"`
	DBPass     string `json:"DBPass"`
	SSHUser    string `json:"SSHUser"`
	SSHKeyPath string `json:"SSHKeyPath"`
	SSHHost    string `json:"SSHHost"`
}
