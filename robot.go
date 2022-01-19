package main

import (
	"fmt"

	"github.com/opensourceways/community-robot-lib/config"
	framework "github.com/opensourceways/community-robot-lib/robot-gitee-framework"
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/sirupsen/logrus"
)

const botName = "welcome"

const welcomeMessage = `Hey ***@%s***, Welcome to %s Community.
All of the projects in %s Community are maintained by ***@%s***.
That means the developers can comment below every pull request or issue to trigger Bot Commands.
Please follow instructions at **[Here](%s)** to find the details.
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

func (bot *robot) NewConfig() config.Config {
	return &configuration{}
}

func (bot *robot) getConfig(cfg config.Config, org, repo string) (*botConfig, error) {
	c, ok := cfg.(*configuration)
	if !ok {
		return nil, fmt.Errorf("can't convert to configuration")
	}
	if bc := c.configFor(org, repo); bc != nil {
		return bc, nil
	}

	return nil, fmt.Errorf("no config for this repo:%s/%s", org, repo)
}

func (bot *robot) RegisterEventHandler(p framework.HandlerRegitster) {
	p.RegisterIssueHandler(bot.handleIssueEvent)
	p.RegisterPullRequestHandler(bot.handlePREvent)
}

func (bot *robot) handlePREvent(e *sdk.PullRequestEvent, c config.Config, log *logrus.Entry) error {
	if sdk.GetPullRequestAction(e) != sdk.PRActionOpened {
		return nil
	}

	org, repo := e.GetOrgRepo()
	cfg, err := bot.getConfig(c, org, repo)
	if err != nil {
		return err
	}

	comment, err := bot.genWelcomeMessage(e.GetPRAuthor(), cfg)
	if err != nil {
		return err
	}

	return bot.cli.CreatePRComment(org, repo, e.GetPRNumber(), comment)
}

func (bot *robot) handleIssueEvent(e *sdk.IssueEvent, c config.Config, log *logrus.Entry) error {
	if sdk.StatusOpen != e.GetAction() {
		return nil
	}

	org, repo := e.GetOrgRepo()
	cfg, err := bot.getConfig(c, org, repo)
	if err != nil {
		return err
	}

	author := e.GetIssueAuthor()
	comment, err := bot.genWelcomeMessage(author, cfg)
	if err != nil {
		return err
	}

	return bot.cli.CreateIssueComment(org, repo, e.GetIssueNumber(), comment)
}

func (bot robot) genWelcomeMessage(author string, cfg *botConfig) (string, error) {
	b, err := bot.cli.GetBot()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(welcomeMessage, author, cfg.CommunityName, cfg.CommunityName, b.Login, cfg.CommandLink), nil
}
