package mysql

import (
	"strings"

	"fmt"

	"github.com/yezihack/gm2m/common"
	"github.com/yezihack/gm2m/conf"
)

//获取表名的列表
func (d *DbTools) GetTableList() (tableResult map[string]string, err error) {
	result, err := GetMasterDB().Query("show tables")
	if err != nil {
		return
	}
	tableList := make([]string, 0)
	//获取配置文件里的table_list设定
	var ConfTableList []string
	ConfTableList, err = common.GetConfTables()
	for _, mapVal := range result {
		for _, tableName := range mapVal {
			if len(ConfTableList) > 0 {
				var okTable bool
				for k := range ConfTableList {
					if ConfTableList[k] == tableName {
						okTable = true
						break
					}
				}
				if okTable {
					tableList = append(tableList, strings.TrimSpace(tableName.(string)))
				} else {
					fmt.Println("过滤的表:", tableName)
				}
			}

		}
	}
	tableResult = make(map[string]string)
	for _, table := range tableList {
		tableInfo, _ := GetMasterDB().Query("select table_comment from information_schema.tables where table_name = ?", table)
		if err != nil {
			return
		}
		tableResult[table] = tableInfo[0]["table_comment"].(string)
	}
	return
}

//获取表结构详情
func (d *DbTools) GetTableDesc(tableName string) (reply []*TableDesc, err error) {
	result, err := GetMasterDB().Query("select column_name,data_type, column_key, is_nullable,column_default,column_type, column_comment from information_schema.columns where table_name = ?", tableName)
	if err != nil {
		return
	}
	reply = make([]*TableDesc, 0)
	i := 0
	for _, row := range result {
		var keyBool bool
		if strings.ToUpper(row["column_key"].(string)) == "PRI" {
			keyBool = true
		}
		oriType := row["data_type"].(string)
		var columnDefault string
		val, ok := row["column_default"].(string)
		if ok {
			columnDefault = val
		}
		reply = append(reply, &TableDesc{
			Index:            i,
			ColumnName:       row["column_name"].(string),
			GoColumnName:     d.T.Capitalize(row["column_name"].(string)),
			OriMysqlType:     oriType,
			UpperMysqlType:   strings.ToUpper(oriType),
			GolangType:       conf.MysqlTypeToGoType[oriType],
			MysqlNullType:    conf.MysqlTypeToGoNullType[oriType],
			ColumnComment:    row["column_comment"].(string),
			IsNull:           row["is_nullable"].(string),
			DefaultValue:     columnDefault,
			ColumnTypeNumber: row["column_type"].(string),
			PrimaryKey:       keyBool,
		})
		i++
	}
	return
}
