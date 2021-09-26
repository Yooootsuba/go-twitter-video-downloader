package twittervideo

import (
	"regexp"
	"strings"
	"github.com/gocolly/colly"
)

type TwitterVideoDownloader struct {
	video_url    string
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

func (self *TwitterVideoDownloader) GetM3U8ListUrl() string {
    var m3u8_list_url string

    c := colly.NewCollector()

    c.OnRequest(func(r *colly.Request) {
        r.Headers.Set("Authorization", self.bearer_token)
        r.Headers.Set("x-guest-token", self.xguest_token)
    })

    c.OnResponse(func(r *colly.Response) {
        pattern, _ := regexp.Compile(`https.*m3u8`)
        m3u8_list_url = pattern.FindString(string(r.Body))
    })

    url := "https://api.twitter.com/1.1/videos/tweet/config/" +
            strings.TrimPrefix(self.video_url, "https://twitter.com/i/status/") +
            ".json"

    c.Visit(url)

    return m3u8_list_url
}

func (self *TwitterVideoDownloader) GetM3U8Url(m3u8_list_url string) string {
    var m3u8_url string

    c := colly.NewCollector()

    c.OnResponse(func(r *colly.Response) {
        pattern, _ := regexp.Compile(`.*m3u8`)
        m3u8_urls  := pattern.FindAllString(string(r.Body), -1)
        m3u8_url    = m3u8_urls[len(m3u8_urls) - 1]
    })

    c.Visit(m3u8_list_url)

    return m3u8_url
}
