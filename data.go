package site

import (
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// ValidateStructure - complete match
// -----------------------------------------------------------------------------

func TestValidateStructure_CompleteMatch(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion(
					"R1",
					validationRegion(
						"R1-CHILD",
						validationRegion("R1-GRANDCHILD"),
					),
				),
				validationRegion("R2"),
			},
			validationReader("READER-1", 1, 2, 3, 4),
			validationReader("READER-2", 5, 6, 7, 8),
		),
	)

	target := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion(
					"R1",
					validationRegion(
						"R1-CHILD",
						validationRegion("R1-GRANDCHILD"),
					),
				),
				validationRegion("R2"),
			},
			validationReader("READER-1", 1, 2, 3, 4),
			validationReader("READER-2", 5, 6, 7, 8),
		),
	)

	result := ValidateStructure(recorded, target)

	if !result.Passed {
		t.Fatalf("expected validation to pass, errors: %v", result.Errors)
	}

	if len(result.Mismatches) != 0 {
		t.Errorf("expected no mismatches, got %v", result.Mismatches)
	}

	assertCategory(
		t,
		"SiteID",
		result.SiteID,
		1,
		1,
		true,
	)

	assertCategory(
		t,
		"Floors",
		result.Floors,
		1,
		1,
		true,
	)

	// R1, R1-CHILD, R1-GRANDCHILD and R2.
	assertCategory(
		t,
		"Regions",
		result.Regions,
		4,
		4,
		true,
	)

	assertCategory(
		t,
		"Readers",
		result.Readers,
		2,
		2,
		true,
	)

	assertCategory(
		t,
		"AntennaPorts",
		result.AntennaPorts,
		8,
		8,
		true,
	)
}

// -----------------------------------------------------------------------------
// Target can contain additional structures
// -----------------------------------------------------------------------------

func TestValidateStructure_TargetExtrasAllowed(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion(
					"R1",
					validationRegion("R1-CHILD"),
				),
			},
			validationReader("READER-1", 1, 2),
		),
	)

	target := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion(
					"R1",
					validationRegion("R1-CHILD"),
					validationRegion("EXTRA-CHILD"),
				),
				validationRegion("EXTRA-REGION"),
			},
			validationReader("READER-1", 1, 2, 3, 4),
			validationReader("EXTRA-READER", 1, 2),
		),
		validationFloor(
			"EXTRA-FLOOR",
			[]ValidationRegion{
				validationRegion("EXTRA-FLOOR-REGION"),
			},
			validationReader("EXTRA-FLOOR-READER", 1),
		),
	)

	result := ValidateStructure(recorded, target)

	if !result.Passed {
		t.Fatalf(
			"expected validation to pass when target has extra structures, errors: %v",
			result.Errors,
		)
	}

	assertCategory(t, "Floors", result.Floors, 1, 1, true)
	assertCategory(t, "Regions", result.Regions, 2, 2, true)
	assertCategory(t, "Readers", result.Readers, 1, 1, true)
	assertCategory(t, "AntennaPorts", result.AntennaPorts, 2, 2, true)
}

// -----------------------------------------------------------------------------
// Site ID mismatch
// -----------------------------------------------------------------------------

func TestValidateStructure_SiteIDMismatch(t *testing.T) {
	recorded := validationSite(
		"RECORDED-SITE",
		validationFloor(
			"F1",
			nil,
			validationReader("READER-1", 1),
		),
	)

	target := validationSite(
		"TARGET-SITE",
		validationFloor(
			"F1",
			nil,
			validationReader("READER-1", 1),
		),
	)

	result := ValidateStructure(recorded, target)

	if result.Passed {
		t.Fatal("expected Site ID mismatch to fail validation")
	}

	assertCategory(t, "SiteID", result.SiteID, 1, 0, false)

	if len(result.Mismatches) != 1 {
		t.Fatalf(
			"expected 1 mismatch, got %d: %v",
			len(result.Mismatches),
			result.Mismatches,
		)
	}

	mismatch := result.Mismatches[0]

	if mismatch.Type != "SiteID" {
		t.Errorf("mismatch Type: want SiteID, got %s", mismatch.Type)
	}

	if mismatch.SiteID != "RECORDED-SITE" {
		t.Errorf(
			"mismatch SiteID: want RECORDED-SITE, got %s",
			mismatch.SiteID,
		)
	}

	assertErrorContains(
		t,
		result,
		"recorded site ID RECORDED-SITE does not match target site ID TARGET-SITE",
	)
}

