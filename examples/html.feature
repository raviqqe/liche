Feature: HTML
  Scenario: Check an empty HTML file
    Given a file named "foo.html" with ""
    When I successfully run `liche foo.html`
    Then the stdout should contain exactly ""

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
    Then the stdout should contain exactly ""
