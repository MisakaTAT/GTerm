# GTerm

A terminal tool developed with Wails + Vue3 + TypeScript.

This is the first development preview version of GTerm. Due to my full-time work commitments, I can only dedicate limited time to this project. Development will progress gradually, and currently, it only implements basic functionality. There are still many bugs and issues that need to be addressed. If you're interested in contributing to the project, pull requests are welcome!

## Screenshots

<table>
<tr>
<td><img src="docs/preview/main-dark.png" alt="Dark Mode Main Interface" width="400"/></td>
<td><img src="docs/preview/main-light.png" alt="Light Mode Main Interface" width="400"/></td>
</tr>
<tr>
<td><img src="docs/preview/session-dark.png" alt="Dark Mode Session View" width="400"/></td>
<td><img src="docs/preview/session-light.png" alt="Light Mode Session View" width="400"/></td>
</tr>
<tr>
<td><img src="docs/preview/cred.png" alt="Credentials Management" width="400"/></td>
<td><img src="docs/preview/new-session.png" alt="Session Management" width="400"/></td>
</tr>
<tr>
<td><img src="docs/preview/perf.png" alt="Performance Monitoring" width="400"/></td>
<td><img src="docs/preview/new-session-perf.png" alt="New Session Performance" width="400"/></td>
</tr>
</table>

# Supported Platforms

- Windows 10/11 AMD64/ARM64
- MacOS 10.13+ AMD64
- MacOS 11.0+ ARM64
- Linux AMD64/ARM64

## Dependencies

- Go 1.24.1
- PNPM (Node 18+)
- Wails CLI v2.10.1

## Development

```bash
wails dev
```

## Building

```bash
wails build
```

