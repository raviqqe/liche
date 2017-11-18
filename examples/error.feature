Feature: Error
  Scenario: Fail with no argument
    Given a file named "foo.md" with ""
    When I run `linkcheck`
    Then the exit status should be 1

  Scenario: Fail with a non-existent file
    Given a file named "foo.md" with ""
    When I run `linkcheck bar.md`
    Then the exit status should be 1
    And the stderr should contain "no such file"
