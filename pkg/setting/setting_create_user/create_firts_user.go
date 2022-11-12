package setting_create_user

import (
	"time"

	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/util"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbPrivilegeuser_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbcompany_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbdepartmentuser_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbrol_user_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbusers_service"
)

var idCompany string = "636a6ab34ec667f66b57ce9d"
var idRolUser string = "636a6e1eb71bb94337dd2d4e"
var idDepartamentUserSys string = "636a6eca83e4d77ab5ac8775"
var idPrivilegeUser string = "636a6ab44ec667f66b57ce9e"
var idUser string = ""

func Create_first_user() string {

	//idCompany = create_company()
	// time.Sleep(1 * time.Second)

	// time.Sleep(1 * time.Second)
	// time.Sleep(1 * time.Second)
	// //idPrivilegeUser = create_privilege_rol_useradmin()
	// time.Sleep(1 * time.Second)
	// //idRolUser = create_rol_users()
	// time.Sleep(1 * time.Second)
	// //idDepartamentUserSys = create_departamentusersys()
	// time.Sleep(1 * time.Second)
	// idUser = create_user()
	// fmt.Println(idCompany)

	return "name"
}

func create_company() string {

	newCompany := authinterfaces.Company_sys{

		NameLong:    "Nombre Largo Empresa #1",
		NameShort:   "Nombre Corto #1",
		Address:     "Direccion empresa  #1 ",
		Slogan:      "Slogan Company #1",
		Phone:       "809-561-2512 / 809-245-5444",
		Status:      1,
		Image:       "logocompany.png",
		Rnc:         "001-0215211-0",
		Others:      "Otros Datos",
		DateAdd:     time.Now(),
		FolderFiles: "9882388812121212121212_3233-2311",
		UrlFiles:    "https://localhost:8000/api/v2/uploads/images",
	}

	result := dbcompany_service.Add(newCompany)

	return result
}

func create_privilege_rol_useradmin() string {
	listPages := []string{""}
	listFunctions := []string{""}
	newPrivilege := authinterfaces.UserPrivileges_sys{
		IdCompany:     idCompany,
		Name:          "Only-USERs",
		Status:        1,
		DefaultPage:   1,
		ListPages:     listPages,
		ListFunctions: listFunctions,
		DateAdd:       time.Now(),
		Note:          "Create New for demo only Users",
	}

	getId := dbPrivilegeuser_service.Add(newPrivilege)

	return getId
	// list1 := []string{{
	// 	Path: "/home1",
	// 	Name: "/home",
	// 	Mode: 1,
	// }, {
	// 	Path: "/home2",
	// 	Name: "/home",
	// 	Mode: 1,
	// }, {
	// 	Path: "/home3",
	// 	Name: "/home",
	// 	Mode: 1,
	// }}

	// // privilege := authinterfaces.UserPrivileges{
	// // 	IdRol:        idRolUser,
	// // 	IdCompany:    idCompany,
	// // 	WebAccess:    true,
	// // 	Config:       1,
	// // 	TypeUser:     2,
	// // 	UrlListblack: list1,
	// // }

	return "getId"
}

func create_rol_users() string {

	//user Rol-admin
	newRolAdmin := authinterfaces.RolUser_sys{
		Name:        "Rol ADMINS",
		Status:      1,
		Note:        "...ADMINS ONLY",
		Date:        time.Now(),
		IdPrivilege: idPrivilegeUser,
		IdCompany:   idCompany,
	}
	getId := dbrol_user_service.Add(newRolAdmin)

	return getId
}

func create_departamentusersys() string {

	newChild := []authinterfaces.DptsUser_sys{}

	newData := authinterfaces.DptsUser_sys{
		Title:     "DPT-ADMIN #1 ",
		OrderNo:   0,
		Status:    1,
		Note:      "DEPARTAMENTO #1 ONLY ADMIN",
		DateAdd:   time.Now(),
		IdCompany: idCompany,
		IdRole:    idRolUser,
		Children:  newChild,
	}

	result := dbdepartmentuser_service.Add(newData)

	return result
}

type Base struct {
	ID string
}

type Child struct {
	Base
	a string
	b string
}

type UrlLiskBlack1 struct {
	path string
	name string
	mode int
}

type UrlLiskBlack struct {
	Path string
	Name string
	Mode int
}

func create_user() string {

	newUser := authinterfaces.User_sys{
		NickName:        "usuario1",
		Name:            "Nombre usuario #1",
		LastName:        "apellido2",
		Contact:         "contact",
		City:            "city",
		Gender:          "male",
		Email:           "usuario1@gmail.com",
		IdDept:          idDepartamentUserSys,
		IdCompany:       idCompany,
		Status:          1,
		Image:           "imagen",
		Note:            "this is  a note for Admin",
		ForcePass:       true,
		Public:          1,
		Password:        util.Encript([]byte("111111")),
		LastLogin:       time.Now(),
		DefaultPathHome: "/dashboard/analysis",
		DateAdd:         time.Now(),
		ToursInit:       true,
	}

	getId := dbusers_service.Add(newUser)
	return getId

}
