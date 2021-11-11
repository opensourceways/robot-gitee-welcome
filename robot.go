package main

import (
	"fmt"
	"github.com/opensourceways/community-robot-lib/giteeclient"

	sdk "gitee.com/openeuler/go-gitee/gitee"
	libconfig "github.com/opensourceways/community-robot-lib/config"
	libplugin "github.com/opensourceways/community-robot-lib/giteeplugin"
	"github.com/sirupsen/logrus"
)

const botName = "welcome"

const welcomeMessage = `Hey ***@%s***, Welcome to %s Community.
All of the projects in %s Community are maintained by ***@%s***.
That means the developers can comment below every pull request or issue to trigger Bot Commands.
Please follow instructions at <%s> to find the details.
`

type iClient interface {
	CreatePRComment(owner, repo string, number int32, comment string) error
	CreateIssueComment(owner, repo string, number string, comment string) error
	GetBot() (sdk.User, error)
}

func newRobot(cli iClient) *robot {
	return &robot{cli: cli}
}

type robot struct {
	cli iClient
}

func (bot *robot) NewPluginConfig() libconfig.PluginConfig {
	return &configuration{}
}

func (bot *robot) getConfig(cfg libconfig.PluginConfig, org, repo string) (*botConfig, error) {
	c, ok := cfg.(*configuration)
	if !ok {
		return nil, fmt.Errorf("can't convert to configuration")
	}
	if bc := c.configFor(org, repo); bc != nil {
		return bc, nil
	}
	return nil, fmt.Errorf("no config for this repo:%s/%s", org, repo)
}

func (bot *robot) RegisterEventHandler(p libplugin.HandlerRegitster) {
	p.RegisterIssueHandler(bot.handleIssueEvent)
	p.RegisterPullRequestHandler(bot.handlePREvent)
}

func (bot *robot) handlePREvent(e *sdk.PullRequestEvent, pc libconfig.PluginConfig, log *logrus.Entry) error {
	action := giteeclient.GetPullRequestAction(e)
	if action != giteeclient.PRActionOpened {
		return nil
	}

	prInfo := giteeclient.GetPRInfoByPREvent(e)
	cfg, err := bot.getConfig(pc, prInfo.Org, prInfo.Repo)
	if err != nil {
		return err
	}

	comment, err := bot.genWelcomeMessage(prInfo.Author, cfg)
	if err != nil {
		return err
	}

	return bot.cli.CreatePRComment(prInfo.Org, prInfo.Repo, prInfo.Number, comment)
}

func (bot *robot) handleIssueEvent(e *sdk.IssueEvent, pc libconfig.PluginConfig, log *logrus.Entry) error {
	ew := giteeclient.NewIssueEventWrapper(e)
	if giteeclient.StatusOpen != ew.GetAction() {
		return nil
	}

	org, repo := ew.GetOrgRep()
	cfg, err := bot.getConfig(pc, org, repo)
	if err != nil {
		return err
	}

	author := ew.GetIssueAuthor()
	number := ew.GetIssueNumber()
	comment, err := bot.genWelcomeMessage(author, cfg)
	if err != nil {
		return err
	}

	return bot.cli.CreateIssueComment(org, repo, number, comment)
}

func (bot robot) genWelcomeMessage(author string, cfg *botConfig) (string, error) {
	b, err := bot.cli.GetBot()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(welcomeMessage, author, cfg.CommunityName, cfg.CommunityName, b.Login, cfg.CommandLink), nil
}
