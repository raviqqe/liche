Feature: Error
  Scenario: Fail with no argument
    Given a file named "foo.md" with ""
    When I run `linkcheck`
    Then the exit status should be 1
