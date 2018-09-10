package controller

import (
	"github.com/gin-gonic/gin"
	"currency-gin/utils"
	"net/http"
	"encoding/json"
)

type Dict utils.D


type YoutubeController struct {
	gin.Context
}

func (ctrl *YoutubeController) Handle(c *gin.Context) {
	method := c.Param("method")
	if method == "" {
		method = "url"
	}
	utils.Info("youtube y2bMethod:", method, " request method:", c.Request.Method)
	url := c.Query("url")
	targetExt := c.Query("ext")

	if url == "" {
		c.String(http.StatusInternalServerError, "invalid url")
		return
	}

	success := false
	var data interface{}

	switch method {
	case "url":
		success, data = ctrl.parseVideoInfo(url, targetExt, false)
	case "url2":
		success, data = ctrl.parseVideoInfo(url, targetExt, true)
	case "download":
	//	url := c.GetString("url")
	//	formatCode := c.GetString("formatCode")
	//
	//	success, resultData := c.parseDownloadUrl(url, formatCode)
	//
	//	data["success"] = success
	//	data["data"] = resultData
	case "subtitle":
		success, data = ctrl.parseSubtitle(url)
	default:
		break
	}

	c.JSON(http.StatusOK, Dict{
		"success": success,
		"data":    data,
	})
}

func (ctrl *YoutubeController) parseVideoInfo(url string, targetExt string, isShort bool) (bool, interface{}) {
	success, resultStr := utils.ExecCmd("youtube-dl", "--dump-json", "--no-warnings", url)

	resultMap := Dict{}
	if !success {
		resultMap["result"] = resultStr
		return success, resultMap
	}
	err := json.Unmarshal([]byte(resultStr), &resultMap)
	if err != nil {
		utils.Error(err.Error())
		resultMap["result"] = err.Error()
		return false, resultMap
	}
	var videoInfo interface{}
	if isShort {
		videoInfo = ctrl.parseShortVideoInfoByResultMap(resultMap, targetExt)
	} else {
		videoInfo = ctrl.parseVideoInfoByResultMap(resultMap, targetExt)
	}
	return success, videoInfo
}


func (ctrl *YoutubeController) parseVideoInfoByResultMap(resultMap Dict, targetExt string) Dict {
	info := Dict{
		"id":			resultMap["id"],
		"title":		resultMap["title"],
		"thumbnail":	resultMap["thumbnail"],
		"webPageUrl":	resultMap["webpage_url"],
	}
	formats := resultMap["formats"].([]interface{})
	var formatList []Dict
	for _, value := range formats {
		formatMap := value.(map[string]interface{})

		format := formatMap["format"]
		formatNote := formatMap["format_note"]
		extension := formatMap["ext"]
		url := formatMap["url"]
		fileSize := formatMap["filesize"]
		if fileSize != nil {
			fileSize = utils.FormatFileSize(fileSize.(float64))
		}

		formatInfo := Dict{
			"format":     	format,
			"formatNote": 	formatNote,
			"extension":  	extension,
			"url":        	url,
			"fileSize":		fileSize,
		}

		insert := true
		if targetExt != "" {
			insert = targetExt == extension
		}
		if insert {
			formatList = append(formatList, formatInfo)
		}
	}
	info["formats"] = formatList
	return info
}


func (ctrl *YoutubeController) parseShortVideoInfoByResultMap(resultMap Dict, targetExt string) Dict {
	info := Dict{}
	info["title"] = resultMap["title"].(string)
	formats := resultMap["formats"].([]interface{})
	totalFormatMap := map[string]string{}
	for _, value := range formats {
		formatMap := value.(map[string]interface{})

		fileSize := formatMap["filesize"]
		if fileSize == nil {
			continue
		}
		fileSize = utils.FormatFileSize(fileSize.(float64))
		utils.Info("fileSize:", fileSize)

		insert := true
		if targetExt != "" {
			ext := formatMap["ext"].(string)
			insert = targetExt == ext
		}
		if !insert {
			continue
		}
		url := formatMap["url"].(string)
		totalFormatMap[fileSize.(string)] = url
	}
	info["formatMap"] = totalFormatMap
	return info
}

func (ctrl *YoutubeController) parseSubtitle(url string) (bool, interface{}) {
	return false, nil
}