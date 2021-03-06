/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package astmodel

import (
	"testing"

	. "github.com/onsi/gomega"
)

/*
 * Shared Test Values
 */

var (
	fullName   = NewPropertyDefinition("FullName", "full-name", StringType)
	familyName = NewPropertyDefinition("FamilyName", "family-name", StringType)
	knownAs    = NewPropertyDefinition("KnownAs", "known-as", StringType)
	gender     = NewPropertyDefinition("Gender", "gender", StringType)
)

/*
 * NewObjectType() tests
 */

func TestNewObjectType_ReturnsEmptyType(t *testing.T) {
	g := NewGomegaWithT(t)

	st := NewObjectType()

	g.Expect(st.properties).To(HaveLen(0))
	g.Expect(st.functions).To(HaveLen(0))
}

/*
 * Properties() tests
 */

func Test_Properties_GivenEmptyObject_ReturnsEmptySlice(t *testing.T) {
	g := NewGomegaWithT(t)
	properties := EmptyObjectType.Properties()
	g.Expect(properties).To(HaveLen(0))
}

func Test_Properties_GivenObjectWithProperties_ReturnsExpectedSortedSlice(t *testing.T) {
	g := NewGomegaWithT(t)
	object := EmptyObjectType.WithProperties(fullName, familyName, knownAs, gender)
	properties := object.Properties()
	g.Expect(properties).To(HaveLen(4))
	g.Expect(properties[0]).To(Equal(familyName))
	g.Expect(properties[1]).To(Equal(fullName))
	g.Expect(properties[2]).To(Equal(gender))
	g.Expect(properties[3]).To(Equal(knownAs))
}

/*
 * Equals() tests
 */

func TestObjectType_Equals_WhenGivenType_ReturnsExpectedResult(t *testing.T) {

	clanName := NewStringPropertyDefinition("Clan")

	personType := NewObjectType().WithProperties(fullName, familyName, knownAs)
	otherPersonType := NewObjectType().WithProperties(fullName, familyName, knownAs)
	reorderedType := NewObjectType().WithProperties(knownAs, familyName, fullName)
	shorterType := NewObjectType().WithProperties(knownAs, fullName)
	longerType := NewObjectType().WithProperties(fullName, familyName, knownAs, gender)
	differentType := NewObjectType().WithProperties(fullName, clanName, knownAs, gender)
	mapType := NewMapType(StringType, personType)

	cases := []struct {
		name      string
		thisType  Type
		otherType Type
		expected  bool
	}{
		// Expect equal to self
		{"Equal to self", personType, personType, true},
		{"Equal to self", otherPersonType, otherPersonType, true},
		// Expect equal to same
		{"Equal to same", personType, otherPersonType, true},
		{"Equal to same", otherPersonType, personType, true},
		// Expect equal when properties are reordered
		{"Equal when properties reordered", personType, reorderedType, true},
		{"Equal when properties reordered", reorderedType, personType, true},
		// Expect not-equal when properties missing
		{"Not-equal when properties missing", personType, shorterType, false},
		{"Not-equal when properties missing", longerType, personType, false},
		// Expect not-equal when properties added
		{"Not-equal when properties added", personType, longerType, false},
		{"Not-equal when properties added", shorterType, personType, false},
		// Expect not-equal for different type
		{"Not-equal when different type", personType, mapType, false},
		{"Not-equal when different type", mapType, personType, false},
		// Expect not-equal for different property (but same property count)
		{"Not-equal when different type", personType, differentType, false},
		{"Not-equal when different type", differentType, personType, false},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			g := NewGomegaWithT(t)

			areEqual := c.thisType.Equals(c.otherType)

			g.Expect(areEqual).To(Equal(c.expected))
		})
	}
}

/*
 * WithProperty() Tests
 */

