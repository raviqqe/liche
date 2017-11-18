Feature: Error
  Scenario: Fail with no argument
    Given a file named "foo.md" with ""
    When I run `liche`
    Then the exit status should be 1

  Scenario: Fail with a non-existent file
    Given a file named "foo.md" with ""
    When I run `liche bar.md`
    Then the exit status should be 1
    And the stderr should contain "no such file"
