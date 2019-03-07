package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"strings"
)

type Value struct {
	Value int `yaml:"value"`
	Unit string `yaml:"unit"`
}

type Carrier struct {
	Name string `yaml:"name"`
	Maxload Value `yaml:"max-load"`
	Speed Value `yaml:"speed"`
}

type Data struct {
	Carriers []Carrier `yaml:"carriers"`
}

var (
	distance  = kingpin.Arg("distance", "How far the data is being transferred in kilometers").Required().Int()
	cardType = kingpin.Flag("card-type", "The card type to use. This will be one of [micro, fullsize].").Default("micro").Short('c').String()
	cardSize = kingpin.Flag("card-size", "The size of card to use in GB").Default("1024").Short('s').Int()
)

func main() {
	kingpin.Parse()
	data := Data{}
	yamlFile, err := ioutil.ReadFile("data.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	carriers := data.Carriers
	//cards :/*/= data.Cards

	//for _, carrier := range(carriers) {
	//	fmt.Printf("%s:\t%d%s\t%d%s\n", carrier.Name, carrier.Maxload.Value, carrier.Maxload.Unit, carrier.Speed.Value, carrier.Speed.Unit)
	//}
	//
	//for _, card := range(cards) {
	//	fmt.Printf("%s:\t%d%s\t%v\n", card.Type, card.Weight.Value, card.Weight.Unit, card.Sizes)
	//}

	fmt.Printf("Transferring data %dkm\nUsing Card type: %s\n", *distance, *cardType)

	for _, carrier := range (carriers) {
		maxCards := carrier.getMaxCardCount(*cardType, data)
		fmt.Printf("%s can carry %d cards for a data transfer of %dGB\n", carrier.Name, maxCards, maxCards * *cardSize)

		travelTime := getTravelTimeInSeconds(carrier.Speed, *distance)
		fmt.Printf("Travel time is %f seconds\n\n", travelTime)
	}
}

func (carrier Carrier) getMaxCardCount(cardType string, data Data) int {
	maxMilligrams := carrier.Maxload.getMilligrams()
	cardWeight := getCardWeight(cardType)

	maxCards := maxMilligrams / cardWeight.getMilligrams()

	return maxCards
}

func (value Value) getMilligrams() int {
	switch value.Unit {
	case "mg":
		return value.Value
	case "g":
		return value.Value * 1000
	case "kg":
		return value.Value * 1000 * 1000
	}

	return -1
}

func getCardWeight(cardType string) Value {
	switch strings.ToLower(cardType) {
	case "fullsize":
		return Value{Value: 2, Unit: "g"}
	case "micro":
		return Value{Value: 500, Unit:"mg"}
	}

	return Value{}
}

func getTravelTimeInSeconds(speed Value, distance int) float64 {
	speedInKmph := speed.Value
	// Todo: have a function to convert speed to kmph if it isnt already

	// Convert to Meters per Hour
	metersPerHour := speedInKmph * 1000

	// Convert to Meters per Minute
	metersPerMinute := float64(metersPerHour / 60.0)

	// Convert to Meters per Second
	metersPerSecond := float64(metersPerMinute / 60.0)

	fmt.Printf("%dkm/hr == %dm/hr == %fm/min == %fm/sec\n", speedInKmph, metersPerHour, metersPerMinute, metersPerSecond)

	distanceInMeters := distance * 1000

	secondsToTravel := float64(distanceInMeters) / metersPerSecond

	return secondsToTravel
}