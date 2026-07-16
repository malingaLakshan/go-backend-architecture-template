package site

import (
	"encoding/json"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// ParseSiteGraph - happy path
// -----------------------------------------------------------------------------

func TestParseSiteGraph_FullHierarchy(t *testing.T) {
	raw := `{
		"id": "SITE-PARSE",
		"name": "Ignored Site Name",
		"floors": [
			{
				"id": "F1",
				"name": "Ground Floor",
				"regions": [
					{
						"id": "R1",
						"name": "Entrance",
						"regions": [
							{
								"id": "R1-CHILD",
								"name": "Entrance Child",
								"regions": [
									{
										"id": "R1-GRANDCHILD",
										"name": "Entrance Grandchild"
									}
								]
							}
						]
					},
					{
						"id": "R2",
						"name": "Checkout"
					}
				],
				"readers": [
					{
						"id": "RDR1",
						"make": "Impinj",
						"antennas": [
							{"port": 1, "gain": 6.0},
							{"port": 2, "gain": 6.0},
							{"port": 3},
							{"port": 4}
						]
					},
					{
						"id": "RDR2",
						"antennas": [
							{"port": 5},
							{"port": 6},
							{"port": 7},
							{"port": 8}
						]
					}
				]
			},
			{
				"id": "F2",
				"regions": [
					{
						"id": "R3",
						"regions": [
							{
								"id": "R3-CHILD"
							}
						]
					}
				],
				"readers": [
					{
						"id": "RDR3",
						"antennas": [
							{"port": 1}
						]
					}
				]
			}
		]
	}`

	vs, err := ParseSiteGraph([]byte(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vs.SiteID != "SITE-PARSE" {
		t.Errorf("SiteID: want SITE-PARSE, got %s", vs.SiteID)
	}

	if len(vs.Floors) != 2 {
		t.Fatalf("Floors: want 2, got %d", len(vs.Floors))
	}

	// First floor.
	f1 := vs.Floors[0]

	if f1.ID != "F1" {
		t.Errorf("Floor[0].ID: want F1, got %s", f1.ID)
	}

	if len(f1.Regions) != 2 {
		t.Fatalf("Floor[0].Regions: want 2, got %d", len(f1.Regions))
	}

	if len(f1.Readers) != 2 {
		t.Fatalf("Floor[0].Readers: want 2, got %d", len(f1.Readers))
	}

	// Recursive Region hierarchy:
	// R1 -> R1-CHILD -> R1-GRANDCHILD.
	r1 := f1.Regions[0]

	if r1.ID != "R1" {
		t.Errorf("Floor[0].Region[0].ID: want R1, got %s", r1.ID)
	}

	if len(r1.Regions) != 1 {
		t.Fatalf(
			"Floor[0].Region[0].Regions: want 1, got %d",
			len(r1.Regions),
		)
	}

	r1Child := r1.Regions[0]

	if r1Child.ID != "R1-CHILD" {
		t.Errorf(
			"nested Region ID: want R1-CHILD, got %s",
			r1Child.ID,
		)
	}

	if len(r1Child.Regions) != 1 {
		t.Fatalf(
			"R1-CHILD.Regions: want 1, got %d",
			len(r1Child.Regions),
		)
	}

	if r1Child.Regions[0].ID != "R1-GRANDCHILD" {
		t.Errorf(
			"grandchild Region ID: want R1-GRANDCHILD, got %s",
			r1Child.Regions[0].ID,
		)
	}

	if f1.Regions[1].ID != "R2" {
		t.Errorf(
			"Floor[0].Region[1].ID: want R2, got %s",
			f1.Regions[1].ID,
		)
	}

	// Readers are directly under the Floor, not under Regions.
	rdr1 := f1.Readers[0]

	if rdr1.ID != "RDR1" {
		t.Errorf("Floor[0].Reader[0].ID: want RDR1, got %s", rdr1.ID)
	}

	wantPorts1 := []int{1, 2, 3, 4}

	if len(rdr1.AntennaPorts) != len(wantPorts1) {
		t.Fatalf(
			"Reader[0].AntennaPorts: want %v, got %v",
			wantPorts1,
			rdr1.AntennaPorts,
		)
	}

	for i, wantPort := range wantPorts1 {
		if rdr1.AntennaPorts[i] != wantPort {
			t.Errorf(
				"Reader[0].AntennaPorts[%d]: want %d, got %d",
				i,
				wantPort,
				rdr1.AntennaPorts[i],
			)
		}
	}

	rdr2 := f1.Readers[1]

	if rdr2.ID != "RDR2" {
		t.Errorf("Floor[0].Reader[1].ID: want RDR2, got %s", rdr2.ID)
	}

	wantPorts2 := []int{5, 6, 7, 8}

	if len(rdr2.AntennaPorts) != len(wantPorts2) {
		t.Fatalf(
			"Reader[1].AntennaPorts: want %v, got %v",
			wantPorts2,
			rdr2.AntennaPorts,
		)
	}

	for i, wantPort := range wantPorts2 {
		if rdr2.AntennaPorts[i] != wantPort {
			t.Errorf(
				"Reader[1].AntennaPorts[%d]: want %d, got %d",
				i,
				wantPort,
				rdr2.AntennaPorts[i],
			)
		}
	}

	// Second floor.
	f2 := vs.Floors[1]

	if f2.ID != "F2" {
		t.Errorf("Floor[1].ID: want F2, got %s", f2.ID)
	}

	if len(f2.Regions) != 1 {
		t.Fatalf("Floor[1].Regions: want 1, got %d", len(f2.Regions))
	}

	if f2.Regions[0].ID != "R3" {
		t.Errorf(
			"Floor[1].Region[0].ID: want R3, got %s",
			f2.Regions[0].ID,
		)
	}

	if len(f2.Regions[0].Regions) != 1 {
		t.Fatalf(
			"Floor[1].Region[0].Regions: want 1, got %d",
			len(f2.Regions[0].Regions),
		)
	}

	if f2.Regions[0].Regions[0].ID != "R3-CHILD" {
		t.Errorf(
			"nested Region ID: want R3-CHILD, got %s",
			f2.Regions[0].Regions[0].ID,
		)
	}

	if len(f2.Readers) != 1 {
		t.Fatalf("Floor[1].Readers: want 1, got %d", len(f2.Readers))
	}

	if f2.Readers[0].ID != "RDR3" {
		t.Errorf(
			"Floor[1].Reader[0].ID: want RDR3, got %s",
			f2.Readers[0].ID,
		)
	}

	if len(f2.Readers[0].AntennaPorts) != 1 ||
		f2.Readers[0].AntennaPorts[0] != 1 {
		t.Errorf(
			"Floor[1].Reader[0].AntennaPorts: want [1], got %v",
			f2.Readers[0].AntennaPorts,
		)
	}
}

// -----------------------------------------------------------------------------
// ParseSiteGraph - unknown fields are ignored
// -----------------------------------------------------------------------------

func TestParseSiteGraph_UnknownFieldsIgnored(t *testing.T) {
	raw := `{
		"id": "SITE-X",
		"completely_unknown": "value",
		"networking": {
			"ip": "1.2.3.4"
		},
		"floors": [
			{
				"id": "F1",
				"bounds": [1, 2, 3, 4],
				"regions": [
					{
						"id": "R1",
						"inventoryType": "retail",
						"regions": [
							{
								"id": "R1-CHILD",
								"behaviors": []
							}
						]
					}
				],
				"readers": [
					{
						"id": "RDR1",
						"model": "R420",
						"timeouts": {},
						"physicality": "ceiling",
						"behaviors": [],
						"antennas": [
							{
								"port": 1,
								"position": [0, 0, 0]
							}
						]
					}
				]
			}
		]
	}`

	vs, err := ParseSiteGraph([]byte(raw))
	if err != nil {
		t.Fatalf(
			"unexpected error parsing SiteGraph with unknown fields: %v",
			err,
		)
	}

	if vs.SiteID != "SITE-X" {
		t.Errorf("SiteID: want SITE-X, got %s", vs.SiteID)
	}

	if len(vs.Floors) != 1 {
		t.Fatalf("Floors: want 1, got %d", len(vs.Floors))
	}

	floor := vs.Floors[0]

	if len(floor.Regions) != 1 {
		t.Fatalf("Regions: want 1, got %d", len(floor.Regions))
	}

	if floor.Regions[0].ID != "R1" {
		t.Errorf("Region ID: want R1, got %s", floor.Regions[0].ID)
	}

	if len(floor.Regions[0].Regions) != 1 {
		t.Fatalf(
			"nested Regions: want 1, got %d",
			len(floor.Regions[0].Regions),
		)
	}

	if floor.Regions[0].Regions[0].ID != "R1-CHILD" {
		t.Errorf(
			"nested Region ID: want R1-CHILD, got %s",
			floor.Regions[0].Regions[0].ID,
		)
	}

	if len(floor.Readers) != 1 {
		t.Fatalf("Readers: want 1, got %d", len(floor.Readers))
	}

	if floor.Readers[0].ID != "RDR1" {
		t.Errorf("Reader ID: want RDR1, got %s", floor.Readers[0].ID)
	}

	if len(floor.Readers[0].AntennaPorts) != 1 ||
		floor.Readers[0].AntennaPorts[0] != 1 {
		t.Errorf(
			"unexpected antenna ports: %v",
			floor.Readers[0].AntennaPorts,
		)
	}
}

// -----------------------------------------------------------------------------
// ParseSiteGraph - error cases
// -----------------------------------------------------------------------------

func TestParseSiteGraph_EmptyInput_Error(t *testing.T) {
	_, err := ParseSiteGraph([]byte{})

	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestParseSiteGraph_MissingRootID_Error(t *testing.T) {
	raw := `{"floors":[]}`

	_, err := ParseSiteGraph([]byte(raw))

	if err == nil {
		t.Fatal("expected error when root id field is missing")
	}

	if !strings.Contains(err.Error(), "root field: id") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseSiteGraph_InvalidJSON_Error(t *testing.T) {
	_, err := ParseSiteGraph([]byte(`{ bad json `))

	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseSiteGraph_MissingFloorID_Error(t *testing.T) {
	raw := `{
		"id": "SITE-X",
		"floors": [
			{
				"regions": []
			}
		]
	}`

	_, err := ParseSiteGraph([]byte(raw))

	if err == nil {
		t.Fatal("expected error when Floor ID is missing")
	}

	if !strings.Contains(err.Error(), "floor[0]") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseSiteGraph_MissingRegionID_Error(t *testing.T) {
	raw := `{
		"id": "SITE-X",
		"floors": [
			{
				"id": "F1",
				"regions": [
					{}
				]
			}
		]
	}`

	_, err := ParseSiteGraph([]byte(raw))

	if err == nil {
		t.Fatal("expected error when Region ID is missing")
	}

	if !strings.Contains(err.Error(), "region[0]") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseSiteGraph_MissingNestedRegionID_Error(t *testing.T) {
	raw := `{
		"id": "SITE-X",
		"floors": [
			{
				"id": "F1",
				"regions": [
					{
						"id": "R1",
						"regions": [
							{}
						]
					}
				]
			}
		]
	}`

	_, err := ParseSiteGraph([]byte(raw))

	if err == nil {
		t.Fatal("expected error when nested Region ID is missing")
	}

	if !strings.Contains(err.Error(), "region[0]") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseSiteGraph_MissingReaderID_Error(t *testing.T) {
	raw := `{
		"id": "SITE-X",
		"floors": [
			{
				"id": "F1",
				"readers": [
					{
						"antennas": [
							{"port": 1}
						]
					}
				]
			}
		]
	}`

	_, err := ParseSiteGraph([]byte(raw))

	if err == nil {
		t.Fatal("expected error when Reader ID is missing")
	}

	if !strings.Contains(err.Error(), "reader[0]") {
		t.Errorf("unexpected error: %v", err)
	}
}

// -----------------------------------------------------------------------------
// ToSiteGraphResponse - round trip
// -----------------------------------------------------------------------------

func TestToSiteGraphResponse_RoundTrip(t *testing.T) {
	original := &ValidationSite{
		SiteID: "SITE-RT",
		Floors: []ValidationFloor{
			{
				ID: "F1",
				Regions: []ValidationRegion{
					{
						ID: "R1",
						Regions: []ValidationRegion{
							{
								ID: "R1-CHILD",
								Regions: []ValidationRegion{
									{
										ID: "R1-GRANDCHILD",
									},
								},
							},
						},
					},
				},
				Readers: []ValidationReader{
					{
						ID:           "RDR1",
						AntennaPorts: []int{1, 2, 3, 4, 5, 6, 7, 8},
					},
				},
			},
		},
	}

	// Convert to SiteGraph-compatible JSON.
	response := ToSiteGraphResponse(original)

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal SiteGraphResponse: %v", err)
	}

	// Parse the converted response back into ValidationSite.
	parsed, err := ParseSiteGraph(data)
	if err != nil {
		t.Fatalf(
			"ParseSiteGraph failed on ToSiteGraphResponse output: %v",
			err,
		)
	}

	if parsed.SiteID != original.SiteID {
		t.Errorf(
			"SiteID: want %s, got %s",
			original.SiteID,
			parsed.SiteID,
		)
	}

	if len(parsed.Floors) != 1 {
		t.Fatalf("Floors: want 1, got %d", len(parsed.Floors))
	}

	floor := parsed.Floors[0]

	if floor.ID != "F1" {
		t.Errorf("Floor ID: want F1, got %s", floor.ID)
	}

	if len(floor.Regions) != 1 {
		t.Fatalf("Regions: want 1, got %d", len(floor.Regions))
	}

	if floor.Regions[0].ID != "R1" {
		t.Errorf("Region ID: want R1, got %s", floor.Regions[0].ID)
	}

	if len(floor.Regions[0].Regions) != 1 {
		t.Fatalf(
			"nested Regions: want 1, got %d",
			len(floor.Regions[0].Regions),
		)
	}

	child := floor.Regions[0].Regions[0]

	if child.ID != "R1-CHILD" {
		t.Errorf("child Region ID: want R1-CHILD, got %s", child.ID)
	}

	if len(child.Regions) != 1 {
		t.Fatalf(
			"grandchild Regions: want 1, got %d",
			len(child.Regions),
		)
	}

	if child.Regions[0].ID != "R1-GRANDCHILD" {
		t.Errorf(
			"grandchild Region ID: want R1-GRANDCHILD, got %s",
			child.Regions[0].ID,
		)
	}

	// Readers must remain directly under the Floor.
	if len(floor.Readers) != 1 {
		t.Fatalf("Readers: want 1, got %d", len(floor.Readers))
	}

	if floor.Readers[0].ID != "RDR1" {
		t.Errorf("Reader ID: want RDR1, got %s", floor.Readers[0].ID)
	}

	ports := floor.Readers[0].AntennaPorts

	if len(ports) != 8 {
		t.Fatalf("AntennaPorts: want 8, got %d: %v", len(ports), ports)
	}

	for i, wantPort := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		if ports[i] != wantPort {
			t.Errorf(
				"AntennaPorts[%d]: want %d, got %d",
				i,
				wantPort,
				ports[i],
			)
		}
	}
}

// -----------------------------------------------------------------------------
// CountValidationSite
// -----------------------------------------------------------------------------

func TestCountValidationSite_MultipleFloors(t *testing.T) {
	vs := makeValidationSiteForCountTest()

	counts := CountValidationSite(vs)

	if counts.Floors != 2 {
		t.Errorf("Floors: want 2, got %d", counts.Floors)
	}

	if counts.Regions != 4 {
		t.Errorf("Regions: want 4, got %d", counts.Regions)
	}

	if counts.Readers != 8 {
		t.Errorf("Readers: want 8, got %d", counts.Readers)
	}

	if counts.AntennaPorts != 32 {
		t.Errorf(
			"AntennaPorts: want 32, got %d",
			counts.AntennaPorts,
		)
	}
}

func TestCountValidationSite_RecursiveRegions(t *testing.T) {
	vs := &ValidationSite{
		SiteID: "SITE-RECURSIVE",
		Floors: []ValidationFloor{
			{
				ID: "F1",
				Regions: []ValidationRegion{
					{
						ID: "R1",
						Regions: []ValidationRegion{
							{
								ID: "R2",
								Regions: []ValidationRegion{
									{
										ID: "R3",
									},
									{
										ID: "R4",
									},
								},
							},
						},
					},
					{
						ID: "R5",
					},
				},
				Readers: []ValidationReader{
					{
						ID:           "RDR1",
						AntennaPorts: []int{1, 2},
					},
				},
			},
		},
	}

	counts := CountValidationSite(vs)

	if counts.Floors != 1 {
		t.Errorf("Floors: want 1, got %d", counts.Floors)
	}

	// R1, R2, R3, R4 and R5.
	if counts.Regions != 5 {
		t.Errorf("Regions: want 5, got %d", counts.Regions)
	}

	if counts.Readers != 1 {
		t.Errorf("Readers: want 1, got %d", counts.Readers)
	}

	if counts.AntennaPorts != 2 {
		t.Errorf(
			"AntennaPorts: want 2, got %d",
			counts.AntennaPorts,
		)
	}
}

func makeValidationSiteForCountTest() *ValidationSite {
	makeReader := func(id string) ValidationReader {
		return ValidationReader{
			ID:           id,
			AntennaPorts: []int{1, 2, 3, 4},
		}
	}

	return &ValidationSite{
		SiteID: "SITE-COUNT",
		Floors: []ValidationFloor{
			{
				ID: "F1",
				Regions: []ValidationRegion{
					{
						ID: "F1-R1",
					},
					{
						ID: "F1-R2",
					},
				},
				Readers: []ValidationReader{
					makeReader("F1-RDR1"),
					makeReader("F1-RDR2"),
					makeReader("F1-RDR3"),
					makeReader("F1-RDR4"),
				},
			},
			{
				ID: "F2",
				Regions: []ValidationRegion{
					{
						ID: "F2-R1",
					},
					{
						ID: "F2-R2",
					},
				},
				Readers: []ValidationReader{
					makeReader("F2-RDR1"),
					makeReader("F2-RDR2"),
					makeReader("F2-RDR3"),
					makeReader("F2-RDR4"),
				},
			},
		},
	}
}