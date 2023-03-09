package utils

import (
	"context"
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"xc7/log"
	"xc7/script"
	"xc7/utils/celpro"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
	"golang.org/x/sync/semaphore"
)

type Task struct {
	Poc    *Ymal
	Target string
}

/*
 *	tips:
 *	已经关闭的channel 能继续读取数据，直到数据读取完
 *  周期性定时器的原理是使用了 go 语言的 channel 机制。NewTicker函数创建了一个周期性定时器，
 *  它会不断的在指定的时间间隔内发送当前时间到channel中。所以当你在程序中调用<-ticker.C时，会一直阻塞，
 *  直到收到channel中的数据，才会继续往下执行。这样就可以在同一时刻只有一个 goroutine 在执行，从而限制了 goroutine 的并发数。
 *  初始版本中WaitGroup并不是用来控制goroutine 的数量，只是用来等待所有goroutine执行完，控制数量由gorou决定，不是很精准。
 *  当然也可是使用可以使用信号量（semaphore），来精确的控制goroutine。
 */

func CheckVuls(pocs []*Ymal, targets []string, gorou int64, proxy string) {

	var wg sync.WaitGroup

	var tasks []Task
	// 将poc和目标加入到结构体数组tasks中
	for _, target := range targets {
		for _, poc := range pocs {
			task := Task{
				poc, target,
			}
			tasks = append(tasks, task)
		}
	}

	// 创建有缓冲通道msg, 用于传递tasks
	msgchan := make(chan Task, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		msgchan <- task
	}
	// 指定信号量的容量
	sem := semaphore.NewWeighted(gorou)

	worker := func(s <-chan Task, wg *sync.WaitGroup) {

		for task := range s {
			sem.Acquire(context.Background(), 1)
			isval, err := execPoc(task, proxy)
			if err != nil {
				fmt.Print(err)
			}
			if isval {
				if task.Poc.Detail.Links != nil {
					log.RedOut("\n" + "[◢◤]Find: " + task.Poc.Name + "\n" + "links: " + task.Poc.Detail.Links[0] + "\n")
				} else {
					log.RedOut("\n" + "[◢◤]Find: " + task.Poc.Name + "\n")
				}
			}
			sem.Release(1)
			wg.Done()
		}
	}

	var i int64
	for i = 0; i < gorou; i++ {
		go worker(msgchan, &wg)
	}

	log.GreenOut("◢◤ " + "Goroutine 数量:" + strconv.Itoa(runtime.NumGoroutine()) + "\n")

	close(msgchan)
	wg.Wait()
}

