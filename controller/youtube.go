package controller

import (
	"github.com/gin-gonic/gin"
	"currency-gin/utils"
	"net/http"
	"encoding/json"
)

type FormatInfo struct {
	Format string `json:"format"`
	FormatNote string `json:"formatNote"`
	Extension string `json:"extension"`
	Url string `json:"url"`
	FileSize string `json:"fileSize"`
}

type VideoInfo struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	WebPageUrl string `json:"webPageUrl"`
	Formats []FormatInfo `json:"formats"`
}

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
		success, data = ctrl.parseVideoInfo(url, targetExt)

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

	c.JSON(http.StatusOK, gin.H{
		"success": success,
		"data":    data,
	})
}

func (ctrl *YoutubeController) parseVideoInfo(url string, targetExt string) (bool, interface{}) {
	success, resultStr := utils.ExecCmd("youtube-dl", "--dump-json", "--no-warnings", url)

	resultMap := map[string]interface{} {}
	if !success {
		return success, resultStr
	}
	err := json.Unmarshal([]byte(resultStr), &resultMap)
	if err != nil {
		utils.Error(err.Error())
		resultMap["result"] = err.Error()
		return false, resultMap
	}
	videoInfo := ctrl.parseVideoInfoByResultMap(resultMap, targetExt)
	return success, videoInfo
}


func (ctrl *YoutubeController) parseVideoInfoByResultMap(resultMap map[string]interface{}, targetExt string) VideoInfo {
	info := VideoInfo{}
	info.Id = resultMap["id"].(string)
	info.Title = resultMap["title"].(string)
	info.Thumbnail = resultMap["thumbnail"].(string)
	info.WebPageUrl = resultMap["webpage_url"].(string)
	formats := resultMap["formats"].([]interface{})
	var formatList []FormatInfo
	for _, value := range formats {
		formatMap := value.(map[string]interface{})
		formatNote := formatMap["format_note"].(string)
		utils.Info("formatNote:", formatNote)

		format := FormatInfo{
			Format:     formatMap["format"].(string),
			FormatNote: formatNote,
			Extension:  formatMap["ext"].(string),
			Url:        formatMap["url"].(string),
		}
		fileSize := formatMap["filesize"]
		if fileSize != nil {
			format.FileSize = utils.FormatFileSize(fileSize.(float64))
		}
		insert := true
		if targetExt != "" {
			insert = targetExt == format.Extension
		}
		if insert {
			formatList = append(formatList, format)
		}
	}
	info.Formats = formatList
	return info
}

func (ctrl *YoutubeController) parseSubtitle(url string) (bool, interface{}) {
	return false, nil
}