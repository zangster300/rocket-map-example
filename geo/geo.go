package geo

import (
	"math"
	"math/rand/v2"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

var landPolygons []orb.Polygon
var landMultiPolygons []orb.MultiPolygon

func LoadLand(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return err
	}

	for _, f := range fc.Features {
		if f.Geometry == nil {
			continue
		}

		switch g := f.Geometry.(type) {
		case orb.Polygon:
			landPolygons = append(landPolygons, g)

		case orb.MultiPolygon:
			landMultiPolygons = append(landMultiPolygons, g)
		}
	}

	return nil
}

func RandomLandPoint() [2]float64 {
	for {
		lon := rand.Float64()*360 - 180
		lat := math.Asin(rand.Float64()*2-1) * 180 / math.Pi

		if lat < -60 {
			continue
		}

		if isOnLand(lon, lat) {
			return [2]float64{lon, lat}
		}
	}
}

func isOnLand(lon, lat float64) bool {
	p := orb.Point{lon, lat}

	for _, poly := range landPolygons {
		if planar.PolygonContains(poly, p) {
			return true
		}
	}

	for _, mp := range landMultiPolygons {
		if planar.MultiPolygonContains(mp, p) {
			return true
		}
	}

	return false
}
