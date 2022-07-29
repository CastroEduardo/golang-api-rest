package conf

const (
	NameUrlApi1 = "/api/v1"
	NameUrlApi2 = "/api/v2"

	//url articles
	ArticlesParms_GET = "/articles/:id"
	Moderator
	Admin

	//error login Auth user
	UserDisabled = "DISABLED"
	UserFailed   = "FAILED"
	UserSuccess  = "SUCCESS"

	//urlUploadFiles
	NameUrlPathFiles = "/uploads"
	//userSys
	ManagedUserSys   = "/managedUserSys"
	UserSysList_Post = "/userlistsys"

	//rol_privileges_departament
	ControllerDepartamentUserSys = "/managedepartment"
	ControllerRolUserSys         = "/managerolsys"
	ControllerPrivilegeUserSys   = "/manageprivilegesys"

	//managed uploadFiles
	ControllerUploadFiles = "/manageduploadfiles"

	//types EVENTS LOGS
	LOGIN_FAILE                     = 111
	LOGIN_SUCCESS                   = 222
	LOGIN_USER_EVENT_ADD            = 333
	LOGIN_USER_EVENT_UPDATE         = 444
	LOGIN_USER_EVENT_REMOVE         = 555
	LOGIN_USER_SEARCH               = 666
	LOGIN_USER_BLOCK                = 777
	LOGIN_USER_GETCLAIM             = 888
	LOGIN_USER_LOGOUT               = 999
	LOGIN_USER_CHECK_TOKEN          = 1010
	USER_CHECK_ISACCOUNT            = 1111
	USER_EVENT_UPDATE_PASSWORD      = 1112
	USER_EVENT_UPDATE_AVATAR        = 1113
	SYSTEM_EVENT_MOVE_FILE          = 1114
	SYSTEM_EVENT_USER_REQUIRED_FILE = 1115
)
