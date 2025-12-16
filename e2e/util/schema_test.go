package util

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSetResourceDataFromMap_EmptyMap(t *testing.T) {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	r := schema.Resource{Schema: s}
	d := r.TestResourceData()

	err := SetResourceDataFromMap(d, map[string]interface{}{})
	if err != nil {
		t.Errorf("SetResourceDataFromMap with empty map should not error, got: %v", err)
	}
}

func TestSetResourceDataFromMap_SingleStringValue(t *testing.T) {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	r := schema.Resource{Schema: s}
	d := r.TestResourceData()

	m := map[string]interface{}{
		"name": "test-name",
	}

	err := SetResourceDataFromMap(d, m)
	if err != nil {
		t.Errorf("SetResourceDataFromMap failed: %v", err)
	}

	if d.Get("name").(string) != "test-name" {
		t.Errorf("Expected name to be 'test-name', got: %v", d.Get("name"))
	}
}

func TestSetResourceDataFromMap_MultipleValues(t *testing.T) {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
	r := schema.Resource{Schema: s}
	d := r.TestResourceData()

	m := map[string]interface{}{
		"name":    "test-resource",
		"count":   42,
		"enabled": true,
	}

	err := SetResourceDataFromMap(d, m)
	if err != nil {
		t.Errorf("SetResourceDataFromMap failed: %v", err)
	}

	if d.Get("name").(string) != "test-resource" {
		t.Errorf("Expected name to be 'test-resource', got: %v", d.Get("name"))
	}
	if d.Get("count").(int) != 42 {
		t.Errorf("Expected count to be 42, got: %v", d.Get("count"))
	}
	if d.Get("enabled").(bool) != true {
		t.Errorf("Expected enabled to be true, got: %v", d.Get("enabled"))
	}
}

func TestSetResourceDataFromMap_InvalidAttribute(t *testing.T) {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	r := schema.Resource{Schema: s}
	d := r.TestResourceData()

	m := map[string]interface{}{
		"invalid_key": "some-value",
	}

	err := SetResourceDataFromMap(d, m)
	if err == nil {
		t.Error("Expected error when setting invalid attribute, got nil")
	}
	if !strings.Contains(err.Error(), "invalid_key") {
		t.Errorf("Expected error message to contain 'invalid_key', got: %v", err)
	}
}

func TestSetResourceDataFromMap_TypeMismatch(t *testing.T) {
	s := map[string]*schema.Schema{
		"count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
	r := schema.Resource{Schema: s}
	d := r.TestResourceData()

	m := map[string]interface{}{
		"count": "not-a-number", // String instead of int
	}

	err := SetResourceDataFromMap(d, m)
	if err == nil {
		t.Error("Expected error when setting wrong type, got nil")
	}
	if !strings.Contains(err.Error(), "count") {
		t.Errorf("Expected error message to contain 'count', got: %v", err)
	}
}

func TestSetResourceDataFromMap_NilValue(t *testing.T) {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	r := schema.Resource{Schema: s}
	d := r.TestResourceData()

	m := map[string]interface{}{
		"name": nil,
	}

	err := SetResourceDataFromMap(d, m)
	if err != nil {
		t.Errorf("SetResourceDataFromMap with nil value should not error, got: %v", err)
	}
}

func TestSetResourceDataFromMap_ListType(t *testing.T) {
	s := map[string]*schema.Schema{
		"tags": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
	r := schema.Resource{Schema: s}
	d := r.TestResourceData()

	m := map[string]interface{}{
		"tags": []interface{}{"tag1", "tag2", "tag3"},
	}

	err := SetResourceDataFromMap(d, m)
	if err != nil {
		t.Errorf("SetResourceDataFromMap with list failed: %v", err)
	}

	tags := d.Get("tags").([]interface{})
	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got: %d", len(tags))
	}
}

func TestSetResourceDataFromMap_SetType(t *testing.T) {
	s := map[string]*schema.Schema{
		"ids": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
		},
	}
	r := schema.Resource{Schema: s}
	d := r.TestResourceData()

	m := map[string]interface{}{
		"ids": schema.NewSet(schema.HashInt, []interface{}{1, 2, 3}),
	}

	err := SetResourceDataFromMap(d, m)
	if err != nil {
		t.Errorf("SetResourceDataFromMap with set failed: %v", err)
	}

	ids := d.Get("ids").(*schema.Set)
	if ids.Len() != 3 {
		t.Errorf("Expected 3 ids in set, got: %d", ids.Len())
	}
}

func TestSetResourceDataFromMap_MapType(t *testing.T) {
	s := map[string]*schema.Schema{
		"metadata": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
	r := schema.Resource{Schema: s}
	d := r.TestResourceData()

	m := map[string]interface{}{
		"metadata": map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}

	err := SetResourceDataFromMap(d, m)
	if err != nil {
		t.Errorf("SetResourceDataFromMap with map failed: %v", err)
	}

	metadata := d.Get("metadata").(map[string]interface{})
	if len(metadata) != 2 {
		t.Errorf("Expected 2 metadata entries, got: %d", len(metadata))
	}
	if metadata["key1"] != "value1" {
		t.Errorf("Expected metadata['key1'] to be 'value1', got: %v", metadata["key1"])
	}
}

func TestSetResourceDataFromMap_PartialSuccess(t *testing.T) {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
	r := schema.Resource{Schema: s}
	d := r.TestResourceData()

	// First key is valid, second is invalid
	m := map[string]interface{}{
		"name":        "test",
		"invalid_key": "value",
	}

	err := SetResourceDataFromMap(d, m)
	if err == nil {
		t.Error("Expected error when one key is invalid")
	}

	// Note: depending on map iteration order, "name" may or may not be set
	// This is acceptable behavior - the function returns error on first failure
}

func TestSetResourceDataFromMap_ZeroValues(t *testing.T) {
	s := map[string]*schema.Schema{
		"count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	r := schema.Resource{Schema: s}
	d := r.TestResourceData()

	m := map[string]interface{}{
		"count":   0,
		"enabled": false,
		"name":    "",
	}

	err := SetResourceDataFromMap(d, m)
	if err != nil {
		t.Errorf("SetResourceDataFromMap with zero values failed: %v", err)
	}

	if d.Get("count").(int) != 0 {
		t.Errorf("Expected count to be 0, got: %v", d.Get("count"))
	}
	if d.Get("enabled").(bool) != false {
		t.Errorf("Expected enabled to be false, got: %v", d.Get("enabled"))
	}
	if d.Get("name").(string) != "" {
		t.Errorf("Expected name to be empty string, got: %v", d.Get("name"))
	}
}
