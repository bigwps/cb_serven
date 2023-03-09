package utils

import (
	"bufio"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync"
	"xc7/log"

	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// 兼容xray poc https://docs.xray.cool/#/guide/poc/v2，暂时没有search规则

type Ymal struct {
	Name       string            `yaml:"name"`
	Set        map[string]string `yaml:"set"`
	Transport  string            `yaml:"transport"`  // 传输方式
	Rules      map[string]Rules  `yaml:"rules"`      // 规则描述
	Expression string            `yaml:"expression"` // 规则表达式
	Detail     Detail            `yaml:"detail"`
}

type Rules struct {
	Request    Request `yaml:"request"`
	Expression string  `yaml:"expression"`
	Search     string  `yaml:"search"`
}

type Request struct {
	Method          string            `yaml:"method"`
	Path            string            `yaml:"path"`
	Headers         map[string]string `yaml:"headers"`
	Body            string            `yaml:"body"`
	FollowRedirects bool              `yaml:"follow_redirects"`
}
type Detail struct {
	Links []string `yaml:"links"`
}

// 配置文件
type Ceye struct {
	Identifier string `yaml:"Identifier"`
	CeyeApi    string `yaml:"CeyeApi"`
}

func ReadLine(filename string) []string {
	var res_target []string
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return res_target
	}

	// 逐行读取文件
	file_scan := bufio.NewScanner(file)
	// 每次调用 input.Scan()，即 读入下一行 ，并移除行末的换行符。
	for file_scan.Scan() {
		val := file_scan.Text()
		if val == "" {
			continue
		}
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
		} else {
			val = "http://" + val
		}
		res_target = append(res_target, val)
	}
	if err := file_scan.Err(); err != nil {
		return res_target
	}
	return res_target
}

func LoadAllPoc() ([]string, error) {
	pocdir := "poc/pocymal/"
	var filesname []string
	// os.Start获取文件属性，os.IsNotExist返回已知错误或目录不存在。
	if _, err := os.Stat(pocdir); os.IsNotExist(err) {
		return nil, err
	}
	// filepath.Abs返回参数的绝对路径
	abspath, _ := filepath.Abs(pocdir)
	err := filepath.Walk(abspath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "yml") || strings.HasSuffix(path, "yaml") {
			filesname = append(filesname, path)
		}
		return err
	})
	allpocs := "◢◤ " + "POC 总数量:" + strconv.Itoa(len(filesname))
	log.GreenOut(allpocs)
	return filesname, err
}

func LoadPoc(pocname string) []string {
	pocfiles, err := LoadAllPoc()
	var custompoc []string
	if err != nil {
		log.ErrorLog("加载poc出现了错误")
	}
	pocnames := strings.Split(pocname, "*")
	for _, name := range pocnames {
		for _, poc := range pocfiles {
			if strings.Contains(strings.ToLower(poc), strings.ToLower(name)) {
				custompoc = append(custompoc, poc)
			}
		}
	}

	log.GreenOut("◢◤ " + pocname + " POC 数量:" + strconv.Itoa(len(custompoc)))
	return custompoc
}

// 创建一个池
var pool = sync.Pool{
	New: func() interface{} {
		return new(Ymal)
	},
}

func ParseYml(pocfiles []string) ([]*Ymal, error) {
	// 结构体数组
	var pocs []*Ymal
	for _, file := range pocfiles {
		// 从池中获取一个结构体实例
		p := pool.Get().(*Ymal)
		// ioutil导入配置文件
		yamlfile, err := ioutil.ReadFile(file)
		if err != nil {
			// 将实例放回池中
			pool.Put(p)
			return nil, err
		}
		// 初始化一个结构体实例[指针],这里使用池来管理初始化结构体提升性能。
		// p := &Ymal{}
		// 使用yaml.Unmarshal将配置文件读取到结构体中
		err = yaml.Unmarshal(yamlfile, p)
		if err != nil {
			pool.Put(p)
			return nil, err
		}
		pocs = append(pocs, p)
	}
	return pocs, nil
}

// 读取config文件
func LoadConfig() (ceyes string, api string) {
	ymalfile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.ErrorLog("导入config.yaml文件出错!")
	}
	c := &Ceye{}
	err = yaml.Unmarshal(ymalfile, c)
	if err != nil {
		log.ErrorLog("读取config.yaml文件出错!")
	}
	return c.Identifier, c.CeyeApi
}
