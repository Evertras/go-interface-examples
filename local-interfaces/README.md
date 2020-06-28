# Defining interfaces locally

Don't even look at the code yet.  Read this first.  You're going to look at the
code in a very specific order with very specific things in mind, so don't cheat!

## Overview

It's not uncommon to see the desire for a "interface package" or shareable
interfaces of some sort.  This is generally how you might approach interfaces
in other languages, and it was how I started doing things when I moved to Go.

It makes sense on paper.  DRY (Don't Repeat Yourself) dictates that you should
move common code into a common area for reuse.  And if you have a lot of pieces
that need to use this interface, why wouldn't you?

Well... a few reasons, it turns out.

## The example code

*Don't look at the code yet!*

This example contains a few super barebones packages that contain some simple code.
The important thing is not what they actually do if you run them, the important
thing is how they're built and what you can understand about them from reading them
at a high level.

These packages contain some functionality to deal with users and a score value
per-user, which can be used for a leaderboard of some sort.  Maybe it's a game,
maybe it's a social networking thing, maybe it's Maybelline.  Doesn't matter.
You can get user info, award points, and get top users from a database.
Additionally, you can send notifications to a user that maybe give them a code
or something for a promotion/reward.

### The "do low level stuff" packages

The `db` package handles all the actual database commands.  It's totally mocked
out because the implementation doesn't matter.  You can imagine writing some actual
database access code here.

The `notifications` package handles sending notifications to users.  Maybe it's an email
or a push notification; again, the implementation doesn't matter here.  Pretend
the inner code does something cool.

### The "do business logic" packages

The `handlers` package contains some functions that build HTTP handlers that could
then be used in a server.  I didn't go far enough to actually implement a server,
but that shouldn't matter here.  In this case we only need the `db` package.

The `leaderboard` package has some functionality to grab top users and send them
notifications.  It's not terribly exciting but it's enough to need `db` and `notifications`.

### The main package

There's a barebones main.go in `cmd` to demonstrate how to pass in `db` and `notifications`
instances to the `handlers` and `leaderboard` bits.

## Ok, time to look at code

You didn't look at the code yet, right?  Great.  I knew I could trust you.

Let's say you've just been hired, and some jerk named Evertras left behind some
legacy code that you're now in charge of maintaining and adding features to.
This is your life now.  Congratulations.

So you sit down at your desk on Day 1, and... where do you even start?

The first question I have for you is: which package would you look at first?

If there was no documentation, [checking the main.go would probably be a good bet.](./cmd/main.go)
You'd then see the main packages in use and how `handlers` and `leaderboard` seem
to be doing some high level stuff with `db` and `notifications` being passed to them.

```golang
leaderboard := leaderboard.New(database, notifier)

leaderboard.NotifyTopPlayers(context.Background(), 3)
```

At this point I would strongly argue you should step into either `handlers` or
`leaderboard`.  Going straight into `db` or `notifications` wouldn't tell you
what the system does, only how it interacts with some lower level resources.

So let's get into some business logic!

## First stop: Leaderboard

Let's go with `leaderboard` because I feel like it.  Fine, you can look at code now.
[Go here.](./leaderboard/leaderboard.go)

So you jump into `leaderboard.New` and find this signature:

```golang
func New(topUserGetter TopUserGetter, topScoreNotifier TopScoreNotifier) *Leaderboard {
    // ...
}
```

Wait, what?  You saw a database and a notifications instance get passed in, what's this?

Joy, that's what.  Because this function is telling you exactly what capabilities are
required to get a Leaderboard working!

```golang
type TopUserGetter interface {
	GetTopUsers(ctx context.Context, count int) ([]*db.User, error)
}

type TopScoreNotifier interface {
	NotifyTopScore(ctx context.Context, id string, score int) error
}

func New(topUserGetter TopUserGetter, topScoreNotifier TopScoreNotifier) *Leaderboard {
    // ...
}
```

### Why is this good?

When you see these interfaces, you now know the following.

1. Leaderboard needs something that can get top users
2. Leaderboard needs something that can notify about a top score
3. The two above things are 100% relevant to Leaderboard
4. Leaderboard cannot do anything else like delete a user

Imagine if instead there was some `IDb` interface from a traitorous C# coder that
contained 50 different function signatures, but all `leaderboard` ever needed was
that one.  How would you know that? You'd have to go through all the code and track
which calls are used.  Not fun.  What if the original intent was to make `leaderboard`
read-only for architectural purposes, but there's this neat `AwardPoints` function
on `IDb` and the temptation is there to use it and the resulting PR makes a senior
dev cry because they never meant for this.  Now who's the jerk?

The point is, these interfaces are documentation in their own right.  They tell you
exactly what, and *only* what, `leaderboard` is going to need to do to the outside world.
This is wonderful for keeping your code clean and vastly reducing the amount of
tribal knowledge required to maintain a project.

It's easy to take for granted what some code is doing while you're writing it and
actively maintaining it.  But you should often take a step back and consider what
your code looks like to someone that's never touched it before.  Whenever you have
tools at your dispoal to reduce the mental load of someone coming in, it's probably
a good idea to use them.

### Implications for testing

[Now look at our tests.](./leaderboard/leaderboard_test.go)

See how simple the mock can be?  When we only need to mock a small subset of a larger
whole, the mock itself is completely manageable and makes it easy for us to clearly
set up whatever scenario we want to test against.

*But wait,* I hear you say.  *What if our interface package also included a mock*
*implementation for testing purposes?  Then we wouldn't ever have to rewrite any mocks!*
*Long live IDb!*

This sounds totally reasonable at first, yeah.  The problem is that the mock will
quickly start to grow, and grow, and grow.  And because different packages will want
to set up specific scenarios for testing, you'll start adding weird configs to set things
up a certain way in the mock.  And then suddenly you realize your mock has gotten
complicated enough to need its own tests because tweaking something suddenly broke a
random actual test that relied on the mock, and... yeah.  I've been down that road.
Learn from my mistakes.

Small, lean, self-contained mocks like this may end up getting copied around to some
extent, and this isn't always as DRY as you could be.  But consider the tradeoffs.

**Simpler mocks make more confident tests.  More confident tests means less friction in
development.**

## The rest of the code

Now that you're thinking in terms of self-contained interfaces, take a look at the rest
of the code.  I've added comments everywhere to preach at you, don't worry.  [Handlers](./handlers/user.go)
is a good next step.

When you get around to [looking at the database code](./db/db.go), notice that there's stuff
in there that isn't used by any of the other packages yet.  You didn't need to know that
`db` could do all these things.  You only had to worry about what those packages needed `db`
to do.  That mindset lets you create much more self-contained and vastly more understandable code.

## Summary

Go interfaces don't work like C# or Java interfaces.  They allow you to very clearly
declare required dependencies and this comes with some great benefits that you can't
easily get from other languages.  Don't fight this by being a DRY zealot.  Embrace it!
