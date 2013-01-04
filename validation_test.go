package validation

import (
  //"fmt"
  "reflect"
  "testing"
)

type RuleTest struct {
  rule  Rule
  in    string
  valid bool
}

var ruleTests = []RuleTest{
  {NotEmpty, "", false},
  {NotEmpty, "test", true},
  {Url, "http://www.google.com/", true},
  {Url, "ftp://john%20doe@www.google.com/", true},
  {Url, "http://www.google.com/?q=go+language", true},
  {Url, "notaurl", false},
  {Url, "www.sortofaurlbutno.com", false},
  {Url, "stillnoturl.com", false},
  {ObjectId, "507f1f77bcf86cd799439011", true},
  {ObjectId, "507f1f77bcf86cd79939011", false},
  {Alpha, "fj3j345kj345kj34", true},
  {Alpha, "_false", false},
  {Email, "test@test.com", true},
  {Email, "test_test@test.com", true},
  {Email, "test.test@test.com", true},
  {Email, "Test_test.com", false},
  {Email, "test%40test.com", false},
  {Numeric, "43545345345", true},
  {Numeric, "a23", false},
  {Numeric, "forty", false},
  {Numeric, "5.0", false},
  {ZipCode, "33145", true},
  {ZipCode, "44245-3456", false},
  {ZipCode, "331456", false},
  {Sha1, "8616fa4f0990c2d6bdd0cfb00789c3e47b9f65d6", true},
  {Sha1, "c3a44c7e06e461adeb0bff1c6c009e4d07d14875", true},
  {Sha1, "z3a44c7e06e461adeb0bff1c6c009e4d07d14875", false},
  {Sha1, "ababcbabcabcbac", false},

  {EqualsAny([]string{"male", "female"}), "male", true},
  {EqualsAny([]string{"male", "female"}), "female", true},
  {EqualsAny([]string{"male", "female"}), "unknown", false},
  {EqualsAny([]string{"male", "female"}), "femal", false},
  {Date("2006-01-02"), "1998-04-20", true},
  {Date("2006-01-02"), "0000-01-01", true},
  {Date("2006-01-02"), "0000-00-00", false},
  {Date("2006-01-02"), "19980420", false},
  {Date("Monday, 02-Jan-06 15:04:05 MST"), "Tuesday, 04-Feb-08 15:05:10 EST", true},
  {Date("Monday, 02-Jan-06 15:04:05 MST"), "Tuesday 04-Feb-08 15:05:10 EST", false},
}

func TestRuleFuncs(t *testing.T) {
  for i, tt := range ruleTests {
    err := tt.rule(tt.in)
    if (err == nil) != tt.valid {
      tot := reflect.TypeOf(tt.rule)
      if tt.valid {
        t.Errorf("%d. %v(%q) => %v. Expected no error.", i, tot.Name(), tt.in, err)
      } else {
        t.Errorf("%d. %v(%q) => %v. Expected error.", i, tot.Name(), tt.in, err)
      }
    }
  }
}

func TestAdd(t *testing.T) {
  v := New()
  v.Add("test", Numeric, "message")

  if _, ok := v.Map["test"]; ok != true {
    t.Errorf("Add(): could not find added rule.")
  }

  if _, ok := v.Map["test2"]; ok != false {
    t.Errorf("Add(): found a rule that wasn't added.")
  }

  if len(v.Map["test"]) != 1 {
    t.Fatalf("Add(): could not find the constraint in the map")
  }

  // Can't directly test if Rule is equal to Numeric in go, since reflection
  // cannot compare functions or tell us their original name.
  // v.Map["test"][0].Rule
  error1 := v.Map["test"][0].Rule("notnumeric")
  error2 := Numeric("notnumeric")
  if (error1 == nil || error2 == nil) || error1.Error() != error2.Error() {
    t.Errorf("Add(): the rule in the constraint doesn't seem to be the one we added.")
  }

  if v.Map["test"][0].Message.Error() != "message" {
    t.Errorf("Add(): found an incorrect error message")
  }

  v.Add("test", Numeric, "message")
  if len(v.Map["test"]) != 2 {
    t.Errorf("Add(): did not properly append a constraint to the same name")
  }
}

func TestValidateKeyValue(t *testing.T) {
  v := New()
  v.Add("test", Numeric, "")
  v.Add("test", NotEmpty, "")

  // should succeed
  valid, _ := v.ValidateKeyValue("test", "123")
  if !valid {
    t.Errorf("ValidateKeyValue(): returned false when expected true.")
  }

  // should error
  valid, messages := v.ValidateKeyValue("test", "fdgd")
  if valid {
    t.Errorf("ValidateKeyValue(): returned true when expected false.")
  }
  if len(messages) != 1 {
    t.Fatalf("ValidateKeyValue(): could not find the test error in messages.")
  }
  if messages[0] != Numeric("dfsdf").Error() {
    t.Errorf("ValidateKeyValue(): returned the wrong error message")
  }
}
