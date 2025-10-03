# Mock Gaming Engine

To run this repo, launch the debugger mode in VS Code.  
It will ask you to create a `launch.json` file. Copy-paste the below configuration:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Run Gaming Engine",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/server",
      "cwd": "${workspaceFolder}",
      "env": {
        "PORT": "8080"
      },
      "args": [],
      "buildFlags": "",
      "showLog": true
    }
  ]
}
