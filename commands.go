package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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

	fmt.Println("All tables have been reset.")
	return nil
}

// addfeed command 입력시 실행되는 함수 : 주어진 이름과 url로 rss feed 수집 및 db에 저장
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 { // addfeed 피드이름 피드url
		return errors.New("the addfeed handler expects two arguments, a new feed name and an url of the feed")
	}

	// // feed를 추가하는 current user 정보를 users 테이블에서 불러오기
	// user, err := s.ptrDB.GetUser(context.Background(), s.ptrCfg.CurrentUserName)
	// if err != nil { // current_user_name 불러오기 실패
	// 	return fmt.Errorf("error getting current user: %w", err)
	// } // @@@@@ 함수 시그니쳐를 변경(User를 입력받도록)해서 이부분 주석처리

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

	// 현재 유저와 추가된 feed pair를 feed_follows 테이블에 추가
	_, err = s.ptrDB.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil { // 테이블 추가 실패
		return fmt.Errorf("error creating feed_follow record: %w", err)
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
		fmt.Printf("* ID:            %s\n", feed.ID)
		fmt.Printf("* Name:          %s\n", feed.Name)
		fmt.Printf("* URL:           %s\n", feed.Url)
		fmt.Printf("* Created:       %v\n", feed.CreatedAt)
		fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
		fmt.Printf("* UserName:      %s\n", feed.UserName)
	}
}

// follow command 입력시 실행되는 함수 : url을 받아서 현재 유저와 url의 피드 pair를 feed_follows 테이블에 저장
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 { // follow url
		return errors.New("the follow handler expects one arguments, a feed url")
	}

	// // feed_follow를 추가하는 current user 정보를 users 테이블에서 불러오기
	// user, err := s.ptrDB.GetUser(context.Background(), s.ptrCfg.CurrentUserName)
	// if err != nil { // current_user_name 불러오기 실패
	// 	return fmt.Errorf("error getting current user: %w", err)
	// } // @@@@@ 함수 시그니쳐를 변경(User를 입력받도록)해서 이부분 주석처리

	feed, err := s.ptrDB.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil { // feed 불러오기 실패
		return fmt.Errorf("error getting stored feed from url: %w", err)
	}

	// feed_follows 테이블에 추가
	now := time.Now()
	feed_follow, err := s.ptrDB.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil { // 테이블 추가 실패
		return fmt.Errorf("error creating feed_follow record: %w", err)
	}

	fmt.Println("=====================================")
	fmt.Printf("Current user %s follows %s", feed_follow.UserName, feed_follow.FeedName)

	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

// unfollow command 입력시 실행되는 함수 : 입력된 url 피드를 unfollow
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 { // unfollow url
		return errors.New("the unfollow handler expects one arguments, a feed url")
	}

	err := s.ptrDB.DeleteFeedFollow(
		context.Background(),
		database.DeleteFeedFollowParams{
			Name: user.Name,
			Url:  cmd.args[0],
		},
	)
	if err != nil {
		return fmt.Errorf("error deleting feed_follow record: %w", err)
	}

	fmt.Println("=====================================")
	fmt.Printf("Current user %s unfollows %s", user.Name, cmd.args[0])

	fmt.Println()
	fmt.Println("=====================================")

	return nil

	// @@@ 해답은 sql 쿼리는 간단하게 user_id, feed_id를 입력받는 걸로 하고 이 handlerUnfollow에서 GetFeedByURL 함수에 url 입력해서 받은 feed의 id를 입력하고 있음
}

// following command 입력시 실행되는 함수 : 현재 유저가 follow중인 feed 리스트 출력
func handlerFollowing(s *state, cmd command, user database.User) error {
	// user, err := s.ptrDB.GetUser(context.Background(), s.ptrCfg.CurrentUserName)
	// if err != nil { // current_user_name 불러오기 실패
	// 	return fmt.Errorf("error getting current user: %w", err)
	// } // @@@@@ 함수 시그니쳐를 변경(User를 입력받도록)해서 이부분 주석처리

	followlist, err := s.ptrDB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil { // current_user_name 불러오기 실패
		return fmt.Errorf("error getting follows for user: %w", err)
	}

	fmt.Printf("Current user %s is following:\n", user.Name)

	for i, follow := range followlist {
		fmt.Printf("%d\n", i)
		fmt.Printf("* FeedName:      %s\n", follow.FeedName)
		fmt.Printf("* URL:           %s\n", follow.Url)
	}

	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

// Middleware - Function transformation : 반복 사용되던 GetUser 코드 부분을 이 함수 한군데로 통일
// Middleware is a way to wrap a function with additional functionality. It is a common pattern that allows us to write DRY code.
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		//@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		user, err := s.ptrDB.GetUser(context.Background(), s.ptrCfg.CurrentUserName)
		if err != nil { // current_user_name 불러오기 실패
			return fmt.Errorf("error getting current user: %w", err)
		}
		// 이 반복 사용되는 코드부분을 handler에서 빼와서 이곳으로 이동
		//@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

		return handler(s, cmd, user)
	}
}

