package main

import (
	"fmt"
	"log"

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
	// user 설정하고 json 파일로 디스크에 저장
	if err := cfg.SetUser("paokimsiwoong"); err != nil {
		// fmt.Println("error setting user : ", err)
		log.Fatalf("error setting user : %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		// fmt.Println("error reading file : ", err)
		log.Fatalf("error reading file : %v", err)
	}
	fmt.Printf("config : %+v\n", cfg) // v대신 +v쓰면 필드명까지 출력
}
