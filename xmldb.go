package x

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var dbcfgPaths = []string{
	`N:\PLATFORM\Common\Config\Cfg\zh-CN\GlobalCfg\PlatformCfg\CommData\DBCfg.xml`,
	`C:\Users\Public\CAXA\CAXA EAP CLIENT\1.0\Cfg\zh-CN\GlobalCfg\PlatformCfg\CommData\DBCfg.xml`,
}

func XmlToViper() (*viper.Viper, error) {
	dbjson, err := DBXMLToJSON(dbcfgPaths)
	if err != nil {
		Xlog.Error(err.Error())
		return nil, err
	}

	// PLM.DATABASESERVER.CONNECTIONPARAM.PARAM.#(-Name=="database_connection_timeout").-Value
	const jsonPath = `PLM.DATABASESERVER.CONNECTIONPARAM.PARAM.#(-Name=="%s").-Value`
	const jsonPath2 = `PLM.DATABASESERVER.CONNECTIONPARAM.PARAM.#(-Name=="%s").%s`

	v := viper.New()
	pwd := gjson.Get(dbjson, fmt.Sprintf(jsonPath, "database_user_pwd")).String()
	encrypt := gjson.Get(dbjson, fmt.Sprintf(jsonPath2, "database_user_pwd", "-Encrypt")).String()
	if encrypt == "TRUE" {
		pwd, err = RC4Decrypt(pwd)
		if err != nil {
			Xlog.Error(err.Error())
			return nil, err
		}
	}

	usr := gjson.Get(dbjson, fmt.Sprintf(jsonPath, "database_user")).String()
	host := gjson.Get(dbjson, fmt.Sprintf(jsonPath, "database_server")).String()
	port := 1433

	// fmt.Print("XML HOST1", host)
	// 如果含有,则取逗号前不含,
	if strings.Contains(host, ",") {
		port, err = strconv.Atoi(strings.Split(host, ",")[1])
		if err != nil {
			Xlog.Error(err.Error())
			return nil, err
		}
		host = strings.Split(host, ",")[0]
	}
	dbname := gjson.Get(dbjson, fmt.Sprintf(jsonPath, "database_name")).String()

	v.SetDefault("usr", usr)
	v.SetDefault("pwd", pwd)
	v.SetDefault("host", host)
	v.SetDefault("port", port)
	v.SetDefault("dbname", dbname)
	v.SetDefault("dbextra", "encrypt=disable")
	v.SetDefault("maxIdleConns", 10)
	v.SetDefault("maxOpenConns", 100)

	return v, nil
}

func DBXMLToJSON(paths []string) (string, error) {
	var xml *os.File
	var err error
	for _, path := range paths { // 遍历 dbcfgPaths ,读出到第一个成功的路径
		xml, err = os.Open(path)
		if err != nil {
			if pathErr, ok := err.(*os.PathError); ok {
				Xlog.Info("File does not exist: %s\n" + pathErr.Path)
				continue
			} else {
				Xlog.Error("An error occurred:" + err.Error())
				return "", err
			}
		}
		if xml != nil {
			break
		}
	}

	//xml := strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?><hello>world</hello>`)
	json, err := xj.Convert(xml)
	if err != nil {
		Xlog.Error("That's embarrassing...")
		return "", err
	}

	Xlog.Info(json.String())

	return json.String(), nil
}
