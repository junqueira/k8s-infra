/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package astmodel

import (
	"fmt"
	"go/ast"
)

// Type represents something that is a Go type
type Type interface {
	// RequiredImports returns a list of packages required by this type
	RequiredImports() []PackageReference

	// References returns the names of all types that this type
	// references. For example, an Array of Persons references a
	// Person.
	References() TypeNameSet

	// AsType renders as a Go abstract syntax tree for a type
	// (yes this says ast.Expr but that is what the Go 'ast' package uses for types)
	AsType(codeGenerationContext *CodeGenerationContext) ast.Expr

	// AsDeclarations renders as a Go abstract syntax tree for a declaration
	AsDeclarations(codeGenerationContext *CodeGenerationContext, name TypeName, description []string) []ast.Decl

	// Equals returns true if the passed type is the same as this one, false otherwise
	Equals(t Type) bool
}

// Types is the set of all types being generated
type Types map[TypeName]TypeDefinition

// Add adds a type to the set, with safety check that it has not already been defined
func (types Types) Add(def TypeDefinition) {
	key := def.Name()
	if _, ok := types[key]; ok {
		panic(fmt.Sprintf("type already defined: %v", key))
	}

	types[key] = def
}

// AddAll adds multiple types to the set, with the same safety check as Add()
func (types Types) AddAll(otherTypes Types) {
	for _, t := range otherTypes {
		types.Add(t)
	}
}

// Where returns a new set of types including only those that satisfy the predicate
func (types Types) Where(predicate func(definition TypeDefinition) bool) Types {
	result := make(Types, len(types))
	for _, t := range types {
		if predicate(t) {
			result[t.Name()] = t
		}
	}

	return result
}

// Contains returns true if the set contains the passed definition
func (types Types) Contains(def TypeDefinition) bool {
	_, ok := types[def.Name()]
	return ok
}

// TypeEquals decides if the types are the same and handles the `nil` case
func TypeEquals(left, right Type) bool {
	if left == nil {
		return right == nil
	}

	return left.Equals(right)
}

// Copy makes an independent copy of this set of types
func (types Types) Copy() Types {
	result := make(Types)
	result.AddAll(types)
	return result
}
