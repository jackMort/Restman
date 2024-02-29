# 🚧 Restman
![GitHub Workflow Status](http://img.shields.io/github/actions/workflow/status/jackMort/Restman/go.yml?branch=main&style=for-the-badge)
![GO](https://img.shields.io/badge/Made%20with%20GO-white.svg?style=for-the-badge&logo=go)

`Restman` is a command-line tool for interacting with RESTful APIs, featuring a TUI (Text-based User Interface). It's designed for developers who prefer to work within the terminal environment, offering a convenient and efficient way to test and debug APIs.

> [!NOTE]
> Restman is currently in active development and is not yet ready for production use.

![preview image](https://github.com/jackMort/Restman/blob/media/preview.png?raw=true)

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Features](#features)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Installation
Restman can be installed in several ways, including using pre-built packages from the releases, Go's package manager, or by building from source.

### Using Pre-built Packages
You can download the pre-built packages for Restman from the [Releases](https://github.com/jackMort/Restman/releases) page on the GitHub repository. Choose the appropriate package for your operating system and architecture.
For example, to download and install Restman on a Unix-like system, you can use the following commands (replace `VERSION` with the actual version you want to install):
```bash
curl -Ls https://api.github.com/repos/jackMort/Restman/releases/latest \
| grep -wo "https.*Linux_x86_64.tar.gz" \
| wget -qi - \
&& tar -xvf restman*.tar.gz restman \
&& chmod +x ./restman \
&& sudo mv restman /usr/local/bin/
```
Make sure to add the executable to your PATH if it's not already included.

### Building from Source
Alternatively, you can clone the repository and build from source:
```bash
git clone https://github.com/jackMort/Restman.git
cd Restman
go build
```
After building, you can run `./restman` to start the application.

### Verifying Installation
After installation, you can verify that Restman is installed correctly by running:
```bash
restman --version
```
This should output the version of Restman that you have installed.

> [!NOTE]
> Make sure to replace `VERSION` with the actual version number and adjust the download URL and file names according to your project's release structure. The instructions should be clear and easy to follow for users who prefer to use pre-built binaries rather than building from source.

## Usage
To start using Restman, navigate to your project directory and run:
```bash
restman
```
You can also pass in an initial URL to work with:
```bash
restman http://example.com/api
```
For a list of commands and options, use the help command:
```bash
restman --help
```

Once Restman is running, you'll be greeted with the TUI where you can configure your requests and view responses.
The interface is designed to be intuitive and easy to navigate. Below is a preview of the TUI with the available keyboard shortcuts:

![preview help image](https://github.com/jackMort/Restman/blob/media/help.png?raw=true)

## Features
- Intuitive Text-based User Interface (TUI)
- Support for various HTTP methods (GET, POST, PUT, DELETE, etc.)
- Ability to save and reuse requests
- Custom headers and body content
- Response highlighting for easy reading
- SSL/TLS support

## Configuration
Restman can be configured using a `.restmanrc` file in your home directory. Here's an example configuration:
```json
{
  "default_headers": {
    "Content-Type": "application/json",
    "User-Agent": "Restman/1.0"
  }
}
```

## Contributing
Contributions are welcome! If you'd like to contribute, please follow these steps:
1. Fork the repository.
2. Create a new branch for your feature or fix.
3. Write your code.
4. Add or update tests as necessary.
5. Ensure your code passes all tests.
6. Submit a pull request against the main branch.
Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## License
Restman is released under the GPL-3.0 License. See the bundled [LICENSE](LICENSE) file for details.
