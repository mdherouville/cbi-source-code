package main

////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

// The following is a sample record from the Taxi Trips dataset retrieved from the City of Chicago Data Portal

////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

// trip_id	"c354c843908537bbf90997917b714f1c63723785"
// trip_start_timestamp	"2021-11-13T22:45:00.000"
// trip_end_timestamp	"2021-11-13T23:00:00.000"
// trip_seconds	"703"
// trip_miles	"6.83"
// pickup_census_tract	"17031840300"
// dropoff_census_tract	"17031081800"
// pickup_community_area	"59"
// dropoff_community_area	"8"
// fare	"27.5"
// tip	"0"
// additional_charges	"1.02"
// trip_total	"28.52"
// shared_trip_authorized	false
// trips_pooled	"1"
// pickup_centroid_latitude	"41.8335178865"
// pickup_centroid_longitude	"-87.6813558293"
// pickup_centroid_location
// type	"Point"
// coordinates
// 		0	-87.6813558293
// 		1	41.8335178865
// dropoff_centroid_latitude	"41.8932163595"
// dropoff_centroid_longitude	"-87.6378442095"
// dropoff_centroid_location
// type	"Point"
// coordinates
// 		0	-87.6378442095
// 		1	41.8932163595
////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"database/sql"
	"encoding/json"

	"github.com/kelvins/geocoder"
	_ "github.com/lib/pq"
)

type TaxiTripsJsonRecords []struct {
	Trip_id                    string `json:"trip_id"`
	Trip_start_timestamp       string `json:"trip_start_timestamp"`
	Trip_end_timestamp         string `json:"trip_end_timestamp"`
	Pickup_centroid_latitude   string `json:"pickup_centroid_latitude"`
	Pickup_centroid_longitude  string `json:"pickup_centroid_longitude"`
	Dropoff_centroid_latitude  string `json:"dropoff_centroid_latitude"`
	Dropoff_centroid_longitude string `json:"dropoff_centroid_longitude"`
}
type TransportationNetworkProviderJsonRecords []struct {
	Trip_id                    string `json:"trip_id"`
	Trip_start_timestamp       string `json:"trip_start_timestamp"`
	Trip_end_timestamp         string `json:"trip_end_timestamp"`
	Pickup_centroid_latitude   string `json:"pickup_centroid_latitude"`
	Pickup_centroid_longitude  string `json:"pickup_centroid_longitude"`
	Dropoff_centroid_latitude  string `json:"dropoff_centroid_latitude"`
	Dropoff_centroid_longitude string `json:"dropoff_centroid_longitude"`
}
type UnemploymentJsonRecords []struct {
	Community_area      string `json:"community_area"`
	Community_area_name string `json:"community_area_name"`
	Unemployment        string `json:"unemployment"`
	Below_poverty_level string `json:"below_poverty_level"`
	Per_capita_income   string `json:"per_capita_income"`
}
type CovidCasesJsonRecords []struct {
	Zip_code     string `json:"zip_code"`
	Week_start   string `json:"week_start"`
	Week_end     string `json:"week_end"`
	Cases_weekly string `json:"cases_weekly"`
	Tests_weekly string `json:"tests_weekly"`
}
type CCVIJsonRecords []struct {
	Community_area_name   string `json:"community_area_name"`
	Community_area_or_zip string `json:"community_area_or_zip"`
	Ccvi_category         string `json:"ccvi_category"`
}
type BuildingPermitsJsonRecords []struct {
	Id             string `json:"id"`
	Permit_        string `json:"permit_"`
	Permit_type    string `json:"permit_type"`
	Community_area string `json:"community_area"`
	Total_fee      string `json:"total_fee"`
	Reported_cost  string `json:"reported_cost"`
}

var limit = "1000"

