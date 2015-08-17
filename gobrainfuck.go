package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

/*
> 增加数据指针 (使其指向当前右边内存单元).
< 减少数据指针(使其指向当前左边内存单元).
+ 对当前内存单元加 1
- 对当前内存单元减 1
. 输出当内存单元
, 接受一个字节的输入，将其放到当前数据指针指向的内存单元
[ 如果当前数据指针指向的单元，值为非0， 进入循环，执行紧靠 [ 后面的指令；否则，向前跳转到与此 [ 匹配的 ] 之后开始执行
] 如果当前数据指针指向的单元，值为非0，向后跳转，回到与此 ] 匹配的 [ 后面执行， 否则，正常流程，继续向下执行


> 右
< 左
+ 上
- 下
[ 始
] 终
. 写
, 读

遇到什么都没写的格子就当里面写了 0
右：向右移动一个格子。嗯你盯着它看就行了，什么都不用做
左：向左移动一个格子
上：给格子里的数字加上 1，擦掉原来的数字再写回去。现在你知道为什么要用铅笔了吧，少年！
下：给格子里的数字减去 1
始：开始重复「始……终」之间的指令，直到你读到「始」之前盯着的那个格子里的数字变成 0 为止。（什么？那个格子里已经是负数了？……不要这么没有下限好不好）
终：如果当前格子里的数字为 0，就跳过，否则回头到「始」那里
写：查当前格子里的数字在 ASCII 表上对应的字母，把它写下来（不，别写在格子里，就写在你买来一直立志想用但是没有用的日记本上吧）
读：随便想一个英文字母，查表找到它对应的数字，写到当前格子里
*/

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