// -----------------------------------------------------------------------------
// Missing Floor
// -----------------------------------------------------------------------------

func TestValidateStructure_MissingFloor(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion("R1"),
			},
			validationReader("READER-1", 1, 2),
		),
		validationFloor(
			"F2",
			[]ValidationRegion{
				validationRegion("R2"),
			},
			validationReader("READER-2", 1),
		),
	)

	target := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion("R1"),
			},
			validationReader("READER-1", 1, 2),
		),
	)

	result := ValidateStructure(recorded, target)

	if result.Passed {
		t.Fatal("expected missing Floor to fail validation")
	}

	assertCategory(t, "Floors", result.Floors, 2, 1, false)

	mismatch := findMismatch(
		t,
		result,
		func(m ValidationMismatch) bool {
			return m.Type == "Floor" && m.FloorID == "F2"
		},
	)

	if mismatch.Message == "" {
		t.Error("expected missing Floor mismatch to contain a message")
	}

	assertErrorContains(
		t,
		result,
		"target site is missing Floor ID F2",
	)
}

// -----------------------------------------------------------------------------
// Recursive Region validation
// -----------------------------------------------------------------------------

func TestValidateStructure_RecursiveRegionsMatch(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion(
					"R1",
					validationRegion(
						"R2",
						validationRegion("R3"),
						validationRegion("R4"),
					),
				),
			},
		),
	)

	target := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion(
					"R1",
					validationRegion(
						"R2",
						validationRegion("R3"),
						validationRegion("R4"),
					),
				),
			},
		),
	)

	result := ValidateStructure(recorded, target)

	if !result.Passed {
		t.Fatalf(
			"expected recursive Region hierarchy to pass, errors: %v",
			result.Errors,
		)
	}

	// R1, R2, R3 and R4.
	assertCategory(t, "Regions", result.Regions, 4, 4, true)
}

func TestValidateStructure_MissingTopLevelRegion(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion("RECORDED-REGION"),
			},
		),
	)

	target := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion("OTHER-REGION"),
			},
		),
	)

	result := ValidateStructure(recorded, target)

	if result.Passed {
		t.Fatal("expected missing top-level Region to fail validation")
	}

	assertCategory(t, "Regions", result.Regions, 1, 0, false)

	mismatch := findMismatch(
		t,
		result,
		func(m ValidationMismatch) bool {
			return m.Type == "Region" &&
				m.FloorID == "F1" &&
				m.ParentRegionID == "" &&
				m.RegionID == "RECORDED-REGION"
		},
	)

	if mismatch.Message == "" {
		t.Error("expected Region mismatch message")
	}

	assertErrorContains(
		t,
		result,
		"target site is missing Region ID RECORDED-REGION under Floor ID F1",
	)
}

func TestValidateStructure_RegionUnderWrongParentFails(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion(
					"PARENT-A",
					validationRegion("CHILD-X"),
				),
				validationRegion("PARENT-B"),
			},
		),
	)

	// CHILD-X exists in the target, but it is under PARENT-B instead of PARENT-A.
	target := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion("PARENT-A"),
				validationRegion(
					"PARENT-B",
					validationRegion("CHILD-X"),
				),
			},
		),
	)

	result := ValidateStructure(recorded, target)

	if result.Passed {
		t.Fatal(
			"expected Region under the wrong parent to fail validation",
		)
	}

	assertCategory(t, "Regions", result.Regions, 3, 2, false)

	mismatch := findMismatch(
		t,
		result,
		func(m ValidationMismatch) bool {
			return m.Type == "Region" &&
				m.FloorID == "F1" &&
				m.ParentRegionID == "PARENT-A" &&
				m.RegionID == "CHILD-X"
		},
	)

	if mismatch.Message == "" {
		t.Error("expected wrong-parent Region mismatch message")
	}

	assertErrorContains(
		t,
		result,
		"target site is missing Region ID CHILD-X under parent Region ID PARENT-A",
	)
}

