package main

import(
    "fmt"
    "log"
    "os/exec"
    "gotwittervideo/twittervideo"
)

func main() {
    // Check FFmpeg is in PATH
    _, err := exec.LookPath("ffmpeg")
    if err != nil {
        log.Fatal("Error: FFmpeg is not found !")
    }

	// Input
    var url string
    fmt.Printf("Enter a Twitter video url : ")
    fmt.Scanln(&url)

    // Output
    fmt.Println("\nDownloading ...")
    downloader := twittervideo.NewTwitterVideoDownloader(url)
    downloader.Download()
    fmt.Println("\nFinished !")
}
