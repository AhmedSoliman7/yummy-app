package scrap

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

type(

	sdasa func () (map[string]string)
)

func crawl(url string, chUrl chan string, chTitle chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close()
	var waitingForTitle bool = false
	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token() 
			if t.Data == "a" {
				var isRestaurant bool = false
			for _, a := range t.Attr {
	  	 		if a.Key == "class" && a.Val == "vendor__inner"{
	        		isRestaurant = true
				}else if(isRestaurant && a.Key == "href"){
					chUrl <- a.Val
					waitingForTitle = true
			}
			}

		}
			
		if(waitingForTitle) {
			if t.Data == "div" {
				for _, a := range t.Attr {
					if a.Key == "class" && a.Val == "vendor__title" {
						z.Next()
						var title string = ""
						for {
							tt := z.Next()
							if(tt == html.TextToken){
								if(len(title) != 0){
									title = title + " "
								}
								title = title + z.Token().Data
							} else if(tt == html.EndTagToken) {
								break
							}
						}
						chTitle <- title
					}
				}
			}
		}

	}
}
}

func Restaurants()(map[string]string) {

	restsUrl := []string{}
	restsTitle := []string{}

	// Channels
	chUrls := make(chan string)
	chTitles := make(chan string)
	chFinished := make(chan bool) 
	go crawl("https://www.otlob.com/restaurants", chUrls, chTitles, chFinished)

	var done bool = false
	for !done {
		select {
		case url := <-chUrls:
			restsUrl = append(restsUrl, url)
		case title := <-chTitles:
			restsTitle = append(restsTitle, title)

		case <-chFinished:
			done = true
		}
	}

	foundRests := make(map[string]string)


	for c := 0; c < len(restsUrl); c++ {
		foundRests[restsTitle[c]] = restsUrl[c]
	}

	close(chUrls)
	close(chTitles)
	return foundRests
}

// func main(){
// 	fmt.Println(getRestaurants())
// }

