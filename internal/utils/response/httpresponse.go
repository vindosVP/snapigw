package response

import "github.com/gin-gonic/gin"

type HttpResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Ok(c *gin.Context, status int, data interface{}) {
	resp := HttpResponse{
		Message: "",
		Data:    data,
	}
	c.JSON(status, resp)
}

func OkMsg(c *gin.Context, status int, data interface{}, msg string) {
	resp := HttpResponse{
		Message: msg,
		Data:    data,
	}
	c.JSON(status, resp)
}

func Err(c *gin.Context, status int, msg string) {
	resp := HttpResponse{
		Message: msg,
		Data:    nil,
	}
	c.JSON(status, resp)
}

func AbortErr(c *gin.Context, status int, msg string) {
	resp := HttpResponse{
		Message: msg,
		Data:    nil,
	}
	c.AbortWithStatusJSON(status, resp)
}

func ErrData(c *gin.Context, status int, msg string, data interface{}) {
	resp := HttpResponse{
		Message: msg,
		Data:    data,
	}
	c.JSON(status, resp)
}
