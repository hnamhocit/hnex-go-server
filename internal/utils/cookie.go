package utils

import "github.com/gin-gonic/gin"

func SetCookies(c *gin.Context, accessToken, refreshToken string) {
	hourInSeconds, sevenDaysInSeconds := 3600, 604800
	c.SetCookie("access_token", accessToken, hourInSeconds, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, sevenDaysInSeconds, "/", "localhost", false, true)
}

func ClearCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
}