func Test_WithProperty_GivenEmptyObject_ReturnsPopulatedObject(t *testing.T) {
	g := NewGomegaWithT(t)
	empty := EmptyObjectType
	object := empty.WithProperty(fullName)
	g.Expect(empty).NotTo(Equal(object)) // Ensure the original wasn't modified
	g.Expect(object.properties).To(HaveLen(1))
	g.Expect(object.Properties()[0]).To(Equal(fullName))
}

func Test_WithProperty_GivenPopulatedObjectAndNewProperty_ReturnsExpectedObject(t *testing.T) {
	g := NewGomegaWithT(t)

	nickname := NewPropertyDefinition("Nickname", "nick-name", StringType)

	original := EmptyObjectType.WithProperties(fullName, familyName, knownAs, gender)
	modified := original.WithProperty(nickname)

	g.Expect(original).NotTo(Equal(modified))
	g.Expect(modified.Properties()).To(ContainElement(fullName))
	g.Expect(modified.Properties()).To(ContainElement(familyName))
	g.Expect(modified.Properties()).To(ContainElement(knownAs))
	g.Expect(modified.Properties()).To(ContainElement(gender))
	g.Expect(modified.Properties()).To(ContainElement(nickname))
}

func Test_WithProperty_GivenPopulatedObjectAndReplacementProperty_ReturnsExpectedObject(t *testing.T) {
	g := NewGomegaWithT(t)

	// Replacement for string typed 'Gender' field
	genderEnum := NewPropertyDefinition("Gender", "gender", &EnumType{})

	original := EmptyObjectType.WithProperties(fullName, familyName, knownAs, gender)
	modified := original.WithProperty(genderEnum)

	g.Expect(original).NotTo(Equal(modified))
	g.Expect(modified.Properties()).To(ContainElement(fullName))
	g.Expect(modified.Properties()).To(ContainElement(familyName))
	g.Expect(modified.Properties()).To(ContainElement(knownAs))
	g.Expect(modified.Properties()).To(ContainElement(genderEnum))
	g.Expect(modified.Properties()).NotTo(ContainElement(gender))
}

/*
 * WithProperties() Tests
 */

func Test_WithProperties_GivenEmptyObject_ReturnsPopulatedObject(t *testing.T) {
	g := NewGomegaWithT(t)
	original := EmptyObjectType
	modified := original.WithProperties(fullName, familyName, knownAs, gender)

	g.Expect(original).NotTo(Equal(modified)) // Ensure the original wasn't modified
	g.Expect(modified.Properties()).To(HaveLen(4))
	g.Expect(modified.Properties()).To(ContainElement(fullName))
	g.Expect(modified.Properties()).To(ContainElement(familyName))
	g.Expect(modified.Properties()).To(ContainElement(knownAs))
	g.Expect(modified.Properties()).To(ContainElement(gender))
}

/*
 * WithoutProperty() Tests
 */

func Test_ObjectWithoutProperty_ReturnsExpectedObject(t *testing.T) {
	g := NewGomegaWithT(t)

	original := EmptyObjectType.WithProperties(fullName, familyName, knownAs, gender)
	modified := original.WithoutProperty(fullName.propertyName)

	g.Expect(original).NotTo(Equal(modified))
	g.Expect(modified.Properties()).NotTo(ContainElement(fullName))
	g.Expect(modified.Properties()).To(ContainElement(familyName))
	g.Expect(modified.Properties()).To(ContainElement(knownAs))
	g.Expect(modified.Properties()).To(ContainElement(gender))
}

/*
 * WithFunction() tests
 */

func Test_WithFunction_GivenEmptyObject_ReturnsPopulatedObject(t *testing.T) {
	g := NewGomegaWithT(t)
	empty := EmptyObjectType
	fn := &FakeFunction{}
	object := empty.WithFunction("Activate", fn)
	g.Expect(empty).NotTo(Equal(object)) // Ensure the original wasn't modified
	g.Expect(object.functions).To(HaveLen(1))
	g.Expect(object.functions["Activate"].Equals(fn)).To(BeTrue())
}
