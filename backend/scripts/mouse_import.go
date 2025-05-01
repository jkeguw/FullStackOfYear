package scripts

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuration for the import script
type Config struct {
	MongoURI      string
	Database      string
	Collection    string
	CsvFile       string
	HasHeader     bool
	DryRun        bool
	TypeMap       map[string]string
	DefaultFields map[string]interface{}
}

// ColumnMapping defines how CSV columns map to MongoDB fields
type ColumnMapping struct {
	ColumnIndex int
	FieldPath   string
	Converter   func(string) (interface{}, error)
}

func ImportMouse() {
	// Command line arguments
	csvFile := flag.String("file", "", "Path to the CSV file to import")
	mongoURI := flag.String("mongo", "mongodb://localhost:27017", "MongoDB connection URI")
	database := flag.String("db", "mouse_db", "MongoDB database name")
	collection := flag.String("collection", "devices", "MongoDB collection name")
	hasHeader := flag.Bool("header", true, "Whether the CSV file has a header row")
	dryRun := flag.Bool("dry-run", false, "Perform a dry run without writing to the database")
	flag.Parse()

	if *csvFile == "" {
		log.Fatal("CSV file path is required. Use -file flag.")
	}

	// Create configuration
	config := &Config{
		MongoURI:   *mongoURI,
		Database:   *database,
		Collection: *collection,
		CsvFile:    *csvFile,
		HasHeader:  *hasHeader,
		DryRun:     *dryRun,
		// Define default types for common values
		TypeMap: map[string]string{
			"mouse":     "mouse",
			"keyboard":  "keyboard",
			"mousepad":  "mousepad",
			"monitor":   "monitor",
			"accessory": "accessory",
		},
		// Default fields to add to all documents
		DefaultFields: map[string]interface{}{
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
	}

	// Run the import
	if err := importMice(config); err != nil {
		log.Fatalf("Error importing mice: %v", err)
	}
}

func importMice(config *Config) error {
	// Open the CSV file
	file, err := os.Open(config.CsvFile)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields per record

	// Read header if it exists
	var header []string
	if config.HasHeader {
		header, err = reader.Read()
		if err != nil {
			return fmt.Errorf("error reading header: %v", err)
		}

		// Clean header names
		for i, h := range header {
			header[i] = strings.TrimSpace(h)
		}
	}

	// Create column mappings based on header or default mapping
	mappings := createColumnMappings(header)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var client *mongo.Client
	var collection *mongo.Collection

	if !config.DryRun {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
		if err != nil {
			return fmt.Errorf("failed to connect to MongoDB: %v", err)
		}
		defer client.Disconnect(ctx)

		// Check the connection
		err = client.Ping(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to ping MongoDB: %v", err)
		}

		collection = client.Database(config.Database).Collection(config.Collection)
	}

	// Process each row
	lineCount := 1
	if config.HasHeader {
		lineCount++
	}

	var importedCount, errorCount int

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errorCount++
			log.Printf("Warning: Error reading line %d: %v", lineCount, err)
			lineCount++
			continue
		}

		// Process the record
		doc, err := processRecord(record, mappings, config)
		if err != nil {
			errorCount++
			log.Printf("Warning: Error processing line %d: %v", lineCount, err)
			lineCount++
			continue
		}

		// Skip if doc is nil
		if doc == nil {
			lineCount++
			continue
		}

		// Print document in dry run mode
		if config.DryRun {
			fmt.Printf("Document from line %d: %+v\n", lineCount, doc)
		} else {
			// Insert the document
			_, err = collection.InsertOne(ctx, doc)
			if err != nil {
				errorCount++
				log.Printf("Warning: Error inserting document from line %d: %v", lineCount, err)
				lineCount++
				continue
			}
		}

		importedCount++
		lineCount++
	}

	// Print summary
	action := "Would import"
	if !config.DryRun {
		action = "Imported"
	}
	fmt.Printf("%s %d of %d records with %d errors\n", action, importedCount, lineCount-1, errorCount)

	return nil
}

