# ChatBotWeatherNotifier

ChatBotWeatherNotifier is a chatbot that provides weather information via a KakaoTalk channel. It is implemented using Go and the Echo framework, running on an AWS EC2 instance.

## Project Structure

my_weatherBot
├── main.go
├── handler/
│   └── weather.go
└── README.md

## Key Features

- Provides weather information for user-specified locations via KakaoTalk.
- Delivers temperature, precipitation type, humidity, precipitation amount, and wind speed information by time.
- Supports simple text and image responses.

## Installation and Execution

### Prerequisites

- Go (version 1.16 or above)
- AWS EC2 instance
- Kakao Developer account and API key

### Installation

1. Clone the repository.

    ```bash
    git clone https://github.com/your-username/ChatBotWeatherNotifier.git
    cd ChatBotWeatherNotifier
    ```

2. Install necessary Go modules.

    ```bash
    go mod tidy
    ```

3. Run the server on your AWS EC2 instance.

    ```bash
    go run main.go
    ```

4. Set up the skill server URL in the KakaoTalk chatbot admin center.

### Example

1. When the user inputs "Weather in Jamsil" in KakaoTalk, the bot returns the weather information for that location.

    ```json
    {
        "version": "2.0",
        "template": {
            "outputs": [
                {
                    "simpleText": {
                        "text": "14:00 Weather\nSky Status: ☀️ Clear\nPrecipitation Type: None\nTemperature: 25℃\nHumidity: 60%\nPrecipitation: 0mm\nWind Speed: 3m/s\n\n"
                    }
                }
            ]
        }
    }
    ```

## File Description

### main.go

- Initializes the server and sets up routing.
- Defines handlers for KakaoTalk chatbot requests.

### handler/weather.go

- Retrieves and formats weather information.
- Uses the Kakao Local API to get location information.
- Uses the Korean Meteorological Administration API to get and process weather data.

## Environment Variable Setup

- Manage your API keys (Kakao and weather) using environment variables or a configuration file.

    ```go
    const (
        KAKAO_API_KEY = "YOUR_KAKAO_API_KEY"
        WEATHER_API_KEY = "YOUR_WEATHER_API_KEY"
    )
    ```

## Contribution

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -am 'Add some feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Create a Pull Request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
