package brickset

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/go-resty/resty/v2"
)

// Marshal JSON in Golang using lower-camel-case object key conventions
// https://gist.github.com/piersy/b9934790a8892db1a603820c0c23e4a7
// Regexp definitions
var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)

type conventionalMarshaller struct {
	Value interface{}
}

func (c conventionalMarshaller) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(c.Value)

	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			// Empty keys are valid JSON, only lowercase if we do not have an
			// empty key.
			if len(match) > 2 {
				// Decode first rune after the double quotes
				r, width := utf8.DecodeRune(match[1:])
				r = unicode.ToLower(r)
				utf8.EncodeRune(match[1:width+1], r)
			}
			return match
		},
	)

	return converted, err
}

type AuthSuccess struct {
}

type BricksetLogin struct {
	ApiKey string
	Hash   string
}

type BricksetConnection struct {
	Status string `json:"status"`
	Hash   string `json:"hash"`
}

type BricksetSetParams struct {
	PageSize  int
	Theme     string
	Subtheme  string
	SetNumber string
	Year      string
	Owned     int
	Wanted    int
	OrderBy   string
}

type BricksetSet struct {
	Status  string `json:"status"`
	Matches int    `json:"matches"`
	Sets    []struct {
		Setid         int    `json:"setID"`
		Number        string `json:"number"`
		Numbervariant int    `json:"numberVariant"`
		Name          string `json:"name"`
		Year          int    `json:"year"`
		Theme         string `json:"theme"`
		Themegroup    string `json:"themeGroup"`
		Subtheme      string `json:"subtheme"`
		Category      string `json:"category"`
		Released      bool   `json:"released"`
		Pieces        int    `json:"pieces"`
		Minifigs      int    `json:"minifigs"`
		Image         struct {
			Thumbnailurl string `json:"thumbnailURL"`
			Imageurl     string `json:"imageURL"`
		} `json:"image"`
		Brickseturl string `json:"bricksetURL"`
		Collection  struct {
			Owned    bool   `json:"owned"`
			Wanted   bool   `json:"wanted"`
			Qtyowned int    `json:"qtyOwned"`
			Rating   int    `json:"rating"`
			Notes    string `json:"notes"`
		} `json:"collection"`
		Collections struct {
			Ownedby  int `json:"ownedBy"`
			Wantedby int `json:"wantedBy"`
		} `json:"collections"`
		Legocom struct {
			Us struct {
				Retailprice        float64   `json:"retailPrice"`
				Datefirstavailable time.Time `json:"dateFirstAvailable"`
				Datelastavailable  time.Time `json:"dateLastAvailable"`
			} `json:"US"`
			Uk struct {
				Retailprice        float64   `json:"retailPrice"`
				Datefirstavailable time.Time `json:"dateFirstAvailable"`
				Datelastavailable  time.Time `json:"dateLastAvailable"`
			} `json:"UK"`
			Ca struct {
			} `json:"CA"`
			De struct {
			} `json:"DE"`
		} `json:"LEGOCom"`
		Rating               float64 `json:"rating"`
		Reviewcount          int     `json:"reviewCount"`
		Packagingtype        string  `json:"packagingType"`
		Availability         string  `json:"availability"`
		Instructionscount    int     `json:"instructionsCount"`
		Additionalimagecount int     `json:"additionalImageCount"`
		Agerange             struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"ageRange,omitempty"`
		Dimensions struct {
			Height float64 `json:"height"`
			Width  float64 `json:"width"`
			Depth  float64 `json:"depth"`
			Weight float64 `json:"weight"`
		} `json:"dimensions,omitempty"`
		Barcode struct {
			Ean string `json:"EAN"`
			Upc string `json:"UPC"`
		} `json:"barcode,omitempty"`
		Extendeddata struct {
		} `json:"extendedData"`
		Lastupdated time.Time `json:"lastUpdated"`
	} `json:"sets"`
}

type BricksetTheme struct {
	Status  string `json:"status"`
	Matches int    `json:"matches"`
	Themes  []struct {
		Theme         string `json:"theme"`
		Setcount      int    `json:"setCount"`
		Subthemecount int    `json:"subthemeCount"`
		Yearfrom      int    `json:"yearFrom"`
		Yearto        int    `json:"yearTo"`
	} `json:"themes"`
}

type BricksetSubtheme struct {
	Status    string `json:"status"`
	Matches   int    `json:"matches"`
	Subthemes []struct {
		Theme    string `json:"theme"`
		Subtheme string `json:"subtheme"`
		Setcount int    `json:"setCount"`
		Yearfrom int    `json:"yearFrom"`
		Yearto   int    `json:"yearTo"`
	} `json:"subthemes"`
}

func sendRequest(apiKey string, url string, body string) (*resty.Response, error) {

	client := resty.New()

	// client.SetDebug(true)

	baseUrl := "https://brickset.com/api/v3.asmx"
	fullUrl := fmt.Sprint(baseUrl, url)

	resp, _ := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetContentLength(true).
		SetBody(body).
		SetResult(AuthSuccess{}).
		Post(fullUrl)

	return resp, nil
}

func Login(apiKey string, username string, password string) (BricksetLogin, error) {

	var login BricksetLogin

	var bricksetResponse BricksetConnection
	body := fmt.Sprint("apiKey=", apiKey, "&username=", username, "&password=", password)

	resp, err := sendRequest(apiKey, "/login", body)

	if err != nil {
		fmt.Println(err)
		return login, err
	}

	err = json.Unmarshal(resp.Body(), &bricksetResponse)

	if err != nil {
		fmt.Println(err)
		return login, err
	}

	login.ApiKey = apiKey
	login.Hash = bricksetResponse.Hash

	return login, nil
}

func GetSets(apiKey string, userHash string, pageSize int, theme string, subtheme string, setNumber string, year string, owned int, wanted int, orderBy string) (BricksetSet, error) {

	if pageSize == 0 {
		pageSize = 500
	}
	var bricksetResponse BricksetSet
	var params BricksetSetParams
	params.PageSize = pageSize
	params.Theme = theme
	params.Subtheme = subtheme
	params.SetNumber = setNumber
	params.Year = year
	params.Owned = owned
	params.Wanted = wanted
	params.OrderBy = orderBy

	encoded, _ := json.MarshalIndent(conventionalMarshaller{params}, "", "  ")

	fmt.Println("Params in JSON are:", string(encoded))

	body := fmt.Sprint("apiKey=", apiKey, "&userHash=", userHash, "&params=", string(encoded))

	resp, err := sendRequest(apiKey, "/getSets", body)

	if err != nil {
		fmt.Println(err)
		return bricksetResponse, err
	}

	err = json.Unmarshal(resp.Body(), &bricksetResponse)

	if err != nil {
		fmt.Println(err)
		return bricksetResponse, err
	}

	return bricksetResponse, nil
}

func GetThemes(apiKey string) (BricksetTheme, error) {

	var bricksetResponse BricksetTheme
	body := fmt.Sprint("apiKey=", apiKey)

	resp, err := sendRequest(apiKey, "/getThemes", body)

	if err != nil {
		fmt.Println(err)
		return bricksetResponse, err
	}

	err = json.Unmarshal(resp.Body(), &bricksetResponse)

	if err != nil {
		fmt.Println(err)
		return bricksetResponse, err
	}

	return bricksetResponse, nil
}

func GetSubthemes(apiKey string, theme string) (BricksetSubtheme, error) {

	var bricksetResponse BricksetSubtheme
	body := fmt.Sprint("apiKey=", apiKey, "&theme=", theme)

	resp, err := sendRequest(apiKey, "/getSubthemes", body)

	if err != nil {
		fmt.Println(err)
		return bricksetResponse, err
	}

	err = json.Unmarshal(resp.Body(), &bricksetResponse)

	if err != nil {
		fmt.Println(err)
		return bricksetResponse, err
	}

	return bricksetResponse, nil
}
