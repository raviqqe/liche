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
