package main

import (
	"encoding/json"
	"fmt"
	"io"

	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

// SofaScore unofficial API: FIFA World Cup 2026 = unique-tournament 16, season 58210.
// Cloudflare blocks plain net/http (TLS JA3 fingerprint) — we mimic the mobile app's
// OkHttp/Android TLS handshake against the api.sofascore.app host, which passes.
const (
	ssBase   = "https://api.sofascore.app/api/v1/unique-tournament/16/season/58210"
	ssUA     = "Sofascore/Android"
	ssMaxPag = 10 // safety cap on pagination pages per direction
)

// ssEvent is the subset of a SofaScore event we care about.
type ssEvent struct {
	ID        int    `json:"id"`
	HomeScore struct {
		Current *int `json:"current"`
	} `json:"homeScore"`
	AwayScore struct {
		Current *int `json:"current"`
	} `json:"awayScore"`
	Status struct {
		Type string `json:"type"` // notstarted | inprogress | finished
	} `json:"status"`
	StartTimestamp int64 `json:"startTimestamp"`
}

type ssResponse struct {
	Events []ssEvent `json:"events"`
}

// outEvent is the normalized shape returned to the frontend.
type outEvent struct {
	ID int    `json:"id"`
	HS *int   `json:"hs"`
	AS *int   `json:"as"`
	St string `json:"st"`
	Ts int64  `json:"ts"`
}

// FetchScores pulls all WC-2026 events from SofaScore and returns a JSON array of
// {id, hs, as, st, ts}. Returns an error on network/parse failure so the frontend
// can show an indicator without touching existing data.
func (a *App) FetchScores() (string, error) {
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(),
		tls_client.WithTimeoutSeconds(15),
		tls_client.WithClientProfile(profiles.Okhttp4Android13),
	)
	if err != nil {
		return "", err
	}

	seen := make(map[int]ssEvent)
	var lastErr error
	for _, dir := range []string{"last", "next"} {
		for n := 0; n < ssMaxPag; n++ {
			url := fmt.Sprintf("%s/events/%s/%d", ssBase, dir, n)
			req, err := fhttp.NewRequest("GET", url, nil)
			if err != nil {
				lastErr = err
				break
			}
			req.Header.Set("User-Agent", ssUA)
			req.Header.Set("Accept", "*/*")

			resp, err := client.Do(req)
			if err != nil {
				lastErr = err
				break
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode != 200 {
				lastErr = fmt.Errorf("sofascore %s/%d: status %d", dir, n, resp.StatusCode)
				break
			}
			var r ssResponse
			if err := json.Unmarshal(body, &r); err != nil || len(r.Events) == 0 {
				break // end of pagination for this direction
			}
			for _, e := range r.Events {
				seen[e.ID] = e
			}
		}
	}

	if len(seen) == 0 {
		if lastErr != nil {
			return "", lastErr
		}
		return "", fmt.Errorf("sofascore: no events returned")
	}

	out := make([]outEvent, 0, len(seen))
	for _, e := range seen {
		out = append(out, outEvent{
			ID: e.ID,
			HS: e.HomeScore.Current,
			AS: e.AwayScore.Current,
			St: e.Status.Type,
			Ts: e.StartTimestamp,
		})
	}
	b, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
