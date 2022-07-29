package authinterfaces

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type LoginUser struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type IUserLogin struct {
	IdUser string `json:"idUser"`
}

type ClaimSession struct {
	User_sys           User_sys           `json:"user_sys" bson:"user_sys"`
	Company_sys        Company_sys        `json:"company_sys" bson:"company_sys"`
	UserPrivileges_sys UserPrivileges_sys `json:"userPrivileges_sys" bson:"user_privilege_sys"`
	RolUser_sys        RolUser_sys        `json:"rolUser_sys" bson:"rol_user_sys"`
	DeptUser_sys       DptsUser_sys       `json:"deptUser_sys" bson:"dept_user_sys"`
}
type Company_sys struct {
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	NameLong    string    `json:"nameLong" bson:"namelong"`
	NameShort   string    `json:"nameShort" bson:"nameshort"`
	Address     string    `json:"address" bson:"address"`
	Slogan      string    `json:"slogan" bson:"slogan"`
	Phone       string    `json:"phone" bson:"phone"`
	Status      int       `json:"status" bson:"status"`
	Image       string    `json:"image" bson:"image"`
	Rnc         string    `json:"rnc" bson:"rnc"`
	Others      string    `json:"others" bson:"others"`
	DateAdd     time.Time `json:"dateAdd" bson:"dateadd"`
	FolderFiles string    `json:"folderFiles" bson:"folderfiles"`
	UrlFiles    string    `json:"urlFiles" bson:"urlfiles"`
}

type User_sys struct {
	ID              string    `json:"id,omitempty" bson:"_id,omitempty"`
	NickName        string    `json:"nickName" bson:"nickname"`
	Name            string    `json:"name" bson:"name"`
	LastName        string    `json:"lastName" bson:"lastName"`
	Contact         string    `json:"contact" bson:"contact"`
	City            string    `json:"city" bson:"city"`
	Gender          string    `json:"gender" bson:"gender"`
	Email           string    `json:"email" bson:"email"`
	Address         string    `json:"address" bson:"address"`
	IdDept          string    `json:"idDept" bson:"iddept"`
	IdCompany       string    `json:"idCompany"  bson:"idcompany"`
	Status          int       `json:"status" bson:"status"`
	Image           string    `json:"image" bson:"image"`
	Note            string    `json:"note" bson:"note"`
	ForcePass       bool      `json:"forcePass" bson:"forcepass"`
	Public          int       `json:"public" bson:"public"`
	Password        string    `json:"password" bson:"password"`
	LastLogin       time.Time `json:"lastLogin" bson:"lastlogin"`
	DefaultPathHome string    `json:"defaultPathHome" bson:"defaultpathhome"`
	DateAdd         time.Time `json:"dateAdd" bson:"dateadd"`
	ToursInit       bool      `json:"toursInit" bson:"toursinit"`
}

type UserPrivileges_sys struct {
	ID            string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string    `json:"name" bson:"name"`
	Status        int       `json:"status" bson:"status"`
	IdCompany     string    `json:"idCompany"  bson:"idcompany"`
	DefaultPage   int       `json:"defaultPage" bson:"default_page"`
	ListPages     []string  `json:"listPages" bson:"list_pages"`
	ListFunctions []string  `json:"listFunctions" bson:"list_functions"`
	DateAdd       time.Time `json:"dateAdd" bson:"dateadd"`
	Note          string    `json:"note"  bson:"note"`
}

// type UrlLiskBlack struct {
// 	Path string
// 	Name string
// 	Mode int
// }

type RolUser_sys struct {
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Status      int       `json:"status" bson:"status"`
	Note        string    `json:"note" bson:"note"`
	Date        time.Time `json:"date" bson:"date"`
	IdPrivilege string    `json:"idPrivilege"  bson:"idprivilege"`
	IdCompany   string    `json:"idCompany"  bson:"idcompany"`
}

//####DEPARTMENTS USER
type DptsUser_sys struct {
	ID        string         `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string         `json:"title" bson:"title"`
	Key       string         `json:"key" bson:"key"`
	OrderNo   int            `json:"orderNo" bson:"orderno"`
	Status    int            `json:"status" bson:"status"`
	Note      string         `json:"note" bson:"note"`
	IdCompany string         `json:"idCompany"  bson:"idcompany"`
	IdRole    string         `json:"idRole"  bson:"idrole"`
	DateAdd   time.Time      `json:"dateAdd" bson:"dateadd"`
	Children  []DptsUser_sys `json:"children" bson:"children"`
}

//####Tree options PRIVILEGES
type TreePrivileges_sys struct {
	ID        string               `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string               `json:"name" bson:"name"`
	Key       string               `json:"key" bson:"key"`
	Status    int                  `json:"status" bson:"status"`
	Note      string               `json:"note" bson:"note"`
	IdCompany string               `json:"idCompany"  bson:"idcompany"`
	DateAdd   time.Time            `json:"dateAdd" bson:"dateadd"`
	Children  []TreePrivileges_sys `json:"children" bson:"children"`
}

// type DtChildrenDptsUser_sys struct {
// 	ID    string `json:"id,omitempty" bson:"_id,omitempty"`
// 	Title string `json:"title" bson:"title"`
// 	Key   string `json:"key" bson:"key"`
// 	Note  string `json:"note" bson:"note"`
// 	// Children
// }

type DepartmentUserSysOld struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Status    int       `json:"status" bson:"status"`
	Note      string    `json:"note" bson:"note"`
	Date      time.Time `json:"date" bson:"date"`
	IdCompany string    `json:"idCompany"  bson:"idcompany"`
}

//########END DEPARTMENTS

type SessionUser struct {
	ID             string    `json:"id,omitempty" bson:"_id,omitempty"`
	IdUser         string    `json:"idUser" bson:"iduser"`
	IdCompany      string    `json:"idCompany"  bson:"idcompany"`
	Token          string    `json:"token" bson:"token"`
	Active         bool      `json:"active" bson:"active"`
	Remember       bool      `json:"remember" bson:"remember"`
	TokenExpire    time.Time `json:"tokenExpire" bson:"tokenExpire"`
	DateAdd        time.Time `json:"dateAdd" bson:"dateadd"`
	DateLogout     time.Time `json:"dateLogout" bson:"datelogout"`
	LastUpdateTime time.Time `json:"lastUpdateTime" bson:"lastupdatetime"`
}

type Token struct {
	IdSession string
	jwt.StandardClaims
}

type RequestUpdatePassUser struct {
	ID             string `json:"id"`
	NewPassword    string `json:"newPassword"`
	OldPassword    string `json:"oldPassword"`
	RepeatPassword string `json:"repeatPassword"`
}
