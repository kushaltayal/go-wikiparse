package wikiparse

import (
	"math"
	"strings"
	"testing"
)

type testinput struct {
	input string
	lat   float64
	lon   float64
	err   string
}

var testdata = []testinput{
	testinput{
		"{{coord|61.1631|-149.9721|type:landmark_globe:earth_region:US-AK_scale:150000_source:gnis|name=Kulis Air National Guard Base}}",
		61.1631,
		-149.9721,
		"",
	},
	testinput{
		"{{coord|29.5734571|N|2.3730469|E|scale:10000000|format=dms|display=title}}",
		29.5734571,
		2.3730469,
		"",
	},
	testinput{
		"{{coord|27|59|16|N|86|56|40|E}}",
		27.98777777,
		86.94444444,
		"",
	},
	testinput{
		"{{coord|27|59|16|S|86|56|40|E}}",
		-27.98777777,
		86.94444444,
		"",
	},
	testinput{
		"{{coord|27|59|16|N|86|56|40|W}}",
		27.98777777,
		-86.94444444,
		"",
	},
	testinput{
		"{{coord|27|59|16|S|86|56|40|W}}",
		-27.98777777,
		-86.94444444,
		"",
	},
	testinput{
		"{{Coord|display=title|45|N|114|W|region:US-ID_type:adm1st_scale:3000000}}",
		45,
		-114,
		"",
	},
	testinput{
		"{{Coord|42||N|82||W|}}",
		42,
		-82,
		"",
	},
	testinput{
		"{{Coord|display=title|41.762736| -72.674286}}",
		41.762736,
		-72.674286,
		"",
	},
	testinput{
		"North Maple in Russell ({{coord|38.895352|-98.861034}}) and it remained his " +
			"official residence throughout his political career." +
			"<ref>{{cite news| url=http://www.time.com/}}",
		38.895352,
		-98.861034,
		"",
	},
	testinput{
		"{{coord|97|59|16|S|86|56|40|W|invalid lat}}",
		-97.98777777,
		-86.94444444,
		"Invalid latitude: -97.98777",
	},
	testinput{
		"{{coord|27|59|16|S|186|56|40|W|invalid long}}",
		-27.98777777,
		-186.94444444,
		"Invalid longitude: -186.9444",
	},
	testinput{
		"<nowiki>{{coord|27|59|16|N|86|56|40|E}}</nowiki>",
		0,
		0,
		"No coord data found.",
	},
	testinput{
		`<nowiki>
{{coord|27|59|16|N|86|56|40|E}}
</nowiki>`,
		0,
		0,
		"No coord data found.",
	},
	testinput{
		"<!-- {{coord|27|59|16|N|86|56|40|E}} -->",
		0,
		0,
		"No coord data found.",
	},
	testinput{
		`<!--
{{coord|27|59|16|N|86|56|40|E}}
-->`,
		0,
		0,
		"No coord data found.",
	},
}

func assertEpsilon(t *testing.T, input, field string, expected, got float64) {
	if math.Abs(got-expected) > 0.00001 {
		t.Fatalf("Expected %v for %v of %v, got %v",
			expected, field, input, got)
	}
}

func testOne(t *testing.T, ti testinput, input string) {
	coord, err := ParseCoords(input)
	switch {
	case err != nil && ti.err == "":
		t.Fatalf("Unexpected error on %v, got %v, wanted %q", input, err, ti.err)
	case err != nil && strings.HasPrefix(err.Error(), ti.err):
		// ok
	case err == nil && ti.err == "":
		// ok
	case err == nil && ti.err != "":
		t.Fatalf("Expected error %q on %v", ti.err, input)
	default:
		t.Fatalf("Wanted %v,%v with error %v, got %#v with error %v",
			ti.lat, ti.lon, ti.err, coord, err)
	}
	t.Logf("Parsed %#v with %v", coord, err)
	assertEpsilon(t, input, "lon", ti.lon, coord.Lon)
	assertEpsilon(t, input, "lat", ti.lat, coord.Lat)
	t.Logf("Results for %s:  %#v", input, coord)
}

func TestCoordSimple(t *testing.T) {
	for _, ti := range testdata {
		testOne(t, ti, ti.input)
	}
}

func TestCoordWithGarbage(t *testing.T) {
	for _, ti := range testdata {
		input := " some random garbage " + ti.input + " and stuff"
		testOne(t, ti, input)
	}
}

func TestCoordMultiline(t *testing.T) {
	for _, ti := range testdata {
		input := " some random garbage\n\nnewlines\n" + ti.input + " and stuff"
		testOne(t, ti, input)
	}
}