func main() {

	// Establish connection to Postgres Database
	db_connection := "user=postgres dbname=chicago_business_intelligence password=root host=/cloudsql/c-b-i-mlbdh:us-central1:mypostgres sslmode=disable port = 5432"

	// Docker image for the microservice - uncomment when deploy
	//db_connection := "user=postgres dbname=chicago_business_intelligence password=root host=host.docker.internal sslmode=disable"

	db, err := sql.Open("postgres", db_connection)
	if err != nil {
		panic(err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Couldn't Connect to database")
		panic(err)
	}

	// Spin in a loop and pull data from the city of chicago data portal
	// Once every hour, day, week, etc.
	// Though, please note that Not all datasets need to be pulled on daily basis
	// fine-tune the following code-snippet as you see necessary
	for {
		// build and fine-tune functions to pull data from different data sources
		// This is a code snippet to show you how to pull data from different data sources.
		GetTaxiTrips(db)
		GetTransportationNetworkProviders(db)
		GetUnemploymentRates(db)
		GetBuildingPermits(db)
		GetCCVI(db)
		GetCovidCasesRecords(db)

		// Pull the data once a day
		// You might need to pull Taxi Trips and COVID data on daily basis
		// but not the unemployment dataset becasue its dataset doesn't change every day
		time.Sleep(24 * time.Hour)
	}

}

func GetTaxiTrips(db *sql.DB) {

	// This function is NOT complete
	// It provides code-snippets for the data source: https://data.cityofchicago.org/Transportation/Taxi-Trips/wrvz-psew
	// You need to complete the implmentation and add the data source: https://data.cityofchicago.org/Transportation/Transportation-Network-Providers-Trips/m6dm-c72p

	// Data Collection needed from two data sources:
	// 1. https://data.cityofchicago.org/Transportation/Taxi-Trips/wrvz-psew
	// 2. https://data.cityofchicago.org/Transportation/Transportation-Network-Providers-Trips/m6dm-c72p

	fmt.Println("GetTaxiTrips: Collecting Taxi Trips Data")

	// Get your geocoder.ApiKey from here :
	// https://developers.google.com/maps/documentation/geocoding/get-api-key?authuser=2
	//AIzaSyAp3-89wiaVO_SixMGATB-zHi566NWkU8o
	//AIzaSyDNlwnDqCj6ksNZv5mp-0ePW-CsFRHqavY
	geocoder.ApiKey = "AIzaSyAp3-89wiaVO_SixMGATB-zHi566NWkU8o"

	drop_table := `drop table if exists taxi_trips`
	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}

	create_table := `CREATE TABLE IF NOT EXISTS "taxi_trips" (
						"id"   SERIAL , 
						"trip_id" VARCHAR(255) UNIQUE, 
						"trip_start" TIMESTAMP WITH TIME ZONE, 
						"trip_end" TIMESTAMP WITH TIME ZONE, 
						"pickup_latitude" DOUBLE PRECISION, 
						"pickup_longitude" DOUBLE PRECISION, 
						"dropoff_latitude" DOUBLE PRECISION, 
						"dropoff_longitude" DOUBLE PRECISION, 
						"pickup_zip" VARCHAR(255), 
						"dropoff_zip" VARCHAR(255), 
						PRIMARY KEY ("id") 
					);`

	_, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

	// While doing unit-testing keep the limit value to 500
	// later you could change it to 1000, 2000, 10,000, etc.
	//fmt.Println("https://data.cityofchicago.org/resource/wrvz-psew.json?$limit=5")
	var url = "https://data.cityofchicago.org/resource/wrvz-psew.json?$limit=" + limit

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	//text := string(body)
	//fmt.Println(text)
	var taxi_trips_list TaxiTripsJsonRecords
	json.Unmarshal(body, &taxi_trips_list)

	for i := 0; i < len(taxi_trips_list); i++ {

		// We will execute definsive coding to check for messy/dirty/missing data values
		// Any record that has messy/dirty/missing data we don't enter it in the data lake/table

		trip_id := taxi_trips_list[i].Trip_id
		if trip_id == "" {
			continue
		}

		// if trip start/end timestamp doesn't have the length of 23 chars in the format "0000-00-00T00:00:00.000"
		// skip this record

		// get Trip_start_timestamp
		trip_start_timestamp := taxi_trips_list[i].Trip_start_timestamp
		if len(trip_start_timestamp) < 23 {
			continue
		}

		// get Trip_end_timestamp
		trip_end_timestamp := taxi_trips_list[i].Trip_end_timestamp
		if len(trip_end_timestamp) < 23 {
			continue
		}

		pickup_centroid_latitude := taxi_trips_list[i].Pickup_centroid_latitude

		if pickup_centroid_latitude == "" {
			continue
		}

		pickup_centroid_longitude := taxi_trips_list[i].Pickup_centroid_longitude
		//pickup_centroid_longitude := taxi_trips_list[i].PICKUP_LONG

		if pickup_centroid_longitude == "" {
			continue
		}

		dropoff_centroid_latitude := taxi_trips_list[i].Dropoff_centroid_latitude
		//dropoff_centroid_latitude := taxi_trips_list[i].DROPOFF_LAT

		if dropoff_centroid_latitude == "" {
			continue
		}

		dropoff_centroid_longitude := taxi_trips_list[i].Dropoff_centroid_longitude
		//dropoff_centroid_longitude := taxi_trips_list[i].DROPOFF_LONG

		if dropoff_centroid_longitude == "" {
			continue
		}

		// Using pickup_centroid_latitude and pickup_centroid_longitude in geocoder.GeocodingReverse
		// we could find the pickup zip-code

		pickup_centroid_latitude_float, _ := strconv.ParseFloat(pickup_centroid_latitude, 64)
		pickup_centroid_longitude_float, _ := strconv.ParseFloat(pickup_centroid_longitude, 64)
		pickup_location := geocoder.Location{
			Latitude:  pickup_centroid_latitude_float,
			Longitude: pickup_centroid_longitude_float,
		}

		pickup_address_list, _ := geocoder.GeocodingReverse(pickup_location)
		//fmt.Println("Address list pickup : ")
		//fmt.Println(pickup_address_list)
		pickup_zip_code := ""
		if len(pickup_address_list) > 0 {
			pickup_address := pickup_address_list[0]
			pickup_zip_code = pickup_address.PostalCode
		} else {
			pickup_zip_code = "zip geocoder error"
		}

		// Using dropoff_centroid_latitude and dropoff_centroid_longitude in geocoder.GeocodingReverse
		// we could find the dropoff zip-code

		dropoff_centroid_latitude_float, _ := strconv.ParseFloat(dropoff_centroid_latitude, 64)
		dropoff_centroid_longitude_float, _ := strconv.ParseFloat(dropoff_centroid_longitude, 64)

		dropoff_location := geocoder.Location{
			Latitude:  dropoff_centroid_latitude_float,
			Longitude: dropoff_centroid_longitude_float,
		}

		dropoff_address_list, _ := geocoder.GeocodingReverse(dropoff_location)
		//fmt.Println("Address list dropoff : ")
		//fmt.Println(dropoff_address_list)
		dropoff_zip_code := ""
		if len(pickup_address_list) > 0 {
			dropoff_address := dropoff_address_list[0]
			dropoff_zip_code = dropoff_address.PostalCode
		} else {
			dropoff_zip_code = "zip geocoder error"
		}

		sql := `INSERT INTO taxi_trips ("trip_id", 
										"trip_start", 
										"trip_end", 
										"pickup_latitude", 
										"pickup_longitude", 
										"dropoff_latitude", 
										"dropoff_longitude", 
										"pickup_zip", 
										"dropoff_zip") 
										values($1, $2, $3, $4, $5, $6, $7, $8, $9)`

		_, err = db.Exec(
			sql,
			trip_id,
			trip_start_timestamp,
			trip_end_timestamp,
			pickup_centroid_latitude,
			pickup_centroid_longitude,
			dropoff_centroid_latitude,
			dropoff_centroid_longitude,
			pickup_zip_code,
			dropoff_zip_code)

		if err != nil {

			panic(err)
		}

	}

}

