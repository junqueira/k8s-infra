/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package jsonast

import (
	"net/url"
)

// Schema abstracts over the exact implementation of
// JSON schema we are using as we need to be able to
// process inputs provided by both 'gojsonschema' and
// 'go-openapi'. It is possible we could switch to
// 'go-openapi' for everything because it could handle
// both JSON Schema and Swagger; but this is not yet done.
type Schema interface {
	url() *url.URL
	title() *string
	description() *string

	hasType(schemaType SchemaType) bool

	// complex things
	hasOneOf() bool
	oneOf() []Schema

	hasAnyOf() bool
	anyOf() []Schema

	hasAllOf() bool
	allOf() []Schema

	// enum things
	enumValues() []string

	// array things
	items() []Schema

	// object things
	requiredProperties() []string
	properties() map[string]Schema
	additionalPropertiesAllowed() bool
	additionalPropertiesSchema() Schema

	// ref things
	isRef() bool
	refGroupName() (string, error)
	refObjectName() (string, error)
	refVersion() (string, error)
	refSchema() Schema
}

// SchemaType defines the type of JSON schema node we are currently processing
type SchemaType string

// Definitions for different kinds of JSON schema
const (
	AnyOf   SchemaType = "anyOf"
	AllOf   SchemaType = "allOf"
	OneOf   SchemaType = "oneOf"
	Ref     SchemaType = "ref"
	Array   SchemaType = "array"
	Bool    SchemaType = "boolean"
	Int     SchemaType = "integer"
	Number  SchemaType = "number"
	Object  SchemaType = "object"
	String  SchemaType = "string"
	Enum    SchemaType = "enum"
	Unknown SchemaType = "unknown"
)
