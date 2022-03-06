package main

import (
	"log"
	"time"

	"github.com/sj14/astral"
)

type AstralResponse struct {
	DawnCivil            int64   `json:"dawnCivil,omitempty"`
	DawnAstronomical     int64   `json:"dawnAstronomical,omitempty"`
	DawnNautical         int64   `json:"dawnNautical,omitempty"`
	GoldenRisingStart    int64   `json:"goldenRisingStart,omitempty"`
	GoldenRisingEnd      int64   `json:"goldenRisingEnd,omitempty"`
	Sunrise              int64   `json:"sunrise,omitempty"`
	SunriseNextDay       int64   `json:"sunriseNextDay,omitempty"`
	Noon                 int64   `json:"noon,omitempty"`
	GoldenSettingStart   int64   `json:"goldenSettingStart,omitempty"`
	GoldenSettingEnd     int64   `json:"goldenSettingEnd,omitempty"`
	Sunset               int64   `json:"sunset,omitempty"`
	DuskCivil            int64   `json:"duskCivil,omitempty"`
	DuskAstronomical     int64   `json:"duskAstronomical,omitempty"`
	DuskNautical         int64   `json:"duskNautical,omitempty"`
	Midnight             int64   `json:"midnight,omitempty"`
	MoonPhase            float64 `json:"moonPhase,omitempty"`
	MoonPhaseDescription string  `json:"moonPhaseDescription,omitempty"`
}

func getAstral(latitude float64, longitude float64) *AstralResponse {
	observer := astral.Observer{
		Latitude:  latitude,
		Longitude: longitude,
	}

	t := time.Now()

	dawnCivil, err := astral.Dawn(observer, t, astral.DepressionCivil)
	if err != nil {
		log.Println(err)
		return nil
	}
	dawnAstronomical, err := astral.Dawn(observer, t, astral.DepressionAstronomical)
	if err != nil {
		log.Println(err)
		return nil
	}
	dawnNautical, err := astral.Dawn(observer, t, astral.DepressionNautical)
	if err != nil {
		log.Println(err)
		return nil
	}

	goldenRisingStart, goldenRisingEnd, err := astral.GoldenHour(observer, t, astral.SunDirectionRising)
	if err != nil {
		log.Println(err)
		return nil
	}

	sunrise, err := astral.Sunrise(observer, t)
	if err != nil {
		log.Println(err)
		return nil
	}

	sunriseNextDay, err := astral.Sunrise(observer, t.Add(24*time.Hour))
	if err != nil {
		log.Println(err)
		return nil
	}

	noon := astral.Noon(observer, t)

	goldenSettingStart, goldenSettingEnd, err := astral.GoldenHour(observer, t, astral.SunDirectionSetting)
	if err != nil {
		log.Println(err)
		return nil
	}

	sunset, err := astral.Sunset(observer, t)
	if err != nil {
		log.Println(err)
		return nil
	}

	duskCivil, err := astral.Dusk(observer, t, astral.DepressionCivil)
	if err != nil {
		log.Println(err)
		return nil
	}
	duskAstronomical, err := astral.Dusk(observer, t, astral.DepressionAstronomical)
	if err != nil {
		log.Println(err)
		return nil
	}
	duskNautical, err := astral.Dusk(observer, t, astral.DepressionNautical)
	if err != nil {
		log.Println(err)
		return nil
	}

	midnight := astral.Midnight(observer, t)

	moonPhase := astral.MoonPhase(t)
	moonDesc, err := astral.MoonPhaseDescription(moonPhase)
	if err != nil {
		log.Fatalf("failed parsing moon phase: %v", err)
	}

	return &AstralResponse{
		DawnCivil:            dawnCivil.Unix(),
		DawnAstronomical:     dawnAstronomical.Unix(),
		DawnNautical:         dawnNautical.Unix(),
		GoldenRisingStart:    goldenRisingStart.Unix(),
		GoldenRisingEnd:      goldenRisingEnd.Unix(),
		Sunrise:              sunrise.Unix(),
		SunriseNextDay:       sunriseNextDay.Unix(),
		Noon:                 noon.Unix(),
		GoldenSettingStart:   goldenSettingStart.Unix(),
		GoldenSettingEnd:     goldenSettingEnd.Unix(),
		Sunset:               sunset.Unix(),
		DuskCivil:            duskCivil.Unix(),
		DuskAstronomical:     duskAstronomical.Unix(),
		DuskNautical:         duskNautical.Unix(),
		Midnight:             midnight.Unix(),
		MoonPhase:            moonPhase,
		MoonPhaseDescription: moonDesc,
	}
}
