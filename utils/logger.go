package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Info(a ...interface{})  {
	fmt.Fprintln(gin.DefaultWriter, "[INFO]", time.Now().Format("[2006/01/02 15:04:05]"), a)
}

func Error(a ...interface{})  {
	fmt.Fprintln(gin.DefaultWriter, "[ERROR]", time.Now().Format("[2006/01/02 15:04:05]"), a)
}