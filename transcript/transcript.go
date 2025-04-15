package transcript

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

var PLAYABILITY_STATUS []string = []string{
	"OK",
	"ERROR",
	"LOGIN_REQUIRED",
}

var PLAYABILITY_FAILED_REASON map[string]string = map[string]string{
	"BOT_DETECTED":      "Sign in to confirm you're not a bot",
	"AGE_RESTRICTED":    "Sign in to confirm your age",
	"VIDEO_UNAVAILABLE": "Video unavailable",
}
