Feature: Options
  Scenario: Check a live link in verbose mode
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    """
    When I successfully run `liche -v foo.md`
    Then the stderr should contain "OK"

  Scenario: Check a live link in verbose mode with a long option
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    """
    When I successfully run `liche --verbose foo.md`
    Then the stderr should contain "OK"

  Scenario: Check live and dead links in verbose mode
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    [The answer](https://some-say-the-answer-is-42.com)
    """
    When I run `liche -v foo.md`
    Then the exit status should be 1
    And the stderr should contain "OK"
    And the stderr should contain "ERROR"

  Scenario: Check a live link with timeout
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    """
    When I successfully run `liche --timeout 10 foo.md`
    Then the stderr should contain exactly ""

  Scenario: Set concurrency
    Given a file named "foo.md" with:
    """
    [Google](https://google.com)
    """
    When I successfully run `liche --concurrency 10 foo.md`
    Then the stderr should contain exactly ""

  Scenario: Search files recursively
    Given a directory named "foo"
    And a file named "foo/bar.md" with:
    """
    [Google](https://google.com)
    """
    When I successfully run `liche --recursive -v foo`
    Then the stderr should contain "OK"
