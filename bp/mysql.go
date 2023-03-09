package bp

import (
	"database/sql"
	"fmt"
	"net"
	"time"
	"xc7/log"

	_ "github.com/go-sql-driver/mysql"
)

func Mysql(bp BpTask, boom string) {
	host, port, _ := net.SplitHostPort(boom)
	// dataSourceName指定数据源
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/mysql?charset=utf8&timeout=%v", bp.User, bp.Pass, host, port, time.Duration(3)*time.Second)
	// Open函数可能只是验证其参数格式是否正确，实际上并不创建与数据库的连接。如果要检查数据源的名称是否真实有效，应该调用Ping方法
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.InfoLog("ERROE:Parameter format error")
	} else {
		// 表示一个连接保持可用的最长时间
		db.SetConnMaxLifetime(time.Duration(3) * time.Second)
		// 在被标记为失效之前一个连接最长空闲时间
		db.SetConnMaxIdleTime(time.Duration(3) * time.Second)
		db.SetMaxIdleConns(0) // 设置池中空闲连接数的上限
		defer db.Close()
		err = db.Ping()
		if err == nil {
			log.WarLog("MYSQL-Weak! " + bp.User + ":" + bp.Pass)
			log.RedOut("\n" + "[◢◤]Find: Mysql " + bp.User + ":" + bp.Pass + "\n")

		} else {
			log.InfoLog("False:" + bp.User + ":" + bp.Pass + " in " + boom)
		}
	}

}
