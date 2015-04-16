package cfg

import (
	"testing"
	"fmt"
)

//测试Ini格式的配置文件解析.
func TestInI(t *testing.T) {
	if iniFile, err := LoadFile("./config.ini"); err != nil {
		fmt.Printf("%s\n", "config.ini error")
	}else{
		fmt.Printf("file %v\n", iniFile)
		if port, status := iniFile.Get("","port"); status != false {
			fmt.Printf("default port %s\n", port)
		}
		if user, status := iniFile.Get("database","user"); status != false {
			fmt.Printf("databbase user %s\n", user)
		}
		if pass, status := iniFile.Get("","port"); status != false {
			fmt.Printf("databbase password %s\n", pass)
		}

		secDb := iniFile.Section("database")
		fmt.Printf("section database %v\n", secDb)

		secNot := iniFile.Section("databasex")
		if secNot == nil {
			fmt.Printf("section database is nil\n")
		}else{
			fmt.Printf("section database %v\n", secNot)
		}

		var s map[int]string
		fmt.Printf("s x %v\n", s[0])

		f := File{}

		if f["f"] == nil {
			fmt.Printf("f f is nil\n")
		} else {
			fmt.Printf("f f %v\n", f["f"])
		}

	}
}
