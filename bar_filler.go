package mpb

import (
	"io"
	"strings"

	"github.com/vbauerster/mpb/v4/decor"
	"github.com/vbauerster/mpb/v4/internal"
)

const (
	rLeft = iota
	rFill
	rTip
	rEmpty
	rRight
	rRefill
)

var defaultBarStyle = []rune("[=>-]+")

type barFiller struct {
	format []rune
	rup    int
}

func (s *barFiller) Fill(w io.Writer, width int, stat *decor.Statistics) {

	str := string(s.format[rLeft])

	// don't count rLeft and rRight [brackets]
	width -= 2

	if width <= 2 {
		io.WriteString(w, str+string(s.format[rRight]))
		return
	}

	progressWidth := internal.Percentage(stat.Total, stat.Current, int64(width))
	needTip := progressWidth < int64(width) && progressWidth > 0

	if needTip {
		progressWidth--
	}

	if s.rup > 0 {
		refillCount := internal.Percentage(stat.Total, int64(s.rup), int64(width))
		rest := progressWidth - refillCount
		str += runeRepeat(s.format[rRefill], int(refillCount)) + runeRepeat(s.format[rFill], int(rest))
	} else {
		str += runeRepeat(s.format[rFill], int(progressWidth))
	}

	if needTip {
		str += string(s.format[rTip])
		progressWidth++
	}

	rest := int64(width) - progressWidth
	str += runeRepeat(s.format[rEmpty], int(rest)) + string(s.format[rRight])
	io.WriteString(w, str)
}

func (s *barFiller) SetRefill(upto int) {
	s.rup = upto
}

func runeRepeat(r rune, count int) string {
	return strings.Repeat(string(r), count)
}
