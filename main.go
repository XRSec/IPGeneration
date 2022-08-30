package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	buildTime, commitId, versionData, author, filePath string
	help, version                                      bool
)

const (
	InfoColor    = "\033[1;34m%s\033[0m\n"
	NoticeColor  = "\033[1;36m%s\033[0m\n"
	WarningColor = "\033[1;33m%s\033[0m\n"
	ErrorColor   = "\033[1;31m%s\033[0m\n"
	DebugColor   = "\033[0;36m%s\033[0m\n"
)

func init() {
	fmt.Printf(InfoColor, "-------------------------------")
	flag.StringVar(&filePath, "f", "ip.txt", "IP-List FilePath")
	flag.BoolVar(&help, "h", false, "Display help information")
	flag.BoolVar(&help, "help", false, "Display help information")
	flag.BoolVar(&version, "v", false, "version")
	flag.BoolVar(&version, "version", false, "version")
	flag.Parse()
	if help {
		flag.Usage()
		fmt.Printf(WarningColor, "-------------------------------")
		os.Exit(0)
	}
	// Version
	if version {
		fmt.Printf("%-15v%v", "Version: ", fmt.Sprintf(NoticeColor, versionData))
		fmt.Printf("%-15v%v", "BuildTime: ", fmt.Sprintf(NoticeColor, buildTime))
		fmt.Printf("%-15v%v", "Author: ", fmt.Sprintf(NoticeColor, author))
		fmt.Printf("%-15v%v", "CommitId: ", fmt.Sprintf(NoticeColor, commitId))
		fmt.Printf(InfoColor, "-------------------------------\n")
		os.Exit(0)
	}
}

//ip到数字
func ip2Int(ip string) uint32 {
	var long uint32
	err := binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	if err != nil {
		fmt.Printf(ErrorColor, fmt.Sprintf("ip2Int error: %v %v", ip, err))
		os.Exit(1)
	}
	return long
}

//数字到IP
func int2IP(ipInt int64) string {
	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt(ipInt&0xff, 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

func main() {
	defer fmt.Printf(InfoColor, "-------------------------------\n")
	if err := resultIPS(); err != nil {
		fmt.Printf(ErrorColor, err)
	}
}

func resultIPS() error {
	fp, err := os.Open(filePath)
	if err != nil {
		return errors.New("文件打开失败: " + err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf(ErrorColor, "关闭文件失败: "+err.Error())
		}
	}(fp) // 关闭文本流

	buf := bufio.NewScanner(fp) // 读取文本数据
	var (
		ipLists []string
		i       = 0
	)
	for {
		if !buf.Scan() {
			break // 文件读完了,退出for
		}
		i++
		fmt.Printf(NoticeColor, fmt.Sprintf("当前处理的是第%v行", i))

		var (
			str1 = strings.Split(buf.Text(), "-")
			ip0  = ip2Int(str1[0])
			ip1  = ip2Int(str1[1])
		)

		if len(str1[1]) <= 3 {
			return errors.New("不支持 /24 /26 格式")
		}

		fmt.Printf(NoticeColor, fmt.Sprintf("当前IP段共有%v个IP", ip1-ip0))

		for l := ip0; l <= ip1; l++ {
			ipLists = append(ipLists, int2IP(int64(l)))
		}
		fmt.Printf(InfoColor, "-----------")
	}

	fmt.Printf(NoticeColor, fmt.Sprintf("总共查询到%v个IP", len(ipLists)))

	ips := strings.Join(ipLists, "\n")
	err = ioutil.WriteFile("new-"+filePath, []byte(ips), 0666)
	return nil
}
