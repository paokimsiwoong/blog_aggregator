package main

import (
	"log"
	"os"

	"github.com/paokimsiwoong/blog_aggregator/internal/config"
)

func main() {
	// json 파일 불러오기
	cfg, err := config.Read()
	if err != nil {
		// fmt.Println("error reading file : ", err)
		// @@@ 해답의 log.Fatalf 사용하기
		log.Fatalf("error reading file : %v", err)
	}

	// state, commands 구조체들 초기화
	stateInstance := state{&cfg}
	cmds := commands{make(map[string]func(*state, command) error)}

	// login command 등록
	cmds.register("login", handlerLogin)

	// 유저 명령어 입력 확인
	if len(os.Args) < 2 {
		// os.Args의 첫번째 arg는 무조건 프로그램 이름이므로 명령어가 포함되어 있으려면 길이가 2 이상이어야 한다
		log.Fatalf("error checking argments : not enough arguments were provided")
	}
	// 유저 명령어 command 구조체에 저장
	cmd := command{
		name: os.Args[1], // 0은 프로그램 이름, 1은 명령어 이름
		args: os.Args[2:],
	}

	// 명령어 실행
	if err := cmds.run(&stateInstance, cmd); err != nil {
		log.Fatalf("%v", err)
	}
}