// agg command 입력시 실행되는 함수 : rss feed 수집
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 { // agg time_between_requests
		return errors.New("the agg handler expects one arguments, a time between requests")
	}

	fmt.Printf("Collecting feeds every %s\n", cmd.args[0])

	// 1s, 1m10s 같은 string 시간간격표현을 time.Duration으로 변환
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error parsing duration: %w", err)
	}

	// ticker.C는 지정된 간격마다 그 시점 시간을 받는 채널
	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return fmt.Errorf("error scraping feed: %w", err)
		}
	}

	// return nil
}

// DB에서 가장 갱신 시점이 오래된 feed를 찾아 갱신하고 내용을 출력하는 함수
func scrapeFeeds(s *state) error {
	// 갱신 안된지 가장 오래된 feed를 받아오기
	feedToFetch, err := s.ptrDB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting oldest feed from db: %w", err)
	}

	// 찾아온 feed는 이제 갱신되므로 last_fetched_at, updated_at 갱신
	now := time.Now()
	if err := s.ptrDB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{UpdatedAt: now, ID: feedToFetch.ID}); err != nil {
		return fmt.Errorf("error marking oldest feed: %w", err)
	}
	// @@@ 해답은 query에 RETURNING *; 을 붙여서 갱신 후 해당 feed를 다시 반환 받는다

	rssFeed, err := rss.FetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Println("=====================================")
	fmt.Printf("Newly fetched feed at %v:\n", now)
	fmt.Printf("Feed title: %s\n", rssFeed.Channel.Title)
	fmt.Printf("Feed url: %s\n", feedToFetch.Url)
	fmt.Printf("Feed description: %s\n", rssFeed.Channel.Description)

	fmt.Println("=====================================")

	for i, item := range rssFeed.Channel.Item {
		fmt.Printf("* item %d title: %s\n", (i + 1), item.Title)
		fmt.Printf("* description: %s\n", item.Description)
		fmt.Printf("* url: %s\n", item.Link)
		fmt.Printf("* PubDate: %s\n", item.PubDate)
		fmt.Println("----------------------------------------")
		fmt.Println()

		now := time.Now()
		if _, err := s.ptrDB.CreatePost(
			context.Background(), database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   now,
				UpdatedAt:   now,
				Title:       item.Title,
				Url:         item.Link,
				Description: convertToNullString(item.Description),
				PublishedAt: convertToNullTime(item.PubDate),
				FeedID:      feedToFetch.ID,
			},
		); err != nil {
			errText := err.Error()
			fmt.Printf("error printing: %s\n", errText)
			if errText == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
				fmt.Println("ignoring duplicate URL error")
			} else {
				return fmt.Errorf("error creating post: %w", err)
			}
		}
	}

	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

// string을 sql.NullString으로 변환하는 함수
func convertToNullString(text string) sql.NullString {
	if len(text) == 0 {
		return sql.NullString{}
		// Valid 필드는 bool이므로 default가 false
	}

	return sql.NullString{
		String: text,
		Valid:  true,
	}
}

// RSSItem 구조체의 PubDate 필드(string)을 받아 sql.NullTime으로 변환하는 함수
func convertToNullTime(pubdata string) sql.NullTime {
	layout := "Mon, 02 Jan 2006 15:04:05 -0700" // 마지막 자리는 시간대 표시 ex: MST는 -0700, UTC는 +0000
	// layout에는 반드시 reference time(2006년 1월 2일 월요일 15:04:05 MST)를 사용해야 한다
	// 이 reference time을 원하는 형식으로 변환해서 사용
	parsed, err := time.Parse(layout, pubdata)
	// @@@ 해답처럼 time 패키지의 time.RFC1123Z const로 대체 가능 RFC1123Z == "Mon, 02 Jan 2006 15:04:05 -0700"
	// parsed, err := time.Parse(time.RFC1123Z, pubdata)
	if err != nil {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  parsed,
		Valid: true,
	}
}

// browse command 입력시 실행되는 함수 : 현재 유저가 follow중인 feed들의 post들을 출력
func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	if len(cmd.args) > 1 { // browse limit
		return errors.New("the browse handler expects at most one arguments, a number of posts to show (if no argument given, number set to default: 2)")
	} else if len(cmd.args) == 0 {
		limit = 2
	} else {
		var err error
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("error converting a to i: %w", err)
		}
	}
	// @@@ 해답 cmd.args 예외처리
	// limit := 2
	// if len(cmd.Args) == 1 {
	// 	if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
	// 		limit = specifiedLimit
	// 	} else {
	// 		return fmt.Errorf("invalid limit: %w", err)
	// 	}
	// }

	// @@@ posts에는 posts.* 와 feeds.name 정보가 들어있음
	posts, err := s.ptrDB.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  int32(limit),
		},
	)
	if err != nil {
		return fmt.Errorf("error getting posts for user: %w", err)
	}

	fmt.Println("=====================================")

	for i, post := range posts {
		fmt.Printf("* post %d title: %s\n--- From feed: %s ---\n", (i + 1), post.Title, post.FeedName)
		fmt.Printf("* description: %s\n", post.Description.String)
		fmt.Printf("* url: %s\n", post.Url)
		fmt.Printf("* PubDate: %v\n", post.PublishedAt.Time)
		fmt.Println("----------------------------------------")
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("=====================================")

	return nil
}
