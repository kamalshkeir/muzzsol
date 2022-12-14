package handlers

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/kamalshkeir/korm"
	"github.com/kamalshkeir/muzzsol/middlewares"
	"github.com/kamalshkeir/muzzsol/services"
	"github.com/kamalshkeir/muzzsol/types"
	"github.com/labstack/echo/v4"
)

const (
	rayon = 6371 // radius of the earth in km
	charset = "abcdefghijklmnopqrstuvwxyz"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

var genders = [2]string{"male", "female"}



// generateRandomLocation generates a random location within a specified radius en Km arround given location (to bookmark, it took me 1 hour)
func GenerateRandomLocation(location *types.Location, radiusKm float64) types.Location {
    // Define the latitude and longitude of the center point
	latitude := 45.766688
	longitude := 4.833756
	if location != nil {
		latitude = location.Latitude
		longitude = location.Longitude
	}

	// Generate a random number seed
	rand.Seed(time.Now().UnixNano())

	// Calculate the minimum and maximum latitude and longitude values
	// that can be generated based on the radius
	minLat := latitude - ((radiusKm/4) / 111.0)
	maxLat := latitude + ((radiusKm/4) / 111.0)
	minLng := longitude - ((radiusKm/4) / (111.0 * math.Cos(latitude)))
	maxLng := longitude + ((radiusKm/4) / (111.0 * math.Cos(latitude)))

	// Generate a random latitude and longitude within the defined range
	randomLat := minLat + (rand.Float64() * (maxLat - minLat))
	randomLng := minLng + (rand.Float64() * (maxLng - minLng))
	return types.Location{
		Latitude: randomLat,
		Longitude: randomLng,
	}
}

// Formule Haversine (this one too)
// https://fr.wikipedia.org/wiki/Formule_de_haversine
func DistanceBetweenLocations(p1, p2 types.Location) float64 {
	lat1 := p1.Latitude* math.Pi / 180 // in Radians
	lon1 := p1.Longitude* math.Pi / 180
	lat2 := p2.Latitude* math.Pi / 180
	lon2 := p2.Longitude* math.Pi / 180

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	// d = 2 * rayon * arcsin(sqrt(sin2((lat2-lat1)/2))+cos(lat1)*cos(lat2)*sin2((lon2-lon1)/2))
	// arcsin = Atan2
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return rayon * c
}

// used to generate random names, passwords, and emails 
func RandomString(length int) string {
	return stringWithCharset(length, charset)
}

// used to generate random age 
func RandomNumberBetween(min, max int) int {
	return seededRand.Intn(max-min) + min
}


func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomUser() *types.User {
	u := types.User{}
	name := RandomString(5)
	u.Email = name + "@email.com"
	var err error

	u.Password, err = services.GenerateHash(name)
	if err != nil {
		fmt.Println("GenerateRandomUser error generating hash:", err)
		return &u
	}
	u.Name = name
	u.Age = uint(RandomNumberBetween(18, 60))
	u.Gender = genders[seededRand.Intn(len(genders))]
	return &u
}

func getSessionUserId(c echo.Context) int {
	if userM,ok := c.Request().Context().Value(middlewares.ContextKey(middlewares.CookieSessionName)).(map[string]any);ok {
		if v,ok := userM["id"];ok {
			switch vv := v.(type) {
			case float64:
				return int(vv)
			case float32:
				return int(vv)
			case int:
				return vv
			case uint:
				return int(vv)
			default:
				return -1
			}
		}
	}
	return -1
}

func calculateAndSortByDistance(userLocation types.Location,profiles []types.User) {
	if len(profiles) == 0 {
		return 
	}
	// get ids from profiles, and do a single query to get all locations for these profiles
	userIdsToGetLocations := []any{}
	for _,p := range profiles {
		userIdsToGetLocations = append(userIdsToGetLocations, p.Id)
	}
	ph := strings.Repeat("?,",len(userIdsToGetLocations))
	ph = ph[:len(ph)-1]
	// get locations
	locations,err := korm.Model[types.Location]().Where("user_id IN ("+ ph + ")",userIdsToGetLocations...).All()
	if err != nil {
		fmt.Println("error getting locations:",err)
		return 
	}
	// fill distance field
	for i := range profiles {
		for _,l := range locations {
			if l.UserId == profiles[i].Id {
				profiles[i].DistanceFromMe = DistanceBetweenLocations(userLocation,l)
			}
		}
	}
	// sort by distance from user
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].DistanceFromMe < profiles[j].DistanceFromMe
	})
}





