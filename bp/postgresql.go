package bp

import (
	"database/sql"
	"fmt"
	"net"
	"time"
	"xc7/log"

	_ "github.com/lib/pq"
)

func Postgresql(bp BpTask, boom string) {
	host, port, _ := net.SplitHostPort(boom)
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", bp.User, bp.Pass, host, port, "postgres", "disable")
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.InfoLog("ERROE:Parameter format error ")
	} else {
		db.SetConnMaxLifetime(time.Duration(3) * time.Second)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			log.WarLog("Postgres-Weak! " + bp.User + ":" + bp.Pass)
			log.RedOut("\n" + "[◢◤]Find: Postgres " + bp.User + ":" + bp.Pass + "\n")

		} else {
			log.InfoLog("False:" + bp.User + ":" + bp.Pass + " in " + boom)
		}
	}

}
