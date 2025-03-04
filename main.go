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
	log.Println("Start ScheduleFunc")

	mu.Lock()
	if isRunning {
		log.Println("Already running")
		mu.Unlock()
		return
	}
	isRunning = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		isRunning = false
		mu.Unlock()
	}()

	posts = make(map[int]*model.Post)
	newPosts := make(map[int]*model.Post)
	postChan := make(chan *model.Post)
	errChan := make(chan error)

	var wg sync.WaitGroup
	maxPage := 2

	if err := database.LoadPosts(posts); err != nil {
		return
	}
	log.Println("Loaded posts count:", len(posts))

	for i := 1; i <= maxPage; i++ {
		wg.Add(1)
		page := strconv.Itoa(i)
		go func(page string) {
			defer wg.Done()
			crawler.GetPostBody(page, postChan, errChan)
		}(page)
	}

	go func() {
		wg.Wait()
		close(postChan)
		close(errChan)
	}()

	for post := range postChan {
		mu.Lock()
		if posts[post.PostNumber] == nil {
			posts[post.PostNumber] = post
			newPosts[post.PostNumber] = post
			log.Println("New post:", post.PostNumber)
		}
		mu.Unlock()
	}

	for err := range errChan {
		log.Println("Error:", err)
	}

	log.Println("Crawler Done")
	if err := database.SavePosts(newPosts); err != nil {
		log.Println("Failed to save posts:", err)
		return
	}
}

func main() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := database.InitDB("test.db"); err != nil {
		log.Println("Failed to init DB:", err)
		return
	}
	defer database.CloseDB()

	ScheduleFunc()
}
