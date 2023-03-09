/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import "xc7/cmd"

func main() {
	cmd.Execute()
}

/*
日志：

  Data         File      Version      改动
--------- | --------- | --------- | ------------------------------------------------------------------
2023/1/6  |  file.go  |    0.1    | 优化解析Yaml函数ParseYml(),使用池管理初始化结构体，避免每次循环都初始化结构体。
2023/1/6  |  exec.go  |    0.1    | 使用信号量（semaphore），来精确的控制goroutine。
2023/1/9  |  log.go   |    0.1    | 使用go-colorable着色编辑器库解决在cmd中颜色无法正常显示的问题
2023/2/3  | request.go|    0.1    | 解决空指针的问题，经过dlv调试发现没有对请求的err进行判断而继续执行造成ParseUrl函数传入空指针。
2023/2/9  |  script   |    1.0    | 对于go文件的动态加载初步方案是利用接口和反射机制，后改为yaml文件定义函数执行，暂时放弃动态加载。
2023/2/11 |  exec.go  |    1.0    | 修复了若干bug
*/
