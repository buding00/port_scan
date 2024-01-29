package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//
//var mu sync.Mutex

//var AllIp chan string

var AllIp = make(chan string, 500)
var wg sync.WaitGroup
var f *os.File
var err error

// 定义常用端口合集列表
var commonPorts = []string{
	"2094", "3306", "3700", "11211", "80", "513", "1521",
	"28017", "2083", "7001", "50030", "4848", "8089", "23",
	"9200", "161", "5901", "2601", "20", "6379", "50070",
	"10000", "5560", "3311", "5432", "139", "8080", "22",
	"50000", "5000", "69", "53", "5984", "8000", "873", "512",
	"5632", "3", "67", "9000", "7002", "2082", "443", "445",
	"21", "8888", "143", "7778", "6082", "8083", "1025", "9080",
	"3389", "110", "50060", "25", "89", "5902", "514", "389",
	"3128", "3312", "4440", "111", "9300", "9090", "68", "81",
	"2604", "27018", "2222", "8649", "9081", "1433", "5900",
	"27017", "22122"}

func initUse(ip string) {
	// 格式化当前时间 年月日时分秒
	now := time.Now().Format("2006-01-02-15-04-05——")
	// 获取当前目录
	dir, _ := os.Getwd()
	fmt.Printf("当前目录:%s\n", dir)
	// 获取当前系统分隔符
	//sep := string(os.PathSeparator)
	// 构造当前文件名
	fileName := fmt.Sprintf("%s%s%s.txt", dir, "\\", now+ip)
	// 创建一个文件对象 追加，写入权限
	f, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("文件创建失败")
		fmt.Println(err.Error())
		return
	}
}
func worker(ports chan string) {
	for address := range ports {
		//fmt.Printf("正在扫描%s", address)
		conn, err := net.DialTimeout("tcp", address, time.Second*1)
		if err != nil {
			wg.Done()
			continue
		} else {
			err = conn.Close()
			wg.Done()
			_, err := f.WriteString(address + "\n")
			if err != nil {
				fmt.Println("写入失败")
			}
			fmt.Println(address + "端口开放!")
		}
	}
}
func main() {
	// 加个时间计算
	start := time.Now()
	var p string
	var PP string
	var ppp string
	var t int
	flag.StringVar(&p, "p", "", "扫描单个ip  例如:10.201.1.10")
	flag.StringVar(&PP, "pp", "", "扫描整个c段  例如:10.201.1.0")
	flag.StringVar(&ppp, "ppp", "", "扫描整个B段  例如:10.201.0.0 非特殊情况不建议使用")
	flag.IntVar(&t, "t", 500, "指定扫描线程数,默认是500")
	flag.Parse()
	if p != "" && ppp != "" {
		fmt.Println("请一次传入单个ip或者c段或者B段")
		return
	}
	if p != "" && PP != "" {
		fmt.Println("请一次传入单个ip或者c段或者B段")
		return
	}
	if ppp != "" && PP != "" {
		fmt.Println("请一次传入单个ip或者c段或者B段")
		return
	}
	if PP != "" && p != "" && ppp != "" {
		fmt.Println("请一次传入单个ip或者c段或者B段")
		return
	}
	if p == "" && PP == "" && ppp == "" {
		fmt.Println("请一次传入单个ip或者c段或者B段")
		return
	}
	if p != "" {
		address := net.ParseIP(p)
		if address == nil {
			fmt.Println("ip地址格式错误")
			return
		} else {
			scanIp(p, t)
			wg.Wait()
			fmt.Printf("扫描完成,文件写入成功 文件名：%s\n", f.Name())
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					fmt.Println("文件关闭失败")
				}
			}(f)
		}
	}
	if PP != "" {
		address := net.ParseIP(PP)
		if address == nil {
			fmt.Println("ip地址格式错误")
			return
		} else {
			scanC(PP, t)
			wg.Wait()
			fmt.Printf("扫描完成,文件写入成功 文件名：%s\n", f.Name())
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					fmt.Println("文件关闭失败")
				}
			}(f)
		}
	}
	if ppp != "" {
		fmt.Println("暂不支持扫描整个B段")
		//fmt.Println("扫描整个B段")
		//address := net.ParseIP(ppp)
		//if address == nil {
		//	fmt.Println("ip地址格式错误")
		//	return
		//} else {
		//	scanB(ppp, t)
		//	wg.Wait()
		//	fmt.Printf("扫描完成,文件写入成功 文件名：%s\n", f.Name())
		//	defer func(f *os.File) {
		//		err := f.Close()
		//		if err != nil {
		//			fmt.Println("文件关闭失败")
		//		}
		//	}(f)
		//}
	}
	// 计算时间
	end := time.Now()
	fmt.Println("耗时：", end.Sub(start))
	// 捕获t参数解析错误
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("参数解析错误")
		}
	}()

}

// 扫描单个ip
func scanIp(ip string, t int) {

	initUse(ip)
	// 50个协程去扫描
	for i := 0; i < t; i++ {
		go worker(AllIp)
	}
	// 单个ip扫描全量端口
	for i := 0; i < 65536; i++ {
		wg.Add(1)
		AllIp <- ip + ":" + strconv.Itoa(i)
	}
	//for _, port := range commonPorts {
	//	wg.Add(1)
	//	AllIp <- ip + ":" + port
	//}
	defer close(AllIp)
}

// 扫描整个c段
func scanC(ip string, t int) {
	var resIp string
	initUse(ip)
	split := strings.Split(ip, ".")[0:3]
	for i := 0; i < len(split); i++ {
		resIp += split[i] + "."
	}
	resIp = resIp[0 : len(resIp)-1]

	for i := 0; i < t; i++ {
		go worker(AllIp)
	}
	// 把c段的ip放入到一个切片中
	var ips []string
	for i := 1; i < 255; i++ {
		ips = append(ips, fmt.Sprintf("%s.%d", resIp, i))
	}

	for _, i2 := range ips {
		for _, port := range commonPorts {
			wg.Add(1)
			AllIp <- i2 + ":" + port
		}
	}
	defer close(AllIp)
}

// 扫描整个B段
func scanB(ip string, t int) {
	var resIp string
	initUse(ip)
	split := strings.Split(ip, ".")[0:2]

	for i := 0; i < len(split); i++ {
		resIp += split[i] + "."
	}
	resIp = resIp[0 : len(resIp)-1]

	for i := 0; i < t; i++ {
		go worker(AllIp)
	}
	// 把c段的ip放入到一个切片中
	var ips []string
	for i := 0; i < 255; i++ {
		ipB := fmt.Sprintf("%s.%d", resIp, i)
		for i := 1; i < 255; i++ {
			ips = append(ips, fmt.Sprintf("%s.%d", ipB, i))
		}
	}

	for _, i2 := range ips {
		for _, port := range commonPorts {
			wg.Add(1)
			AllIp <- i2 + ":" + port
		}
	}
	close(AllIp)
}