func GetUnemploymentRates(db *sql.DB) {
	fmt.Println("GetUnemploymentRates: Collecting Unemployment Rates Data")

	drop_table := `drop table if exists unemployment_rates`
	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}

	create_table := `CREATE TABLE IF NOT EXISTS "unemployment_rates" (
		"id"   SERIAL , 
		"com_area_number" VARCHAR(255), 
		"com_area_name" VARCHAR(255), 
		"zip" VARCHAR(255), 
		"unemployment" DOUBLE PRECISION, 
		"below_poverty_level" DOUBLE PRECISION, 
		"per_capita_income" DOUBLE PRECISION, 
		PRIMARY KEY ("id") 
	);`

	_, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

	var url = "https://data.cityofchicago.org/resource/iqnk-2tcu.json?$limit=" + limit
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	//text := string(body)
	//fmt.Println(text)
	var unemployment_rates_list UnemploymentJsonRecords
	json.Unmarshal(body, &unemployment_rates_list)

	for i := 0; i < len(unemployment_rates_list); i++ {

		community_area := unemployment_rates_list[i].Community_area
		if community_area == "" {
			continue
		}

		community_area_name := unemployment_rates_list[i].Community_area_name
		if community_area_name == "" {
			continue
		}

		unemployment := unemployment_rates_list[i].Unemployment
		if unemployment == "" {
			continue
		}

		below_poverty_level := unemployment_rates_list[i].Below_poverty_level
		if below_poverty_level == "" {
			continue
		}

		per_capita_income := unemployment_rates_list[i].Per_capita_income
		if per_capita_income == "" {
			continue
		}

		zip := GetZipFromCommunityArea(&community_area)

		sql := `INSERT INTO unemployment_rates ("com_area_number", 
												"com_area_name", 
												"zip", 
												"unemployment", 
												"below_poverty_level", 
												"per_capita_income") 
												values($1, $2, $3, $4, $5, $6)`

		_, err = db.Exec(
			sql,
			community_area,
			community_area_name,
			zip,
			unemployment,
			below_poverty_level,
			per_capita_income)

		if err != nil {

			panic(err)
		}
	}
}

