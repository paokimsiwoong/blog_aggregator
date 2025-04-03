package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

// 주어진 url에서 rss feed를 HTTP GET request로 가져오는 함수
// 추가로 가져온 텍스트 데이터에 escaped 된 문자들 원래 문자로 변환
func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// http.NewRequest 에 context 추가한 버전 함수
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error creating request: %w", err)
	}

	// setting User-Agent header : This is a common practice to identify your program to the server.
	req.Header.Set("User-Agent", "gator")

	client := http.Client{
		Timeout: 10 * time.Second, // @@@ 해답처럼 Timeout 필드 설정
	}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error getting response: %w", err)
	}
	defer res.Body.Close()

	// response에서 []byte 얻어내기
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error reading response body: %w", err)
	}

	// []byte unmarshal로 RSSFeed 구조체에 저장하기
	var rssFeed RSSFeed
	if err := xml.Unmarshal(data, &rssFeed); err != nil { // xml data에 RSSFeed 구조체에 없는 field가 있어도 unmarshal 과정에서 알아서 버린다
		return &RSSFeed{}, fmt.Errorf("error unmarshalling data: %w", err)
	}

	// html.UnescapeString 함수로 &ldquo; ==> ' 와 같이 변환하기
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
	// for _, item := range rssFeed.Channel.Item {
	// 	// fmt.Println(item.Title)
	// 	// item.Title = html.UnescapeString(item.Title)
	// 	item.Description = ""
	//  이 item은 원본이 아닌 복사 ==> 변경해도 원본 그대로
	// }
	// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}
	// @@@ 해답 비교
	// for i, item := range rssFeed.Channel.Item {
	// 	item.Title = html.UnescapeString(item.Title)
	// 	item.Description = html.UnescapeString(item.Description)
	// 	rssFeed.Channel.Item[i] = item
	// }

	return &rssFeed, nil
}
