package twittervideo

import (
    // "fmt"
    "regexp"
    "strings"
    "github.com/gocolly/colly"
)

type TwitterVideoDownloader struct {
    m3u8_url	 string
	video_url    string
	video_url_id string
    headers      string
}

func NewTwitterVideoDownloader(url string) *TwitterVideoDownloader {
    self := new(TwitterVideoDownloader)
    self.video_url = url
    return self
}

func (self *TwitterVideoDownloader) GetBearerToken() string {
    var token string

    c := colly.NewCollector()

    c.OnResponse(func(r *colly.Response) {
        pattern, _ := regexp.Compile(`"Bearer.*?"`)
        token       = pattern.FindString(string(r.Body))
        token       = strings.Trim(token, `"`)
    })

    c.Visit("https://abs.twimg.com/web-video-player/TwitterVideoPlayerIframe.cefd459559024bfb.js")

    return token
}
