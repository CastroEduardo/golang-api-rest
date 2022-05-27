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
	User           User           `json:"user" bson:"user"`
	Company        Company        `json:"company" bson:"company"`
	UserPrivileges UserPrivileges `json:"userPrivileges" bson:"userPrivileges"`
	RolUser        RolUser        `json:"rolUser" bson:"rolUser"`
}

type Company struct {
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

type User struct {
	ID              string    `json:"id,omitempty" bson:"_id,omitempty"`
	NickName        string    `json:"nickName" bson:"nickname"`
	Name            string    `json:"name" bson:"name"`
	LastName        string    `json:"lastName" bson:"lastName"`
	Contact         string    `json:"contact" bson:"contact"`
	City            string    `json:"city" bson:"city"`
	Gender          string    `json:"gender" bson:"gender"`
	Email           string    `json:"email" bson:"email"`
	IdRol           string    `json:"idRol" bson:"idrol"`
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
	ToursInit       bool      `json:"toursinit" bson:"toursinit"`
}

type UserPrivileges struct {
	ID           string         `json:"id,omitempty" bson:"_id,omitempty"`
	IdRol        string         `json:"idRol" bson:"idrol"`
	IdCompany    string         `json:"idCompany"  bson:"idcompany"`
	WebAccess    bool           `json:"webAccess" bson:"webaccess"`
	Config       int            `json:"config" bson:"config"`
	TypeUser     int            `json:"typeUser" bson:"typeUser"`
	UrlListblack []UrlLiskBlack `json:"urlListBlack" bson:"urllistblack"`
}

type RolUser struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Status    int       `json:"status" bson:"status"`
	Note      string    `json:"note" bson:"note"`
	Date      time.Time `json:"date" bson:"date"`
	IdCompany string    `json:"idCompany"  bson:"idcompany"`
}

type UrlLiskBlack struct {
	Path string
	Name string
	Mode int
}

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
