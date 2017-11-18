Feature: Markdown
  Scenario: Check an empty markdown file
    Given a file named "foo.md" with ""
    When I successfully run `linkcheck foo.md`
    Then the stdout should contain exactly ""

  Scenario: Check a markdown file
    Given a file named "foo.md" with:
    """
    # Title

    ## Section

    - List 1
        - Item 1
        - Item 2
    - List 2

    ```
    #!/bin/sh

    echo Hello, world!
    ```
    """
    When I successfully run `linkcheck foo.md`
    Then the stdout should contain exactly ""

  Scenario: Check a markdown file which contains a live link
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    """
    When I successfully run `linkcheck foo.md`
    Then the stdout should contain exactly ""

  Scenario: Check a markdown file which contains a dead link
    Given a file named "foo.md" with:
    """
    [The answer](https://some-say-the-answer-is-42.com)
    """
    When I run `linkcheck foo.md`
    Then the exit status should be 1
    And the stderr should contain "ERROR"

  Scenario: Check a markdown file which contains live and dead links
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    [The answer](https://some-say-the-answer-is-42.com)
    """
    When I run `linkcheck foo.md`
    Then the exit status should be 1
    And the stderr should not contain "OK"
    And the stderr should contain "ERROR"

  Scenario: Check a markdown file which contains a live link in verbose mode
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    """
    When I successfully run `linkcheck -v foo.md`
    Then the stderr should contain "OK"

  Scenario: Check a markdown file which contains a live link in verbose mode with a long option
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    """
    When I successfully run `linkcheck --verbose foo.md`
    Then the stderr should contain "OK"

  Scenario: Check a markdown file which contains live and dead links in verbose mode
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    [The answer](https://some-say-the-answer-is-42.com)
    """
    When I run `linkcheck -v foo.md`
    Then the exit status should be 1
    And the stderr should contain "OK"
    And the stderr should contain "ERROR"

  Scenario: Check 2 markdown files
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    """
    And a file named "bar.md" with:
    """
    [Yahoo](https://yahoo.com)
    """
    When I successfully run `linkcheck foo.md bar.md`
    Then the stdout should contain exactly ""

  Scenario: Check 2 markdown files
    Given a file named "foo.md" with:
    """
    [![Circle CI](https://img.shields.io/circleci/project/github/raviqqe/linkcheck.svg?style=flat-square)](https://circleci.com/gh/raviqqe/linkcheck)
    [![Go Report Card](https://goreportcard.com/badge/github.com/raviqqe/linkcheck?style=flat-square)](https://goreportcard.com/report/github.com/raviqqe/linkcheck)
    """
    And a file named "foo.sh" with:
    """
    linkcheck -v foo.md 2>&1 | wc -l
    """
    When I successfully run `sh foo.sh`
    Then the stdout should contain exactly "5"
