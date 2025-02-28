package main

import (
	"log"
	"strconv"
	"sync"

	"github.com/jhphon0730/crawler_auto_dcdc/pkg/crawler"
	"github.com/jhphon0730/crawler_auto_dcdc/pkg/database"
	"github.com/jhphon0730/crawler_auto_dcdc/pkg/model"
)

var (
	// key is post number
	posts     map[int]*model.Post = make(map[int]*model.Post)
	isRunning bool                = false
	mu        sync.Mutex
)

// 해당 함수가 매일 새벽 6시에 실행 된다고 가정 (cron)
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
	newPosts := make(map[int]*model.Post)
	postChan := make(chan *model.Post)
	errChan := make(chan error)

	var wg sync.WaitGroup
	maxPage := 10

	wg.Add(1)
	if err := database.LoadPosts(posts); err != nil {
		wg.Done()
		return
	}
	wg.Done()
	wg.Wait()

	for i := 1; i <= maxPage; i++ {
		wg.Add(1)
		go func(page string) {
			crawler.GetPostBody(page, postChan, errChan, &wg)
		}(strconv.Itoa(i))
	}

	for post := range postChan {
		mu.Lock()
		if posts[post.PostNumber] == nil {
			posts[post.PostNumber] = post
			newPosts[post.PostNumber] = post
		}
		mu.Unlock()
	}

	wg.Wait()
	// 에러를 직접 처리
	for err := range errChan {
		log.Println("Error:", err)
	}
	close(postChan)
	close(errChan)

	log.Println("Crawler Done")
	if err := database.SavePosts(newPosts); err != nil {
		log.Println("Failed to save posts:", err)
		return
	}
	log.Println("Saved posts count:", len(newPosts))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := database.InitDB("db.db"); err != nil {
		log.Println("Failed to init DB:", err)
		return
	}
	defer database.CloseDB()

	ScheduleFunc()

	select {}
}
