Feature: HTML
  Scenario: Check an empty HTML file
    Given a file named "foo.html" with ""
    When I successfully run `liche foo.html`
    Then the stderr should contain exactly ""

  Scenario: Check a HTML file
    Given a file named "foo.html" with:
    """
    <!DOCTYPE html>
    <html>
    <head>
      <title>My title</title>
    </head>
    <body>
      <div>
        <a href="https://google.com">Google</a>
        <a href="https://yahoo.com">Yahoo</a>
      </div>
    </body>
    </html>
    """
    When I successfully run `liche foo.html`
    Then the stderr should contain exactly ""

  Scenario: Ignore id reference
    Given a file named "foo.html" with:
    """
    <!DOCTYPE html>
    <html>
    <head>
      <title>My title</title>
    </head>
    <body>
      <div id="foo">
        <a href="#foo">Google</a>
      </div>
    </body>
    </html>
    """
    When I successfully run `liche foo.html`
    Then the stderr should contain exactly ""

  Scenario: Set document root
    Given a file named "foo.html" with:
    """
    <!DOCTYPE html>
    <html>
    <head>
      <title>My title</title>
    </head>
    <body>
      <div>
        <a href="/foo.html">Google</a>
      </div>
    </body>
    </html>
    """
    When I successfully run `liche --document-root . foo.html`
    Then the stderr should contain exactly ""

  Scenario: Fail without document root
    Given a file named "foo.html" with:
    """
    <!DOCTYPE html>
    <html>
    <head>
      <title>My title</title>
    </head>
    <body>
      <div>
        <a href="/foo.html">Google</a>
      </div>
    </body>
    </html>
    """
    When I run `liche foo.html`
    Then the exit status should be 1

  Scenario: Set document root to a sub directory
    Given a directory named "sub"
    And a file named "sub/foo.html" with:
    """
    <!DOCTYPE html>
    <html>
    <head>
      <title>My title</title>
    </head>
    <body>
      <div>
        <a href="/foo.html">Google</a>
      </div>
    </body>
    </html>
    """
    When I successfully run `liche --document-root sub sub/foo.html`
    Then the stderr should contain exactly ""
