package handler

import (
	"fmt"
	"gin/model"
	"gin/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 流式传输处理器 - 支持 Server-Sent Events (SSE)
// @Summary 发送消息到LLM并获取流式响应
// @Description 发送用户消息到LLM并通过流式传输获取响应
// @Tags LLM API
// @Accept json
// @Produce text/event-stream
// @Param request body model.Message true "用户消息"
// @Success 200 {string} string "流式响应数据"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/llm-message/deepseek-stream [post]
func SendMessageToLLMStreamHandler(c *gin.Context) {
	var req model.Message
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": 400, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// 获取数据流通道
	dataChan, errChan := service.SendMessageToLLMStream(req.Content)

	// 获取响应写入器
	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "不支持流式传输", "code": 500})
		return
	}

	// 标记流式传输已开始
	c.Status(http.StatusOK)

	// 循环处理数据和错误
	for {
		select {
		case data, ok := <-dataChan:
			if !ok {
				// 数据通道已关闭
				fmt.Fprint(w, "data: [DONE]\n\n")
				flusher.Flush()
				return
			}

			// 发送数据块
			if data != "" {
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}

		case err := <-errChan:
			if err != nil {
				// 发送错误信息
				fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
				flusher.Flush()
			}
			return

		case <-c.Request.Context().Done():
			// 客户端断开连接
			fmt.Println("客户端断开连接")
			return
		}
	}
}
