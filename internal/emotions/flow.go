package emotions

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
)

type interpretationFlow = core.Flow[WindowSummary, InterpretationOut, struct{}]

func defineInterpretationFlow(g *genkit.Genkit) *interpretationFlow {
	return genkit.DefineFlow(g, "interpretWindow",
		func(ctx context.Context, window WindowSummary) (InterpretationOut, error) {
			prompt := fmt.Sprintf(
				"You are analyzing facial expressions automatically detected via webcam "+
					"from an individual, over a %d-second window, with %d samples.\n"+
					"Distribution of observed expressions (count by label): %v\n"+
					"Predominant expression: %s.\n\n"+
					"Write a brief, neutral, and cautious summary of the observed pattern. "+
					"Do not state with certainty what the person is feeling—always frame your observations "+
					"in terms of 'apparent facial expression.' Do not provide medical or psychological advice.",
				(window.EndMs-window.BeginMs)/1000, window.Samples, window.Distribution, window.Predominant,
			)

			out, _, err := genkit.GenerateData[InterpretationOut](ctx, g,
				ai.WithPrompt(prompt),
			)
			if err != nil {
				return InterpretationOut{}, err
			}
			return *out, nil
		},
	)
}
