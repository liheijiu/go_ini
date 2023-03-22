package bin

import (
	"fmt"
)

// mysql  confg
type Mysqlconfg struct {
	Address string `ini:"address"`
	Port    int    `ini:"port"`
}

// redis  confg
type Redisconfg struct {
	Host string `ini:"host"`
	Port int    `ini:"port"`
	Test bool   `ini:"test"`
}

// merge  confg
type Conf struct {
	Mysqlconfg `ini:"mysql"`
	Redisconfg `ini:"redis"`
}

func Run() {
	var cfg Conf
	//var mx = new(int)

	fileRead, sTypeof, sValueof, err := loadInfo("/go_ini/conf/conf.ini", &cfg)
	if err != nil {
		fmt.Printf("load ini failed, err:%v\n", err)
		return
	}
	err = FileRead(fileRead, sTypeof, sValueof)
	if err != nil {
		fmt.Printf("read ini failed, err:%v\n", err)
		return
	}
	fmt.Printf("%#v\n", cfg)
}
