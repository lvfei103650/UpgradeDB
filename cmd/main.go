package main

import (
	"UpgradeWhenDisconnected/midprocess"
	"fmt"
)


func main() {
	fmt.Println("before start, need to do some tasks")
	//1. 解析config.yaml, 初始化参数
	PodName, ImageTagName := midprocess.InitConfig()

	//2. 注册还是直接打开edgecore/.db
	midprocess.InitDBAccess()

	//3. 一系列操作，解决断网升级的task
	midprocess.StopEdgecore()
	midprocess.RemoveTargetContainers(PodName)

	errProcessDB := midprocess.ProcessDB(PodName, ImageTagName)
	if errProcessDB != nil {
		fmt.Printf("error : %v", errProcessDB)
	}

	//4. 过程当中有一项注意： 检查解压后的本地images的tag是否和要更新的tag一致，如果一致...,否则...。

	//5. 重新启动edgecore
	//midprocess.RestartEdgecore()

	//6. 本地验证是否程序生效，edgecore是否重新拉起新的pod
}
