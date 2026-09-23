Unreleased
==========

### Upgrade note: new request header and proxy allowlists

This release sends an `X-Retry-Count` request header on retries. If your
traffic to Segment goes through a proxy, gateway or WAF that allowlists
request headers, add it before upgrading or retried uploads will be
rejected. The `Authorization` header is unchanged: this client has always
sent the write key as HTTP Basic credentials.

* Send `X-Retry-Count` on retries, so the server can distinguish a retry from a first attempt. Omitted on the first attempt.
* Unified retry handling: 429, 408, 410, 460 and 5xx (except 501, 505 and 511) are retried. `Retry-After` is honoured on all of them, not just 429, which brings 529 in through the generic 5xx rule.
* `Retry-After` accepts numeric seconds and the RFC 7231 HTTP-date formats, capped at 300s.
* Rate-limited retries are bounded by elapsed time rather than counted against the retry limit, so a long `Retry-After` no longer exhausts the budget.
* New `Config.MaxTotalBackoffDuration` and `Config.MaxRateLimitDuration` (default 12 hours each) bound the two waits, reported as `ErrBackoffBudgetExceeded` and `ErrRateLimitBudgetExceeded`.
* New `Config.ShutdownTimeout` (default 75s) bounds how long `Close` waits for in-flight retries, so shutdown neither discards a batch the server asked us to resend nor blocks for the full rate-limit budget. The final attempt carries it as a request deadline, so the bound covers the in-flight request too.
* Negative `MaxRetries`, `MaxTotalBackoffDuration`, `MaxRateLimitDuration` and `ShutdownTimeout` are rejected at construction. A negative retry count previously dropped every batch after its first failure. Zero still means "use the default", per the zero-value convention on `Config`.
* `Retry-After` is read before the response body, so a mid-read I/O error no longer loses it and push the attempt onto the counted-backoff budget.
* Only 2xx responses count as a successful upload. A 3xx is now reported as a failed upload rather than silently treated as delivered. It is not retried: a redirect `net/http` already declined to follow will not succeed on a retry. The Segment endpoint does not redirect, so this only affects custom `Endpoint` values.

v3.3.0 / 2023-10-31
===================

* Add groupId to context so the track events can related to both a user and distinct id and group
  * Note: When updating to this version, verify the groupId is not being found using the Extra map.  The new groupId field will now take precedence.

v3.1.0 / 2019-09-20
===================

  * add consistent panic error message
  * Expose the Message interface Validate method
  * return error if a custom type is enqueued
  * Handle pointer types in Enqueue()
  * message: update maxMessageBytes to 32KB

v3.0.1 / 2018-10-02
===================

* Migrate from Circle V1 format to Circle V2
* Adds CLI for sending segment events
* Vendor packages back-go and uuid instead of using gitsubmodules


v3.0.0 / 2016-06-02
===================

 * 3.0 is a significant rewrite with multiple breaking changes.
 * [Quickstart](https://segment.com/docs/sources/server/go/quickstart/).
 * [Documentation](https://segment.com/docs/sources/server/go/).
 * [GoDocs](https://godoc.org/gopkg.in/segmentio/analytics-go.v3).
 * [What's New in v3](https://segment.com/docs/sources/server/go/#what-s-new-in-v3).


v2.1.0 / 2015-12-28
===================

 * Add ability to set custom timestamps for messages.
 * Add ability to set a custom `net/http` client.
 * Add ability to set a custom logger.
 * Fix edge case when client would try to upload no messages.
 * Properly upload in-flight messages when client is asked to shutdown.
 * Add ability to set `.integrations` field on messages.
 * Fix resource leak with interval ticker after shutdown.
 * Add retries and back-off when uploading messages.
 * Add ability to set  custom flush interval.

v2.0.0 / 2015-02-03
===================

 * rewrite with breaking API changes

v1.2.0 / 2014-09-03
==================

 * add public .Flush() method
 * rename .Stop() to .Close()

v1.1.0 / 2014-09-02
==================

 * add client.Stop() to flash/wait. Closes #7

v1.0.0 / 2014-08-26
==================

 * fix response close
 * change comments to be more go-like
 * change uuid libraries

0.1.2 / 2014-06-11
==================

 * add runnable example
 * fix: close body

0.1.1 / 2014-05-31
==================

 * refactor locking

0.1.0 / 2014-05-22
==================

 * replace Debug option with debug package

0.0.2 / 2014-05-20
==================

 * add .Start()
 * add mutexes
 * rename BufferSize to FlushAt and FlushInterval to FlushAfter
 * lower FlushInterval to 5 seconds
 * lower BufferSize to 20 to match other clients
