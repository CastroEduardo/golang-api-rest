package setting_create_user

import (
	"fmt"
	"time"

	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/util"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbcompany_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbdepartamentuser_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbprivilege_rol_user_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbrol_user_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbusers_service"
)

var idCompany string = "629aa1f72c2b0990fdde5631"
var idRolUser string = ""
var idDepartamentUserSys string = ""
var idPrivilegeUser string = ""
var idUser string = ""

func Create_first_user() string {

	// idCompany = create_company()
	time.Sleep(1 * time.Second)
	idDepartamentUserSys = create_departamentusersys()
	time.Sleep(1 * time.Second)
	idRolUser = create_rol_useradmin()
	time.Sleep(1 * time.Second)
	idPrivilegeUser = create_privileges_rol_useradmin()
	time.Sleep(1 * time.Second)
	idUser = create_user()

	fmt.Println(idRolUser)
	//idPrivilegeUser = create_privileges_rol_useradmin()

	return "name"
}

func create_company() string {

	newCompany := authinterfaces.Company{

		NameLong:    "Nombre Largo Empresa #1",
		NameShort:   "Nombre Corto #1",
		Address:     "Direccion #1 ",
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

func create_departamentusersys() string {

	newData := authinterfaces.DepartamentUserSys{
		Name:      "Departamento 2",
		Status:    1,
		Note:      "DEPARTAMENTO #2 X",
		Date:      time.Now(),
		IdCompany: idCompany,
	}

	result := dbdepartamentuser_service.Add(newData)

	return result
}

func create_rol_useradmin() string {

	//user Rol-admin
	newRolAdmin := authinterfaces.RolUser{
		Name:          "Rol USUARIO",
		Status:        1,
		Note:          "... User ONLY",
		Date:          time.Now(),
		IdCompany:     idCompany,
		IdDepartament: idDepartamentUserSys,
	}
	getId := dbrol_user_service.Add(newRolAdmin)

	return getId
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

func create_privileges_rol_useradmin() string {

	list1 := []authinterfaces.UrlLiskBlack{{
		Path: "/home1",
		Name: "/home",
		Mode: 1,
	}, {
		Path: "/home2",
		Name: "/home",
		Mode: 1,
	}, {
		Path: "/home3",
		Name: "/home",
		Mode: 1,
	}}

	privilege := authinterfaces.UserPrivileges{
		IdRol:        idRolUser,
		IdCompany:    idCompany,
		WebAccess:    true,
		Config:       1,
		TypeUser:     2,
		UrlListblack: list1,
	}

	getId := dbprivilege_rol_user_service.Add(privilege)

	return getId
}

func create_user() string {

	newUser := authinterfaces.User{
		NickName:        "usuario2",
		Name:            "Nombre usuario #2",
		LastName:        "apellido2",
		Contact:         "contact",
		City:            "city",
		Gender:          "male",
		Email:           "usuario2@gmail.com",
		IdRol:           idRolUser,
		IdCompany:       idCompany,
		Status:          1,
		Image:           "imagen",
		Note:            "this is  a note for user",
		ForcePass:       true,
		Public:          1,
		Password:        util.Encript([]byte("22222")),
		LastLogin:       time.Now(),
		DefaultPathHome: "/dashboard/analysis",
		DateAdd:         time.Now(),
		ToursInit:       true,
	}

	getId := dbusers_service.Add(newUser)
	return getId

}
