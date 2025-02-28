package main

import (
	"log"
	"sync"
	"strconv"

	"github.com/jhphon0730/crawler_auto_dcdc/pkg/model"
	"github.com/jhphon0730/crawler_auto_dcdc/pkg/crawler"
)

var (
	// key is post number
	posts map[int]*model.Post = make(map[int]*model.Post)

	isRunning bool = false

	mu sync.Mutex
)

// 해당 함수가 매일 새벽 6시에 실행 된다고 가정 ( cron )
func ScheduleFunc() {
	if isRunning {
		log.Println("Already running")
		return
	}
	isRunning = true

	defer func() {
		isRunning = false
	}()

	// posts 초기화 및 post, err 채널 생성
	posts = make(map[int]*model.Post)
	postChan := make(chan *model.Post)
	errChan := make(chan error)

	var wg sync.WaitGroup
	maxPage := 10

	for i := 1; i <= maxPage; i++ {
		wg.Add(1)
		go func(page string) {
			crawler.GetPostBody(page, postChan, errChan, &wg)
		}(strconv.Itoa(i))
	}

	go func() {
		wg.Wait()
		close(postChan)
	}()

	for post := range postChan {
		mu.Lock()
		if posts[post.PostNumber] == nil {
			posts[post.PostNumber] = post
			log.Println("Insert post number:", post.PostNumber)
		} else {
			log.Println("Already exist post number:", post.PostNumber)
		}
		mu.Unlock()
	}

	go func() {
		for err := range errChan {
			log.Println("Error:", err)
		}
	}()

	log.Println("Finish")
	log.Println("Posts Count:", len(posts))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ScheduleFunc()

	select {}
}
