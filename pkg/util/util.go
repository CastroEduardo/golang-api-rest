package util

import "github.com/CastroEduardo/golang-api-rest/pkg/setting"

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}
