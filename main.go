package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // _ "github.com/lib/pq" 는 postgres driver를 사용한다고 알리는 것. main.go 내부에서 직접 코드 작성할 때 쓰이지는 않음
	"github.com/paokimsiwoong/blog_aggregator/internal/config"
	"github.com/paokimsiwoong/blog_aggregator/internal/database"
)

func main() {
	// json 파일 불러오기
	cfg, err := config.Read()
	if err != nil {
		// fmt.Println("error reading file : ", err)
		// @@@ 해답의 log.Fatalf 사용하기
		log.Fatalf("error reading file : %v", err)
	}

	// cfg.DBURL에 저장된 connection string(postgres://username:password@localhost:5432/dbname?sslmode=disable 형태)로 database 연결
	db, errr := sql.Open("postgres", cfg.DBURL)
	// db는 *sql.DB 타입
	if errr != nil {
		log.Fatalf("error connecting to db : %v", err)
	}
	// @@@ 해답처럼 db.Close() defer 걸어두기
	defer db.Close()

	// sqlc가 생성한 database 패키지 사용
	dbQueries := database.New(db)

	// state, commands 구조체들 초기화
	stateInstance := state{
		ptrCfg: &cfg,
		ptrDB:  dbQueries,
	}
	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	// login, register command 등록
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)

	// 유저 명령어 입력 확인
	if len(os.Args) < 2 {
		// os.Args의 첫번째 arg는 무조건 프로그램 이름이므로 명령어가 포함되어 있으려면 길이가 2 이상이어야 한다
		log.Fatalf("error checking arguments : not enough arguments were provided")
	}
	// 유저 명령어 command 구조체에 저장
	cmd := command{
		name: os.Args[1], // 0은 프로그램 이름, 1은 명령어 이름
		args: os.Args[2:],
	}

	// 명령어 실행
	if err := cmds.run(&stateInstance, cmd); err != nil {
		// log.Fatalf("%v", err)
		// @@@ 해답 log.Fatal 함수 사용 반영
		log.Fatal(err)
	}
}
