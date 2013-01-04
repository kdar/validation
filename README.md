A very simple validation framework that is not dependent on any other web/http framework.

Example
-------

    package main

    import (
      "fmt"
      "github.com/kdar/validation"
    )

    func main() {
      v := validation.New()
      v.Add("test", validation.Numeric, "Test was not a numeric.")

      data := map[string]string{
        "test": "5",
      }

      valid, messages := v.Validate(data)
      fmt.Printf("%#v: %#v\n", data, valid)
      fmt.Println(messages)

      data["test"] = "hello"

      valid, messages = v.Validate(data)
      fmt.Printf("%#v: %#v\n", data, valid)
      fmt.Println(messages)
    }
