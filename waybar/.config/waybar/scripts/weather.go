package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var weatherCodes = map[string]string{
	"113": "☀️",
	"116": "⛅️",
	"119": "☁️",
	"122": "☁️",
	"143": "🌫",
	"176": "🌦",
	"179": "🌧",
	"182": "🌧",
	"185": "🌧",
	"200": "⛈",
	"227": "🌨",
	"230": "❄️",
	"248": "🌫",
	"260": "🌫",
	"263": "🌦",
	"266": "🌦",
	"281": "🌧",
	"284": "🌧",
	"293": "🌦",
	"296": "🌦",
	"299": "🌧",
	"302": "🌧",
	"305": "🌧",
	"308": "🌧",
	"311": "🌧",
	"314": "🌧",
	"317": "🌧",
	"320": "🌨",
	"323": "🌨",
	"326": "🌨",
	"329": "❄️",
	"332": "❄️",
	"335": "❄️",
	"338": "❄️",
	"350": "🌧",
	"353": "🌦",
	"356": "🌧",
	"359": "🌧",
	"362": "🌧",
	"365": "🌧",
	"368": "🌨",
	"371": "❄️",
	"374": "🌧",
	"377": "🌧",
	"386": "⛈",
	"389": "🌩",
	"392": "⛈",
	"395": "❄️",
}

type WeatherData struct {
	CurrentCondition []struct {
		WeatherCode     string `json:"weatherCode"`
		FeelsLikeC      string `json:"FeelsLikeC"`
		TempC           string `json:"temp_C"`
		WindspeedMiles  string `json:"windspeedMiles"`
		Humidity        string `json:"humidity"`
		WeatherDesc     []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
	} `json:"current_condition"`
	Weather []struct {
		Date     string `json:"date"`
		MaxtempC string `json:"maxtempC"`
		MintempC string `json:"mintempC"`
		Astronomy []struct {
			Sunrise string `json:"sunrise"`
			Sunset  string `json:"sunset"`
		} `json:"astronomy"`
		Hourly []struct {
			Time               string `json:"time"`
			WeatherCode        string `json:"weatherCode"`
			FeelsLikeC         string `json:"FeelsLikeC"`
			ChanceOfFog        string `json:"chanceoffog"`
			ChanceOfFrost      string `json:"chanceoffrost"`
			ChanceOfOvercast   string `json:"chanceofovercast"`
			ChanceOfRain       string `json:"chanceofrain"`
			ChanceOfSnow       string `json:"chanceofsnow"`
			ChanceOfSunshine   string `json:"chanceofsunshine"`
			ChanceOfThunder    string `json:"chanceofthunder"`
			ChanceOfWindy      string `json:"chanceofwindy"`
			WeatherDesc        []struct {
				Value string `json:"value"`
			} `json:"weatherDesc"`
		} `json:"hourly"`
	} `json:"weather"`
}

type OutputData struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

func formatTime(timeStr string) string {
	cleaned := strings.Replace(timeStr, "00", "", -1)
	if len(cleaned) == 0 {
		return "24"
	}
	if len(cleaned) == 1 {
		return "0" + cleaned
	}
	return cleaned
}

func formatTemp(temp string) string {
	return fmt.Sprintf("%-3s", temp+"°")
}

func formatChances(hour struct {
	ChanceOfFog      string `json:"chanceoffog"`
	ChanceOfFrost    string `json:"chanceoffrost"`
	ChanceOfOvercast string `json:"chanceofovercast"`
	ChanceOfRain     string `json:"chanceofrain"`
	ChanceOfSnow     string `json:"chanceofsnow"`
	ChanceOfSunshine string `json:"chanceofsunshine"`
	ChanceOfThunder  string `json:"chanceofthunder"`
	ChanceOfWindy    string `json:"chanceofwindy"`
}) string {
	chanceValues := []string{
		hour.ChanceOfFog, hour.ChanceOfFrost, hour.ChanceOfOvercast,
		hour.ChanceOfRain, hour.ChanceOfSnow, hour.ChanceOfSunshine,
		hour.ChanceOfThunder, hour.ChanceOfWindy,
	}

	chanceNames := []string{
		"Fog", "Frost", "Overcast", "Rain", "Snow", "Sunshine", "Thunder", "Wind",
	}

	var conditions []string
	for i, chanceStr := range chanceValues {
		if chance, err := strconv.Atoi(chanceStr); err == nil && chance > 0 {
			conditions = append(conditions, fmt.Sprintf("%s %s%%", chanceNames[i], chanceStr))
		}
	}

	return strings.Join(conditions, ", ")
}

func main() {
	// Fetch weather data
	resp, err := http.Get("https://wttr.in/Amager?format=j1")
	if err != nil {
		fmt.Printf("Error fetching weather data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var weather WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		fmt.Printf("Error decoding weather data: %v\n", err)
		return
	}

	var data OutputData

	// Set text
	currentCondition := weather.CurrentCondition[0]
	data.Text = weatherCodes[currentCondition.WeatherCode] + " " + currentCondition.FeelsLikeC + "°C"

	// Build tooltip
	data.Tooltip = fmt.Sprintf("<b>%s %s°C</b>\n",
		currentCondition.WeatherDesc[0].Value,
		currentCondition.TempC)
	data.Tooltip += fmt.Sprintf("Feels like: %s°C\n", currentCondition.FeelsLikeC)
	data.Tooltip += fmt.Sprintf("Wind: %smi/h\n", currentCondition.WindspeedMiles)
	data.Tooltip += fmt.Sprintf("Humidity: %s%%\n", currentCondition.Humidity)

	now := time.Now()

	for i, day := range weather.Weather {
		data.Tooltip += "\n<b>"
		if i == 0 {
			data.Tooltip += "Today, "
		} else if i == 1 {
			data.Tooltip += "Tomorrow, "
		}
		data.Tooltip += fmt.Sprintf("%s</b>\n", day.Date)
		data.Tooltip += fmt.Sprintf("⬆️ %s° ⬇️ %s° ", day.MaxtempC, day.MintempC)
		data.Tooltip += fmt.Sprintf("🌅 %s 🌇 %s\n", day.Astronomy[0].Sunrise, day.Astronomy[0].Sunset)

		for _, hour := range day.Hourly {
			// Skip past hours for today
			if i == 0 {
				hourTime, err := strconv.Atoi(formatTime(hour.Time))
				if err == nil && hourTime < now.Hour()-2 {
					continue
				}
			}

			data.Tooltip += fmt.Sprintf("%s %s %s %s, %s\n",
				formatTime(hour.Time),
				weatherCodes[hour.WeatherCode],
				formatTemp(hour.FeelsLikeC),
				hour.WeatherDesc[0].Value,
				formatChances(struct {
					ChanceOfFog      string `json:"chanceoffog"`
					ChanceOfFrost    string `json:"chanceoffrost"`
					ChanceOfOvercast string `json:"chanceofovercast"`
					ChanceOfRain     string `json:"chanceofrain"`
					ChanceOfSnow     string `json:"chanceofsnow"`
					ChanceOfSunshine string `json:"chanceofsunshine"`
					ChanceOfThunder  string `json:"chanceofthunder"`
					ChanceOfWindy    string `json:"chanceofwindy"`
				}{
					hour.ChanceOfFog, hour.ChanceOfFrost, hour.ChanceOfOvercast,
					hour.ChanceOfRain, hour.ChanceOfSnow, hour.ChanceOfSunshine,
					hour.ChanceOfThunder, hour.ChanceOfWindy,
				}))
		}
	}

	// Output JSON
	output, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling output: %v\n", err)
		return
	}

	fmt.Println(string(output))
}