func TestValidateStructure_RegionUnderWrongFloorFails(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion("REGION-X"),
			},
		),
		validationFloor("F2", nil),
	)

	// REGION-X exists, but only under F2.
	target := validationSite(
		"SITE-1",
		validationFloor("F1", nil),
		validationFloor(
			"F2",
			[]ValidationRegion{
				validationRegion("REGION-X"),
			},
		),
	)

	result := ValidateStructure(recorded, target)

	if result.Passed {
		t.Fatal(
			"expected Region under the wrong Floor to fail validation",
		)
	}

	findMismatch(
		t,
		result,
		func(m ValidationMismatch) bool {
			return m.Type == "Region" &&
				m.FloorID == "F1" &&
				m.ParentRegionID == "" &&
				m.RegionID == "REGION-X"
		},
	)
}

// -----------------------------------------------------------------------------
// Reader validation
// -----------------------------------------------------------------------------

func TestValidateStructure_MissingReaderUnderFloor(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion("R1"),
			},
			validationReader("READER-1", 1, 2),
			validationReader("READER-2", 1),
		),
	)

	target := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion("R1"),
			},
			validationReader("READER-1", 1, 2),
		),
	)

	result := ValidateStructure(recorded, target)

	if result.Passed {
		t.Fatal("expected missing Reader to fail validation")
	}

	assertCategory(t, "Readers", result.Readers, 2, 1, false)

	mismatch := findMismatch(
		t,
		result,
		func(m ValidationMismatch) bool {
			return m.Type == "Reader" &&
				m.FloorID == "F1" &&
				m.ReaderID == "READER-2"
		},
	)

	if mismatch.Message == "" {
		t.Error("expected missing Reader mismatch message")
	}

	assertErrorContains(
		t,
		result,
		"target site is missing Reader ID READER-2 under Floor ID F1",
	)
}

func TestValidateStructure_ReaderUnderWrongFloorFails(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			nil,
			validationReader("READER-X", 1),
		),
		validationFloor("F2", nil),
	)

	// READER-X exists in the target but is under F2, not F1.
	target := validationSite(
		"SITE-1",
		validationFloor("F1", nil),
		validationFloor(
			"F2",
			nil,
			validationReader("READER-X", 1),
		),
	)

	result := ValidateStructure(recorded, target)

	if result.Passed {
		t.Fatal(
			"expected Reader under the wrong Floor to fail validation",
		)
	}

	findMismatch(
		t,
		result,
		func(m ValidationMismatch) bool {
			return m.Type == "Reader" &&
				m.FloorID == "F1" &&
				m.ReaderID == "READER-X"
		},
	)
}

// -----------------------------------------------------------------------------
// Antenna-port validation
// -----------------------------------------------------------------------------

func TestValidateStructure_MissingAntennaPort(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			nil,
			validationReader("READER-1", 1, 2, 3, 4),
		),
	)

	target := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			nil,
			validationReader("READER-1", 1, 2, 4),
		),
	)

	result := ValidateStructure(recorded, target)

	if result.Passed {
		t.Fatal("expected missing antenna port to fail validation")
	}

	assertCategory(t, "Readers", result.Readers, 1, 1, true)
	assertCategory(t, "AntennaPorts", result.AntennaPorts, 4, 3, false)

	mismatch := findMismatch(
		t,
		result,
		func(m ValidationMismatch) bool {
			return m.Type == "AntennaPort" &&
				m.FloorID == "F1" &&
				m.ReaderID == "READER-1" &&
				m.AntennaPort == 3
		},
	)

	if mismatch.Message == "" {
		t.Error("expected missing antenna-port mismatch message")
	}

	assertErrorContains(
		t,
		result,
		"target site is missing antenna port 3 under Reader ID READER-1",
	)
}

