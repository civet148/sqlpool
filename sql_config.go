package sqlpool

type SqlConfig struct {
	Mysql Mysql `toml:"mysql"`
	Queue Queue `toml:"queue"`
	Log   Log   `toml:"log"`
	Redis Redis `toml:"redis"`
}

type Log struct {
	LogPath  string `toml:"logPath"`  //log file to write
	LogLevel string `toml:"logLevel"` // "debug" "info" "warn" "error"
}

type Mysql struct {
	Masters []Masters `toml:"masters"`
	Slaves  []Slaves  `toml:"slaves"`
}

type Masters struct {
	Name   string `toml:"name"`
	Dsn    string `toml:"dsn"`
	Active int    `toml:"active"`
	Idle   int    `toml:"idle"`
}

type Slaves struct {
	Name   string `toml:"name"`
	Dsn    string `toml:"dsn"` // mysql://root:123456@127.0.0.1/test?charset=utf8mb4
	Active int    `toml:"active"`
	Idle   int    `toml:"idle"`
	Slave  bool   `toml:"slave"`
}

type Queue struct {
	Routines int  `toml:"routines"`
	Timeout  int  `toml:"timeout"`
	Cap      int  `toml:"cap"`
	Debug    bool `toml:"debug"`
}

type Redis struct {
	Password       string   `toml:"password"`       //auth password
	Index          int      `toml:"index"`          //db index
	MasterHost     string   `toml:"masterHost"`     //master node
	ReplicateHosts []string `toml:"replicateHosts"` //replicate node
	ConnTimeout    int      `toml:"connTimeout"`    //millisecond
	ReadTimeout    int      `toml:"readTimeout"`    //millisecond
	WriteTimeout   int      `toml:"writeTimeout"`   //millisecond
	KeepAlive      int      `toml:"keepAlive"`      //second
	AliveTime      int      `toml:"aliveTime"`      //second
}
