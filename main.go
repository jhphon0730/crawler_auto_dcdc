package main

import (
	"log"
	"strconv"
	"sync"

	"github.com/jhphon0730/crawler_auto_dcdc/pkg/crawler"
	"github.com/jhphon0730/crawler_auto_dcdc/internal/model"
	"github.com/jhphon0730/crawler_auto_dcdc/internal/server"
	"github.com/jhphon0730/crawler_auto_dcdc/internal/database"
)

var (
	// key is post number
	posts     map[int]*model.Post = make(map[int]*model.Post)
	isRunning bool                = false
	mu        sync.Mutex

	maxPage int = 10
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

	mu.Lock()
	if err := database.LoadPosts(posts); err != nil {
		return
	}
	mu.Unlock()
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

	if len(newPosts) > 0 {
		log.Println("New posts count:", len(newPosts))
	} else {
		log.Println("No new posts")
	}
	return
}

func main() {
	if err := database.InitDB("test.db"); err != nil {
		log.Fatalln("Failed to init DB:", err)
		return
	}
	defer database.CloseDB()

	// Initial & Run server
	server.InitialServer()
}