func TestValidateStructure_AntennaPortUnderWrongReaderFails(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			nil,
			validationReader("READER-1", 1, 2),
			validationReader("READER-2", 3),
		),
	)

	// Port 2 exists in the target, but under READER-2 rather than READER-1.
	target := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			nil,
			validationReader("READER-1", 1),
			validationReader("READER-2", 2, 3),
		),
	)

	result := ValidateStructure(recorded, target)

	if result.Passed {
		t.Fatal(
			"expected antenna port under the wrong Reader to fail validation",
		)
	}

	findMismatch(
		t,
		result,
		func(m ValidationMismatch) bool {
			return m.Type == "AntennaPort" &&
				m.FloorID == "F1" &&
				m.ReaderID == "READER-1" &&
				m.AntennaPort == 2
		},
	)
}

// -----------------------------------------------------------------------------
// Multiple mismatches and backward-compatible Errors
// -----------------------------------------------------------------------------

func TestValidateStructure_MultipleMismatches(t *testing.T) {
	recorded := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion(
					"R1",
					validationRegion("R1-CHILD"),
				),
			},
			validationReader("READER-1", 1, 2),
			validationReader("READER-2", 1),
		),
	)

	target := validationSite(
		"SITE-1",
		validationFloor(
			"F1",
			[]ValidationRegion{
				validationRegion("R1"),
			},
			validationReader("READER-1", 1),
		),
	)

	result := ValidateStructure(recorded, target)

	if result.Passed {
		t.Fatal("expected validation with multiple mismatches to fail")
	}

	// Missing:
	// - R1-CHILD
	// - READER-1 antenna port 2
	// - READER-2
	if len(result.Mismatches) != 3 {
		t.Fatalf(
			"expected 3 mismatches, got %d: %v",
			len(result.Mismatches),
			result.Mismatches,
		)
	}

	if len(result.Errors) != len(result.Mismatches) {
		t.Fatalf(
			"Errors and Mismatches lengths differ: Errors=%d Mismatches=%d",
			len(result.Errors),
			len(result.Mismatches),
		)
	}

	for i, mismatch := range result.Mismatches {
		if result.Errors[i] != mismatch.Message {
			t.Errorf(
				"Errors[%d] does not mirror mismatch message: want %q, got %q",
				i,
				mismatch.Message,
				result.Errors[i],
			)
		}
	}
}

// -----------------------------------------------------------------------------
// Test helpers
// -----------------------------------------------------------------------------

func validationSite(
	siteID string,
	floors ...ValidationFloor,
) *ValidationSite {
	return &ValidationSite{
		SiteID: siteID,
		Floors: floors,
	}
}

func validationFloor(
	id string,
	regions []ValidationRegion,
	readers ...ValidationReader,
) ValidationFloor {
	return ValidationFloor{
		ID:      id,
		Regions: regions,
		Readers: readers,
	}
}

func validationRegion(
	id string,
	children ...ValidationRegion,
) ValidationRegion {
	return ValidationRegion{
		ID:      id,
		Regions: children,
	}
}

func validationReader(
	id string,
	ports ...int,
) ValidationReader {
	return ValidationReader{
		ID:           id,
		AntennaPorts: ports,
	}
}

func assertCategory(
	t *testing.T,
	name string,
	actual ValidationCategoryResult,
	required int,
	matched int,
	passed bool,
) {
	t.Helper()

	if actual.Required != required {
		t.Errorf(
			"%s.Required: want %d, got %d",
			name,
			required,
			actual.Required,
		)
	}

	if actual.Matched != matched {
		t.Errorf(
			"%s.Matched: want %d, got %d",
			name,
			matched,
			actual.Matched,
		)
	}

	if actual.Passed != passed {
		t.Errorf(
			"%s.Passed: want %t, got %t",
			name,
			passed,
			actual.Passed,
		)
	}
}

func assertErrorContains(
	t *testing.T,
	result *ValidationResult,
	expected string,
) {
	t.Helper()

	for _, message := range result.Errors {
		if strings.Contains(message, expected) {
			return
		}
	}

	t.Errorf(
		"expected an error containing %q, got %v",
		expected,
		result.Errors,
	)
}

func findMismatch(
	t *testing.T,
	result *ValidationResult,
	predicate func(ValidationMismatch) bool,
) ValidationMismatch {
	t.Helper()

	for _, mismatch := range result.Mismatches {
		if predicate(mismatch) {
			return mismatch
		}
	}

	t.Fatalf(
		"expected matching validation mismatch, got %v",
		result.Mismatches,
	)

	return ValidationMismatch{}
}