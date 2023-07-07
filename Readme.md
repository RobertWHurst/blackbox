
<p align="center">
  <a href="https://pkg.go.dev/github.com/RobertWHurst/blackbox">
    <img src="https://img.shields.io/github/go-mod/go-version/RobertWHurst/blackbox">
  </a>
  <a href="https://github.com/RobertWHurst/blackbox/actions/workflows/ci.yml">
    <img src="https://github.com/RobertWHurst/blackbox/actions/workflows/ci.yml/badge.svg">
  </a>
  <a href="https://github.com/sponsors/RobertWHurst">
    <img src="https://img.shields.io/static/v1?label=Sponsor&message=%E2%9D%A4&logo=GitHub&color=%23fe8e86">
  </a>
</p>

# blackbox - A flight recorder for Go

__If you encounter a bug please [report it][bug-report].__

Blackbox is a simple logger that will do exactly what you want a logger
to do - log messages, messages that are clear and contextual.

```go
logger := blackbox.New()

logger.AddTarget(blackbox.NewPrettyTarget(os.Stdout, os.Stderr))

logger.Info("Hello world", blackbox.Ctx{ "type": "greeting" })
```

## Getting started

To add blackbox to your project use go get:

```sh
go get github.com/RobertWHurst/blackbox
```

Next you'll need to create a logger in your project.

```go
logger := blackbox.New()
```

This logger will not output anything until you add a target to it. So let's do
that now. We'll add a pretty target that will output to stdout and stderr.

```go
logger.AddTarget(blackbox.NewPrettyTarget(os.Stdout, os.Stderr))
```

The pretty target will output messages in a human readable format.

```sh
2000-01-01T12:00:00Z00:00 info    Hello world type=greeting
```

## Providing Context

blackbox has a concept called contexts. Contexts are a way to provide
additional information associated with a message. For example, if you're
logging a message about a user, you might want to include the user's id in the
message. You could do this by logging a string containing the user's id, but
that might be more difficult to parse, and search for. Instead you can use
contexts which can be formatted in a searchable way.

Contexts are also useful for associating things like request ids, component
names, and other information that might be useful to have in a log message.

```go
logger.Info("Hello world", blackbox.Ctx{ "type": "greeting" })
```

## Levels

blackbox has 6 levels. Trace, Debug, Info, Warn, Error, and Fatal. Each level
has a purpose, and a meaning. Some have special behavior such as Fatal which
will cause the program to exit.

### Trace

Trace should be used for overly detailed information to help debug the flow of
a program.

### Debug

Debug should be used to aid in debugging, but should be less detailed than
Trace.

### Info

Info should be used for general information about the flow of a program.
These messages should be expected to be recorded in production.

### Warn

Warn should be used to indicate that something unexpected happened, but the
program is still able to continue.

### Error

Error should be used to indicate that something unexpected happened, and it
caused a failure, but the program is still able to continue.

### Fatal

Fatal should be used to indicate that a critical failure has occurred, and the
program needs to exit. Fatal will call os.Exit(1) after logging the message.

## Targets

Targets handle logger output. Loggers can have more than one target. There are
two targets included with blackbox. A pretty target, and a JSON target. Both
will write to any pair of io.Writer. This allows you to write to files, stdout,
stderr, or any other io.Writer.

Let's take a look at these two targets.

### Pretty

The pretty target will output messages in a human readable format.

```sh
2000-01-01T12:00:00Z00:00 info    Hello world type=greeting
```

Not shown in the example above, but the pretty target will also colorize the
level of the message, and the keys in the context using ansi color codes.

You can also customize the format of the message by way of SetLevel,
ShowTimestamp, ShowContext, and UseColor.

```go
logger.AddTarget(blackbox.NewPrettyTarget(os.Stdout, os.Stderr).
    SetLevel(blackbox.Trace).
    ShowTimestamp(false).
    ShowContext(false).
    UseColor(false))
```

### JSON

The json target will output messages in a JSON format.

```json
{"time":"2000-01-01T12:00:00Z00:00","level":"info","message":"Hello world"}
```

With the exception of color, the json target has the same customization
options as the pretty target.

```go
logger.AddTarget(blackbox.NewPrettyTarget(os.Stdout, os.Stderr).
    SetLevel(blackbox.Trace).
    ShowTimestamp(false).
    ShowContext(false))
```

## Implementing Targets

Targets are simple to implement. They only need to implement the Target
interface - a single Log Method.

```go
type Target interface {
    Log(level Level, values []interface{}, context Context) error
}
```

Please note that if synchronization is needed, it should be handled by the
target. The logger will not handle synchronization.

Let's go over the arguments to the Log method.

The first argument is the level of the message, exactly as it was passed to the
logger.

The second argument is a slice of values. The reason this isn't a single string
is because the logger will accept any number of arguments of any type. By
passing the values to the target, it allows the target to decide how to
format and/or consume them respective of value type.

The third argument is the context of the message, a map with string keys and
interface{} as the values.

With these values targets can to a wide range of things with maximum flexibility
and control.

## Help Welcome

If you want to support this project by throwing be some coffee money It's
greatly appreciated.

[![sponsor](https://img.shields.io/static/v1?label=Sponsor&message=%E2%9D%A4&logo=GitHub&color=%23fe8e86)](https://github.com/sponsors/RobertWHurst)

If your interested in providing feedback or would like to contribute please feel
free to do so. I recommend first [opening an issue][feature-request] expressing
your feedback or intent to contribute a change, from there we can consider your
feedback or guide your contribution efforts. Any and all help is greatly
appreciated since this is an open source effort after all.

Thank you!

[bug-report]: https://github.com/RobertWHurst/blackbox/issues/new?template=bug_report.md
[feature-request]: https://github.com/RobertWHurst/blackbox/issues/new?template=feature_request.md
