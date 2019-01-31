Service-Level Agreement
=======================

*originally by https://github.com/shurcooL/SLA/blob/master/README.md*

This service-level agreement (SLA) describes the level of support that you should expect.

Uptime
------

-	**[95%+ uptime](https://uptime.is/95)**

	If the packages have tests enabled, then the packages will build successfully, and have their tests pass at least 95% of the time in a year.

	Check a summary of building state at the [dashboard](BUILD-DASHBOARD.md).

Response Time
-------------

-	**≤ 7 days response time**

	Issues, pull requests and comments are responded to within 7 days.

API Stability Changes
---------------------

-	**No backwards compatibility**

	I do not guarantee backwards compatibility. I do check from time to time if there are importers using my packages and I shall try to reach them through issue/bug trackers about the breaking changes.

-	**Vanity URL stability**

	Once I publish a package in a URL, it shall remain in the same URL until I see no one is importing it anymore. I reserve the right to fork the package into another URL in case of conceptual changes.

Go Version
----------

-	**Current stable Go version only**

	Quoting @shurcooL _ipsis literis_:
	> Current stable Go version is supported. Previous versions may work, but aren't guaranteed to. I don't go out of my way to break support for previous versions, but I don't hold back on using new features from current stable Go version.

Applicability
-------------

This SLA applies to all Go packages under the following namespaces:

-	[`cirello.io/...`](https://cirello.io/...)

with the exception of:

-	`cirello.io/exp/...`