// Simple validation that isn't dependent on any framework.
// Idea and code taken from the Tango project: https://github.com/Astrata/tango/
// Modified for my purposes.
package validation

import (
  "errors"
  "fmt"
  "reflect"
  "regexp"
  "strings"
  "time"
  //"strconv"
)

// Validation function.
type Rule func(string) error

// A set of rules to be applied on a single variable.
type Constraint struct {
  Rule    Rule
  Message error
}

// A set of validation rules.
type Rules struct {
  Map map[string][]Constraint
}

// Returns a new set of validation rules.
func New() *Rules {
  self := &Rules{}
  self.Map = make(map[string][]Constraint)
  return self
}

// Adds a rule to a set of constraints.
// func (self *Constraint) Add(rule Rule) {
//   self.All = append(self.All, rule)
// }

// Validates a key/value pair. Returns wether it's valid
// and a slice of errors.
func (self *Rules) ValidateKeyValue(key, value string) (bool, []string) {
  errors := []string{}
  passed := true
  if constraints, ok := self.Map[key]; ok == true {
    for _, constraint := range constraints {
      test := constraint.Rule(value)
      if test != nil {
        passed = false
        error := test.Error()
        if constraint.Message != nil {
          error = constraint.Message.Error()
        }
        errors = append(errors, error)
      }
    }
  }

  return passed, errors
}

// Validates input data.
func (self *Rules) Validate(params map[string]string) (bool, map[string][]string) {
  valid := true
  messages := map[string][]string{}

  for key, _ := range params {
    value := params[key]
    passed, errors := self.ValidateKeyValue(key, value)
    if passed == false {
      messages[key] = errors
      valid = false
    }
  }

  return valid, messages
}

// Validates values in a structure.
// If a rule name has 
func (self *Rules) ValidateStruct(s interface{}) (bool, map[string][]string) {
  return self.validateStructWithKeyPrefix("", s)
}

// Validates values in a structure with a key prefix
// If a rule name has 
func (self *Rules) validateStructWithKeyPrefix(prefix string, s interface{}) (bool, map[string][]string) {
  valid := true
  messages := map[string][]string{}

  valueOfT := reflect.ValueOf(s)
  if valueOfT.Kind() == reflect.Ptr {
    valueOfT = valueOfT.Elem()
  }

  typeOfT := valueOfT.Type()
  for i := 0; i < valueOfT.NumField(); i++ {
    f := valueOfT.Field(i)
    key := prefix + typeOfT.Field(i).Name

    passed := true
    var errors []string
    switch f.Kind() {
    case reflect.String:
      passed, errors = self.ValidateKeyValue(key, f.Interface().(string))
      if passed == false {
        messages[key] = errors
        valid = false
      }
    case reflect.Struct:
      var errormap map[string][]string
      passed, errormap = self.validateStructWithKeyPrefix(key+".", f.Interface())
      if passed == false {
        for key, value := range errormap {
          messages[key] = value
        }
        valid = false
      }
    }
  }

  return valid, messages
}

// Adds a new rule
func (self *Rules) Add(name string, rule Rule, message string) {
  constraint := Constraint{Rule: rule}
  var constraints []Constraint
  var ok bool

  if constraints, ok = self.Map[name]; ok == false {
    constraints = []Constraint{}
  }

  if len(message) > 0 {
    constraint.Message = errors.New(message)
  }
  constraints = append(constraints, constraint)
  self.Map[name] = constraints
}

// Adds a new rule that is required
func (self *Rules) AddRequired(name string, rule Rule, message string) {
  self.Add(name, NotEmpty, "")
  self.Add(name, rule, message)
}

// A rule that returns error if the value is empty.
func NotEmpty(value string) error {
  if value == "" {
    return fmt.Errorf("This value is required")
  }
  return nil
}

// A rule that returns error if the value is not an URL.
func Url(value string) error {
  match := MatchExpr(value, `(?i)^([a-z]+:\/\/[a-z0-9][a-z0-9\-\.]*.+)?$`)
  if match == nil {
    return nil
  }
  return fmt.Errorf("Value must be an URL.")
}

// A rule that returns error if the value is not a BSON ObjectId.
func ObjectId(value string) error {
  match := MatchExpr(value, `^([a-f0-9]{24})?$`)
  if match == nil {
    return nil
  }
  return fmt.Errorf("Expecting an ObjectId.")
}

// A rule that returns error if the value is not a-zA-Z0-9.
func Alpha(value string) error {
  match := MatchExpr(value, `(?i)^([a-z0-9]+)?$`)
  if match == nil {
    return nil
  }
  return fmt.Errorf("Value must be a number or a letter from A to Z (case does not matter).")
}

// A rule that returns error if the value is not an e-mail.
func Email(value string) error {
  passed := MatchExpr(value, `(?i)^([a-z0-9][a-z0-9\.\-+_]*@[a-z0-9\-\.]+.[a-z]+)?$`)
  if passed == nil {
    return nil
  }
  return fmt.Errorf("Value must be an e-mail address.")
}

// A rule that returns error if the value is not numeric.
func Numeric(value string) error {
  passed := MatchExpr(value, `(?i)^([0-9]+)?$`)
  if passed == nil {
    return nil
  }
  return fmt.Errorf("Value must be a number.")
}

// A rule that returns error if the value is not a zip code.
func ZipCode(value string) error {
  passed := MatchExpr(value, `(?i)^([0-9]{5})?$`)
  if passed == nil {
    return nil
  }
  return fmt.Errorf("Value must be a zipcode (XXXXX).")
}

// A rule that returns a func that returns error if the value does not match 
// any of the strings passed in the slice.
func EqualsAny(any []string) Rule {
  return func(value string) error {
    for _, x := range any {
      if x == value {
        return nil
      }
    }

    return fmt.Errorf("Did not match the following: %s", strings.Join(any, ", "))
  }
}

// A rule that returns a func that returns error if the value does not match 
// the date format.
func Date(format string) Rule {
  return func(value string) error {
    _, err := time.Parse(format, value)
    return err
  }
}

// A rule that returns error if the value is not a SHA1 hash.
func Sha1(value string) error {
  passed := MatchExpr(value, `(?i)^([a-f0-9]{40})$`)
  if passed == nil {
    return nil
  }
  return fmt.Errorf("Value must be a SHA1 hash.")
}

// A rule that returns error if the value does not match a pattern.
func MatchExpr(value string, expr string) error {
  match, _ := regexp.MatchString(expr, value)
  if match == true {
    return nil
  }
  return fmt.Errorf("Value does not match pattern %s.", expr)
}
