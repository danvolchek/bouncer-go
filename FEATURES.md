# Features

<!-- TOC -->
* [Direct message forwarding](#direct-message-forwarding)
* [User note logging](#user-note-logging)
* [User watching](#user-watching)
* [Spam prevention](#spam-prevention)
* [Event tracking](#event-tracking)
* [User reports](#user-reports)
* [Backups](#backups)
* [Misc](#misc)
  * [Thread improvements](#thread-improvements)
  * [Uptime tracking](#uptime-tracking)
  * [Help](#help)
  * [Speak as bouncer in the server](#speak-as-bouncer-in-the-server)
  * [Sync bot commands](#sync-bot-commands)
<!-- TOC -->

## Direct message forwarding

*Description*

DMs to bouncer are forwarded to threads underneath a channel for each user.

Staff send messages back to the user through bouncer.

*Purpose*

This allows all staff to handle moderation questions/acts as a shared mod queue.

*Associated commands*

Send messages to users. (`$reply`)

Preview what a sent message would look like. (`$preview`)

Users can be blocked/unblocked, causing their messages to be not forwarded/forwarded any more. (`$block`, `$unblock`)

Users who have not been replied to are tracked in a wait list. The list can be checked/emptied by staff. (`$waiting`, `$clear`)

## User note logging

*Description*

Allow staff to store/retrieve moderation events/relevant information for moderation about users.

Track which staff members performed which actions.

*Purpose*

Track historical actions/relevant information about users so staff don't need to remember.

*Associated commands*

Store info: (`$ban`, `$kick`, `$note`, `$unban`, `$warn`, `$spam`)
* Log category
  * Either ban (the user was banned), kick (the user was kicked), note (general messages), unban (the user was unbanned), warn (the user was issued a warning), spam (the user was banned for spamming)
* Message containing the note.

If the category isn't note, also send a message to the user indicating what happened.

Also send the message to the log channel, as well as the user's reply thread.

Event history can be looked up given a user reference. (`$search`)

Delete/edit logged info by index. (`$remove`, `$edit`)

Display graphs of ban/warn count by user as well as bans/warns by month. (`$graph`)

## User watching

*Description*

Copy all messages certain users send to a channel.

*Purpose*

Track users who are potentially breaking the rules more easily.

*Associated commands*

Add/remove users to/from the watch list. (`$watch`, `$unwatch`)

List watched users. (`$watchlist`)

## Spam prevention

*Description*

Watch all server messages and timeout/delete messages by users who are spamming messages automatically.

*Purpose*

Prevent spam.

*Associated commands*

Manually unmark a user as a spammer. (`$unmute`)

## Event tracking

*Description*

Track and log the following events to a channel:
* Name changes
* Role changes
* Timeout status changes
* Ban events
* User joining the server events
* User leaving the server events
* User joining/leaving audio channels
* User removing reactions
* User message deletions
* User message edits
* Bulk user message deletions

*Purpose*

Record actions relevant for moderation.

## User reports

*Description*

Allow users to report messages through a discord level app/slash command.

*Associated commands*

`/Report`

## Backups

*Description*

Regularly upload configuration/database to the server.

*Purpose*

Have data backed up in case something goes wrong.

## Misc

### Thread improvements

Joins all threads; sets all threads to maximum timeout.

### Uptime tracking

Print the amount of time the bot has been running. (`$uptime`)

### Help

Print a help message describing available commands. (`$help`)

### Speak as bouncer in the server

Print a message in a server channel. (`$say`)

### Sync bot commands

Syncs app/slash commands with the server. (`$sync`)