## robot-gitee-welcome

[中文README](README_zh_CN.md)

###  Overview

The welcome bot adds the following prompt when a user creates an issue or pull request：

`Hey xx, Welcome to xx Community.
All of the projects in xx Community are maintained by xx.
That means the developers can comment below every pull request or issue to trigger Bot Commands.
Please follow instructions at <xx> to find the details.`

###  Configuration

```yaml
config_items:
  - repos:  #list of repositories to be managed by robot
     -  owner/repo
     -  owner1
    excluded_repos: #Robot manages the list of repositories to be excluded
     - owner1/repo1
    community_name: opensourceways #community Name (required field)
    command_link: http://opensourceways.cn/command_link #community robot command instructions url (required field)
```

