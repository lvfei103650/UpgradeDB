package pkg

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"k8s.io/api/core/v1"
)
const (
	TypeName = "pod"
)

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
	cmd := exec.Command("systemctl stop edgecore")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

//step 2:
//	docker stop `docker ps | grep xxx | awk'{print $1}'`
//	docker rm  `docker ps | grep xxx | awk'{print $1}'`
func RemoveTargetContainers() {
	cmd := exec.Command("docker stop `docker ps | grep xxx | awk'{print $1}'`")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	cmd = exec.Command("docker rm  `docker ps | grep xxx | awk'{print $1}'`")
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
	var podstructs  v1.PodSpec
	content, err1 := json.Marshal(podMeta)
	if err1 != nil {
		fmt.Printf("err: %v", err)
	}
	json.Unmarshal(content, &podstructs)
	podstructs.Containers[0].Image = imagesTag

	contentAfter, err := json.Marshal(podstructs)
	json.Unmarshal(contentAfter, &podMeta)

	//3. update db
	err = InsertOrUpdate(&podMeta)
	if err != nil {
		fmt.Printf("errpr : %v", err)
		return err
	}
	return err
}

//step 4:
func RestartEdgecore() {
	cmd := exec.Command("systemctl restart edgecore")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}