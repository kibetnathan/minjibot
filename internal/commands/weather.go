package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	weatherGeocodeAPI  = "https://geocoding-api.open-meteo.com/v1/search"
	weatherForecastAPI = "https://api.open-meteo.com/v1/forecast"
)

type weatherGeocodingResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Country   string  `json:"country"`
		Admin1    string  `json:"admin1"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type weatherForecastResponse struct {
	Current struct {
		Temperature      float64 `json:"temperature_2m"`
		RelativeHumidity int     `json:"relative_humidity_2m"`
		WeatherCode      int     `json:"weather_code"`
		WindSpeed        float64 `json:"wind_speed_10m"`
	} `json:"current"`
	Daily struct {
		TemperatureMax []float64 `json:"temperature_2m_max"`
		TemperatureMin []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

func weatherMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-weather <location>`")
		return err
	}
	embed, err := buildWeatherEmbed(strings.Join(args, " "))
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func weatherSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	location := strings.TrimSpace(OptString(opts, "location"))
	if location == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/weather location:<location>`"},
		})
	}
	embed, err := buildWeatherEmbed(location)
	if err != nil {
		return err
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

func buildWeatherEmbed(location string) (*discordgo.MessageEmbed, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	geoURL := fmt.Sprintf("%s?name=%s&count=1&language=en&format=json", weatherGeocodeAPI, url.QueryEscape(location))
	geoResp, err := client.Get(geoURL)
	if err != nil {
		return nil, err
	}
	defer geoResp.Body.Close()

	var geo weatherGeocodingResponse
	if err := json.NewDecoder(geoResp.Body).Decode(&geo); err != nil {
		return nil, err
	}
	if len(geo.Results) == 0 {
		return &discordgo.MessageEmbed{
			Color:       0xED4245,
			Title:       "Weather",
			Description: fmt.Sprintf("No location found for `%s`.", location),
		}, nil
	}

	place := geo.Results[0]
	placeName := place.Name
	if place.Admin1 != "" && place.Admin1 != place.Name {
		placeName = fmt.Sprintf("%s, %s", place.Name, place.Admin1)
	}
	if place.Country != "" {
		placeName = fmt.Sprintf("%s, %s", placeName, place.Country)
	}

	forecastURL := fmt.Sprintf(
		"%s?latitude=%.4f&longitude=%.4f&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m&daily=temperature_2m_max,temperature_2m_min&timezone=auto&forecast_days=1",
		weatherForecastAPI, place.Latitude, place.Longitude,
	)
	fResp, err := client.Get(forecastURL)
	if err != nil {
		return nil, err
	}
	defer fResp.Body.Close()

	var fc weatherForecastResponse
	if err := json.NewDecoder(fResp.Body).Decode(&fc); err != nil {
		return nil, err
	}

	var high, low string
	if len(fc.Daily.TemperatureMax) > 0 {
		high = fmt.Sprintf("%.1f°C", fc.Daily.TemperatureMax[0])
	}
	if len(fc.Daily.TemperatureMin) > 0 {
		low = fmt.Sprintf("%.1f°C", fc.Daily.TemperatureMin[0])
	}

	emoji := weatherCodeEmoji(fc.Current.WeatherCode)
	return &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: fmt.Sprintf("%s %s", emoji, placeName),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Conditions", Value: weatherCodeLabel(fc.Current.WeatherCode), Inline: true},
			{Name: "Temperature", Value: fmt.Sprintf("%.1f°C", fc.Current.Temperature), Inline: true},
			{Name: "Humidity", Value: fmt.Sprintf("%d%%", fc.Current.RelativeHumidity), Inline: true},
			{Name: "Wind", Value: fmt.Sprintf("%.1f km/h", fc.Current.WindSpeed), Inline: true},
			{Name: "Today's High", Value: high, Inline: true},
			{Name: "Today's Low", Value: low, Inline: true},
		},
	}, nil
}

func weatherCodeLabel(code int) string {
	labels := map[int]string{
		0: "Clear sky", 1: "Mainly clear", 2: "Partly cloudy", 3: "Overcast",
		45: "Fog", 48: "Depositing rime fog",
		51: "Light drizzle", 53: "Drizzle", 55: "Dense drizzle",
		56: "Light freezing drizzle", 57: "Freezing drizzle",
		61: "Light rain", 63: "Rain", 65: "Heavy rain",
		66: "Light freezing rain", 67: "Freezing rain",
		71: "Light snow", 73: "Snow", 75: "Heavy snow",
		77: "Snow grains",
		80: "Light showers", 81: "Showers", 82: "Violent showers",
		85: "Light snow showers", 86: "Snow showers",
		95: "Thunderstorm", 96: "Thunderstorm with hail", 99: "Thunderstorm with heavy hail",
	}
	if v, ok := labels[code]; ok {
		return v
	}
	return "Unknown"
}

func weatherCodeEmoji(code int) string {
	switch {
	case code == 0:
		return "☀️"
	case code == 1 || code == 2:
		return "🌤️"
	case code == 3:
		return "☁️"
	case code >= 45 && code <= 48:
		return "🌫️"
	case code >= 51 && code <= 67:
		return "🌧️"
	case code >= 71 && code <= 77:
		return "🌨️"
	case code >= 80 && code <= 82:
		return "🌦️"
	case code >= 85 && code <= 86:
		return "🌨️"
	case code >= 95:
		return "⛈️"
	default:
		return "🌡️"
	}
}