func execPoc(task Task, proxy string) (bool, error) {

	// 用于存放set,rules
	setMap := make(map[string]interface{})

	target := task.Target
	u, _ := url.Parse(target)
	poc := *task.Poc
	log.InfoLog(poc.Name)

	// NewEnvOption
	options := NewEnvOption()

	// options传入Set全局变量
	options.Update_Set_Options(poc.Set)

	/*
	 *  配置自定义环境NewEnv
	 *  *CelOption需要实现cel.Library（方法CompileOptions）。
	 *  语言规范：https://github.com/google/cel-spec/blob/master/doc/langdef.md
	 */

	env, err := cel.NewEnv(cel.Lib(options))

	if err != nil {
		log.ErrorLog("cel NewEnv 错误!")
		return false, err
	}
	for k, v := range poc.Set {
		var out ref.Val

		switch {
		case v == "newReverse()":
			setMap[k] = NewReverse()

			// 忽略剩余的循环体而直接进入下一次循环，这里获得 {reverse:ceye的解析}，跳出后v = reverse.domain
			continue
		case v == "RedisUnauth()":
			if script.RedisUnauth(u.Hostname(), u.Port()) == true {
				return true, err
			} else {
				return false, err
			}
		case v == "MongodbUnauth()":
			if script.MongodbUnauth(u.Hostname(), u.Port()) == true {
				return true, err
			} else {
				return false, err
			}
		case v == "RsyncUnauth()":
			if script.RsyncUnauth(u.Hostname(), u.Port()) == true {
				return true, err
			} else {
				return false, err
			}
		case v == "Memcache()":
			if script.Memcache(u.Hostname(), u.Port()) == true {
				return true, err
			} else {
				return false, err
			}
		case v == "Zookeeper()":
			if script.Zookeeper(u.Hostname(), u.Port()) == true {
				return true, err
			} else {
				return false, err
			}
		case strings.Contains(v, "target"):
			if v == "target.host" {
				setMap[k] = u.Host
				continue
			}
			if v == "target.url" {
				setMap[k] = u.String()
				continue
			}
		default:
			break

		}

		// 处理其他的函数，比如randomLowercase()
		out, err = Evaluate(env, v, setMap)
		if err != nil {
			return false, nil
		}
		// 一个接口类型的变量中是任何类型的值，要检测动态类型，
		switch value := out.Value().(type) {
		case *celpro.UrlType_C:
			setMap[k] = UrltypeCTostring(value)
			//log.WarLog("UrlType_C_Tostring")
		case int64:
			setMap[k] = int(value)
			//log.WarLog("Expression is  int")
		default:
			setMap[k] = fmt.Sprintf("%v", out)
			//log.WarLog("Expression is  string")
		}
	}

	for rn, rule := range poc.Rules {
		for k1, v1 := range setMap {

			/*
			 * 有时候，map[string]interface{} 有可能存储的是map，也可能是数组等等，那么在取值的时候需要做类型判断，例如：.(类型)
			 * 例如setmap里面第一个reverse是map类型的存放的是ceye的解析，第二个参数是是ceye string
			 * 其他函数中暂时都默认为string,randomInt是int
			 */

			// 这里判断一下除map类型的都替换
			_, ismap := v1.(map[string]string)
			if ismap {
				continue
			}
			// 替换reverseURL
			value := fmt.Sprintf("%v", v1)

			// 这里用for循环是因为有些poc有多个header
			for k2, v2 := range rule.Request.Headers {
				rule.Request.Headers[k2] = strings.ReplaceAll(v2, "{{"+k1+"}}", value)
			}
			rule.Request.Body = strings.ReplaceAll(rule.Request.Body, "{{"+k1+"}}", value)
			rule.Request.Path = strings.ReplaceAll(rule.Request.Path, "{{"+k1+"}}", value)
		}

		URL := target + rule.Request.Path

		switch rule.Request.Method {
		case "GET":
			errr, res := SendRequest("GET", URL, strings.NewReader(rule.Request.Body), rule.Request.Headers, proxy, rule.Request.FollowRedirects)
			if errr != nil {
				log.ErrorLog("GET请求无响应")
				return false, err
			}
			setMap["response"] = res
		case "POST":
			errr, res := SendRequest("POST", URL, strings.NewReader(rule.Request.Body), rule.Request.Headers, proxy, rule.Request.FollowRedirects)
			if errr != nil {
				log.ErrorLog("POST请求无响应")
				return false, err
			}
			setMap["response"] = res
		case "PUT":
			errr, res := SendRequest("PUT", URL, strings.NewReader(rule.Request.Body), rule.Request.Headers, proxy, rule.Request.FollowRedirects)
			if errr != nil {
				log.ErrorLog("PUT请求无响应")
				return false, err
			}
			setMap["response"] = res
		case "DELETE":
			errr, res := SendRequest("DELETE", URL, strings.NewReader(rule.Request.Body), rule.Request.Headers, proxy, rule.Request.FollowRedirects)
			if errr != nil {
				log.ErrorLog("DELETE请求无响应")
				return false, err
			}
			setMap["response"] = res
		case "PATCH":
			errr, res := SendRequest("DELETE", URL, strings.NewReader(rule.Request.Body), rule.Request.Headers, proxy, rule.Request.FollowRedirects)
			if errr != nil {
				log.ErrorLog("PATCH请求无响应")
				return false, err
			}
			setMap["response"] = res
		case "OPTIONS  ":
			log.GreenOut("OPTIONS方法未启用")
		}

		// 判断expression
		var out1 ref.Val
		out1, err = Evaluate(env, rule.Expression, setMap)
		if err != nil {
			log.ErrorLog("Evaluate err")
			return false, err
		}

		// 将r0等结果注册到env
		options.UpdateFunctionOptions(rn, out1)
	}
	// 加载r0&r1等表达式
	env, err = cel.NewEnv(cel.Lib(options))

	//只关心Expression
	out, err := Evaluate(env, poc.Expression, setMap)
	if err != nil {
		return false, err
	}
	if out.Value() == false {
		//log.InfoLog("未发现" + poc.Name)
		return false, err
	}

	log.WarLog(poc.Name + ":" + target)
	return true, nil
}
