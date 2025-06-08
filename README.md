# Scribe

> Scribe is currently in active development and not all planned functionality is currently implemented.

Scribe is a local-only, TUI project management tool built for managing multiple projects. Rather than being
a typical project management or todo list application focused on "what needs to be done" and "what has been done"
Scribe is focused on "what **was being done**" by making it easy to track what was work was planned for the
day ([Session Planning](#session-planning)) and through [Note Taking](#note-taking) for a session. These features let Scribe track the
important context to what and why you were working on a project.

- [Features](#features)
    - [JSON Database](#json-database)
    - [Session Planning](#session-planning)
    - [Note Taking](#note-taking)
    - [Reporting](#reporting)
- [Installing](#installing-scribe)
- [Using Scribe](#using-scribe)
    - [Running Scribe Locally](#running-scribe-locally)
    - [Running Scribe Globally](#running-scribe-globally)
    - [Reports](#reports)
- [Keybindings](#keybindings)
    - [Navigation](#navigation)
    - [Interaction](#interaction)
- [Features](#features)

## Features

### JSON Database
Scribe uses a simple JSON file as a database that is stored locally in a project, or globally, allowing for simple management
of many projects in progress.

### Session Planning
Scribe lets you select a subset of tasks as "planned tasks" for a session (a day) to easily see what work was planned on being done.

### Note Taking
Take simple text-only notes that are tied to a session providing context to the planned and completed tasks of a session.

### Reporting
Scribe provides a `report` sub-command that provides a view of the planned and completed tasks for a session as well as
the notes for that session. These reports can provide context when picking up a project after a long period of time away
or when trying to refresh the context on why something was done.

## Installing Scribe

Scribe can be installed using the `go install` tool by running the command `go install github.com/darwinfroese/scribe@latest`

Scribe can also be installed by downloading the binary for you respective operating system from the [releases](https://github.com/darwinfroese/scribe/releases/)
tab and putting it in a folder on your `PATH`.

## Using Scribe

### Running Scribe Locally
Scribe can be run on a local database by running the command `scribe` which will create and interact with a `.scribe`
file in the current folder.

### Running Scribe Globally
Scribe can be run on a global database by running the command `scribe --global` which will create and interact with a `.scribe`
file in your `HOME` directory.

### Reports
Reports can be run on either the local or global database by adding the `--global` flag to the report command. The following
commands are available:

- **report**: will output a report for the last session
- **report all**: will output a report for all sessions
- **report last #**: will output a report for the last # of sesssion
- **report list**: will output a list of all sessions
- **report start YYYY-MM-DD end YYYY-MM-DD**: will ouput a report for all sessions between the start and end date

## Keybindings
The following keybinds are available in Scribe:

### Navigation
- **arrow keys**: navigates between items in the lists
- **hjkl**: navigates between items in the list
- **ctrl+hjkl**: navigates between panes
- **tab/shift+tab**: navigates between fields/buttons in dialogs
- **enter**: interacts with buttons or dropdowns

### Interaction
- **a**: opens the "add task" dialog
- **A (shift+a)**: opens the "add child task" dialog (defaults to the current task as the parent)
- **e**: opens the "edit task" dialog for the current task
- **t**: toggles a task as a child of the task above it (un-toggling not implemented yet, can only nest one level)
- **spacebar**: completes a task (or un-completes a task if it's completed)
- **s**: sorts the tasks in descending order of priority
- **S (shift+s)**: sorts the tasks in ascending order of priority
- **p**: marks a task as "planned" for the session
- **n**: opens the notes editor dialog for the session
- **x**: deletes a task

