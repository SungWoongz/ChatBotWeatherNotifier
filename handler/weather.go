package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

var (
	kakaoAPIKey   string
	weatherAPIKey string
)

// Response represents the entire JSON structure.
type Response struct {
	Documents []Document `json:"documents"`
	Meta      Meta       `json:"meta"`
}

type Document struct {
	AddressName       string `json:"address_name"`
	CategoryGroupCode string `json:"category_group_code"`
	CategoryGroupName string `json:"category_group_name"`
	CategoryName      string `json:"category_name"`
	Distance          string `json:"distance"`
	ID                string `json:"id"`
	Phone             string `json:"phone"`
	PlaceName         string `json:"place_name"`
	PlaceURL          string `json:"place_url"`
	RoadAddressName   string `json:"road_address_name"`
	X                 string `json:"x"`
	Y                 string `json:"y"`
}

// Meta represents the meta structure in the JSON.
type Meta struct {
	IsEnd         bool `json:"is_end"`
	PageableCount int  `json:"pageable_count"`
	SameName      struct {
		Keyword        string   `json:"keyword"`
		Region         []string `json:"region"`
		SelectedRegion string   `json:"selected_region"`
	} `json:"same_name"`
	TotalCount int `json:"total_count"`
}

var codeToWeather = map[string]string{
	"T1H": "기온(℃)",
	"RN1": "1시간 강수량(mm)",
	"SKY": "하늘상태",
	"UUU": "동서바람(m/s)",
	"VVV": "남북바람(m/s)",
	"REH": "습도(%)",
	"PTY": "날씨", //강수형태
	"LGT": "낙뢰(kA)",
	"VEC": "풍향(deg)",
	"WSD": "풍속(m/s)",
}

var ptyMap = map[string]string{
	"0": "☀️ 맑음",
	"1": "🌧️ 비",
	"2": "🌧️ 🌨️ 비/눈",
	"3": "🌨️ 눈",
	"5": "🌧️ 빗방울",
	"6": "🌧️ 🌨️ 빗방울눈날림",
	"7": "🌨️ 눈날림",
}

var skyMap = map[string]string{
	"1": "☀️ 맑음",
	"3": "⛅ 구름많음",
	"4": "☁️ 흐림",
}

func SetAPIKeys(kakaoKey, weatherKey string) {
	kakaoAPIKey = kakaoKey
	weatherAPIKey = weatherKey
}

func GetswyWeather(location string) (string, error) {
	resp, err := getLocationInfo(location)
	if err != nil {
		fmt.Printf("error while getLocationInfo: %v", err)
		return "", err
	}
	if len(resp.Documents) == 0 {
		err := fmt.Errorf("There is no data")
		return "", err
	}

	lon, err := strconv.ParseFloat(resp.Documents[0].X, 64)
	if err != nil {
		return "", fmt.Errorf("error parsing longitude: %v", err)
	}
	lat, err := strconv.ParseFloat(resp.Documents[0].Y, 64)
	if err != nil {
		return "", fmt.Errorf("error parsing latitude: %v", err)
	}

	// Convert the coordinates to grid x, y
	x, y := convertLonLatToGrid(lon, lat)
	weather, err := getWeatherForecast(x, y)
	if err != nil {
		return "", err
	}

	return weather, nil
}

func getLocationInfo(location string) (Response, error) {
	// Define the API endpoint and query parameters
	baseURL := "https://dapi.kakao.com/v2/local/search/keyword"
	params := url.Values{}
	params.Add("query", location)
	params.Add("size", "1")

	// Create the request URL
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Create a new HTTP request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the required headers
	req.Header.Set("Authorization", kakaoAPIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v", err)
		return Response{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading request: %v", err)
		return Response{}, err
	}

	// Print the response
	fmt.Println(string(body))

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v", err)
		return Response{}, err
	}

	return response, nil
}

