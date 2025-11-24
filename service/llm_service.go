package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"gin/model"
	"io"
	"net/http"
	"strings"
	"time"
	"gin/config"
)

// 流式请求 - 返回一个通道用于接收数据块
func SendMessageToLLMStream(message string) (<-chan string, <-chan error) {
	dataChan := make(chan string, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(dataChan)
		defer close(errChan)

		url := "https://api.deepseek.com/chat/completions"

		// 构建请求体，启用流式传输
		requestBody := map[string]interface{}{
			"model": "deepseek-chat",
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": message,
				},
			},
			"stream": true, // 启用流式传输
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			errChan <- fmt.Errorf("无法编码请求体: %v", err)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			errChan <- fmt.Errorf("无法创建请求: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer " + config.DeepseekAPIKey)

		client := &http.Client{Timeout: 5 * time.Minute} // 流式传输需要更长的超时
		resp, err := client.Do(req)
		if err != nil {
			errChan <- fmt.Errorf("请求失败: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errChan <- fmt.Errorf("API 错误 (状态码 %d): %s", resp.StatusCode, string(body))
			return
		}

		// 使用 bufio.Scanner 逐行读取流式数据
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			// 跳过空行
			if line == "" {
				continue
			}

			// DeepSeek API 返回的格式：data: {...json...}
			if strings.HasPrefix(line, "data: ") {
				jsonStr := strings.TrimPrefix(line, "data: ")

				// 检查是否是结束标记
				if jsonStr == "[DONE]" {
					break
				}

				// 解析 JSON
				var streamResponse model.ChatStreamResponse
				if err := json.Unmarshal([]byte(jsonStr), &streamResponse); err != nil {
					// 跳过无法解析的行，继续处理下一行
					fmt.Printf("警告: 无法解析流数据: %v\n", err)
					continue
				}

				// 提取文本内容
				if len(streamResponse.Choices) > 0 && streamResponse.Choices[0].Delta.Content != "" {
					dataChan <- streamResponse.Choices[0].Delta.Content
				}
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("读取流数据失败: %v", err)
		}
	}()

	return dataChan, errChan
}
