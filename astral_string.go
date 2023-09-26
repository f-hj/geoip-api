package main

import (
	"log"
	"time"

	"github.com/sj14/astral/pkg/astral"
)

type AstralStringResponse struct {
	Golden   *AstralRisingSetting   `json:"golden,omitempty"`
	Blue     *AstralRisingSetting   `json:"blue,omitempty"`
	Sunrise  string                 `json:"sunrise,omitempty"`
	Sunset   string                 `json:"sunset,omitempty"`
	Dawn     *AstralCivilAstroNauti `json:"dawn,omitempty"`
	Dusk     *AstralCivilAstroNauti `json:"dusk,omitempty"`
	Day      *AstralStartEnd        `json:"day,omitempty"`
	Noon     string                 `json:"noon,omitempty"`
	Night    *AstralStartEnd        `json:"night,omitempty"`
	Midnight string                 `json:"midnight,omitempty"`
	Moon     *AstralPhase           `json:"moon,omitempty"`
}

type AstralCivilAstroNauti struct {
	Civil        string `json:"civil,omitempty"`
	Astronomical string `json:"astronomical,omitempty"`
	Nautical     string `json:"nautical,omitempty"`
}

type AstralRisingSetting struct {
	Rising  *AstralStartEnd `json:"rising,omitempty"`
	Setting *AstralStartEnd `json:"setting,omitempty"`
}

type AstralStartEnd struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

type AstralPhase struct {
	Phase       float64 `json:"phase,omitempty"`
	Description string  `json:"description,omitempty"`
}

func getAstralString(t time.Time, latitude float64, longitude float64, timezone string) *AstralStringResponse {
	observer := astral.Observer{
		Latitude:  latitude,
		Longitude: longitude,
	}

	if t.IsZero() {
		t = time.Now()
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc, _ = time.LoadLocation("UTC")
	}

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

	blueRisingStart, blueRisingEnd, err := astral.BlueHour(observer, t, astral.SunDirectionRising)
	if err != nil {
		log.Println(err)
		return nil
	}

	blueSettingStart, blueSettingEnd, err := astral.BlueHour(observer, t, astral.SunDirectionSetting)
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

	dayFrom, dayTo, err := astral.Daylight(observer, t)
	if err != nil {
		log.Println(err)
		return nil
	}

	nightFrom, nightTo, err := astral.Night(observer, t)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &AstralStringResponse{
		Dawn: &AstralCivilAstroNauti{
			Astronomical: dawnAstronomical.In(loc).Format(time.RFC3339),
			Civil:        dawnCivil.In(loc).Format(time.RFC3339),
			Nautical:     dawnNautical.In(loc).Format(time.RFC3339),
		},
		Blue: &AstralRisingSetting{
			Rising: &AstralStartEnd{
				Start: blueRisingStart.In(loc).Format(time.RFC3339),
				End:   blueRisingEnd.In(loc).Format(time.RFC3339),
			},
			Setting: &AstralStartEnd{
				Start: blueSettingStart.In(loc).Format(time.RFC3339),
				End:   blueSettingEnd.In(loc).Format(time.RFC3339),
			},
		},
		Golden: &AstralRisingSetting{
			Rising: &AstralStartEnd{
				Start: goldenRisingStart.In(loc).Format(time.RFC3339),
				End:   goldenRisingEnd.In(loc).Format(time.RFC3339),
			},
			Setting: &AstralStartEnd{
				Start: goldenSettingStart.In(loc).Format(time.RFC3339),
				End:   goldenSettingEnd.In(loc).Format(time.RFC3339),
			},
		},
		Sunrise: sunrise.In(loc).Format(time.RFC3339),
		Noon:    noon.In(loc).Format(time.RFC3339),
		Sunset:  sunset.In(loc).Format(time.RFC3339),
		Dusk: &AstralCivilAstroNauti{
			Astronomical: duskAstronomical.In(loc).Format(time.RFC3339),
			Civil:        duskCivil.In(loc).Format(time.RFC3339),
			Nautical:     duskNautical.In(loc).Format(time.RFC3339),
		},
		Day: &AstralStartEnd{
			Start: dayFrom.In(loc).Format(time.RFC3339),
			End:   dayTo.In(loc).Format(time.RFC3339),
		},
		Midnight: midnight.In(loc).Format(time.RFC3339),
		Moon: &AstralPhase{
			Phase:       moonPhase,
			Description: moonDesc,
		},
		Night: &AstralStartEnd{
			Start: nightFrom.In(loc).Format(time.RFC3339),
			End:   nightTo.In(loc).Format(time.RFC3339),
		},
	}
}
