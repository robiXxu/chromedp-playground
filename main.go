package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func navAndShot(ctx context.Context, queryAction chromedp.QueryAction, url string, label string) {
	var buf []byte
	chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		queryAction,
		chromedp.FullScreenshot(&buf, 90),
	)

	err := ioutil.WriteFile(fmt.Sprintf("%s.png", label), buf, 0o644)
	if err != nil {
		log.Fatalf("Failed to take screenshot: %+v", err)
	}
}

func main() {
	fmt.Println("Just a chromedp playground")

	args := append(
		chromedp.DefaultExecAllocatorOptions[:],

		chromedp.Headless,
		chromedp.DisableGPU,

		chromedp.NoSandbox,
		//chromedp.Flag("font-render-hinting", "none"),
		//chromedp.Flag("ignore-gpu-blocklist", true),
		//chromedp.Flag("enable-accelerated-video-decode", true),
		//chromedp.Flag("enable-gpu-rasterization", true),

		chromedp.Flag("use-gl", "angle"),
		//chromedp.Flag("use-gl", "swiftshader-webgl"),
		//chromedp.Flag("use-angle", "swiftshader"),
		chromedp.Flag("use-angle", "swiftshader"),

		//chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("enable-threaded-compositing", true),

		chromedp.Flag("allow-insecure-localhost", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("allow-file-access-from-files", true),

		//chromedp.Flag("disable-gpu-compositing", true),
		chromedp.Flag("incognito", true),

		//chromedp.Flag("disable-gpu-watchdog", true),
		//chromedp.Flag("disable-hang-monitor", true),
		//chromedp.Flag("single-process", true),
	)

	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), args...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorCtx,
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	dirPath, _ := os.Getwd()
	log.Println(dirPath)

	webgl := fmt.Sprintf("file://%s/index.html", dirPath)

	//navAndShot(ctx, chromedp.WaitReady("body"), "chrome://gpu", "gpu")
	navAndShot(ctx, chromedp.Poll("window.ready===true", nil, chromedp.WithPollingTimeout(2*time.Second)), webgl, "webgl")

}