func GetBuildingPermits(db *sql.DB) {
	fmt.Println("GetBuildingPermits: Collecting Building Permits Data")

	drop_table := `drop table if exists building_permits`
	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}

	create_table := `CREATE TABLE IF NOT EXISTS "building_permits" (
		"id"   SERIAL , 
		"permit_id" VARCHAR(255), 
		"permit_number" VARCHAR(255), 
		"permi_zip" VARCHAR(255), 
		"permit_com_area" DOUBLE PRECISION, 
		"permit_fee" DOUBLE PRECISION, 
		"permit_reported_cost" DOUBLE PRECISION, 
		PRIMARY KEY ("id") 
	);`

	_, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

	var url = "https://data.cityofchicago.org/resource/ydr8-5enu.json?$limit=" + limit
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	//text := string(body)
	//fmt.Println(text)
	var building_permits_list BuildingPermitsJsonRecords
	json.Unmarshal(body, &building_permits_list)

	for i := 0; i < len(building_permits_list); i++ {

		permit_id := building_permits_list[i].Id
		if permit_id == "" {
			continue
		}

		permit_number := building_permits_list[i].Permit_
		if permit_number == "" {
			continue
		}

		permit_com_area := building_permits_list[i].Community_area
		if permit_com_area == "" {
			continue
		}

		permit_fee := building_permits_list[i].Total_fee
		if permit_fee == "" {
			continue
		}

		permit_reported_cost := building_permits_list[i].Reported_cost
		if permit_reported_cost == "" {
			continue
		}

		permit_zip := GetZipFromCommunityArea(&permit_com_area)

		sql := `INSERT INTO building_permits ("permit_id", 
												"permit_number", 
												"permi_zip", 
												"permit_com_area", 
												"permit_fee", 
												"permit_reported_cost") 
												values($1, $2, $3, $4, $5, $6)`

		_, err = db.Exec(
			sql,
			permit_id,
			permit_number,
			permit_zip,
			permit_com_area,
			permit_fee,
			permit_reported_cost)

		if err != nil {

			panic(err)
		}

	}
}

