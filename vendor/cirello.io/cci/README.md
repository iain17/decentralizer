# cirello.io CI service

## Design goals

cirello.io CI service implement a continuous integration service in a single
binary that is capable running standalone.

* No retries - on failure, you should have the option to try again.
* No worker filter - all workers attached to a repository should be able to run its CD/CI steps.
* No global queue - every worker is attached to one queue only, and every build target is its own queue.
* Only support Github and Slack.