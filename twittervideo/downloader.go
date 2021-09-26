package twittervideo

import (
	"regexp"
	"strings"
	"github.com/gocolly/colly"
)

type TwitterVideoDownloader struct {
    m3u8_url	 string
	video_url    string
	video_url_id string
    bearer_token string
    xguest_token string
}

func NewTwitterVideoDownloader(url string) *TwitterVideoDownloader {
    self := new(TwitterVideoDownloader)
    self.video_url = url
    return self
}

func (self *TwitterVideoDownloader) GetBearerToken() string {
    c := colly.NewCollector()

    c.OnResponse(func(r *colly.Response) {
        pattern, _ := regexp.Compile(`"Bearer.*?"`)
        self.bearer_token = strings.Trim(pattern.FindString(string(r.Body)), `"`)
    })

    c.Visit("https://abs.twimg.com/web-video-player/TwitterVideoPlayerIframe.cefd459559024bfb.js")

    return self.bearer_token
}

func (self *TwitterVideoDownloader) GetXGuestToken() string {
    c := colly.NewCollector()

    c.OnRequest(func(r *colly.Request) {
        r.Headers.Set("Authorization", self.bearer_token)
    })

    c.OnResponse(func(r *colly.Response) {
        pattern, _ := regexp.Compile(`[0-9]+`)
        self.xguest_token = pattern.FindString(string(r.Body))
    })

    c.Post("https://api.twitter.com/1.1/guest/activate.json", nil)

    return self.xguest_token
}
