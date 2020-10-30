package pkg

import (
	"UpgradeWhenDisconnected/src/common/dbm"
	"fmt"
	"strings"
)


const (
	MetaTableName = "meta"
)

type Meta struct {
	//ID int64 `orm:"pk; auto; column(id)"`
	Key string `orm:"column(key); size(256); pk"`
	Type string `orm:"column(type); size(32)"`
	Value string `orm:"column(value); null; type(text)"`
}

func SaveMeta(meta *Meta) error {
	num, err := dbm.DBAccess.Insert(meta)
	fmt.Printf("Insert affected Num : %d, %v", num, err)
	if err == nil || IsNonUniqueNameError(err) {
		return nil
	}
	return nil
}

func IsNonUniqueNameError(err error) bool {
	str := err.Error()
	if strings.HasSuffix(str, "are not unique") || strings.Contains(str, "UNIQUE constraint failed") || strings.HasSuffix(str, "constraint failed") {
		return true
	}
	return false
}

func DeleteMetaByKey(key string) error {
	num ,err := dbm.DBAccess.QueryTable(MetaTableName).Filter("key", key).Delete()
	fmt.Printf("Delete affected num: %d, %v", num, err)
	return err
}

func UpdateMeta(meta *Meta) error {
	num ,err := dbm.DBAccess.Update(meta)
	fmt.Printf("Update affected num: %d, %v", num, err)
	return err
}

func InsertOrUpdate(meta *Meta) error {
	_, err := dbm.DBAccess.Raw("INSERT OR REPLACE INTO meta (key, type, value) VALUE (?,?,?)", meta.Key, meta.Type, meta.Value).Exec()
	fmt.Printf("Update result %v", err)
	return err
}

func UpdateMetaField(key string, col string, value interface{}) error{
	num, err := dbm.DBAccess.QueryTable(MetaTableName).Filter("key", key).Update(map[string]interface{}{col: value})
	fmt.Printf("Update affected Num: %d, %v", num, err)
	return err
}

func UpdateMetaFields(key string, cols map[string]interface{}) error {
	num, err := dbm.DBAccess.QueryTable(MetaTableName).Filter("key", key).Update(cols)
	fmt.Printf("Update affected Num: %d ,%v", num, err)
	return err
}

//return only meta's value
func QueryMeta(key string, condition string) (*[]string, error) {
	meta := new([]Meta)
	_, err := dbm.DBAccess.QueryTable(MetaTableName).Filter(key, condition).All(meta)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, v := range *meta {
		result = append(result, v.Value)
	}
	return &result, nil
}

//return all meta
func QueryAllMeta(key string, condition string) (*[]Meta, error) {
	meta := new([]Meta)
	_, err := dbm.DBAccess.QueryTable(MetaTableName).Filter(key, condition).All(meta)
	if err != nil {
		return nil, err
	}
	return meta, nil
}



