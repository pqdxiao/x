package x

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"regexp"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// // ServeOnlyRequest 仅用于适配构建EVENT请求,但不返回,需自行返回
// func ServeOnlyRequest(requestBody string, c *gin.Context) (s string, e error) {
// 	// 使用标准库解析JSON请求
// 	var reqBody map[string]interface{}
// 	if err := json.Unmarshal([]byte(requestBody), &reqBody); err != nil {
// 		Xlog.ErrErr("json.Unmarshal() Error", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return "", err
// 	}

// 	// 构建请求
// 	req := cppeapcore.NewRequestRef()
// 	defer req.Delete()
// 	if err := req.Unmarshal([]byte(requestBody)); err != nil {
// 		logger.ErrErr("req.Unmarshal() Error", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to unmarshal request"})
// 		return "", err
// 	}

// 	// 验证token
// 	if claims, exist := c.Get("jwt-claims"); exist {
// 		cl := claims.(map[string]interface{})
// 		uid := cl["uid"].(float64)
// 		req.SetUserID(int64(uid))
// 	} else {
// 		err := errors.New("get token error")
// 		logger.ErrErr("c.Get(jwt-claims) Error", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return "", err
// 	}

// 	rsp, err := cppeapcore.Serve(req)
// 	if err != nil {
// 		logger.ErrErr("cppeapcore.Serve:"+requestBody, err)
// 		// 错误号转换
// 		strErr := err.Error()
// 		reErrorCode := `C.EAPRespond Error Code :([-]?\d+$)`
// 		re := regexp.MustCompile(reErrorCode)
// 		results := re.FindAllStringSubmatch(strErr, -1)
// 		if len(results) > 0 {
// 			result := results[0]
// 			if len(result) > 1 {
// 				errorCode, _ := strconv.ParseInt(result[1], 10, 32)
// 				errorCode = 0xffffffff + errorCode + 1
// 				errorCodeHex := fmt.Sprintf("0x%X", errorCode)
// 				if _, ok := errorConfig[errorCodeHex]; ok {
// 					errInfo := errorConfig[errorCodeHex]
// 					errInfoStr := errInfo.Zh_cn
// 					lang := c.GetHeader("X-CAXA-Lang")
// 					switch lang {
// 					case "zh-cn":
// 						errInfoStr = errInfo.Zh_cn
// 					case "zh-tw":
// 						errInfoStr = errInfo.Zh_tw
// 					case "en-us":
// 						errInfoStr = errInfo.En_us
// 					}
// 					c.JSON(http.StatusInternalServerError, gin.H{
// 						"code":  errorCodeHex,
// 						"error": errInfoStr,
// 					})
// 					return "", err
// 				} else {
// 					c.JSON(http.StatusInternalServerError, gin.H{
// 						"code":  errorCodeHex,
// 						"error": "unknown error",
// 					})
// 					return "", err
// 				}
// 			}
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return "", err
// 	}
// 	defer rsp.Delete()

// 	//c.Data(http.StatusOK, "application/json", rsp.Body)
// 	return string(rsp.Body), nil
// }
