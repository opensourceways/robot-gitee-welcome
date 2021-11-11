## robot-gitee-welcome

### 概述

welcome 机器人在用户创建issue 或 pull request的时候会添加如下提示语：

`Hey xx, Welcome to xx Community.
All of the projects in xx Community are maintained by xx.
That means the developers can comment below every pull request or issue to trigger Bot Commands.
Please follow instructions at <xx> to find the details.
`

### 配置

例子：

```yaml
config_items:
  - repos:  #robot需管理的仓库列表
     -  owner/repo
     -  owner1
    excluded_repos: #robot 管理列表中需排除的仓库
     - owner1/repo1
    community_name: opensourceways #社区名(必填项)
    command_link: http://opensourceways.cn/command_link #社区机器指令说明连接(必填项)
```

