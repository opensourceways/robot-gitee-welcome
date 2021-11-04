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
	return nil, fmt.Errorf("no %s robot config for this repo:%s/%s", botName, org, repo)
}

func (bot *robot) RegisterEventHandler(p libplugin.HandlerRegitster) {
	p.RegisterIssueHandler(bot.handleIssueEvent)
	p.RegisterPullRequestHandler(bot.handlePREvent)
}

func (bot *robot) handlePREvent(e *sdk.PullRequestEvent, cfg libconfig.PluginConfig, log *logrus.Entry) error {
	action := giteeclient.GetPullRequestAction(e)
	if action != giteeclient.PRActionOpened {
		return nil
	}

	prInfo := giteeclient.GetPRInfoByPREvent(e)
	botConfig, err := bot.getConfig(cfg, prInfo.Org, prInfo.Repo)
	if err == nil {
		return err
	}

	comment, err := bot.genWelcomeMessage(prInfo.Author, botConfig)
	if err != nil {
		return err
	}

	return bot.cli.CreatePRComment(prInfo.Org, prInfo.Repo, prInfo.Number, comment)
}

func (bot *robot) handleIssueEvent(e *sdk.IssueEvent, cfg libconfig.PluginConfig, log *logrus.Entry) error {
	if giteeclient.StatusOpen != *e.Action {
		return nil
	}

	org, repo := giteeclient.GetOwnerAndRepoByIssueEvent(e)
	bCfg, err := bot.getConfig(cfg, org, repo)
	if err == nil {
		return err
	}

	author := e.Issue.User.Login
	number := e.Issue.Number
	comment, err := bot.genWelcomeMessage(author, bCfg)
	if err != nil {
		return err
	}

	return bot.cli.CreateIssueComment(org, repo, number, comment)
}

func (bot robot) genWelcomeMessage(author string, bCfg *botConfig) (string, error) {
	b, err := bot.cli.GetBot()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(welcomeMessage, author, bCfg.CommunityName, bCfg.CommunityName, b.Login, bCfg.CommandLink), nil
}
