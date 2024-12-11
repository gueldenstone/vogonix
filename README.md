# vogonix

## About

vogonix is a jira client for time tracking. The aim is to provide an offline-first app to track time for a certain issue assigned to the user.

The app is in early development and still unstable.

The project uses [wails](https://wails.io/) to provide a cross-platform desktop application. The backend is written in [go](https://go.dev/), the frontend uses [svelte](https://svelte.dev/).

## Usage

1. Create a folder called `.vogonix` in you home directory.
2. Copy the `config.tmpl.yml` into `~/.vogonix/config.yml`
3. Add you jira token, user name and instance url
4. Build and start vogonix

## Build

1. [Install wails](https://wails.io/docs/gettingstarted/installation)
2. [Install bun](https://bun.sh/docs/installation)
3. run `wails build`
4. open the executable under `build/bin/...`

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.
