# Docker-compose CLI manager

### A CLI tool for working with multiple docker-compose files.

---

## Installation

Compile or download a binary file, and place it anywhere in your PATH.

Run `dccm completion --help` and follow instructions displayed for shell autocompletion.

This is required to be performed after each update.

## Usage

<p align="center">
    <img src="/docs/animated_gif/dccm-demo.gif?raw=true"/>
</p>

## Usage details

All commands can be followed with `--help` for detailed instructions.

| Command                                         | Description                                                                                                                                     |
|-------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| `dccm`, `dccm --help`                           | Displays help.                                                                                                                                  |
| `dccm completion --help`                        | Displays instructions on enabling shell autocompletion.                                                                                         |
| `dccm list`                                     | Lists all saved projects.                                                                                                                       |
| `dccm add [projectName]`                        | Adds `docker-compose.yml` file from current directory to specified project. `[projectName]` defaults to current directory name..                |
| `dccm add [projectName] [file]`                 | Adds specified file to the specified project.                                                                                                   |
| `dccm status`                                   | Prints out statuses of all projects.                                                                                                            |
| `dccm status [projectName]`                     | Prints out status of a specified project.                                                                                                       |
| `dccm remove [projectName]`                     | Removes a project from saved projects.                                                                                                          |
| `dccm up [projectNames ...]`                    | Runs `docker-compose up -d` command on a project(s). `[projectName]` defaults to current directory name.                                        |
| `dccm down [projectNames ...]`                  | Runs `docker-compose down --remove-orphans --volumes` command on a project(s). `[projectName]` defaults to current directory name.              |
| `dccm start [projectNames ...]`                 | Runs `docker-compose start` command on a project(s). `[projectName]` defaults to current directory name.                                        |
| `dccm stop [projectNames ...]`                  | Runs `docker-compose stop` command on a project(s). `[projectName]` defaults to current directory name.                                         |
| `dccm exec [projectName] [container] [command]` | Runs `docker-compose exec container command` command on a project. `[container]` and `[command]` values are persisted and can be later ommited. |
