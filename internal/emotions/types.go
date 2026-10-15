package emotions

type EmotionSample struct {
	Emotion     string  `json:"emotion"`
	Trust       float64 `json:"trust"`
	TimestampMs int64   `json:"timestampMs"`
}

type WindowSummary struct {
	BeginMs        int64          `json:"beginMs"`
	EndMs          int64          `json:"endMs"`
	Samples        int            `json:"samples"`
	Distribution   map[string]int `json:"distribution"`
	Predominant    string         `json:"predominant"`
	Interpretation string         `json:"interpretation,omitempty"`
}

type InterpretationOut struct {
	Summary string `json:"summary" jsonschema_description:"One to two short paragraphs describing the pattern of expressions observed in the window, using a neutral and cautious tone."`
	Caveat  string `json:"caveat" jsonschema_description:"A brief reminder that facial expression does not equate to one’s actual emotional state."`
}
