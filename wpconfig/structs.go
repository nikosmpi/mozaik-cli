package wpconfig

type WPConfig struct {
	MySQLPath   string    `json:"MySQLPath"`
	DBHost      string    `json:"DBHost"`
	DBUser      string    `json:"DBUser"`
	DBPass      string    `json:"DBPass"`
	DBName      string    `json:"DBName"`
	DBPrefix    string    `json:"DBPrefix"`
	Staging     Remote    `json:"Staging"`
	ReplaceList []Replace `json:"ReplaceList"`
}
type Remote struct {
	DBHost     string `json:"DBHost"`
	DBUser     string `json:"DBUser"`
	DBPass     string `json:"DBPass"`
	DBName     string `json:"DBName"`
	DBPrefix   string `json:"DBPrefix"`
	SSHUser    string `json:"SSHUser"`
	SSHKeyPath string `json:"SSHKeyPath"`
	SSHHost    string `json:"SSHHost"`
}

type Replace struct {
	Old string `json:"old"`
	New string `json:"new"`
}
