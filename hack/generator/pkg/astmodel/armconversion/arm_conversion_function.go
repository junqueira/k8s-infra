/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package armconversion

import (
	"fmt"
	"github.com/Azure/k8s-infra/hack/generator/pkg/astmodel"
	"go/ast"
)

type ConversionDirection string

const (
	ConversionDirectionToArm   = ConversionDirection("ToArm")
	ConversionDirectionFromArm = ConversionDirection("FromArm")
)

// ArmConversionFunction represents an ARM conversion function for converting between a Kubernetes resource
// and an ARM resource.
type ArmConversionFunction struct {
	armTypeName astmodel.TypeName
	armType     *astmodel.ObjectType
	idFactory   astmodel.IdentifierFactory
	direction   ConversionDirection
	isResource  bool
}

var _ astmodel.Function = &ArmConversionFunction{}

// RequiredImports returns the imports required for this conversion function
func (c *ArmConversionFunction) RequiredImports() []astmodel.PackageReference {
	var result []astmodel.PackageReference

	// Because this interface is attached to an object, by definition that object will specify
	// its own required imports. We don't want to call the objects required imports here
	// because then we're in infinite recursion (object delegates to us, we delegate back to it)

	// We need these because we're going to be constructing/casting to the types
	// of the properties in the ARM object, so we need to import those.
	result = append(result, c.armType.RequiredImports()...)
	result = append(result, astmodel.MakeGenRuntimePackageReference())
	result = append(result, astmodel.MakePackageReference("fmt"))

	return result
}

// References this type has to the given type
func (c *ArmConversionFunction) References() astmodel.TypeNameSet {
	return c.armType.References()
}

// AsFunc returns the function as a Go AST
func (c *ArmConversionFunction) AsFunc(codeGenerationContext *astmodel.CodeGenerationContext, receiver astmodel.TypeName, methodName string) *ast.FuncDecl {
	switch c.direction {
	case ConversionDirectionToArm:
		return c.asConvertToArmFunc(codeGenerationContext, receiver, methodName)
	case ConversionDirectionFromArm:
		return c.asConvertFromArmFunc(codeGenerationContext, receiver, methodName)
	default:
		panic(fmt.Sprintf("Unknown conversion direction %s", c.direction))
	}
}

func (c *ArmConversionFunction) asConvertToArmFunc(
	codeGenerationContext *astmodel.CodeGenerationContext,
	receiver astmodel.TypeName,
	methodName string) *ast.FuncDecl {

	builder := newConvertToArmFunctionBuilder(
		c,
		codeGenerationContext,
		receiver,
		methodName)
	return builder.functionDeclaration()
}

func (c *ArmConversionFunction) asConvertFromArmFunc(
	codeGenerationContext *astmodel.CodeGenerationContext,
	receiver astmodel.TypeName,
	methodName string) *ast.FuncDecl {

	builder := newConvertFromArmFunctionBuilder(
		c,
		codeGenerationContext,
		receiver,
		methodName)
	return builder.functionDeclaration()
}

// Equals determines if this function is equal to the passed in function
func (c *ArmConversionFunction) Equals(other astmodel.Function) bool {
	// TODO: Equality on functions is currently awkward because we can't easily pass
	// TODO: a reference to the object the function is on to the function (since both
	// TODO: are immutable and it's impossible to have two immutable objects with
	// TODO: references to each other). Really this equality is always in the context
	// TODO: of comparing objects (since you can't have a free-floating function
	// TODO: with a receiver), and as such we just need to compare things that aren't
	// TODO: the receiver type.

	if o, ok := other.(*ArmConversionFunction); ok {
		return c.armType.Equals(o.armType) &&
			c.armTypeName.Equals(o.armTypeName) &&
			c.direction == o.direction
	}

	return false
}
