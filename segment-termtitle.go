package main

// Port of set_term_title segment from powerine-shell:
// https://github.com/b-ryan/powerline-shell/blob/master/powerline_shell/segments/set_term_title.py

import (
	"fmt"
	"os"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentTermTitle(p *powerline) []pwl.Segment {
	var title string

	term := os.Getenv("TERM")
	if !(strings.Contains(term, "xterm") || strings.Contains(term, "rxvt")) {
		return []pwl.Segment{}
	}

	if p.cfg.Shell == "bash" {
		title = "\\[\\e]0;\\w\\a\\]"
	} else if p.cfg.Shell == "zsh" {
		title = "%{\033]0;%~\007%}"
	} else {
		cwd := p.cwd
		title = fmt.Sprintf("\033]0;%s@%s: %s\007", p.username, p.hostname, cwd)
	}

	return []pwl.Segment{{
		Name:           "termtitle",
		Content:        title,
		Priority:       MaxInteger, // do not truncate
		HideSeparators: true,       // do not draw separators
	}}
}