func GetCCVI(db *sql.DB) {

	fmt.Println("GetCCVI: Collecting CCVI Data")

	drop_table := `drop table if exists ccvi_table`
	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}

	create_table := `CREATE TABLE IF NOT EXISTS "ccvi_table" (
		"id"   SERIAL , 
		"com_area_name" VARCHAR(255), 
		"com_area_number" VARCHAR(255), 
		"zip" VARCHAR(255), 
		"ccvi_category" VARCHAR(255), 
		PRIMARY KEY ("id") 
	);`

	_, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

	var url = "https://data.cityofchicago.org/resource/xhc6-88s9.json?$limit=" + limit
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	//text := string(body)
	//fmt.Println(text)
	var ccvi_list CCVIJsonRecords
	json.Unmarshal(body, &ccvi_list)

	for i := 0; i < len(ccvi_list); i++ {

		com_area_name := ccvi_list[i].Community_area_name
		if com_area_name == "" {
			continue
		}

		com_area_number := ccvi_list[i].Community_area_or_zip
		if com_area_number == "" {
			continue
		}

		ccvi_category := ccvi_list[i].Ccvi_category
		if ccvi_category == "" {
			continue
		}
		zip := GetZipFromCommunityArea(&com_area_number)

		sql := `INSERT INTO ccvi_table (		"com_area_name", 
												"com_area_number", 
												"zip", 
												"ccvi_category") 
												values($1, $2, $3, $4)`

		_, err = db.Exec(
			sql,
			com_area_name,
			com_area_number,
			zip,
			ccvi_category)

		if err != nil {

			panic(err)
		}
	}
}

func GetCovidCasesRecords(db *sql.DB) {

	fmt.Println("GetCovidCasesRecords: Collecting covid cases Data")

	drop_table := `drop table if exists covid_cases_table`
	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}

	create_table := `CREATE TABLE IF NOT EXISTS "covid_cases_table" (
		"id"   SERIAL , 
		"zip" VARCHAR(255), 
		"week_start" TIMESTAMP WITH TIME ZONE, 
		"week_end" TIMESTAMP WITH TIME ZONE, 
		"cases_weekly" VARCHAR(255), 
		"tests_weekly" VARCHAR(255),
		PRIMARY KEY ("id") 
	);`

	_, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

	var url = "https://data.cityofchicago.org/resource/yhhz-zm2v.json?$limit=" + limit
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	//text := string(body)
	//fmt.Println(text)
	var covid_cases_list CovidCasesJsonRecords
	json.Unmarshal(body, &covid_cases_list)

	for i := 0; i < len(covid_cases_list); i++ {

		zip := covid_cases_list[i].Zip_code
		if zip == "" {
			continue
		}

		week_start := covid_cases_list[i].Week_start
		if len(week_start) < 23 {
			continue
		}
		week_end := covid_cases_list[i].Week_end
		if len(week_end) < 23 {
			continue
		}

		cases_weekly := covid_cases_list[i].Cases_weekly
		if cases_weekly == "" {
			continue
		}

		tests_weekly := covid_cases_list[i].Tests_weekly
		if tests_weekly == "" {
			continue
		}

		sql := `INSERT INTO covid_cases_table ("zip", 
												"week_start", 
												"week_end", 
												"cases_weekly",
												"tests_weekly") 
												values($1, $2, $3, $4, $5)`

		_, err = db.Exec(
			sql,
			zip,
			week_start,
			week_end,
			cases_weekly,
			tests_weekly)

		if err != nil {

			panic(err)
		}

	}

}

