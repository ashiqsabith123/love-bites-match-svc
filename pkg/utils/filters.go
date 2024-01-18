package utils

import (
	"fmt"
	"math"
	"strconv"

	authPb "github.com/ashiqsabith123/love-bytes-proto/auth/pb"
)

func (F *Utils) FilterWithDistance(user *authPb.UserRepsonse, matchUsersList []authPb.UserRepsonse) {

	lat1, _ := strconv.ParseFloat(user.Lat, 64)
	lng1, _ := strconv.ParseFloat(user.Log, 64)

	for i := 0; i < len(matchUsersList)-1; i++ {
		lat2, _ := strconv.ParseFloat(matchUsersList[i].Lat, 64)
		lng2, _ := strconv.ParseFloat(matchUsersList[i].Log, 64)
		if dist := distance(lat1, lng1, lat2, lng2); dist > 50 {
			matchUsersList = append(matchUsersList[:i], matchUsersList[i+1:]...)

		}
	}

	for _, v := range matchUsersList {
		fmt.Println(v.Fullname)
	}

}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	dist = dist * 1.609344

	return dist
}
