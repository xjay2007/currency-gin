package controller

import (
	"github.com/gin-gonic/gin"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"currency-gin/utils"
	"fmt"
)

type ConverterController struct {
	gin.Context
}

func (ctrl *ConverterController) Handle(c *gin.Context) {
	fromUnit := c.Param("fromUnit")
	toUnit := c.Param("toUnit")
	if fromUnit == "" {
		fromUnit = "USD"
	}
	if toUnit == "" {
		toUnit = "CNY"
	}
	fromUnit = strings.ToUpper(fromUnit)
	toUnit = strings.ToUpper(toUnit)
	utils.Info("fromUnit:", fromUnit, " toUnit:", toUnit)

	key := fromUnit + "_" + toUnit
	url := "https://free.currencyconverterapi.com/api/v6/convert?q=" + key + "&compact=y"
	utils.Info("request url: ", url)

	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		c.String(http.StatusServiceUnavailable, err.Error())
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.String(http.StatusServiceUnavailable, err.Error())
		return
	}
	utils.Info("response:", string(body))

	resultMap := map[string]map[string]interface{} {}
	err = json.Unmarshal(body, &resultMap)
	if err != nil {
		c.String(http.StatusServiceUnavailable, err.Error())
		return
	}
	rate := resultMap[key]["val"].(float64)

	c.JSON(http.StatusOK, gin.H{
		"rate":		fmt.Sprintf("%.6f", rate),
		"from":		fromUnit,
		"to":		toUnit,
	})
}