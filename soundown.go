package main

import (
	// "fmt"
	"log"
	"flag"
	"net/http"
	"io/ioutil"
	"io"
	"os"
	"regexp"
	"encoding/json"
)

var soundcloud_url string

type Track struct {
	StreamUrl string `json:"streamUrl"`
}

func main() {

	flag.StringVar(&soundcloud_url, "url", "","url to soundcloud.com song")
	flag.Parse()

	if (soundcloud_url != ""){

		// Fetch songs page, so we couldn't scrape it for json string
		soundcloud_resp, err := http.Get(soundcloud_url)
		if (err != nil){ log.Fatal(err) }

		soundcloud_body, err := ioutil.ReadAll(soundcloud_resp.Body)
		defer soundcloud_resp.Body.Close()

		if (err != nil){ log.Fatal(err) }

		// Convert response body to string
		var soundcloud_html = string(soundcloud_body)

		// Scrape soundcloud song page for track json string
		buffer_tracks, err := regexp.Compile(`window.SC.bufferTracks.push\((.*?)\);\n`)
		if (err != nil){ log.Fatal(err) }

		// Store only first track on page
		var find_tracks = buffer_tracks.FindAllStringSubmatch(soundcloud_html, 1)

		// Do we have a match?
		if (find_tracks == nil){
			log.Fatal("couldn't find json string")
		}

		// Create JSON struc for track
		var track Track
		err = json.Unmarshal([]byte(find_tracks[0][1]), &track)
		if (err != nil){ log.Fatal(err) }

		// If everything OK, stream song to output
		stream_resp, err := http.Get(track.StreamUrl)
		if (err != nil){ log.Fatal(err) }

		io.Copy(os.Stdout, stream_resp.Body)
		defer stream_resp.Body.Close()

	}else{
		// Show arguments missing 
		flag.PrintDefaults()
	}

}