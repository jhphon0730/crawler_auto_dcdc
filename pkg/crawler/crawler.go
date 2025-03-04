package crawler

import (
	"bytes"
	"strings"
	"strconv"

	"github.com/PuerkitoBio/goquery"

	"github.com/jhphon0730/crawler_auto_dcdc/internal/model"
	"github.com/jhphon0730/crawler_auto_dcdc/pkg/network"
)

const (
	BASE_URL = "https://gall.dcinside.com/board/lists/?id=ohmygirl&page="

	POST_WRAPPER = "tr.us-post"
	POST_ATTR_DATA_TYPE = "data-type"
	POST_ATTR_NUMBER  = "data-no"
	POST_TITLE   = "td.gall_tit > a"
	POST_WRITER  = "td.gall_writer > span > em"
	POST_WRITE_DATE = "td.gall_date"
	POST_WRITE_DATE_TITLE = "title"
)

func GetPostBody(pageNumber string, postChan chan *model.Post, errChan chan error) {
	body, err := network.GetRequest(BASE_URL + pageNumber, nil)
	if err != nil {
		errChan <- err
		return 
	}

	parsePostBody(body, postChan, errChan)
	return
}

func parsePostBody(body []byte, postChan chan *model.Post, errChan chan error) {

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		errChan <- err
		return
	}

	doc.Find(POST_WRAPPER).Each(func(i int, s *goquery.Selection) {
		var post *model.Post = &model.Post{}

		// post number
		post_num_str, ok := s.Attr(POST_ATTR_NUMBER)
		if ok {
			post_num, err := strconv.Atoi(post_num_str)
			if err != nil {
				errChan <- err
				return
			}
			post.PostNumber = post_num
		}

		// post title
		post_title := s.Find(POST_TITLE).Text()
		post.Title = strings.TrimSpace(post_title)
		if strings.Contains(post.Title, "[") && strings.Contains(post.Title, "]") {
			post.Title = post.Title[:strings.Index(post.Title, "[")]
		}

		// post writer
		post_writer := s.Find(POST_WRITER).Text()
		post.Writer = post_writer

		// post write date
		post_write_date := s.Find(POST_WRITE_DATE).AttrOr(POST_WRITE_DATE_TITLE, "")
		post.WriteDate = post_write_date

		// post ( data-type )
		post_type := s.AttrOr(POST_ATTR_DATA_TYPE, "")
		post.DataType = post_type

		postChan <- post
	})

	return
}