func createColumnMappings(header []string) []ColumnMapping {
	mappings := make([]ColumnMapping, 0)

	if len(header) > 0 {
		// Map based on header
		for i, h := range header {
			// Skip empty headers
			if h == "" {
				continue
			}

			mapping := ColumnMapping{
				ColumnIndex: i,
				FieldPath:   mapHeaderToField(h),
				Converter:   getConverter(h),
			}
			mappings = append(mappings, mapping)
		}
	} else {
		// Default mapping if no header
		defaultMappings := []struct {
			index int
			field string
		}{
			{0, "name"},
			{1, "brand"},
			{2, "dimensions.length"},
			{3, "dimensions.width"},
			{4, "dimensions.height"},
			{5, "dimensions.weight"},
			{6, "shape.type"},
			{7, "technical.sensor"},
			{8, "technical.maxDPI"},
			{9, "technical.pollingRate"},
			{10, "connection_type"},
		}

		for _, m := range defaultMappings {
			mapping := ColumnMapping{
				ColumnIndex: m.index,
				FieldPath:   m.field,
				Converter:   getConverter(m.field),
			}
			mappings = append(mappings, mapping)
		}
	}

	return mappings
}

func mapHeaderToField(header string) string {
	// Clean and normalize header
	header = strings.TrimSpace(strings.ToLower(header))

	// Map common headers to fields
	fieldMap := map[string]string{
		"name":                 "name",
		"brand":                "brand",
		"length":               "dimensions.length",
		"length (mm)":          "dimensions.length",
		"width":                "dimensions.width",
		"width (mm)":           "dimensions.width",
		"height":               "dimensions.height",
		"height (mm)":          "dimensions.height",
		"weight":               "dimensions.weight",
		"weight (g)":           "dimensions.weight",
		"shape":                "shape.type",
		"hump placement":       "shape.humpPlacement",
		"front flare":          "shape.frontFlare",
		"side curvature":       "shape.sideCurvature",
		"hand compatibility":   "shape.handCompatibility",
		"thumb rest":           "shape.thumbRest",
		"ring finger rest":     "shape.ringFingerRest",
		"material":             "material",
		"connectivity":         "connection_type",
		"sensor":               "technical.sensor",
		"sensor technology":    "technical.sensorTechnology",
		"sensor position":      "technical.sensorPosition",
		"dpi":                  "technical.maxDPI",
		"polling rate":         "technical.pollingRate",
		"tracking speed":       "technical.trackingSpeed",
		"tracking speed (ips)": "technical.trackingSpeed",
		"acceleration":         "technical.acceleration",
		"acceleration (g)":     "technical.acceleration",
		"side buttons":         "technical.sideButtons",
		"middle buttons":       "technical.middleButtons",
	}

	if field, ok := fieldMap[header]; ok {
		return field
	}

	// If not in map, use the header itself with camelCase conversion
	return toCamelCase(header)
}

func toCamelCase(s string) string {
	// Replace non-alphanumeric with spaces
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return ' '
	}, s)

	// Split into words
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	// Convert to camelCase
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			words[i] = strings.ToUpper(words[i][0:1]) + strings.ToLower(words[i][1:])
		}
	}

	// Join and ensure first character is lowercase
	result := strings.Join(words, "")
	if len(result) > 0 {
		result = strings.ToLower(result[0:1]) + result[1:]
	}

	return result
}

func getConverter(field string) func(string) (interface{}, error) {
	// Determine converter based on field type
	field = strings.ToLower(field)

	// Numeric fields
	numericFields := []string{
		"dimensions.length", "length", "length (mm)",
		"dimensions.width", "width", "width (mm)",
		"dimensions.height", "height", "height (mm)",
		"dimensions.weight", "weight", "weight (g)",
		"technical.maxdpi", "dpi",
		"technical.pollingrate", "polling rate",
		"technical.trackingspeed", "tracking speed", "tracking speed (ips)",
		"technical.acceleration", "acceleration", "acceleration (g)",
		"technical.sidebuttons", "side buttons",
		"technical.middlebuttons", "middle buttons",
	}

	// Boolean fields
	booleanFields := []string{
		"shape.thumbrest", "thumb rest",
		"shape.ringfingerrest", "ring finger rest",
	}

	// Check for numeric fields
	for _, nf := range numericFields {
		if strings.Contains(field, nf) {
			return func(val string) (interface{}, error) {
				val = strings.TrimSpace(val)
				if val == "" {
					return 0.0, nil
				}

				// Remove any non-numeric characters except decimal point
				val = strings.Map(func(r rune) rune {
					if (r >= '0' && r <= '9') || r == '.' {
						return r
					}
					return -1
				}, val)

				// Parse as float64
				return strconv.ParseFloat(val, 64)
			}
		}
	}

	// Check for boolean fields
	for _, bf := range booleanFields {
		if strings.Contains(field, bf) {
			return func(val string) (interface{}, error) {
				val = strings.ToLower(strings.TrimSpace(val))
				return val == "yes" || val == "true" || val == "1", nil
			}
		}
	}

	// Default string converter
	return func(val string) (interface{}, error) {
		return strings.TrimSpace(val), nil
	}
}