func convertLonLatToGrid(lon, lat float64) (int, int) {
	const (
		Re    = 6371.00877   // 지도반경
		Grid  = 5.0          // 격자간격 (km)
		Slat1 = 30.0         // 표준위도 1
		Slat2 = 60.0         // 표준위도 2
		Olon  = 126.0        // 기준점 경도
		Olat  = 38.0         // 기준점 위도
		Xo    = 210.0 / Grid // 기준점 X좌표
		Yo    = 675.0 / Grid // 기준점 Y좌표
	)
	const PI = math.Pi
	const DEGRAD = PI / 180.0

	re := Re / Grid
	slat1 := Slat1 * DEGRAD
	slat2 := Slat2 * DEGRAD
	olon := Olon * DEGRAD
	olat := Olat * DEGRAD

	sn := math.Tan(PI*0.25+slat2*0.5) / math.Tan(PI*0.25+slat1*0.5)
	sn = math.Log(math.Cos(slat1)/math.Cos(slat2)) / math.Log(sn)
	sf := math.Tan(PI*0.25 + slat1*0.5)
	sf = math.Pow(sf, sn) * math.Cos(slat1) / sn
	ro := math.Tan(PI*0.25 + olat*0.5)
	ro = re * sf / math.Pow(ro, sn)

	ra := math.Tan(PI*0.25 + (lat * DEGRAD * 0.5))
	ra = re * sf / math.Pow(ra, sn)
	theta := lon*DEGRAD - olon
	if theta > PI {
		theta -= 2.0 * PI
	}
	if theta < -PI {
		theta += 2.0 * PI
	}
	theta *= sn
	x := ra*math.Sin(theta) + Xo
	y := ro - ra*math.Cos(theta) + Yo

	return int(math.Ceil(x)), int(math.Ceil(y))
}

func getWeatherForecast(nx, ny int) (string, error) {
	baseURL := "http://apis.data.go.kr/1360000/VilageFcstInfoService_2.0/getUltraSrtFcst"
	params := url.Values{}
	params.Add("serviceKey", weatherAPIKey)
	params.Add("numOfRows", "60")
	params.Add("pageNo", "1")

	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return "", fmt.Errorf("error loading location: %v", err)
	}

	now := time.Now().In(location)
	if now.Minute() < 30 {
		now = now.Add(-1 * time.Hour)
	}
	baseDate := now.Format("20060102")
	baseTime := fmt.Sprintf("%02d00", now.Hour())

	params.Add("base_date", baseDate)
	params.Add("base_time", baseTime)
	params.Add("nx", fmt.Sprintf("%d", nx))
	params.Add("ny", fmt.Sprintf("%d", ny))
	params.Add("dataType", "JSON")

	// Create the request URL
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Create a new HTTP request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}
	fmt.Println("Response Body:", string(body))

	// Parse the response JSON
	var weatherResponse map[string]interface{}
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response JSON: %v", err)
	}

	items := weatherResponse["response"].(map[string]interface{})["body"].(map[string]interface{})["items"].(map[string]interface{})["item"].([]interface{})
	if len(items) == 0 {
		return "", fmt.Errorf("no weather data found")
	}

	// Process and format the forecast data
	forecastMap := make(map[string]map[string]string)

	for _, item := range items {
		it := item.(map[string]interface{})
		category := it["category"].(string)
		fcstValue := it["fcstValue"].(string)
		fcstTime := it["fcstTime"].(string)

		if _, exists := forecastMap[fcstTime]; !exists {
			forecastMap[fcstTime] = make(map[string]string)
		}
		forecastMap[fcstTime][category] = fcstValue
	}

	times := make([]string, 0, len(forecastMap))
	for time := range forecastMap {
		times = append(times, time)
	}

	sort.Strings(times)

	var forecast string
	for i, time := range times {
		values := forecastMap[time]
		temp := values["T1H"]
		humidity := values["REH"]
		precipitation := values["RN1"]
		windSpeed := values["WSD"]
		skyState := skyMap[values["SKY"]]
		precipitationType := ptyMap[values["PTY"]]
		formattedTime := formatTime(time)
		precipitationStr := "강수없음"
		if precipitation != "강수없음" && precipitation != "-" && precipitation != "0" && precipitation != "null" {
			precipitationStr = precipitation + "mm"
		}

		forecast += fmt.Sprintf("%s 날씨\n하늘상태: %s\n강수형태: %s\n기온: %s℃\n습도: %s%%\n강수량: %s\n풍속: %sm/s",
			formattedTime, skyState, precipitationType, temp, humidity, precipitationStr, windSpeed)
		if i < len(times)-1 {
			forecast += "\n\n"
		}
	}

	return forecast, nil
}

func formatPrecipitation(value string) string {
	switch value {
	case "강수없음", "-", "null", "0":
		return "강수 없음"
	case "PCP = 6.2":
		return "강수량: 6.2mm"
	case "PCP = 30":
		return "강수량: 30.0 ~ 50.0mm"
	default:
		return value
	}
}

func formatTime(timeStr string) string {
	hour := timeStr[:2]
	return fmt.Sprintf("%s:00", hour)
}
