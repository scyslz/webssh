package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
	"webssh/core"

	"github.com/gin-gonic/gin"
)

// ResponseBody 响应信息结构体
type ResponseBody struct {
	Duration string
	Data     interface{}
	Msg      string
}

// TimeCost 计算方法执行耗时
func TimeCost(start time.Time, body *ResponseBody) {
	body.Duration = time.Since(start).String()
}

// CheckSSH 检查ssh连接是否能连接
func CheckSSH(c *gin.Context) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	sshInfo := c.DefaultQuery("sshInfo", "")
	sshClient, err := core.DecodedMsgToSSHClient(sshInfo)
	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}

	err = sshClient.GenerateClient()
	defer sshClient.Close()

	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
	}
	return &responseBody
}

// SaveSsh saves an SSH connection configuration to a JSON file.
var fileMutex sync.Mutex
var filePath = "ssh_list.json"

func SaveSsh(c *gin.Context) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// Decode new SSH configurations from request body
	var newSshConfigs []core.SSHInfo
	if err := c.BindJSON(&newSshConfigs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Write updated data back to file
	updatedData, err := json.MarshalIndent(newSshConfigs, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to marshal JSON: %v", err)})
		return
	}

	if err := ioutil.WriteFile(filePath, updatedData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to write file: %v", err)})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "SSH configuration saved successfully"})
	}
}

// SshList retrieves the list of saved SSH configurations.
func SshList(c *gin.Context) {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return an empty array
			c.JSON(http.StatusOK, []core.SSHInfo{})
			return
		}
		// Other read errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read file: %v", err)})
		return
	}

	var sshList []core.SSHInfo
	if len(data) > 0 {
		if err := json.Unmarshal(data, &sshList); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unmarshal JSON: %v", err)})
			return
		}
	}

	c.JSON(http.StatusOK, sshList)
}