func processRecord(record []string, mappings []ColumnMapping, config *Config) (bson.M, error) {
	// Create a new document
	doc := bson.M{}

	// Skip empty rows
	allEmpty := true
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			allEmpty = false
			break
		}
	}
	if allEmpty {
		return nil, nil
	}

	// Populate document according to mappings
	for _, mapping := range mappings {
		// Skip if column index is out of range
		if mapping.ColumnIndex >= len(record) {
			continue
		}

		// Get field value
		value := record[mapping.ColumnIndex]
		if value == "" {
			continue
		}

		// Convert value
		converted, err := mapping.Converter(value)
		if err != nil {
			return nil, fmt.Errorf("error converting field %s: %v", mapping.FieldPath, err)
		}

		// Skip empty strings
		if strVal, ok := converted.(string); ok && strVal == "" {
			continue
		}

		// Set field in document
		setNestedField(doc, mapping.FieldPath, converted)
	}

	// Validate and fix document
	if err := validateDocument(doc, config); err != nil {
		return nil, err
	}

	// Add default fields
	for k, v := range config.DefaultFields {
		if _, exists := doc[k]; !exists {
			doc[k] = v
		}
	}

	// Add ID if not present
	if _, exists := doc["_id"]; !exists {
		doc["_id"] = primitive.NewObjectID()
	}

	// Ensure type is always "mouse" for this data
	doc["type"] = "mouse"

	return doc, nil
}

func setNestedField(doc bson.M, path string, value interface{}) {
	parts := strings.Split(path, ".")

	// Handle simple case
	if len(parts) == 1 {
		doc[parts[0]] = value
		return
	}

	// Handle nested fields
	current := doc
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]

		// Check if this part exists already
		if _, exists := current[part]; !exists {
			current[part] = bson.M{}
		}

		// If it's not a map/document, make it one
		if subDoc, ok := current[part].(bson.M); ok {
			current = subDoc
		} else {
			// Replace with a map if it's not one
			newSubDoc := bson.M{}
			current[part] = newSubDoc
			current = newSubDoc
		}
	}

	// Set the final value
	current[parts[len(parts)-1]] = value
}

func validateDocument(doc bson.M, config *Config) error {
	// Check required fields
	if _, exists := doc["name"]; !exists || doc["name"] == "" {
		return fmt.Errorf("document is missing required field 'name'")
	}

	// Try to ensure brand exists
	if _, exists := doc["brand"]; !exists || doc["brand"] == "" {
		// Try to extract brand from name
		if name, ok := doc["name"].(string); ok {
			nameParts := strings.SplitN(name, " ", 2)
			if len(nameParts) > 0 {
				doc["brand"] = nameParts[0]
			} else {
				doc["brand"] = "Unknown"
			}
		} else {
			doc["brand"] = "Unknown"
		}
	}

	// Ensure dimensions exist
	if _, exists := doc["dimensions"]; !exists {
		doc["dimensions"] = bson.M{}
	}

	// Ensure technical specs exist
	if _, exists := doc["technical"]; !exists {
		doc["technical"] = bson.M{}
	}

	// Ensure shape info exists
	if _, exists := doc["shape"]; !exists {
		doc["shape"] = bson.M{
			"type": "unknown",
		}
	}

	return nil
}

// Helper function to update nested fields with specified separator
func getNestedValue(doc bson.M, path string) interface{} {
	parts := strings.Split(path, ".")

	var current interface{} = doc

	for _, part := range parts {
		switch v := current.(type) {
		case bson.M:
			var exists bool
			current, exists = v[part]
			if !exists {
				return nil
			}
		default:
			return nil
		}
	}

	return current
}
