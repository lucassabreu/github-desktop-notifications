package daemon

import (
	"context"
	"log"
	"time"

	"github.com/lucassabreu/github-desktop-notifications/notify"

	"github.com/google/go-github/github"
	"github.com/lucassabreu/github-desktop-notifications/githubclient"
)

type comment struct {
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
}

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

		log.Println("Ping...")
		for _, n := range ns {
			go processNotification(n, client, errorc)
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

func processNotification(n *github.Notification, client *github.Client, errorc chan error) {
	req, err := client.NewRequest("GET", n.GetSubject().GetLatestCommentURL(), nil)
	if err != nil {
		errorc <- err
		return
	}

	var o comment

	if _, err = client.Do(context.Background(), req, &o); err != nil {
		errorc <- err
		return
	}

	notify.Notify(
		n.GetReason(),
		o.Body,
		o.HTMLURL,
	)
}
