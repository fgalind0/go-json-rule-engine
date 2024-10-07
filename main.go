package main

import (
	"errors"
	"fmt"
	"strings"
)

var operationsMap = map[string]func(string, string) bool{
	"StringEquals":      StringEquals,
	"StringNotEquals":   StringNotEquals,
	"StringContains":    StringContains,
	"StringNotContains": StringNotContains,
}

func main() {
	fmt.Println("it's working")
	// operationsMap["Rando"] = StringNotContains
	// operationsLoop := []string{"StringEquals", "StringNotEquals", "StringContains", "StringNotContains", "souldntExist"}
	nameContainsCondition := StringCondition{
		Subject:  "tom",
		Operator: "StringContains",
		Value:    "fernando",
	}
	searchContainsCondition := StringCondition{
		Subject:  "brox",
		Operator: "StringContains",
		Value:    "NYC,big apple, manhattan",
	}
	cityEqualsCondition := StringCondition{
		Subject:  "nyc",
		Operator: "StringEquals",
		Value:    "New York",
	}
	stateEqualsCondition := StringCondition{
		Subject:  "NJ",
		Operator: "StringEquals",
		Value:    "NY",
	}
	andBlock:= AndGroup{
		Conditions: []Condition{&stateEqualsCondition, &cityEqualsCondition},
	}

	orBlock := OrGroup{
		Conditions: []Condition{&nameContainsCondition, &searchContainsCondition},
		NestedGroups: []*AndGroup{&andBlock},
	}

	rule := Rule{
		ConditionGroup: &orBlock,
	}

	result, err := rule.ConditionGroup.Resolve()
	fmt.Printf("\nthe rule evaluated to \n%v\nand had the following error %v\n",result,err)
	
	// un commenting both below will result to true
	// cityEqualsCondition.Subject = "New York"
	// stateEqualsCondition.Subject = "NY"

	//uncommenting next line result to true
	// nameContainsCondition.Subject = "fern"

	//uncommenting next line result to true
	// searchContainsCondition.Subject = "big apple"
	result, err = rule.ConditionGroup.Resolve()
	fmt.Printf("\nthe rule evaluated to \n%v\nand had the following error %v\n",result,err)
	
}

type Condition interface {
	Check() (bool, error)
}

type StringCondition struct {
	Subject  string // this would dynamically change on each rule evaluation
	Operator string
	Value    string
	Evaluate func(subject, value string) bool
}

type Rule struct {
	ConditionGroup
}

type ConditionGroup interface {
	Resolve() (bool, error)
}

type AndGroup struct {
	Conditions   []Condition
	NestedGroups []*OrGroup
}

type OrGroup struct {
	Conditions   []Condition
	NestedGroups []*AndGroup
}

// and groups require all to be true. At first false can return false
func (ag *AndGroup) Resolve() (bool, error) {
	// checking individual conditions will generally be faster than going into nested group
	// first verify that this and group has conditions to check
	if ag.Conditions != nil {
		for _, condition := range ag.Conditions {
			result, err := condition.Check()
			if !result || err != nil {
				return result, err
			}
		}
	}
	// only need to call resolve on nested group if it exists
	if ag.NestedGroups != nil {
		for _, group := range ag.NestedGroups {
			groupResult, err := group.Resolve()
			// because this is "AND" group any false makes whole thing false
			// also any errors will also default bool to false for saftey
			if !groupResult || err != nil {
				return groupResult, err
			}
		}
	}
	// if we make it out of looping over all conditions and groups and are here then
	// we did not get a false so can return true and nil
	return true, nil
}

// Or group means first true can return
func (og *OrGroup) Resolve() (bool, error) {
	// generally speaking checking condition first is fastest
	if og.Conditions != nil {
		for _, condition := range og.Conditions {
			result, err := condition.Check()
			if result {
				// true case
				return result, err
			} else if err != nil {
				// error so return false and err
				return result, err
			}
		}
	}
	// check NestedGroups
	if og.NestedGroups != nil {
		for _, group := range og.NestedGroups {
			groupResult, err := group.Resolve()
			if groupResult {
				// true case
				return groupResult, err
			} else if err != nil {
				// error so return false and err
				return groupResult, err
			}
		}
	}
	// if we have checked all possible things and did not get true 
	// then we will return false and nil
	return false, nil
}

// to do comeback and move operation to a func
func (c *StringCondition) Check() (bool, error) {
	op, ok := operationsMap[c.Operator]
	if ok {
		c.Evaluate = op
		return c.Run()
	}
	fmt.Println("Operation not found.")
	return false, errors.New("operation not found")
}
func (c *StringCondition) Run() (bool, error) {
	result := c.Evaluate(c.Subject, c.Value)
	fmt.Printf("The following results are \nValue: %v\nSubject: %v\nOperation: %v\nResult: %v\n\n", c.Value, c.Subject, c.Operator, result)
	return result, nil
}
func StringEquals(subject, value string) bool {
	return subject == value
}

func StringNotEquals(subject, value string) bool {
	return subject != value
}

func StringContains(subject, value string) bool {
	return strings.Contains(value, subject)
}
func StringNotContains(subject, value string) bool {
	return !strings.Contains(value, subject)
}
