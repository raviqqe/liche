Feature: Relative paths
  Scenario: Check a markdown file
    Given a file named "foo.md" with:
    """
    [bar](bar.md)
    """
    And a file named "bar.md" with ""
    When I successfully run `liche -v foo.md`
    Then the stderr should contain "OK"

  Scenario: Check a directory
    Given a file named "foo.md" with:
    """
    [bar](bar)
    """
    And a directory named "bar"
    When I successfully run `liche -v foo.md`
    Then the stderr should contain "OK"

  Scenario: Check an image
    Given a file named "foo.md" with:
    """
    ![foo](foo.png)
    """
    And a file named "foo.png" with ""
    When I successfully run `liche -v foo.md`
    Then the stderr should contain "OK"

  Scenario: Check a non-existent markdown file
    Given a file named "foo.md" with:
    """
    [bar](bar.md)
    """
    When I run `liche foo.md`
    Then the exit status should be 1
    And the stderr should contain "ERROR"
