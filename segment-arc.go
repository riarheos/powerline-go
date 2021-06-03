package main

import (
	"encoding/json"
	"fmt"
	pwl "github.com/justjanne/powerline-go/powerline"
	"os/exec"
)

type arcStatusInfo struct {
	BranchInfo struct {
		Ahead  int
		Behind int
		Local  struct {
			Name   string
			Commit struct {
				Id string
			}
		}
		Detached bool
	} `json:"branch_info"`
	Status struct {
		Untracked []struct{}
		Changed   []struct{}
		Staged    []struct{}
	}
}

func runArcCommand(args ...string) ([]byte, error) {
	command := exec.Command("arc", args...)
	command.Env = gitProcessEnv
	out, err := command.Output()
	return out, err
}

func isArcTree() (bool, error) {
	out, err := runArcCommand("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false, err
	}

	return string(out) == "true\n", nil
}

func arcStatus() (arcStatusInfo, error) {
	out, err := runArcCommand("status", "--branch", "--json", "--no-sync-status")
	if err != nil {
		return arcStatusInfo{}, err
	}

	info := arcStatusInfo{}
	err = json.Unmarshal(out, &info)
	if err != nil {
		return arcStatusInfo{}, err
	}

	return info, nil
}

func makeSegment(p *powerline, segments *[]pwl.Segment, count int, symbol string, fg uint8, bg uint8) {
	if count > 0 {
		*segments = append(*segments, pwl.Segment{
			Name:       "arc-status",
			Content:    fmt.Sprintf("%d%s", count, symbol),
			Foreground: fg,
			Background: bg,
		})
		(*segments)[0].Background = p.theme.RepoDirtyBg
	}
}

func segmentArc(p *powerline) []pwl.Segment {
	inTree, err := isArcTree()
	if err != nil || !inTree {
		return []pwl.Segment{}
	}

	status, err := arcStatus()
	if err != nil {
		return []pwl.Segment{}
	}

	var branchName string
	if status.BranchInfo.Detached {
		branchName = fmt.Sprintf("%s %s", p.symbols.RepoDetached, status.BranchInfo.Local.Commit.Id[:10])
	} else {
		branchName = fmt.Sprintf("%s %s", p.symbols.RepoBranch, status.BranchInfo.Local.Name)
	}

	segments := []pwl.Segment{{
		Name:       "arc-branch",
		Content:    branchName,
		Foreground: p.theme.RepoCleanFg,
		Background: p.theme.RepoCleanBg,
	}}

	makeSegment(p, &segments,
		status.BranchInfo.Ahead, p.symbols.RepoAhead, p.theme.GitAheadFg, p.theme.GitAheadBg)

	makeSegment(p, &segments,
		status.BranchInfo.Behind, p.symbols.RepoBehind, p.theme.GitBehindFg, p.theme.GitBehindBg)

	makeSegment(p, &segments,
		len(status.Status.Untracked), p.symbols.RepoUntracked, p.theme.GitUntrackedFg, p.theme.GitUntrackedBg)

	makeSegment(p, &segments,
		len(status.Status.Changed), p.symbols.RepoNotStaged, p.theme.GitNotStagedFg, p.theme.GitNotStagedBg)

	makeSegment(p, &segments,
		len(status.Status.Staged), p.symbols.RepoStaged, p.theme.GitStagedFg, p.theme.GitStagedBg)

	return segments
}
