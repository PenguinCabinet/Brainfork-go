package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	broadcastchannel "github.com/Maki-Daisuke/go-broadcast-channel"
)

var chars = map[rune]string{
	'<': "<",
	'>': ">",
	'+': "+",
	'-': "-",
	'.': ".",
	',': ",",
	'[': "[",
	']': "]",
	'{': "{",
	'}': "}",
	'(': "(",
	')': ")",
	'~': "~",
	'^': "^",
	'v': "v",
	';': ";",
	'|': "|",
}

func program_split(program []rune) []string {
	var A []string
	var stack_temp string
	for len(program) > 0 {
		check_match := false
		for _, e := range chars {
			if len(program) >= len([]rune(e)) {
				if string(program[len(program)-len([]rune(e)):]) == e {
					if stack_temp != "" {
						A = append(A, stack_temp)
					}
					A = append(A, string(program[len(program)-len([]rune(e)):]))
					program = program[:len(program)-len([]rune(e))]
					check_match = true
					stack_temp = ""
				}
			}
		}
		if !check_match {
			stack_temp += string([]rune{program[len(program)-1]})
			program = program[:len(program)-1]
		}
	}
	if stack_temp != "" {
		A = append(A, stack_temp)
	}
	for i := 0; i < len(A)/2; i++ {
		A[i], A[len(A)-i-1] = A[len(A)-i-1], A[i]
	}
	return A
}

func program_check(program []string) []string {
	var A []string
	A_index := 0
	nest1 := 0
	nest2 := 0
	nest3 := 0
	comment_flag := false
	for _, e := range program {
		if e[0] == '\n' || e[0] == '\r' {
			comment_flag = false
			continue
		}
		if comment_flag {
			continue
		}
		if e == chars[';'] {
			comment_flag = true
			continue
		}
		if e == chars['<'] ||
			e == chars['>'] ||
			e == chars['+'] ||
			e == chars['-'] ||
			e == chars['.'] ||
			e == chars[','] ||
			e == chars['['] ||
			e == chars[']'] ||
			e == chars['{'] ||
			e == chars['}'] ||
			e == chars['|'] ||
			e == chars['('] ||
			e == chars[')'] ||
			e == chars['^'] ||
			e == chars['v'] ||
			e == chars['~'] {
			if e == chars['['] {
				nest1++
			}
			if e == chars[']'] {
				nest1--
			}
			if e == chars['{'] {
				nest2++
			}
			if e == chars['}'] {
				nest2--
			}
			if e == chars['('] {
				nest3++
			}
			if e == chars[')'] {
				nest3--
			}

			A = append(A, e)
			A_index++
		}
		if nest1 < 0 || nest2 < 0 || nest3 < 0 {
			log.Fatalln("ERROR:Incorrect block nesting.")
			os.Exit(-1)
		}
	}

	if nest1 != 0 || nest2 != 0 || nest3 != 0 {
		log.Fatalln("ERROR:Incorrect block nesting.")
	}

	return A
}

func getchar() rune {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	return []rune(string([]byte(input)[0]))[0]
}

var VM_mem_size int = 1024

var broadcasts []*broadcastchannel.Broadcaster[bool]

func program_run(program []string, mem *[]byte, mem_Mutex *[]sync.Mutex, mem_count int) {

	for program_count := 0; program_count < len(program); program_count++ {
		switch program[program_count] {
		case chars['>']:
			mem_count++
		case chars['<']:
			mem_count--
		case chars['+']:
			(*mem)[mem_count]++
		case chars['-']:
			(*mem)[mem_count]--
		case chars['.']:
			fmt.Printf("%c", rune((*mem)[mem_count]))
		case chars[',']:
			temp := int(getchar())
			(*mem)[mem_count] = byte(temp)
		case chars['[']:
			if (*mem)[mem_count] == 0 {
				nest := 0
				for true {
					program_count++
					if program[program_count] == chars['['] {
						nest++
					}
					if program[program_count] == chars[']'] && nest == 0 {
						break
					}
					if program[program_count] == chars[']'] {
						nest--
					}
				}
			}
		case chars[']']:
			if (*mem)[mem_count] != 0 {
				nest := 0
				for true {
					program_count--
					if program[program_count] == chars[']'] {
						nest++
					}
					if program[program_count] == chars['['] && nest == 0 {
						break
					}
					if program[program_count] == chars['['] {
						nest--
					}
				}
			}
		case chars['{']:
			var thread_programs [][]string

			program_count_thread_start := program_count
			for true {
				program_count++
				if program[program_count] == "}" {
					thread_programs = append(thread_programs, program[(program_count_thread_start+1):program_count])
					break
				}
				if program[program_count] == "|" {
					thread_programs = append(thread_programs, program[(program_count_thread_start+1):program_count])
					program_count_thread_start = program_count
				}
			}

			var wg sync.WaitGroup
			for _, e := range thread_programs {

				e2 := e
				mem_count_temp := mem_count

				wg.Add(1)
				go (func() {
					program_run(e2, mem, mem_Mutex, mem_count_temp)
					wg.Done()
				})()
			}
			wg.Wait()

		case chars['~']:
			time.Sleep(100 * time.Millisecond)
		case chars['(']:
			program_count_start := program_count
			nest := 0
			for true {
				program_count++
				if program[program_count] == "(" {
					nest++
				} else if program[program_count] == ")" && nest == 0 {
					break
				} else if program[program_count] == ")" {
					nest--
				}
			}
			(*mem_Mutex)[mem_count].Lock()
			program_run(program[(program_count_start+1):program_count], mem, mem_Mutex, mem_count)
			(*mem_Mutex)[mem_count].Unlock()

		case chars['^']:
			ch := make(chan bool)
			broadcasts[mem_count].Subscribe(ch)
			<-ch
		case chars['v']:
			broadcasts[mem_count].Chan() <- true
		}
	}
}

func interpreter_main(args []string) {

	if len(args) < 2 {
		log.Fatalln("Please give the input Brainfork file.")
	}
	program := ""
	var splited_program []string
	program_bytes, err := ioutil.ReadFile(args[1])
	if err != nil {
		log.Fatalln("The input Brainfork file does not exist.")
	}
	program = string(program_bytes)

	splited_program = program_split([]rune(program))

	splited_program = program_check(splited_program)

	mem := make([]byte, VM_mem_size)
	mem_Mutex := make([]sync.Mutex, VM_mem_size)

	broadcasts = make([]*broadcastchannel.Broadcaster[bool], VM_mem_size)

	for i := 0; i < VM_mem_size; i++ {
		broadcasts[i] = broadcastchannel.New[bool](0).WithTimeout(100 * time.Microsecond)
		defer broadcasts[i].Close()
	}

	program_run(splited_program, &mem, &mem_Mutex, 0)
}

func main() {

	interpreter_main(os.Args)
}
