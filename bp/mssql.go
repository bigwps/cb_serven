package bp

import (
	"database/sql"
	"fmt"
	"net"
	"time"
	"xc7/log"

	_ "github.com/denisenkom/go-mssqldb"
)

func Msssql(bp BpTask, boom string) {
	host, port, _ := net.SplitHostPort(boom)
	dataSourceName := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%v;encrypt=disable;timeout=%v", host, bp.User, bp.Pass, port, time.Duration(3)*time.Second)
	db, err := sql.Open("mssql", dataSourceName)
	if err != nil {
		log.InfoLog("ERROE:Parameter format error")
	} else {
		db.SetConnMaxLifetime(time.Duration(3) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(3) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			log.WarLog("SQL SERVER-Weak! " + bp.User + ":" + bp.Pass)
			log.RedOut("\n" + "[◢◤]Find: Mssql " + bp.User + ":" + bp.Pass + "\n")

		} else {
			log.InfoLog("False:" + bp.User + ":" + bp.Pass + " in " + boom)
		}
	}
}
