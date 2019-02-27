package daemon

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/github"
	"github.com/lucassabreu/github-desktop-notifications/githubclient"
)

func lookForNotifications(token string, log *log.Logger, errorc chan error) {
	stdlog.Println("Looking for notifications")

	client := githubclient.NewGithub(token)
	c := context.Background()

	since := time.Now()

	for {
		req, err := client.NewRequest("GET", "notifications", nil)
		if err != nil {
			errorc <- err
			return
		}

		req.Header.Add("If-Modified-Since", since.Format(time.RFC1123))

		var ns []*github.Notification
		r, err := client.Do(c, req, &ns)

		if err != nil && r.StatusCode != 304 {
			errorc <- err
			return
		}

		log.Println(fmt.Sprintf("Ping..."))
		for _, n := range ns {

			req, err = client.NewRequest("GET", n.GetSubject().GetLatestCommentURL(), nil)
			if err != nil {
				errorc <- err
				return
			}

			var o map[string]interface{}
			_, err := client.Do(c, req, &o)
			if err != nil {
				errorc <- err
				return
			}
			log.Println(fmt.Sprintf("%s - %s", o["body"], o["html_url"]))
		}

		if r.StatusCode != 304 {
			since, err = time.Parse(time.RFC1123, r.Header.Get("Last-Modified"))
			if err != nil {
				errorc <- err
				return
			}
		}

		poolInterval := r.Header.Get("X-Poll-Interval") + "s"
		poolInterval = "5s"

		d, err := time.ParseDuration(poolInterval)
		if err != nil {
			errorc <- err
			return
		}
		time.Sleep(d)
	}
}
