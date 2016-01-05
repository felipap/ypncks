package pancakeFinder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"regexp"
)

type Menu []Meal

const (
	PathLocations = "http://www.yaledining.org/fasttrack/locations.cfm?version=3"
	PathMenu      = "http://www.yaledining.org/fasttrack/menus.cfm?location=%d&version=3"
)

type Location struct {
	Id            int     `field:"ID_LOCATION"`
	LocationCode  float64 `field:"LOCATIONCODE"`
	Name          string  `field:"DININGLOCATIONNAME"`
	Type          string  `field:"TYPE"`
	Capacity      int     `field:"CAPACITY"`
	GeoLocation   string  `field:"GEOLOCATION"`
	IsClosed      int     `field:"ISCLOSED"`
	Address       string  `field:"ADDRESS"`
	Phone         string  `field:"PHONE"`
	Manager1Name  string  `field:"MANAGER1NAME"`
	Manager1Email string  `field:"MANAGER1EMAIL"`
	Manager2Name  string  `field:"MANAGER2NAME"`
	Manager2Email string  `field:"MANAGER2EMAIL"`
	Manager3Name  string  `field:"MANAGER3NAME"`
	Manager3Email string  `field:"MANAGER3EMAIL"`
	Manager4Name  string  `field:"MANAGER4NAME"`
	Manager4Email string  `field:"MANAGER4EMAIL"`
}

type Meal struct {
	IdLocation    int    `field:"ID_LOCATION"`
	LocationCode  int    `field:"LOCATIONCODE"`
	Location      string `field:"LOCATION"`
	MealName      string `field:"MEALNAME"`
	MealCode      int    `field:"MEALCODE"`
	MenuDate      string `field:"MENUDATE"`
	Id            int    `field:"ID"`
	Course        string `field:"COURSE"`
	CourseCode    int    `field:"COURSECODE"`
	MenuItemId    int    `field:"MENUITEMID"`
	Name          string `field:"MENUITEM"`
	IsPar         bool   `field:"ISPAR"`
	MealOpens     string `field:"MEALOPENS"`
	MealCloses    string `field:"MEALCLOSES"`
	IsDefaultMeal int    `field:"ISDEFAULTMEAL"`
	IsMenu        bool   `field:"ISMENU"`
}

func (m Meal) String() string {
	return fmt.Sprintf("[%s (%d, %s)]", m.Name, m.Id, m.MealName)
}

func (l Location) String() string {
	return fmt.Sprintf("[%s (%d, %s)]", l.Name, l.Id, l.Name)
}

func unpackStruct(m map[string]interface{}, iObj interface{}) error {
	rt := reflect.TypeOf(iObj).Elem()
	rv := reflect.ValueOf(iObj)

	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("Object passed is not a pointer.")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("Object passed is not a pointer to a struct.")
	}

	if !rv.CanSet() {
		log.Fatal("Can't set attribute in unpack.")
	}

	for i := 0; i < rt.NumField(); i++ {
		tfield := rt.Field(i)
		jsonCol := tfield.Tag.Get("field")

		if _jsonVal, ok := m[jsonCol]; ok {
			jsonVal := reflect.ValueOf(_jsonVal)

			if _jsonVal != nil {
				//fmt.Printf("> %s %v\n", jsonCol, jsonVal)
				// Convert jsonVal to tfield's type
				// See jimmyfrasche.github.io/go-reflection-codex
				if jsonVal.Type().ConvertibleTo(tfield.Type) {
					rv.Field(i).Set(jsonVal.Convert(tfield.Type))
				} else {
					log.Fatalf("Can't convert json field %v of type %s "+
						"to struct field %v of type %s.", jsonCol, jsonVal.Type(),
						tfield.Name, tfield.Type)
				}
			} else {
			}
		}
	}

	return nil
}

// This parses data coming from both endpoints.
func parseYaleResp(body []byte) ([]map[string]interface{}, error) {
	var dataLists map[string][]interface{}
	err := json.Unmarshal(body, &dataLists)
	if err != nil {
		return nil, err
	}

	// Parse and copy column names.
	cols := make([]string, len(dataLists["COLUMNS"]))
	for x, v := range dataLists["COLUMNS"] {
		col, ok := v.(string)
		if ok == false {
			return nil, fmt.Errorf("Unexpected format for COLUMNS")
		}
		cols[x] = col
	}

	var dataFinal []map[string]interface{}

	// Zip object [{column:value,...},...]
	for _, _row := range dataLists["DATA"] {
		row, ok := _row.([]interface{})
		if ok == false {
			return nil, fmt.Errorf("Unexpected format for DATA.")
		}
		fmtRow := make(map[string]interface{})
		for key, attr := range row {
			// fmt.Printf("%+v: %v\n", cols[key], attr)
			fmtRow[cols[key]] = attr
		}
		dataFinal = append(dataFinal, fmtRow)
	}

	return dataFinal, nil
}

func fetchDiningHalls() ([]Location, error) {
	log.Printf("Fetching dining halls at %s.\n", PathLocations)
	resp, err := http.Get(PathLocations)
	if err != nil {
		log.Fatal("Error getting locations.\n", err)
		return nil, fmt.Errorf("Failed to fetch locations.")
	}
	defer resp.Body.Close()
	bodyLocs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Failed to fetch locations.")
	}

	dataLocs, err := parseYaleResp(bodyLocs)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Failed to fetch locations.")
	}

	locs := make([]Location, len(dataLocs))
	for i := range dataLocs {
		unpackStruct(dataLocs[i], &locs[i])
	}

	// Keep only residential colleges dining halls.
	for i := 0; i < len(locs); {
		if locs[i].Type != "Residential" {
			locs = append(locs[:i], locs[i+1:]...)
		} else {
			i++
		}
	}

	return locs, nil
}

func fetchMenu(hall Location, out chan<- Menu) {
	var menu Menu = nil
	log.Printf("Fetching menu for %s.\n", hall.Name)

	defer func() { out <- menu }()

	url := fmt.Sprintf(PathMenu, hall.Id)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting menu for %s (%f).\n", hall.Name, hall.Id, err)
		return
	}
	defer resp.Body.Close()
	bodyMenu, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return
	}

	dataMenu, err := parseYaleResp(bodyMenu)
	if err != nil {
		log.Fatal(err)
		return
	}
	menu = make([]Meal, len(dataMenu))
	for i := range dataMenu {
		unpackStruct(dataMenu[i], &menu[i])
	}
}

func Find() []Meal {
	dhalls, err := fetchDiningHalls()
	if err != nil {
		log.Fatal("Failed fetching dining halls.")
		log.Fatal(err)
		return nil
	}

	cheff := make(chan Menu)
	limiter := make(chan bool)
	for _, loc := range dhalls {
		go func(loc Location) {
			<-limiter
			fetchMenu(loc, cheff)
			limiter <- true
		}(loc)
	}
	fmt.Printf("Integer value.\n")

	for i := 0; i < 1; i++ { // Limit # of simultaneous requests.
		limiter <- true
	}


	var isPancake = func(meal Meal) bool {
		match, err := regexp.Match("(?i)pancakes?", []byte(meal.Name))
		if err != nil {
			log.Fatal("Regexp failed.")
		}
		return match
	}
	

	pancakes := make([]Meal, 0)
	for _ = range dhalls {
		for _, meal := range <-cheff {
			if isPancake(meal) {
				pancakes = append(pancakes, meal)
				//fmt.Printf("Meal: %s \n", meal)
			}
		}
	}

	return pancakes
}
