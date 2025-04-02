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

	// sql.Open의 첫번째 인자로 사용하는 sql 드라이버를 지정(_ "github.com/lib/pq"이 postgres)
	// 두번째 인자로는 cfg.DBURL에 저장된 connection string(postgres://username:password@localhost:5432/dbname?sslmode=disable 형태)로 database 연결
	db, errr := sql.Open("postgres", cfg.DBURL)
	// db는 *sql.DB 타입
	if errr != nil {
		log.Fatalf("error connecting to db : %v", err)
	}
	// @@@ 해답처럼 db.Close() defer 걸어두기
	defer db.Close()
	// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
	// http://golang.site/go/article/106-SQL-DB-%ED%99%9C%EC%9A%A9
	// http://golang.site/go/article/107-MySql-%EC%82%AC%EC%9A%A9---%EC%BF%BC%EB%A6%AC
	// sql.DB 타입에는 하나의 Row(record)만 반환할 때 쓰이는 QueryRow(), 복수 rows의 경우에는 Query() 사용
	// func (db *sql.DB) QueryRow(query string, args ...any) *sql.Row
	// func (db *sql.DB) Query(query string, args ...any) (*sql.Rows, error)
	//
	// ex: row := db.QueryRow("SELECT * FROM users WHERE id = 1;")는 id=1인 row를 sql.Row에 담아서 반환
	// sql.Row의 method인 func (r *sql.Row) Scan(dest ...any) error 를 이용해 sql.Row에 담긴 정보를 dest에 입력한 변수들에 전달 가능
	// var i User
	// err := row.Scan(
	// 	&i.ID,
	// 	&i.CreatedAt,
	// 	&i.UpdatedAt,
	// 	&i.Name,
	// )
	//
	// ex2: Query를 사용해 여러 row의 정보를 sql.Rows에 받은 경우 sql.Rows의 메소드 Next()를 for rows.Next() {]와 같이 써서 사용
	// for rows.Next() {
	// 	var i User
	// 	if err := rows.Scan(
	// 		&i.ID,
	// 		&i.CreatedAt,
	// 		&i.UpdatedAt,
	// 		&i.Name,
	// 	); err != nil {
	// 		return nil, err
	// 	}
	// 	items = append(items, i)
	// }
	//
	// ex3: QueryRow()로 데이터 읽기뿐 아니라 쓰기도 가능
	// row := db.QueryRow("INSERT INTO users (id, created_at, updated_at, name) Values ($1, $2, $3, $4) RETURNING id, created_at, updated_at, name", 변수1, 변수2, 변수3, 변수4)
	// @@@ 첫번째 인자 query string에는 formatted string 사용 가능 => 두번째 인자 args ...any 에 파라메터들 입력
	// @@@ 이 formatted string을 사용할 때 sql 종류에 따라 placeholder가 다르다. postgreSQL은 $1, $2, ... 이지만 MySQL은 ?, Oracle은 :val1, :val2 등을 사용한다
	// @@@@@@ insert 한 이후에 반환이 필요없을 경우 (*sql.DB).Exec() 사용
	// @@@@@@ (func (db *sql.DB) Exec(query string, args ...any) (sql.Result, error))
	//
	// http://golang.site/go/article/108-MySql-%EC%82%AC%EC%9A%A9---DML
	// ex4: Exec() : INSERT, UPDATE, DELETE등 리턴 데이터가 없는 쿼리 사용시 사용
	// ex5: Prepare() : Prepared Statement는 데이타베이스 서버에 Placeholder를 가진 SQL문을 미리 준비시키는 것으로, 차후 해당 Statement를 호출할 때 준비된 SQL문을 빠르게 실행하도록 하는 기법
	// // Prepared Statement 생성
	// stmt, err := db.Prepare("UPDATE test1 SET name=$1 WHERE id=$2")
	// stmt는 *sql.Stmt
	// checkError(err)
	// defer stmt.Close()
	//
	// // Prepared Statement 실행
	// _, err = stmt.Exec("Tom", 1) //Placeholder 파라미터 순서대로 전달
	// checkError(err)
	// _, err = stmt.Exec("Jack", 2)
	// checkError(err)
	// _, err = stmt.Exec("Shawn", 3)
	// checkError(err)
	// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
	// (*sql.DB).QueryRow(), (*sql.DB).Query(), (*sql.DB).Exec(), (*sql.DB).Prepare() 등은 모두 기본 context.Background를 사용한다
	// context를 자유롭게 활용하려면 각 함수들의 이름 끝에 Context를 붙이면 context.Context를 추가 인자로 입력가능
	// ===> (*sql.DB).QueryRowContext(), (*sql.DB).QueryContext(), (*sql.DB).ExecContext(), (*sql.DB).PrepareContext() 등
	// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
	// http://golang.site/go/article/109-MySql-%EC%82%AC%EC%9A%A9---%ED%8A%B8%EB%9E%9C%EC%9E%AD%EC%85%98
	// ex6: transaction : 트랜잭션은 복수 개의 SQL 문을 실행하다 중간에 어떤 한 SQL문에서라도 에러가 발생하면 전체 SQL문을 취소하게 되고 (이를 롤백이라 한다), 모두 성공적으로 실행되어야 전체를 커밋하게 된다.
	// 복수 개의 SQL 문을 하나의 트랜잭션으로 묶기 위하여 sql.DB의 Begin() 메서드를 사용한다.
	// Begin() 메서드는 sql.Tx 객체를 리턴하는데, 이 Tx 객체로부터 Tx.Exec() 등을 실행하여 트랜잭션을 수행한 후, 마지막에 최종 Commit을 위해 Tx.Commit() 메서드를 호출한다.
	// 트랜잭션을 취소하는 롤백을 위해서는 Tx.Rollback() 메서드를 호출한다.
	// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

	// sqlc가 생성한 database 패키지 사용
	dbQueries := database.New(db)
	// db(*sql.DB) 타입은 internal/database/db.go의 DBTX 인터페이스를 만족
	// 반환값으로는 *Queries 반환 type Queries struct { db DBTX}

	// Queries 타입은 GetUser, GetUsers, CreateUser, ResetUsers, WithTx등 다양한 method들 존재
	// @@@ WithTx는 transaction에 쓰임. 입력인자인 *sql.Tx는 (*sql.DB).Begin()로 생성 ((*sql.DB).Begin()은 복수의 sql쿼리를 하나의 transaction으로 묶어서 실행할 때 쓰이는 함수)

	// state, commands 구조체들 초기화
	stateInstance := state{
		ptrCfg: &cfg,
		ptrDB:  dbQueries,
	}
	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	// command 등록
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handlerUsers)
	cmds.register("reset", handlerReset)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)

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
