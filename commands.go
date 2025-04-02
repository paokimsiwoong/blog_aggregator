package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/paokimsiwoong/blog_aggregator/internal/database"
	"github.com/paokimsiwoong/blog_aggregator/internal/rss"
)

// login command 입력 시 실행되는 함수
func handlerLogin(s *state, cmd command) error {
	// login 뒤의 추가 명령어들은 cmd.args에 저장되어 있음
	// args의 길이가 1이 아니면 login <userName> 형태가 아니므로 에러
	if len(cmd.args) != 1 {
		return errors.New("the login handler expects a single argument, the username")
	}

	// db에 없는 유저로 로그인하려하면
	if _, err := s.ptrDB.GetUser(context.Background(), cmd.args[0]); err != nil {
		return fmt.Errorf("error getting user : %w", err)
	}

	// config 파일의 current_user_name 필드를 cmd.args[0]로 수정
	if err := s.ptrCfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error setting user : %w", err)
	}

	fmt.Printf("the user %s has logged in.\n", cmd.args[0])
	return nil
}

/*
register command 입력 시 실행되는 함수 : db의 users 테이블에 신규 유저 insert
*/
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("the register handler expects a single argument, a new username")
	}

	// 이미 있는 유저를 또 등록하려 하는 경우
	if _, err := s.ptrDB.GetUser(context.Background(), cmd.args[0]); err == nil { // 기존에 있는 유저라 get하는데 문제없어서 err == nil 이면
		return errors.New("error creating user : can not register existing user")
	}
	// @@@ 해답에서는 이부분 딱히 없음 (∵ 기존 유저가 있을 경우 CreateUser의 err != nil ==> error creating user : pq: duplicate key value violates unique constraint "users_name_key")

	now := time.Now()

	user, err := s.ptrDB.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(), // func uuid.New() uuid.UUID : New creates a new random UUID or panics. New is equivalent to the expression
			CreatedAt: now,
			UpdatedAt: now,
			Name:      cmd.args[0],
		},
	)
	if err != nil {
		return fmt.Errorf("error creating user : %w", err)
	}

	// config 파일의 current_user_name 필드를 cmd.args[0]로 수정
	if err := s.ptrCfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error setting user : %w", err)
	}

	fmt.Printf("user %s was created : %+v\n", cmd.args[0], user)
	return nil
}

// users command 입력 시 실행되는 함수 : 모든 유저 리스트 출력
func handlerUsers(s *state, cmd command) error {
	curName := s.ptrCfg.CurrentUserName
	users, err := s.ptrDB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users table : %w", err)
	}

	fmt.Printf("total users : %d\n", len(users))

	for _, user := range users {
		if user.Name == curName {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}

	return nil
}

// reset command 입력 시 실행되는 함수 : users 테이블 리셋
func handlerReset(s *state, cmd command) error {
	err := s.ptrDB.ResetUsers(context.Background())

	if err != nil {
		return fmt.Errorf("error deleting users table : %w", err)
	}

	fmt.Println("users table and feeds table have been reset.")
	return nil
}

// agg command 입력시 실행되는 함수 : rss feed 수집
func handlerAgg(s *state, cmd command) error {
	// rss feed 가져올 url
	url := "https://www.wagslane.dev/index.xml"

	// _, err := rss.FetchFeed(context.Background(), url)
	rssFeed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Printf("Fetched feed: %+v\n", *rssFeed)
	// fmt.Printf("Fetched feed: %+v\n", rssFeed)
	// @@@ 해답처럼 포인터를 바로 넣어도 문제없이 출력(단지 출력 맨 앞에 &가 추가되는 차이)
	return nil
}

// addfeed command 입력시 실행되는 함수 : 주어진 이름과 url로 rss feed 수집 및 db에 저장
func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 { // addfeed 피드이름 피드url
		return errors.New("the addfeed handler expects two arguments, a new feed name and an url of the feed")
	}

	// feed를 추가하는 current user 정보를 users 테이블에서 불러오기
	user, err := s.ptrDB.GetUser(context.Background(), s.ptrCfg.CurrentUserName)
	if err != nil { // current_user_name 불러오는 데 실패
		return fmt.Errorf("error getting current user: %w", err)
	}

	// feeds 테이블에 새 feed 추가
	now := time.Now()
	feed, err := s.ptrDB.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Name:      cmd.args[0],
			Url:       cmd.args[1],
			UserID:    user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating feed : %w", err)
	}

	// fmt.Printf("Added feed: %+v\n", feed)
	// @@@ 해답의 깔끔한 출력 방식으로 대체
	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

// feed 출력 함수
func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

// feeds command 입력시 실행되는 함수 : feeds 테이블의 정보를 모두 불러오는 함수
func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.ptrDB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}

	fmt.Println("Stored feeds:")

	printFeeds(feeds)

	fmt.Println()
	fmt.Println("=====================================")
	return nil
}

// feeds 출력 함수
func printFeeds(feeds []database.GetFeedsRow) {
	for _, feed := range feeds {
		fmt.Printf("* Name:          %s\n", feed.Name)
		fmt.Printf("* URL:           %s\n", feed.Url)
		fmt.Printf("* Created:       %v\n", feed.CreatedAt)
		fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
		fmt.Printf("* UserName:      %s\n", feed.UserName)
	}
}
