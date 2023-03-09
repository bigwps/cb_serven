package bp

import (
	"database/sql"
	"fmt"
	"net"
	"time"
	"xc7/log"

	_ "github.com/sijms/go-ora/v2"
)

func Oracle(bp BpTask, boom string) {
	host, port, _ := net.SplitHostPort(boom)
	dataSourceName := fmt.Sprintf("oracle://%s:%s@%s:%s/orcl", bp.User, bp.Pass, host, port)
	db, err := sql.Open("oracle", dataSourceName)
	if err != nil {
		log.InfoLog("ERROE:Parameter format error")
	} else {
		db.SetConnMaxLifetime(time.Duration(3) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(3) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			log.WarLog("Oracle-Weak! " + bp.User + ":" + bp.Pass)
			log.RedOut("\n" + "[◢◤]Find: Mssql " + bp.User + ":" + bp.Pass + "\n")

		} else {
			log.InfoLog("False:" + bp.User + ":" + bp.Pass + " in " + boom)
		}
	}

}
