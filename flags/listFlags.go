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
	listFlags.BoolVarP(&booleanFlags.Approved, "approved", "a", false, "Filter output to include MRs approved by the configured me user ID.")
	listFlags.BoolVarP(&booleanFlags.Browser, "browser", "b", false, "Open merge requests in the browser.")
	listFlags.BoolVarP(&booleanFlags.Draft, "draft", "d", false, "Filter output to include draft merge requests.")
	listFlags.BoolVarP(&booleanFlags.Group, "group", "g", false, "Filter output to the usernames configuration.")
	listFlags.BoolVarP(&booleanFlags.Projects, "projects", "p", false, "Filter output to the projects configuration.")
	listFlags.BoolVarP(&booleanFlags.Ready, "ready", "r", false, "Filter output to include merge requests that are ready to merge.")
	listFlags.StringVarP(&valueFlags.GroupId, "group-id", "i", "", "Override the configured groud ID.")
	listFlags.IntVarP(&valueFlags.Me, "me", "m", 0, "Override the configured me user ID.")
	listFlags.StringVarP(&valueFlags.AccessToken, "access-token", "t", "", "Override the configured access token.")
	listFlags.StringVarP(&valueFlags.UsernamesRaw, "users", "u", "", "Filter output to the specified usernames.")
}

func DescribeListFlags() string {
	return `
    - Use the '-a, --approved' flag to filter output to include MRs approved by the configured 'me' user ID.
    - Use the '-b, --browser' flag to open MRs in the browser.
    - Use the '-d, --drafts' flag to include draft MRs.
    - Use the '-g, --group' flag to filter output to the usernames configuration.
    - Use the '-i, --group-id' flag to override the configured group ID.
    - Use the '-m, --me' flag to override the configured me user ID.
    - Use the '-p, --projects' flag to filter output to the projects configuration.
    - Use the '-r, --ready' flag to filter output to include merge requests that are ready to merge.
    - Use the '-t, --access-token' flag to override the configured access token.
    - Use the '-u, --users' flag to override configured usernames and only filter on usernames you provided. Accepts a CSV string of usernames.`
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
