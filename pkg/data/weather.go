package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type WeatherRequestConfig struct {
	Key string
}

const (
	weatherForecastURL = "https://api.weatherapi.com/v1/forecast.json?key=%v&q=%v&days=%v&aqi=no&alerts=no"
)

var (
	WeatherClientConfig = &WeatherRequestConfig{}
)

func makeWeatherAPIRequest(url string) []byte {
	log.Printf("Making API Requst to %v", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
		return []byte("")
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return []byte("")
	}

	defer res.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(res.Body)
	bodyString := string(bodyBytes)
	if 200 > res.StatusCode && res.StatusCode > 299 {
		log.Fatalln("Request failed:", res.StatusCode, bodyString)
		return []byte("")
	}

	return bodyBytes

}

type WeatherForecastResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Date      string `json:"date"`
			DateEpoch int    `json:"date_epoch"`
			Day       struct {
				MaxtempC          float64 `json:"maxtemp_c"`
				MaxtempF          float64 `json:"maxtemp_f"`
				MintempC          float64 `json:"mintemp_c"`
				MintempF          float64 `json:"mintemp_f"`
				AvgtempC          float64 `json:"avgtemp_c"`
				AvgtempF          float64 `json:"avgtemp_f"`
				MaxwindMph        float64 `json:"maxwind_mph"`
				MaxwindKph        float64 `json:"maxwind_kph"`
				TotalprecipMm     float64 `json:"totalprecip_mm"`
				TotalprecipIn     float64 `json:"totalprecip_in"`
				TotalsnowCm       float64 `json:"totalsnow_cm"`
				AvgvisKm          float64 `json:"avgvis_km"`
				AvgvisMiles       float64 `json:"avgvis_miles"`
				Avghumidity       float64 `json:"avghumidity"`
				DailyWillItRain   int     `json:"daily_will_it_rain"`
				DailyChanceOfRain int     `json:"daily_chance_of_rain"`
				DailyWillItSnow   int     `json:"daily_will_it_snow"`
				DailyChanceOfSnow int     `json:"daily_chance_of_snow"`
				Condition         struct {
					Text string `json:"text"`
					Icon string `json:"icon"`
					Code int    `json:"code"`
				} `json:"condition"`
				Uv float64 `json:"uv"`
			} `json:"day"`
			Astro struct {
				Sunrise          string `json:"sunrise"`
				Sunset           string `json:"sunset"`
				Moonrise         string `json:"moonrise"`
				Moonset          string `json:"moonset"`
				MoonPhase        string `json:"moon_phase"`
				MoonIllumination string `json:"moon_illumination"`
				IsMoonUp         int    `json:"is_moon_up"`
				IsSunUp          int    `json:"is_sun_up"`
			} `json:"astro"`
			Hour []struct {
				TimeEpoch int     `json:"time_epoch"`
				Time      string  `json:"time"`
				TempC     float64 `json:"temp_c"`
				TempF     float64 `json:"temp_f"`
				IsDay     int     `json:"is_day"`
				Condition struct {
					Text string `json:"text"`
					Icon string `json:"icon"`
					Code int    `json:"code"`
				} `json:"condition"`
				WindMph      float64 `json:"wind_mph"`
				WindKph      float64 `json:"wind_kph"`
				WindDegree   int     `json:"wind_degree"`
				WindDir      string  `json:"wind_dir"`
				PressureMb   float64 `json:"pressure_mb"`
				PressureIn   float64 `json:"pressure_in"`
				PrecipMm     float64 `json:"precip_mm"`
				PrecipIn     float64 `json:"precip_in"`
				Humidity     int     `json:"humidity"`
				Cloud        int     `json:"cloud"`
				FeelslikeC   float64 `json:"feelslike_c"`
				FeelslikeF   float64 `json:"feelslike_f"`
				WindchillC   float64 `json:"windchill_c"`
				WindchillF   float64 `json:"windchill_f"`
				HeatindexC   float64 `json:"heatindex_c"`
				HeatindexF   float64 `json:"heatindex_f"`
				DewpointC    float64 `json:"dewpoint_c"`
				DewpointF    float64 `json:"dewpoint_f"`
				WillItRain   int     `json:"will_it_rain"`
				ChanceOfRain int     `json:"chance_of_rain"`
				WillItSnow   int     `json:"will_it_snow"`
				ChanceOfSnow int     `json:"chance_of_snow"`
				VisKm        float64 `json:"vis_km"`
				VisMiles     float64 `json:"vis_miles"`
				GustMph      float64 `json:"gust_mph"`
				GustKph      float64 `json:"gust_kph"`
				Uv           float64 `json:"uv"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func FetchWeatherForecast(location string, days string) WeatherForecastResponse {
	url := fmt.Sprintf(weatherForecastURL, WeatherClientConfig.Key, location, days)
	fmt.Sprintf(url)
	// res := makeWeatherAPIRequest(url)
	var respStruct WeatherForecastResponse
	// json.Unmarshal(res, &respStruct)

	j := `{"location":{"name":"Hartford","region":"Connecticut","country":"USA","lat":41.77,"lon":-72.7,"tz_id":"America/New_York","localtime_epoch":1677969333,"localtime":"2023-03-04 17:35"},"current":{"last_updated_epoch":1677969000,"last_updated":"2023-03-04 17:30","temp_c":4.4,"temp_f":39.9,"is_day":1,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/day/122.png","code":1009},"wind_mph":6.9,"wind_kph":11.2,"wind_degree":350,"wind_dir":"N","pressure_mb":1009.0,"pressure_in":29.78,"precip_mm":0.0,"precip_in":0.0,"humidity":70,"cloud":100,"feelslike_c":0.9,"feelslike_f":33.7,"vis_km":16.0,"vis_miles":9.0,"uv":1.0,"gust_mph":11.2,"gust_kph":18.0},"forecast":{"forecastday":[{"date":"2023-03-04","date_epoch":1677888000,"day":{"maxtemp_c":2.4,"maxtemp_f":36.3,"mintemp_c":0.2,"mintemp_f":32.4,"avgtemp_c":1.3,"avgtemp_f":34.4,"maxwind_mph":16.3,"maxwind_kph":26.3,"totalprecip_mm":19.7,"totalprecip_in":0.78,"totalsnow_cm":12.5,"avgvis_km":5.3,"avgvis_miles":3.0,"avghumidity":91.0,"daily_will_it_rain":1,"daily_chance_of_rain":83,"daily_will_it_snow":0,"daily_chance_of_snow":3,"condition":{"text":"Moderate rain","icon":"//cdn.weatherapi.com/weather/64x64/day/302.png","code":1189},"uv":1.0},"astro":{"sunrise":"06:21 AM","sunset":"05:45 PM","moonrise":"02:53 PM","moonset":"05:24 AM","moon_phase":"Waxing Gibbous","moon_illumination":"89","is_moon_up":1,"is_sun_up":0},"hour":[{"time_epoch":1677906000,"time":"2023-03-04 00:00","temp_c":1.8,"temp_f":35.2,"is_day":0,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/night/122.png","code":1009},"wind_mph":10.7,"wind_kph":17.3,"wind_degree":109,"wind_dir":"ESE","pressure_mb":1010.0,"pressure_in":29.83,"precip_mm":0.0,"precip_in":0.0,"humidity":87,"cloud":100,"feelslike_c":-2.4,"feelslike_f":27.7,"windchill_c":-2.4,"windchill_f":27.7,"heatindex_c":1.8,"heatindex_f":35.2,"dewpoint_c":-0.2,"dewpoint_f":31.6,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":10.0,"vis_miles":6.0,"gust_mph":14.5,"gust_kph":23.4,"uv":1.0},{"time_epoch":1677909600,"time":"2023-03-04 01:00","temp_c":1.0,"temp_f":33.8,"is_day":0,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/night/122.png","code":1009},"wind_mph":10.7,"wind_kph":17.3,"wind_degree":103,"wind_dir":"ESE","pressure_mb":1008.0,"pressure_in":29.77,"precip_mm":0.0,"precip_in":0.0,"humidity":96,"cloud":100,"feelslike_c":-3.4,"feelslike_f":25.9,"windchill_c":-3.4,"windchill_f":25.9,"heatindex_c":1.0,"heatindex_f":33.8,"dewpoint_c":0.4,"dewpoint_f":32.7,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":2.0,"vis_miles":1.0,"gust_mph":14.8,"gust_kph":23.8,"uv":1.0},{"time_epoch":1677913200,"time":"2023-03-04 02:00","temp_c":0.8,"temp_f":33.4,"is_day":0,"condition":{"text":"Heavy snow","icon":"//cdn.weatherapi.com/weather/64x64/night/338.png","code":1225},"wind_mph":12.8,"wind_kph":20.5,"wind_degree":94,"wind_dir":"E","pressure_mb":1007.0,"pressure_in":29.74,"precip_mm":4.6,"precip_in":0.18,"humidity":97,"cloud":100,"feelslike_c":-3.8,"feelslike_f":25.2,"windchill_c":-3.8,"windchill_f":25.2,"heatindex_c":0.8,"heatindex_f":33.4,"dewpoint_c":0.4,"dewpoint_f":32.7,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":61,"vis_km":10.0,"vis_miles":6.0,"gust_mph":18.1,"gust_kph":29.2,"uv":1.0},{"time_epoch":1677916800,"time":"2023-03-04 03:00","temp_c":0.9,"temp_f":33.6,"is_day":0,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/night/122.png","code":1009},"wind_mph":16.3,"wind_kph":26.3,"wind_degree":97,"wind_dir":"E","pressure_mb":1004.0,"pressure_in":29.66,"precip_mm":0.0,"precip_in":0.0,"humidity":97,"cloud":100,"feelslike_c":-4.5,"feelslike_f":23.9,"windchill_c":-4.5,"windchill_f":23.9,"heatindex_c":0.9,"heatindex_f":33.6,"dewpoint_c":0.4,"dewpoint_f":32.7,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":2.0,"vis_miles":1.0,"gust_mph":23.0,"gust_kph":37.1,"uv":1.0},{"time_epoch":1677920400,"time":"2023-03-04 04:00","temp_c":1.1,"temp_f":34.0,"is_day":0,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/night/122.png","code":1009},"wind_mph":14.8,"wind_kph":23.8,"wind_degree":95,"wind_dir":"E","pressure_mb":1000.0,"pressure_in":29.52,"precip_mm":0.0,"precip_in":0.0,"humidity":93,"cloud":100,"feelslike_c":-5.1,"feelslike_f":22.8,"windchill_c":-5.1,"windchill_f":22.8,"heatindex_c":1.1,"heatindex_f":34.0,"dewpoint_c":0.2,"dewpoint_f":32.4,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":10.0,"vis_miles":6.0,"gust_mph":21.0,"gust_kph":33.8,"uv":1.0},{"time_epoch":1677924000,"time":"2023-03-04 05:00","temp_c":1.1,"temp_f":34.0,"is_day":0,"condition":{"text":"Heavy snow","icon":"//cdn.weatherapi.com/weather/64x64/night/338.png","code":1225},"wind_mph":6.9,"wind_kph":11.2,"wind_degree":55,"wind_dir":"NE","pressure_mb":998.0,"pressure_in":29.47,"precip_mm":8.4,"precip_in":0.33,"humidity":96,"cloud":100,"feelslike_c":-4.6,"feelslike_f":23.7,"windchill_c":-4.6,"windchill_f":23.7,"heatindex_c":1.1,"heatindex_f":34.0,"dewpoint_c":0.5,"dewpoint_f":32.9,"will_it_rain":0,"chance_of_rain":67,"will_it_snow":0,"chance_of_snow":0,"vis_km":2.0,"vis_miles":1.0,"gust_mph":10.3,"gust_kph":16.6,"uv":1.0},{"time_epoch":1677927600,"time":"2023-03-04 06:00","temp_c":1.0,"temp_f":33.8,"is_day":0,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/night/122.png","code":1009},"wind_mph":11.6,"wind_kph":18.7,"wind_degree":67,"wind_dir":"ENE","pressure_mb":997.0,"pressure_in":29.44,"precip_mm":0.0,"precip_in":0.0,"humidity":96,"cloud":100,"feelslike_c":-4.0,"feelslike_f":24.8,"windchill_c":-4.0,"windchill_f":24.8,"heatindex_c":1.0,"heatindex_f":33.8,"dewpoint_c":0.5,"dewpoint_f":32.9,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":10.0,"vis_miles":6.0,"gust_mph":17.2,"gust_kph":27.7,"uv":1.0},{"time_epoch":1677931200,"time":"2023-03-04 07:00","temp_c":1.2,"temp_f":34.2,"is_day":1,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/day/122.png","code":1009},"wind_mph":11.6,"wind_kph":18.7,"wind_degree":59,"wind_dir":"ENE","pressure_mb":995.0,"pressure_in":29.39,"precip_mm":0.0,"precip_in":0.0,"humidity":95,"cloud":100,"feelslike_c":-4.2,"feelslike_f":24.4,"windchill_c":-4.2,"windchill_f":24.4,"heatindex_c":1.2,"heatindex_f":34.2,"dewpoint_c":0.5,"dewpoint_f":32.9,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":2.0,"vis_miles":1.0,"gust_mph":16.8,"gust_kph":27.0,"uv":1.0},{"time_epoch":1677934800,"time":"2023-03-04 08:00","temp_c":1.2,"temp_f":34.2,"is_day":1,"condition":{"text":"Moderate or heavy sleet","icon":"//cdn.weatherapi.com/weather/64x64/day/320.png","code":1207},"wind_mph":11.6,"wind_kph":18.7,"wind_degree":51,"wind_dir":"NE","pressure_mb":993.0,"pressure_in":29.33,"precip_mm":5.0,"precip_in":0.2,"humidity":95,"cloud":100,"feelslike_c":-4.1,"feelslike_f":24.6,"windchill_c":-4.1,"windchill_f":24.6,"heatindex_c":1.2,"heatindex_f":34.2,"dewpoint_c":0.4,"dewpoint_f":32.7,"will_it_rain":0,"chance_of_rain":70,"will_it_snow":0,"chance_of_snow":0,"vis_km":10.0,"vis_miles":6.0,"gust_mph":16.6,"gust_kph":26.6,"uv":1.0},{"time_epoch":1677938400,"time":"2023-03-04 09:00","temp_c":1.0,"temp_f":33.8,"is_day":1,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/day/122.png","code":1009},"wind_mph":11.4,"wind_kph":18.4,"wind_degree":34,"wind_dir":"NE","pressure_mb":993.0,"pressure_in":29.32,"precip_mm":0.0,"precip_in":0.0,"humidity":96,"cloud":100,"feelslike_c":-3.8,"feelslike_f":25.2,"windchill_c":-3.8,"windchill_f":25.2,"heatindex_c":1.0,"heatindex_f":33.8,"dewpoint_c":0.5,"dewpoint_f":32.9,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":2.0,"vis_miles":1.0,"gust_mph":16.1,"gust_kph":25.9,"uv":1.0},{"time_epoch":1677942000,"time":"2023-03-04 10:00","temp_c":1.1,"temp_f":34.0,"is_day":1,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/day/122.png","code":1009},"wind_mph":11.2,"wind_kph":18.0,"wind_degree":25,"wind_dir":"NNE","pressure_mb":994.0,"pressure_in":29.34,"precip_mm":0.0,"precip_in":0.0,"humidity":95,"cloud":100,"feelslike_c":-3.5,"feelslike_f":25.7,"windchill_c":-3.5,"windchill_f":25.7,"heatindex_c":1.1,"heatindex_f":34.0,"dewpoint_c":0.5,"dewpoint_f":32.9,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":2.0,"vis_miles":1.0,"gust_mph":15.7,"gust_kph":25.2,"uv":1.0},{"time_epoch":1677945600,"time":"2023-03-04 11:00","temp_c":1.4,"temp_f":34.5,"is_day":1,"condition":{"text":"Moderate snow","icon":"//cdn.weatherapi.com/weather/64x64/day/332.png","code":1219},"wind_mph":11.6,"wind_kph":18.7,"wind_degree":23,"wind_dir":"NNE","pressure_mb":995.0,"pressure_in":29.37,"precip_mm":1.7,"precip_in":0.07,"humidity":94,"cloud":100,"feelslike_c":-3.2,"feelslike_f":26.2,"windchill_c":-3.2,"windchill_f":26.2,"heatindex_c":1.4,"heatindex_f":34.5,"dewpoint_c":0.6,"dewpoint_f":33.1,"will_it_rain":1,"chance_of_rain":83,"will_it_snow":0,"chance_of_snow":0,"vis_km":2.0,"vis_miles":1.0,"gust_mph":15.4,"gust_kph":24.8,"uv":1.0},{"time_epoch":1677949200,"time":"2023-03-04 12:00","temp_c":1.6,"temp_f":34.9,"is_day":1,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/day/122.png","code":1009},"wind_mph":11.4,"wind_kph":18.4,"wind_degree":7,"wind_dir":"N","pressure_mb":996.0,"pressure_in":29.41,"precip_mm":0.0,"precip_in":0.0,"humidity":93,"cloud":100,"feelslike_c":-3.2,"feelslike_f":26.2,"windchill_c":-3.2,"windchill_f":26.2,"heatindex_c":1.6,"heatindex_f":34.9,"dewpoint_c":0.6,"dewpoint_f":33.1,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":5.0,"vis_miles":3.0,"gust_mph":14.8,"gust_kph":23.8,"uv":1.0},{"time_epoch":1677952800,"time":"2023-03-04 13:00","temp_c":1.6,"temp_f":34.9,"is_day":1,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/day/122.png","code":1009},"wind_mph":11.4,"wind_kph":18.4,"wind_degree":5,"wind_dir":"N","pressure_mb":997.0,"pressure_in":29.44,"precip_mm":0.0,"precip_in":0.0,"humidity":92,"cloud":100,"feelslike_c":-3.4,"feelslike_f":25.9,"windchill_c":-3.4,"windchill_f":25.9,"heatindex_c":1.6,"heatindex_f":34.9,"dewpoint_c":0.4,"dewpoint_f":32.7,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":2.0,"vis_miles":1.0,"gust_mph":14.8,"gust_kph":23.8,"uv":1.0},{"time_epoch":1677956400,"time":"2023-03-04 14:00","temp_c":2.0,"temp_f":35.6,"is_day":1,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/day/122.png","code":1009},"wind_mph":11.2,"wind_kph":18.0,"wind_degree":0,"wind_dir":"N","pressure_mb":999.0,"pressure_in":29.5,"precip_mm":0.0,"precip_in":0.0,"humidity":90,"cloud":100,"feelslike_c":-2.9,"feelslike_f":26.8,"windchill_c":-2.9,"windchill_f":26.8,"heatindex_c":2.0,"heatindex_f":35.6,"dewpoint_c":0.5,"dewpoint_f":32.9,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":2.0,"vis_miles":1.0,"gust_mph":13.9,"gust_kph":22.3,"uv":1.0},{"time_epoch":1677960000,"time":"2023-03-04 15:00","temp_c":2.3,"temp_f":36.1,"is_day":1,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/day/122.png","code":1009},"wind_mph":11.6,"wind_kph":18.7,"wind_degree":8,"wind_dir":"N","pressure_mb":1001.0,"pressure_in":29.57,"precip_mm":0.0,"precip_in":0.0,"humidity":88,"cloud":100,"feelslike_c":-2.6,"feelslike_f":27.3,"windchill_c":-2.6,"windchill_f":27.3,"heatindex_c":2.3,"heatindex_f":36.1,"dewpoint_c":0.5,"dewpoint_f":32.9,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":5.0,"vis_miles":3.0,"gust_mph":14.1,"gust_kph":22.7,"uv":1.0},{"time_epoch":1677963600,"time":"2023-03-04 16:00","temp_c":2.4,"temp_f":36.3,"is_day":1,"condition":{"text":"Overcast","icon":"//cdn.weatherapi.com/weather/64x64/day/122.png","code":1009},"wind_mph":11.4,"wind_kph":18.4,"wind_degree":13,"wind_dir":"NNE","pressure_mb":1003.0,"pressure_in":29.63,"precip_mm":0.0,"precip_in":0.0,"humidity":87,"cloud":99,"feelslike_c":-2.4,"feelslike_f":27.7,"windchill_c":-2.4,"windchill_f":27.7,"heatindex_c":2.4,"heatindex_f":36.3,"dewpoint_c":0.4,"dewpoint_f":32.7,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":7.0,"vis_miles":4.0,"gust_mph":13.4,"gust_kph":21.6,"uv":1.0},{"time_epoch":1677967200,"time":"2023-03-04 17:00","temp_c":2.3,"temp_f":36.1,"is_day":1,"condition":{"text":"Cloudy","icon":"//cdn.weatherapi.com/weather/64x64/day/119.png","code":1006},"wind_mph":9.6,"wind_kph":15.5,"wind_degree":11,"wind_dir":"NNE","pressure_mb":1005.0,"pressure_in":29.69,"precip_mm":0.0,"precip_in":0.0,"humidity":86,"cloud":85,"feelslike_c":-2.3,"feelslike_f":27.9,"windchill_c":-2.3,"windchill_f":27.9,"heatindex_c":2.3,"heatindex_f":36.1,"dewpoint_c":0.2,"dewpoint_f":32.4,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":5.0,"vis_miles":3.0,"gust_mph":11.2,"gust_kph":18.0,"uv":1.0},{"time_epoch":1677970800,"time":"2023-03-04 18:00","temp_c":2.2,"temp_f":36.0,"is_day":0,"condition":{"text":"Partly cloudy","icon":"//cdn.weatherapi.com/weather/64x64/night/116.png","code":1003},"wind_mph":7.8,"wind_kph":12.6,"wind_degree":358,"wind_dir":"N","pressure_mb":1008.0,"pressure_in":29.75,"precip_mm":0.0,"precip_in":0.0,"humidity":86,"cloud":62,"feelslike_c":-2.2,"feelslike_f":28.0,"windchill_c":-2.2,"windchill_f":28.0,"heatindex_c":2.2,"heatindex_f":36.0,"dewpoint_c":0.0,"dewpoint_f":32.0,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":7.0,"vis_miles":4.0,"gust_mph":11.0,"gust_kph":17.6,"uv":1.0},{"time_epoch":1677974400,"time":"2023-03-04 19:00","temp_c":2.3,"temp_f":36.1,"is_day":0,"condition":{"text":"Clear","icon":"//cdn.weatherapi.com/weather/64x64/night/113.png","code":1000},"wind_mph":7.8,"wind_kph":12.6,"wind_degree":359,"wind_dir":"N","pressure_mb":1009.0,"pressure_in":29.81,"precip_mm":0.0,"precip_in":0.0,"humidity":85,"cloud":22,"feelslike_c":-1.8,"feelslike_f":28.8,"windchill_c":-1.8,"windchill_f":28.8,"heatindex_c":2.3,"heatindex_f":36.1,"dewpoint_c":0.1,"dewpoint_f":32.2,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":10.0,"vis_miles":6.0,"gust_mph":11.2,"gust_kph":18.0,"uv":1.0},{"time_epoch":1677978000,"time":"2023-03-04 20:00","temp_c":0.5,"temp_f":32.9,"is_day":0,"condition":{"text":"Partly cloudy","icon":"//cdn.weatherapi.com/weather/64x64/night/116.png","code":1003},"wind_mph":7.2,"wind_kph":11.5,"wind_degree":6,"wind_dir":"N","pressure_mb":1011.0,"pressure_in":29.85,"precip_mm":0.0,"precip_in":0.0,"humidity":86,"cloud":28,"feelslike_c":-3.6,"feelslike_f":25.5,"windchill_c":-3.6,"windchill_f":25.5,"heatindex_c":0.5,"heatindex_f":32.9,"dewpoint_c":-1.6,"dewpoint_f":29.1,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":7.0,"vis_miles":4.0,"gust_mph":10.5,"gust_kph":16.9,"uv":1.0},{"time_epoch":1677981600,"time":"2023-03-04 21:00","temp_c":0.4,"temp_f":32.7,"is_day":0,"condition":{"text":"Partly cloudy","icon":"//cdn.weatherapi.com/weather/64x64/night/116.png","code":1003},"wind_mph":4.5,"wind_kph":7.2,"wind_degree":357,"wind_dir":"N","pressure_mb":1012.0,"pressure_in":29.88,"precip_mm":0.0,"precip_in":0.0,"humidity":86,"cloud":38,"feelslike_c":-3.5,"feelslike_f":25.7,"windchill_c":-3.5,"windchill_f":25.7,"heatindex_c":0.4,"heatindex_f":32.7,"dewpoint_c":-1.7,"dewpoint_f":28.9,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":10.0,"vis_miles":6.0,"gust_mph":7.6,"gust_kph":12.2,"uv":1.0},{"time_epoch":1677985200,"time":"2023-03-04 22:00","temp_c":0.3,"temp_f":32.5,"is_day":0,"condition":{"text":"Partly cloudy","icon":"//cdn.weatherapi.com/weather/64x64/night/116.png","code":1003},"wind_mph":3.1,"wind_kph":5.0,"wind_degree":352,"wind_dir":"N","pressure_mb":1013.0,"pressure_in":29.92,"precip_mm":0.0,"precip_in":0.0,"humidity":86,"cloud":28,"feelslike_c":-3.2,"feelslike_f":26.2,"windchill_c":-3.2,"windchill_f":26.2,"heatindex_c":0.3,"heatindex_f":32.5,"dewpoint_c":-1.8,"dewpoint_f":28.8,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":2.0,"vis_miles":1.0,"gust_mph":5.8,"gust_kph":9.4,"uv":1.0},{"time_epoch":1677988800,"time":"2023-03-04 23:00","temp_c":0.2,"temp_f":32.4,"is_day":0,"condition":{"text":"Clear","icon":"//cdn.weatherapi.com/weather/64x64/night/113.png","code":1000},"wind_mph":2.5,"wind_kph":4.0,"wind_degree":311,"wind_dir":"NW","pressure_mb":1015.0,"pressure_in":29.96,"precip_mm":0.0,"precip_in":0.0,"humidity":85,"cloud":6,"feelslike_c":-3.1,"feelslike_f":26.4,"windchill_c":-3.1,"windchill_f":26.4,"heatindex_c":0.2,"heatindex_f":32.4,"dewpoint_c":-2.0,"dewpoint_f":28.4,"will_it_rain":0,"chance_of_rain":0,"will_it_snow":0,"chance_of_snow":0,"vis_km":2.0,"vis_miles":1.0,"gust_mph":4.7,"gust_kph":7.6,"uv":1.0}]}]}}`
	json.Unmarshal([]byte(j), &respStruct)

	return respStruct
}
