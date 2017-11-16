Feature: Markdown
  Scenario: Check an empty markdown file
    Given a file named "foo.md" with ""
    When I successfully run `linkcheck foo.md`
    Then the stdout should contain exactly ""
