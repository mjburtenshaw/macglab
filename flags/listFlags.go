package flags

import (
	"strings"

	"github.com/mjburtenshaw/macglab/config"
	"github.com/spf13/cobra"
)

type BooleanFlags struct {
	Approved bool
	Browser  bool
	Draft    bool
	Group    bool
	Projects bool
	Ready    bool
}

type ResolvedFlags struct {
	AccessToken string
	GroupId     string
	Me          int
	Usernames   []string
}

type TrueUpFlags map[string]bool

type RawValueFlags struct {
	AccessToken  string
	GroupId      string
	Me           int
	UsernamesRaw string
}

type ListFlags struct {
	Boolean  BooleanFlags
	RawValue RawValueFlags
	Resolved ResolvedFlags
	TrueUp   TrueUpFlags
}

var booleanFlags = BooleanFlags{
	Approved: false,
	Browser:  false,
	Draft:    false,
	Group:    false,
	Projects: false,
	Ready:    false,
}

var valueFlags = RawValueFlags{
	AccessToken:  "",
	GroupId:      "",
	Me:           0,
	UsernamesRaw: "",
}

func AddListFlags(listCmd *cobra.Command) {
	listFlags := listCmd.PersistentFlags()
	listFlags.BoolVarP(&booleanFlags.Approved, "approved", "a", false, "Include MRs you approved.")
	listFlags.BoolVarP(&booleanFlags.Browser, "browser", "b", false, "Open MRs in the browser.")
	listFlags.BoolVarP(&booleanFlags.Draft, "draft", "d", false, "Include draft MRs.")
	listFlags.BoolVarP(&booleanFlags.Group, "group", "g", false, "ONLY include MRs where the author is listed in the provided users (*see -u, --users*) or the configured usernames.")
	listFlags.BoolVarP(&booleanFlags.Projects, "projects", "p", false, "ONLY include MRs where the author is listed in ANY of the configured projects; but it only returns MRs for projects the author is listed under.")
	listFlags.BoolVarP(&booleanFlags.Ready, "ready", "r", false, "Include mergeable MRs.")
	listFlags.StringVarP(&valueFlags.GroupId, "group-id", "i", "", "Override the configured groud ID.")
	listFlags.IntVarP(&valueFlags.Me, "me", "m", 0, "Override the configured me user ID with the given number.")
	listFlags.StringVarP(&valueFlags.AccessToken, "access-token", "t", "", "Override the configured access token.")
	listFlags.StringVarP(&valueFlags.UsernamesRaw, "users", "u", "", "Override configured usernames and ONLY filter on usernames you provided. Accepts a CSV of usernames.")
}

func GetListFlags(conf *config.Config) (listFlags ListFlags) {
	resolvedFlags, trueUpFlags := resolveListFlags(conf)
	listFlags = ListFlags{
		Boolean:  booleanFlags,
		RawValue: valueFlags,
		Resolved: resolvedFlags,
		TrueUp:   trueUpFlags,
	}
	return listFlags
}

func resolveListFlags(conf *config.Config) (resolvedFlags ResolvedFlags, trueUpFlags TrueUpFlags) {
	resolvedFlags = ResolvedFlags{
		AccessToken: conf.AccessToken,
		GroupId:     conf.GroupId,
		Me:          conf.Me,
		Usernames:   []string{},
	}

	trueUpFlags = make(TrueUpFlags)

	if valueFlags.AccessToken != "" {
		resolvedFlags.AccessToken = valueFlags.AccessToken
		trueUpFlags["shouldAskToUpdateAccessToken"] = true
	}

	if valueFlags.GroupId != "" {
		resolvedFlags.GroupId = valueFlags.GroupId
		trueUpFlags["shouldAskToUpdateGroupId"] = true
	}

	if valueFlags.Me != 0 {
		resolvedFlags.Me = valueFlags.Me
		trueUpFlags["shouldAskToUpdateMe"] = true
	}

	valueFlags.UsernamesRaw = strings.ReplaceAll(valueFlags.UsernamesRaw, " ", "")
	if valueFlags.UsernamesRaw != "" {
		resolvedFlags.Usernames = strings.Split(valueFlags.UsernamesRaw, ",")
	}

	return resolvedFlags, trueUpFlags
}
