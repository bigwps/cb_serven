package bp

import (
	"context"
	"sync"
	"xc7/log"

	"golang.org/x/sync/semaphore"
)

type BpTask struct {
	User string
	Pass string
}

func CheckWeak(boom string, name string, gogou int64) {

	mysql := []string{"root", "mysql"}
	mssql := []string{"sa", "sql"}
	postgresql := []string{"postgres", "admin", "vulhub"}
	oracle := []string{"sys", "system", "admin", "test", "web", "orcl"}
	Passwords := []string{
		"vulhub",
		"123456",
		"123.com",
		"admin",
		"admin123",
		"admin1234",
		"admin666",
		"admin888",
		"root",
		"toor",
		"root123",
		"abc123",
		"abc1234",
		"",
		"pass123",
		"pass@123",
		"password",
		"123123",
		"654321",
		"111111",
		"123",
		"1",
		"mysql123",
		"mysql",
		"sqlserver",
		"mssql",
		"oracle",
		"oracle123",
		"postgresql",
		"admin@123",
		"admin123!@#",
		"{user}",
		"{user}1",
		"{user}111",
		"{user}123",
		"{user}@123",
		"{user}_123",
		"{user}#123",
		"{user}@111",
		"{user}@2019",
		"{user}@123#4",
		"P@ssw0rd!",
		"P@ssw0rd",
		"Passw0rd",
		"qwe123",
		"12345678",
		"administrator",
		"test",
		"test123",
		"123qwe",
		"123qwe!@#",
		"123456789",
		"123321",
		"666666",
		"a123456.",
		"123456~a",
		"123456!a",
		"000000",
		"1234567890",
		"8888888",
		"!QAZ2wsx",
		"1qaz2wsx",
		"abc123",
		"abc123456",
		"1qaz@WSX",
		"a11111",
		"a12345",
		"Aa1234",
		"Aa1234.",
		"Aa12345",
		"a123456",
		"a123123",
		"Aa123123",
		"Aa123456",
		"Aa12345.",
		"sysadmin",
		"system",
		"1qaz!QAZ",
		"2wsx@WSX",
		"qwe123!@#",
		"Aa123456!",
		"A123456s!",
		"sa123456",
		"1q2w3e",
		"Charge123",
		"123456",
		"zaq1@WSX",
		"mongodb123",
		"qweasdzxc",
		"12345",
		"1234",
		"qwerty",
		"1q2w3e4r",
		"qazwsx",
		"123qwe",
		"123qaz",
		"0000",
		"1234567",
		"123456qwerty",
		"password123",
		"q1w2e3r4",
		"okmnji",
		"aucma",
		"sinop",
		"%user%",
		"%user%123",
		"%user%1234",
		"%user%123456",
		"%user%12345",
		"%user%@123",
		"%user%@123456",
		"%user%@12345",
		"%user%#123",
		"%user%#123456",
		"%user%#12345",
		"%user%_123",
		"%user%_123456",
		"%user%_12345",
		"%user%123!@#",
		"%user%!@#$",
		"%user%!@#",
		"%user%~!@",
		"%user%!@#123",
		"%user%2022",
		"%user%2021",
		"%user%2020",
		"%user%2019",
		"%user%2018",
		"%user%2017",
		"%user%2016",
		"%user%2015",
		"%user%@2017",
		"%user%@2016",
		"%user%@2015",
	}

	vulnmap := map[string]func(BpTask, string){
		"mssql":   Msssql,
		"mysql":   Mysql,
		"oracle":  Oracle,
		"postgre": Postgresql,
	}

	var wg sync.WaitGroup
	var poc []BpTask
	switch {
	case name == "mssql":
		poc = addscan(mssql, Passwords)
	case name == "mysql":
		poc = addscan(mysql, Passwords)
	case name == "oracle":
		poc = addscan(oracle, Passwords)
	case name == "postgre":
		poc = addscan(postgresql, Passwords)
	default:
		log.WarLog("不支持爆破此协议")
		return
	}
	msgchan := make(chan BpTask, len(poc))
	for _, bp := range poc {
		wg.Add(1)
		msgchan <- bp
	}
	sem := semaphore.NewWeighted(gogou)
	worker := func(s <-chan BpTask, wg *sync.WaitGroup, name string) {
		for bp := range s {
			sem.Acquire(context.Background(), 1)
			vulnmap[name](bp, boom)
			sem.Release(1)
			wg.Done()
		}

	}
	var i int64
	for i = 0; i < gogou; i++ {
		go worker(msgchan, &wg, name)

	}
	wg.Wait()

}

func addscan(sqlname []string, pass []string) []BpTask {
	var bptasks []BpTask
	for _, sqlname := range sqlname {
		for _, pass := range pass {
			bptask := BpTask{
				sqlname, pass,
			}
			bptasks = append(bptasks, bptask)
		}
	}
	return bptasks
}
