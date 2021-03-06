/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package astmodel

import (
	"go/ast"
	"go/token"
)

// PackageImport represents an import of a name from a package
type PackageImport struct {
	PackageReference PackageReference
	name             string
}

// NewPackageImport creates a new package import from a reference
func NewPackageImport(packageReference PackageReference) PackageImport {
	return PackageImport{
		PackageReference: packageReference,
	}
}

// WithName creates a new package reference with a friendly name
func (pi PackageImport) WithName(name string) PackageImport {
	pi.name = name
	return pi
}

func (pi PackageImport) AsImportSpec() *ast.ImportSpec {
	var name *ast.Ident
	if pi.name != "" {
		name = ast.NewIdent(pi.name)
	}

	return &ast.ImportSpec{
		Name: name,
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"" + pi.PackageReference.PackagePath() + "\"",
		},
	}
}

// PackageName is the package name of the package reference
func (pi PackageImport) PackageName() string {
	if pi.name != "" {
		return pi.name
	}

	return pi.PackageReference.PackageName()
}

// Equals returns true if the passed package reference references the same package, false otherwise
func (pi PackageImport) Equals(ref PackageImport) bool {
	packagesEqual := pi.PackageReference.Equals(ref.PackageReference)
	namesEqual := pi.name == ref.name

	return packagesEqual && namesEqual
}