func GetTransportationNetworkProviders(db *sql.DB) {

	fmt.Println("GetTransportationNetworkProviders: Collecting transportation network providers Data")

	geocoder.ApiKey = "AIzaSyAp3-89wiaVO_SixMGATB-zHi566NWkU8o"

	drop_table := `drop table if exists transportation_network_providers`
	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}

	create_table := `CREATE TABLE IF NOT EXISTS "transportation_network_providers" (
						"id"   SERIAL , 
						"trip_id" VARCHAR(255) UNIQUE, 
						"trip_start" TIMESTAMP WITH TIME ZONE, 
						"trip_end" TIMESTAMP WITH TIME ZONE, 
						"pickup_latitude" DOUBLE PRECISION, 
						"pickup_longitude" DOUBLE PRECISION, 
						"dropoff_latitude" DOUBLE PRECISION, 
						"dropoff_longitude" DOUBLE PRECISION, 
						"pickup_zip" VARCHAR(255), 
						"dropoff_zip" VARCHAR(255), 
						PRIMARY KEY ("id") 
					);`

	_, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}
	var url = "https://data.cityofchicago.org/resource/m6dm-c72p.json?$limit=" + limit

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	//text := string(body)
	//fmt.Println(text)
	var transportation_network_providers_list TransportationNetworkProviderJsonRecords
	json.Unmarshal(body, &transportation_network_providers_list)

	for i := 0; i < len(transportation_network_providers_list); i++ {

		trip_id := transportation_network_providers_list[i].Trip_id
		if trip_id == "" {
			continue
		}

		trip_start_timestamp := transportation_network_providers_list[i].Trip_start_timestamp
		if len(trip_start_timestamp) < 23 {
			continue
		}

		trip_end_timestamp := transportation_network_providers_list[i].Trip_end_timestamp
		if len(trip_end_timestamp) < 23 {
			continue
		}

		pickup_centroid_latitude := transportation_network_providers_list[i].Pickup_centroid_latitude

		if pickup_centroid_latitude == "" {
			continue
		}

		pickup_centroid_longitude := transportation_network_providers_list[i].Pickup_centroid_longitude

		if pickup_centroid_longitude == "" {
			continue
		}

		dropoff_centroid_latitude := transportation_network_providers_list[i].Dropoff_centroid_latitude

		if dropoff_centroid_latitude == "" {
			continue
		}

		dropoff_centroid_longitude := transportation_network_providers_list[i].Dropoff_centroid_longitude

		if dropoff_centroid_longitude == "" {
			continue
		}

		// Using pickup_centroid_latitude and pickup_centroid_longitude in geocoder.GeocodingReverse
		// we could find the pickup zip-code

		pickup_centroid_latitude_float, _ := strconv.ParseFloat(pickup_centroid_latitude, 64)
		pickup_centroid_longitude_float, _ := strconv.ParseFloat(pickup_centroid_longitude, 64)
		pickup_location := geocoder.Location{
			Latitude:  pickup_centroid_latitude_float,
			Longitude: pickup_centroid_longitude_float,
		}

		pickup_address_list, _ := geocoder.GeocodingReverse(pickup_location)
		//fmt.Println("Address list pickup : ")
		//fmt.Println(pickup_address_list)
		pickup_zip_code := ""
		if len(pickup_address_list) > 0 {
			pickup_address := pickup_address_list[0]
			pickup_zip_code = pickup_address.PostalCode
		} else {
			pickup_zip_code = "zip geocoder error"
		}

		// Using dropoff_centroid_latitude and dropoff_centroid_longitude in geocoder.GeocodingReverse
		// we could find the dropoff zip-code

		dropoff_centroid_latitude_float, _ := strconv.ParseFloat(dropoff_centroid_latitude, 64)
		dropoff_centroid_longitude_float, _ := strconv.ParseFloat(dropoff_centroid_longitude, 64)

		dropoff_location := geocoder.Location{
			Latitude:  dropoff_centroid_latitude_float,
			Longitude: dropoff_centroid_longitude_float,
		}

		dropoff_address_list, _ := geocoder.GeocodingReverse(dropoff_location)
		//fmt.Println("Address list dropoff : ")
		//fmt.Println(dropoff_address_list)
		dropoff_zip_code := ""
		if len(pickup_address_list) > 0 {
			dropoff_address := dropoff_address_list[0]
			dropoff_zip_code = dropoff_address.PostalCode
		} else {
			dropoff_zip_code = "zip geocoder error"
		}

		sql := `INSERT INTO transportation_network_providers ("trip_id", 
															  "trip_start", 
															  "trip_end", 
															  "pickup_latitude",
															  "pickup_longitude", 
															  "dropoff_latitude", 
															  "dropoff_longitude", 
															  "pickup_zip", 
															  "dropoff_zip") 
															  values($1, $2, $3, $4, $5, $6, $7, $8, $9)`

		_, err = db.Exec(
			sql,
			trip_id,
			trip_start_timestamp,
			trip_end_timestamp,
			pickup_centroid_latitude,
			pickup_centroid_longitude,
			dropoff_centroid_latitude,
			dropoff_centroid_longitude,
			pickup_zip_code,
			dropoff_zip_code)

		if err != nil {

			panic(err)
		}

	}

}

func GetZipFromCommunityArea(com_area *string) string {

	return "zip/area crossreference not implemented"
}
