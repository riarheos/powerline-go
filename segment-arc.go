package main

import (
	"encoding/json"
	"fmt"
	pwl "github.com/justjanne/powerline-go/powerline"
	"os/exec"
)

type counts struct {
	Untracked int
	Changed int
	Staged int
}

func runArcCommand(args ...string) ([]byte, error) {
	command := exec.Command("arc", args...)
	command.Env = gitProcessEnv
	out, err := command.Output()
	return out, err
}

func isArcTree(path string) (bool, error) {
	out, err := runArcCommand("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false, err
	}

	return string(out) == "true\n", nil
}

func arcBranch() (string, error) {
	out, err := runArcCommand("info", "--json")
	if err != nil {
		return "", err
	}

	var info struct {
		Branch string `json:"branch"`
	}

	err = json.Unmarshal(out, &info)
	if err != nil {
		return "", err
	}

	return info.Branch, nil
}

func arcStatus() (counts, error) {
	out, err := runArcCommand("st", "--json", "--no-sync-status")
	if err != nil {
		return counts{}, err
	}

	var info struct {
		Status struct {
			Untracked []struct {}
			Changed []struct {}
			Staged []struct {}
		}
	}
	err = json.Unmarshal(out, &info)
	if err != nil {
		return counts{}, err
	}

	c := counts{}

	c.Untracked = len(info.Status.Untracked)
	c.Changed = len(info.Status.Changed)
	c.Staged = len(info.Status.Staged)

	return c, nil
}

func segmentArc(p *powerline) []pwl.Segment {
	inTree, err := isArcTree(p.cwd)
	if err != nil || !inTree {
		return []pwl.Segment{}
	}


	branch, err := arcBranch()
	if err != nil {
		return []pwl.Segment{}
	}

	status, err := arcStatus()
	if err != nil {
		return []pwl.Segment{}
	}

	segments := []pwl.Segment{{
		Name: "arc-branch",
		Content: fmt.Sprintf("%s %s", p.symbols.RepoBranch, branch),
		Foreground: p.theme.RepoCleanFg,
		Background: p.theme.RepoCleanBg,
	}}

	if status.Untracked > 0 {
		segments = append(segments, pwl.Segment{
			Name: "arc-untracked",
			Content: fmt.Sprintf("%d%s", status.Untracked, p.symbols.RepoUntracked),
			Foreground: p.theme.GitUntrackedFg,
			Background: p.theme.GitUntrackedBg,
		})
		segments[0].Background = p.theme.RepoDirtyBg
	}

	if status.Changed > 0 {
		segments = append(segments, pwl.Segment{
			Name: "arc-changed",
			Content: fmt.Sprintf("%d%s", status.Changed, p.symbols.RepoNotStaged),
			Foreground: p.theme.GitNotStagedFg,
			Background: p.theme.GitNotStagedBg,
		})
		segments[0].Background = p.theme.RepoDirtyBg
	}

	if status.Staged > 0 {
		segments = append(segments, pwl.Segment{
			Name: "arc-staged",
			Content: fmt.Sprintf("%d%s", status.Staged, p.symbols.RepoStaged),
			Foreground: p.theme.GitStagedFg,
			Background: p.theme.GitStagedBg,
		})
		segments[0].Background = p.theme.RepoDirtyBg
	}

	return segments
}
