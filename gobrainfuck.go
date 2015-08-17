package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)



/**
 * FKVM
 * 表达brainfuck使用的机器模型，连续字节内存块
 */

type FKVM struct {
	code     string      //执行的代码
	mem      []byte      //模拟内存
	pos      int         //记录当前操作格子位置
	whilemap map[int]int //记录循环的下标
}

func NewFKVM() *FKVM {
	t := new(FKVM)
	t.mem = make([]byte, 4) //4096
	t.whilemap = make(map[int]int, 128)
	return t
}

func (this *FKVM) run() {
	var pc int = 0 //执行下标

	for pc < len(this.code) {
		switch this.code[pc] {
		case '>':
			this.pos++
			if len(this.mem) <= this.pos {
				this.mem = append(this.mem, 0)
			}

		case '<':
			this.pos--
		case '+':
			this.mem[this.pos]++
		case '-':
			this.mem[this.pos]--
		case '.':
			fmt.Printf("%c", this.mem[this.pos])
		case ',':
			fmt.Scanf("%c", this.mem[this.pos])
		case '[':
			if this.mem[this.pos] == 0 {
				pc = this.whilemap[pc]
			}
		case ']':
			if this.mem[this.pos] != 0 {
				pc = this.whilemap[pc]
			}
		}
		pc++
	}
}

func (this *FKVM) parse(code string) *FKVM {
	codes := make([]byte, 0) //解析后的代码
	pcstack := make([]int, 0)

	//记录[,对应的],索引(指令)位置
	whilemap := make(map[int]int, 128)
	pc := 0

	for _, char := range code {

		switch char {
		case '>', '<', '+', '-', '.', ',', '[', ']':
			codes = append(codes, byte(char))
			if char == '[' {
				pcstack = append(pcstack, pc)
			} else if char == ']' {

				last := len(pcstack) - 1
				left := pcstack[last]
				pcstack = pcstack[:last]
				right := pc
				whilemap[right] = left
				whilemap[left] = right

			}
			pc++
		}
	}

	this.code = string(codes)
	this.whilemap = whilemap
	return this
}

func main() {

	flag.Parse()
	var path string = flag.Arg(0)

	if f, err := os.Stat(path); path == "" || err != nil || f.IsDir() {
		fmt.Println("File not found...")
		os.Exit(1)
	}

	//获取文本code
	code := readFile(path)

	//启动VM开始执行
	NewFKVM().parse(code).run()

}

/*读取文件*/
func readFile(path string) string {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic(err)
	}
	return string(f)
}
