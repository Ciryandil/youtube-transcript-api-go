package transcript

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/Ciryandil/youtube-transcripts-api-go/api/constants"
	"github.com/Ciryandil/youtube-transcripts-api-go/api/utils"
	"github.com/Ciryandil/youtube-transcripts-api-go/proxy"
)

type TranscriptSnippet struct {
	Text     string
	Start    float64
	Duration float64
}

type Transcript struct {
	Snippets     []TranscriptSnippet
	VideoId      string
	Language     string
	LanguageCode string // TODO: can there be a definitive map across languages and language codes -> i18n languages api route youtube data
	//bool isGenerated not included -> make separate arrays/maps instead
}

type TranslationLanguage struct {
	Language     string
	LanguageCode string
}

type TranscriptDownloader struct {
	client      *http.Client
	proxyConfig *proxy.ProxyConfig
}

var PLAYABILITY_FAILED_REASON map[string]string = map[string]string{
	"BOT_DETECTED":      "Sign in to confirm you're not a bot",
	"AGE_RESTRICTED":    "Sign in to confirm your age",
	"VIDEO_UNAVAILABLE": "Video unavailable",
}

func constructVideoUnavailabilitySubreasons(playabilityStatusData map[string]interface{}) []string {
	if errorScreen, ok := playabilityStatusData["errorScreen"].(map[string]interface{}); ok {
		if renderer, ok := errorScreen["playerErrorMessageRenderer"].(map[string]interface{}); ok {
			if subreason, ok := renderer["subreason"].(map[string]interface{}); ok {
				if runs, ok := subreason["runs"].([]map[string]interface{}); ok {
					textList := make([]string, 0)
					for _, run := range runs {
						runText, ok := run["text"].(string)
						if !ok {
							continue
						}
						textList = append(textList, runText)
					}
					return textList
				}
			}
		}
	}
	return nil
}

func assertPlayability(playabilityStatusData map[string]interface{}) error {
	playabilityStatus, ok := playabilityStatusData["status"].(string)
	if !ok {
		return fmt.Errorf("playability status not found")
	}
	if playabilityStatus != "OK" { // TODO: Maybe find a better way to represent this than raw strings
		reason, ok := playabilityStatusData["reason"].(string)
		if !ok {
			return fmt.Errorf("reason not found")
		}
		if playabilityStatus == "LOGIN_REQUIRED" {
			if reason == PLAYABILITY_FAILED_REASON["BOT_DETECTED"] {
				return fmt.Errorf("Request blocked")
			} else if reason == PLAYABILITY_FAILED_REASON["AGE_RESTRICTED"] {
				return fmt.Errorf("Video is age restricted")
			}
		} else if playabilityStatus == "ERROR" && reason == PLAYABILITY_FAILED_REASON["VIDEO_UNAVAILABLE"] {
			subReasons := constructVideoUnavailabilitySubreasons(playabilityStatusData)
			return fmt.Errorf("video unplayable: reason : %s : subreasons: %v", reason, subReasons)
		}
	}
	return nil
}

func getCaptionsJsonFromVideoData(videoData map[string]interface{}) (map[string]interface{}, error) {
	if captions, ok := videoData["captions"].(map[string]interface{}); ok {
		if renderer, ok := captions["playerCaptionsTracklistRenderer"].(map[string]interface{}); ok {
			if captionTracks, ok := renderer["captionTracks"].(map[string]interface{}); ok {
				return captionTracks, nil
			}
		}
	}
	return nil, fmt.Errorf("transcript disabled")
}

func extractCaptionsJson(html string) (map[string]interface{}, error) {
	videoData, err := parseJSVars("ytInitialPlayerResponse", html)
	if err != nil {
		if strings.Contains(html, "class=\"g-recaptcha\"") {
			return nil, fmt.Errorf("your IP has been blocked")
		}
		return nil, fmt.Errorf("failed to extract captions: %v", err)
	}
	playabilityStatusData, ok := videoData["playabilityStatus"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("playability status data not found")
	}
	err = assertPlayability(playabilityStatusData)
	if err != nil {
		return nil, err
	}
	return getCaptionsJsonFromVideoData(videoData)
}

// TODO
func fetchVideoHtml(videoId string, downloader TranscriptDownloader) (string, error) {
	return "", nil
}

func fetchHtml(videoId string, downloader TranscriptDownloader) (string, error) {
	endpoint := fmt.Sprintf("%s%s", constants.WATCH_URL, videoId)
	resp, statusCode, _, err := utils.HTTPRequest(downloader.client, "GET", endpoint, nil, nil, nil)
	if err != nil {
		return "", err
	}
	respString := string(resp)
	if statusCode >= 400 {
		return "", fmt.Errorf("error fetching video: Code: %d, Response: %s", statusCode, respString)
	}
	return html.UnescapeString(respString), nil
}
