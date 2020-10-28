package main

import (
	"UpgradeWhenDisconnected/src/common/dbm"
	"UpgradeWhenDisconnected/src/pkg"
	"fmt"
)

const (
	// DataBaseDriverName is sqlite3
	DataBaseDriverName = "sqlite3"
	// DataBaseAliasName is default
	DataBaseAliasName = "default"
	// DataBaseDataSource is edge.db
	DataBaseDataSource = "/var/lib/kubeedge/edgecore.db"
)

func main() {
	fmt.Println("before start, need to do some tasks")
	//1. 解析config.yaml, 初始化参数
	var c pkg.Conf
	c.GetConf()
	fmt.Printf("c podName: %s, imageTagName: %s", c.PodName, c.ImageTagName)

	//2. 注册还是直接打开edgecore/.db
	dbm.InitDBConfig(DataBaseDriverName, DataBaseAliasName, DataBaseDataSource)


	//3. 一系列操作，解决
	pkg.StopEdgecore()
	pkg.RemoveTargetContainers()

	errProcessDB := pkg.ProcessDB(c.PodName, c.ImageTagName)
	if errProcessDB != nil {
		fmt.Printf("error : %v", errProcessDB)
	}

	//4. 过程当中有一项注意： 检查解压后的本地images的tag是否和要更新的tag一致，如果一致...,否则...。


	pkg.RestartEdgecore()
}