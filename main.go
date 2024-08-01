package main

import (
	"my_weatherBot/handler"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	kakaoAPIKey := os.Getenv("KAKAO_API_KEY")
	if kakaoAPIKey == "" {
		panic("KAKAO_API_KEY environment variable is not set")
	}

	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		panic("WEATHER_API_KEY environment variable is not set")
	}

	handler.SetAPIKeys(kakaoAPIKey, weatherAPIKey)

	e.POST("/contextcontrol", func(c echo.Context) error {
		data := map[string]interface{}{
			"version": "2.0",
			"context": map[string]interface{}{
				"values": []map[string]interface{}{
					{
						"name":     "abc",
						"lifeSpan": 10,
						"ttl":      60,
						"params": map[string]interface{}{
							"key1": "val1",
							"key2": "val2",
						},
					},
					{
						"name":     "def",
						"lifeSpan": 5,
						"params": map[string]interface{}{
							"key3": "1",
							"key4": "true",
							"key5": "{\"jsonKey\": \"jsonVal\"}",
						},
					},
					{
						"name":     "ghi",
						"lifeSpan": 0,
					},
				},
			},
		}
		return c.JSON(http.StatusOK, data)
	})

	e.POST("/weather", weather)

	e.POST("/simpletext", func(c echo.Context) error {
		data := map[string]interface{}{
			"version": "2.0",
			"template": map[string]interface{}{
				"outputs": []map[string]interface{}{
					{
						"simpleText": map[string]interface{}{
							"text": "간단한 텍스트 요소입니다.",
						},
					},
				},
			},
		}
		return c.JSON(http.StatusOK, data)
	})

	e.POST("/simpleimage", func(c echo.Context) error {
		data := map[string]interface{}{
			"version": "2.0",
			"template": map[string]interface{}{
				"outputs": []map[string]interface{}{
					{
						"simpleImage": map[string]interface{}{
							"imageUrl": "https://t1.kakaocdn.net/openbuilder/sample/lj3JUcmrzC53YIjNDkqbWK.jpg",
							"altText":  "보물상자입니다",
						},
					},
				},
			},
		}
		return c.JSON(http.StatusOK, data)
	})

	e.POST("/basiccard", func(c echo.Context) error {
		data := map[string]interface{}{
			"version": "2.0",
			"template": map[string]interface{}{
				"outputs": []map[string]interface{}{
					{
						"textCard": map[string]interface{}{
							"title":       "챗봇 관리자센터에 오신 것을 환영합니다.",
							"description": "챗봇 관리자센터로 챗봇을 제작해 보세요. \n카카오톡 채널과 연결하여, 이용자에게 챗봇 서비스를 제공할 수 있습니다.",
							"buttons": []map[string]interface{}{
								{
									"action":     "webLink",
									"label":      "소개 보러가기",
									"webLinkUrl": "https://chatbot.kakao.com/docs/getting-started-overview/",
								},
								{
									"action":     "webLink",
									"label":      "챗봇 만들러 가기",
									"webLinkUrl": "https://chatbot.kakao.com/",
								},
							},
						},
					},
				},
			},
		}

		return c.JSON(http.StatusOK, data)
	})

	e.POST("/commercecard", func(c echo.Context) error {
		data := map[string]interface{}{
			"version": "2.0",
			"template": map[string]interface{}{
				"outputs": []map[string]interface{}{
					{
						"basicCard": map[string]interface{}{
							"title":       "보물상자",
							"description": "보물상자 안에는 뭐가 있을까",
							"thumbnail": map[string]interface{}{
								"imageUrl": "https://t1.kakaocdn.net/openbuilder/sample/lj3JUcmrzC53YIjNDkqbWK.jpg",
							},
							"buttons": []map[string]interface{}{
								{
									"action":      "message",
									"label":       "열어보기",
									"messageText": "짜잔! 우리가 찾던 보물입니다",
								},
								{
									"action":     "webLink",
									"label":      "구경하기",
									"webLinkUrl": "https://e.kakao.com/t/hello-ryan",
								},
							},
						},
					},
				},
			},
		}

		return c.JSON(http.StatusOK, data)
	})

	e.Start(":8080")
}

// 핸들러 함수
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func weather(c echo.Context) error {
	var req struct {
		Action struct {
			Params struct {
				Location string `json:"location"`
			} `json:"params"`
		} `json:"action"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"version": "2.0",
			"template": map[string]interface{}{
				"outputs": []map[string]interface{}{
					{
						"simpleText": map[string]interface{}{
							"text": "파라미터를 읽는 데 실패했습니다.",
						},
					},
				},
			},
		})
	}

	location := req.Action.Params.Location

	resp, err := handler.GetswyWeather(location)
	if err != nil {
		data := map[string]interface{}{
			"version": "2.0",
			"template": map[string]interface{}{
				"outputs": []map[string]interface{}{
					{
						"simpleText": map[string]interface{}{
							"text": err.Error(),
						},
					},
				},
			},
		}
		return c.JSON(http.StatusBadRequest, data)
	}

	data := map[string]interface{}{
		"version": "2.0",
		"template": map[string]interface{}{
			"outputs": []map[string]interface{}{
				{
					"simpleText": map[string]interface{}{
						"text": resp,
					},
				},
			},
		},
	}

	return c.JSON(http.StatusOK, data)
}
