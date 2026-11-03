package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func SendJSONHelper(client *http.Client,body []byte, url string) error {
        // Создаём HTTP-запрос
        httpReq, err := http.NewRequest(
            http.MethodPost,
            url,
            bytes.NewReader(body),
        )
        if err != nil {
            return fmt.Errorf("failed to create request: %w", err)
        }

        httpReq.Header.Set("Content-Type", "application/json")

        // Отправляем запрос
        resp, err := client.Do(httpReq)
        if err != nil {
            return fmt.Errorf("failed to send request: %w", err)
        }
        defer resp.Body.Close()

        // Проверяем HTTP-статус
        if resp.StatusCode < 200 || resp.StatusCode >= 300 {
            responseBody, _ := io.ReadAll(resp.Body)
            return fmt.Errorf(
                "server returned status %d: %s",
                resp.StatusCode,
                string(responseBody),
            )
        }

        return nil
}
