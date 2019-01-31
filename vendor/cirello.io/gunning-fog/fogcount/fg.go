package fogcount // import "cirello.io/gunning-fog/fogcount"

import (
	"bufio"
	"io"
)

// Analyze processes an English text and produces its Gunning Fox index score.
// Refer to its logic in https://en.wikipedia.org/wiki/Gunning_fog_index - it
// does not analyse word endings (-es, -ed, or -ing), or discriminate proper
// nouns, familiar jargon or compound words.
func Analyze(rdr io.Reader) float64 {
	var (
		phraseSize          int
		phraseCount         int
		hardWords           int
		words               int
		totalSentenceLength int
	)

	scanner := bufio.NewScanner(rdr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		words++
		phraseSize++
		w := scanner.Text()

		lastRune := w[len(w)-1]
		switch lastRune {
		case '.', ',', ';':
			w = w[:len(w)-1]
			totalSentenceLength += phraseSize
			phraseCount++
			phraseSize = 0
		}

		if len(w) > 6 {
			hardWords++
		}
	}

	return 0.4 * (float64(totalSentenceLength/phraseCount + 100*hardWords/words))
}
