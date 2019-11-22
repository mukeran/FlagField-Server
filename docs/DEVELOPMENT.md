# FlagField-Server Development Document

FlagField is a contest-based "Capture The Flag" system.

We intended to build a full-functional open-source CTF platform.

[TOC]

## Changelog

- v0.0.1 Initialized by mukeran. FlagField is now open sourced.



## Techniques Used

`golang>=1.13`  `gin>=1.4.0`   `gorm>=1.9.10`



## Code Style

Standard Go style.



## Project Layout

Partly obey standards from https://github.com/golang-standards/project-layout.

```
.
├── cmd
│   ├── manager
│   ├── migrator
│   ├── server
│   └── setup
│       └── main.go // Main entries for each app.
├── docs
│   └── DEVELOPMENT.md
├── internal // Main sources
│   ├── manager
│   ├── migrator
│   ├── server
│   │   ├── controllers
│   │   ├── instance
│   │   ├── middlewares
│   │   └── server.go
│   └── setup
├── test // Units test
├── config.example.json // Example config file
├── docker-compose.yml
├── docker-entrypoint.sh // Docker entrypoint
├── Dockerfile
├── go.mod
├── Makefile
└── README.md
```



## API Structure /v1

FlagField is using RESTful-like API。

### /admin

```
.
├── GET    Get admin list  List()
├── POST   Add admin       Add()
└── DELETE Delete admin    Delete()
```

### /config

```
.
├── GET     Get configure list  List()
├── PATCH   Edit configure      Set()
└── :configKey
    └── GET Get configure       Get()
```

### /contest

```
.
├── GET  Get contest list   List()
├── POST Add a new contest  Create()
└── :contestID (__practice__ means default practice contest)
    ├── GET    Get contest info           Show()
    ├── PATCH  Edit contest info          Modify()
    ├── DELETE Delete particular contest  SingleDelete()
    ├── admin
    │   ├── GET    Get contest admin list  AdminList()
    │   ├── POST   Add contest admin       AdminAdd()
    │   └── DELETE Delete contest admin    AdminDelete()
    ├── team
    │   ├── GET    Get team snapshots        TeamList()
    │   ├── POST   Add team to contest       TeamAdd()
    │   └── DELETE Delete team from contest  TeamDelete()
    ├── problem
    │   ├── GET    Get contest problems  ProblemHandler{}.List()
    │   ├── POST   Add a problem         ProblemHandler{}.Create()
    │   └── :problemAlias
    │       ├── GET    Get problem info   ProblemHandler{}.Show()
    │       ├── PATCH  Edit problem info  ProblemHandler{}.Modify()
    │       ├── DELETE Delete problem     ProblemHandler{}.SingleDelete()
    │       ├── flag
    │       │   ├── GET    Get flag list      FlagHandler{}.List()
    │       │   ├── POST   Add a flag         FlagHandler{}.Create()
    │       │   ├── DELETE Batch delete flag  FlagHandler{}.BatchDelete()
    │       │   └── :flagOrder
    │       │       ├── GET    Get flag info       FlagHandler{}.Show()
    │       │       ├── PATCH  Edit flag info      FlagHandler{}.Modify()
    │       │       └── DELETE Delete single flag  FlagHandler{}.SingleDelete()
    │       ├── hint
    │       │   ├── GET    Get hint list  HintHandler{}.List()
    │       │   ├── POST   Add a hint     HintHandler{}.Create()
    │       │   └── :hintOrder
    │       │       ├── GET    Get hint info       HintHandler{}.Show()
    │       │       ├── PATCH  Edit hint info      HintHandler{}.Modify()
    │       │       └── DELETE Delete single hint  HintHandler{}.SingleDelete()
    │       ├── submission
    │       │   └── POST Submit flag  ProblemHandler{}.SubmissionCreate()
    │       └── tag
    │           ├── POST   Add tags     ProblemHandler{}.TagAdd()
    │           └── DELETE Delete tags  ProblemHandler{}.TagDelete()
    ├── notification
    │   ├── GET  Get contest notification list  NotificationList()
    │   ├── POST Add a contest notification     NotificationCreate()
    │   └── :notificationOrder
    │       └── DELETE Delete a notification    NotificationSingleDelete()
    └── statistic
        └── GET Get contest statistic  StatisticShow()
```

### /notification

```
.
├── GET    Get personal notification    List()
├── POST   Add a notification           New()
├── DELETE Delete notifications         Delete()
└── :notificationID
    └── PATCH Mark notification as read MarkRead()
```

### /resource

```
.
├── GET  Get resource list List()
├── POST Upload resource   Upload()
└── /:resourceUUID
    ├── GET    Download resource           Download()
    ├── PATCH  Edit resource               Modify()
    └── DELETE Delete particular resource  SingleDelete()
```

### /session

```
.
├── GET    Get session list           Show()
├── POST   Create a session           Create()
├── DELETE Destroy current session    Destroy()
└── /__current__/
    └── GET Get current session info  ViewCurrent()
```

### /statistic

```
.
└── GET Get index statistics Show()
```

### /submission

```
.
└── GET Get submission list List()
```

### /user

```
.
├── GET  Get user list       List()
├── POST Add/Register a user Register()
└── :userID
    ├── GET    Get user info    Show()
    ├── DELETE Delete this user SingleDelete()
    ├── password
    │   └── PUT Edit password   PasswordModify()
    ├── email
    │   └── PUT Edit email      EmailModify()
    ├── profile
    │   └── PUT Edit profile    ProfileModify()
    └── team
        ├── GET             Get teams this user joined     TeamList()
        ├── GET invitation  Get team invitations received  TeamInvitationList()
        └── GET application Get team applications sent     TeamApplicationList()
```

### /team

```
.
├── GET  Get team list  List()
├── POST Add a team     Create()
└── :teamID
    ├── GET    Get team info     Show()
    ├── PATCH  Edit team info    Modify()
    ├── DELETE Delete this team  SingleDelete()
    ├── admin
    │   ├── GET    Get team admin list  AdminList()
    │   ├── POST   Add team admin       AdminAdd()
    │   ├── DELETE Delete team admin    AdminDelete()
    ├── user
    │   ├── GET    Get team members    UserShow()
    │   ├── POST   Add team member     UserAdd()
    │   └── DELETE Delete team member  UserDelete()
    ├── statistic
    │   └── GET Get team statistics  StatisticShow()
    ├── invitation
    │   ├── GET           Get invitations sent by team   InvitationList()
    │   ├── POST          Send invitation                InvitationNew()
    │   ├── DELETE        Cancel invitation              InvitationCancel()
    │   ├── GET accept    Accept invitation              InvitationAccept()
    │   ├── POST accept   Accept invitation using token  InvitationAcceptByToken()
    │   ├── GET reject    Reject invitation              InvitationReject()
    │   ├── GET token     Get team invitation token      InvitationTokenShow()
    │   └── DELETE token  Refresh team invitation token  InvitationTokenRefresh()
    └── application
        ├── GET         Get team applications  ApplicationList()
        ├── POST        Send team application  ApplicationNew()
        ├── DELETE      Cancel application     ApplicationCancel()
        ├── POST accept Accept application     ApplicationAccept()
        └── POST reject Reject application     ApplicationReject()
```

### /time

```
.
└── GET Get system time  TimeShow()
```

