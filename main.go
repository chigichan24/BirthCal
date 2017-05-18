package BirthCal

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"regexp"
	"strconv"
)

type User struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Month int    `json:"month"`
	Day   int    `json:"day"`
}

func getBirthDay(lists ...User) []User {
	var result []User
	for _, data := range lists {
		doc, _ := goquery.NewDocument(data.Url)
		doc.Find(".ProfileHeaderCard-birthdate").Each(func(_ int, s *goquery.Selection) {
			date := s.Find(".js-tooltip").Text()
			r := regexp.MustCompile("[0-9][0-9]*")
			t := r.FindAllStringSubmatch(date, -1)

			if len(t) > 0 {
				data.Month, _ = strconv.Atoi(t[0][0])
				data.Day, _ = strconv.Atoi(t[1][0])
				fmt.Println(data.Name, " ", data.Month, " ", data.Day)
				result = append(result, data)
			}
		})
	}
	return result
}

func getFollowersURL(api *anaconda.TwitterApi, v url.Values) []User {
	flists, err := api.GetFollowersList(v)
	if err != nil {
		panic(err)
	}
	var lists []User
	for _, user := range flists.Users {
		t := "https://twitter.com/" + user.ScreenName
		lists = append(lists, User{user.ScreenName, t, -1, -1})
	}
	return lists
}

func outputJSON(list []User) string {

	result, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}
	s := string(result)
	return s
}

func main() {
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api := anaconda.NewTwitterApi(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)

	v := url.Values{}
	v.Set("chigichan24", "")
	var users = getBirthDay(getFollowersURL(api, v)...)

	for _, tmp := range users {
		fmt.Println(tmp)
	}

	fmt.Printf(outputJSON(users))
}
