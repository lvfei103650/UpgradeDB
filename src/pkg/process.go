package pkg

import (
	"UpgradeWhenDisconnected/src/common/dbm"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"k8s.io/api/core/v1"
	"os/exec"
	"strings"
)
const (
	TypeName = "pod"
	// DataBaseDriverName is sqlite3
	DataBaseDriverName = "sqlite3"
	// DataBaseAliasName is default
	DataBaseAliasName = "default"
	// DataBaseDataSource is edge.db
	DataBaseDataSource = "/var/lib/kubeedge/edgecore.db"
)

func InitConfig() (podName string, imageTagName string) {
	var c Conf
	c.GetConf()
	fmt.Printf("c podName: %s, imageTagName: %s", c.PodName, c.ImageTagName)
	return c.PodName, c.ImageTagName
}

func InitDBAccess() {
	orm.RegisterModel(new(Meta))
	dbm.InitDBConfig(DataBaseDriverName, DataBaseAliasName, DataBaseDataSource)
}

func iscontainSubString(key string, subKey string) bool {
	return strings.Contains(key, subKey)
}

func queryByFuzzyString(key string) (Meta, error){
	podMetasRecord, err:= QueryAllMeta("type", TypeName)
	if err != nil {
		fmt.Printf("list pods failed, error: %v", err)
		return Meta{}, err
	}
	for _, v := range *podMetasRecord {
		if iscontainSubString(v.Key, key) {
			return v, nil
		}
	}
	return Meta{}, err
}

func checkSysStopEdgecoreExists() {
	path, err := exec.LookPath("systemctl stop edgecore")
	if err != nil {
		fmt.Printf("didn't find 'systemctl stop edgecore' executable\n")
	} else {
		fmt.Printf("'systemctl stop edgecore' executable is in '%s'\n", path)
	}
}

//step 1 :
func StopEdgecore() {
	//checkSysStopEdgecoreExists()
	cmd := exec.Command("sh", "-c", "systemctl stop edgecore")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

//step 2:
//	docker stop `docker ps | grep xxx | awk'{print $1}'`
//	docker rm  `docker ps | grep xxx | awk'{print $1}'`
func RemoveTargetContainers(key string) {

	cmdStopStr := fmt.Sprintf("%s %s %s","docker stop `docker ps | grep ", key, "| awk '{print $1}'`")

	cmd := exec.Command("sh", "-c", cmdStopStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	cmdRemoveStr := fmt.Sprintf("%s %s %s","docker rm `docker ps | grep ", key, " | awk '{print $1}'`")
	cmd = exec.Command("sh", "-c", cmdRemoveStr)
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

//step 3:
func ProcessDB(key string, imagesTag string) error{
	//1. queryByFuzzyString get podname
	podMeta, err := queryByFuzzyString(key)
	if err != nil {
		fmt.Printf("err: %v", err)
		return err
	}

	//2. modify imageTag in meta.value
	var podstructs  v1.Pod
	json.Unmarshal([]byte(podMeta.Value), podstructs)
	podstructs.Spec.Containers[0].Image = imagesTag
	contentAfter, _ := json.Marshal(podstructs)

	//3. update db
	meta := &Meta{
		Key: podMeta.Key,
		Type: podMeta.Type,
		Value: string(contentAfter)}
	err2 := UpdateMeta(meta)
	if err != nil {
		fmt.Printf("errpr : %v", err2)
		return err
	}
	return err
}

//step 4:
func RestartEdgecore() {
	cmd := exec.Command("sh", "-c", "systemctl restart edgecore")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}