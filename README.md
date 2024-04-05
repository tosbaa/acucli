# Acucli

Acucli is a command-line tool developed in Go, designed to interact with Acunetix scans efficiently. It allows users to manage their Acunetix scans directly from the terminal, providing a streamlined and accessible way to handle web application security assessments.

## Features

- Create, Delete, List and Set/Get Configuration to targets
- Create, Delete, List and Add Targets to Target Group
- Create, Delete, List and Import/Export Scan Profiles
- Trigger Scans

## Installation

You can install Acucli directly from the source code hosted on GitHub. Ensure you have Go installed on your system before proceeding with the installation. Also grab a copy of the .acucli.yaml file to work with configuration setup and setting env variables. You can find it on the repository

```bash
go install github.com/tosbaa/acucli@latest

```

## Usage

After installation, you can start using Acucli to interact with your Acunetix scans. For detailed usage instructions and command options, refer to the [documentation](https://github.com/tosbaa/acucli) or use the help command:

```bash
acucli --help
```

### Target

```bash
acucli target list # Lists the target with their corresponding ids

echo "https://target.com" | acucli target add # Adds the target from stdin

echo "<TARGET-ID>" | acucli target remove # Removed the target with the given id

acucli target --id <TARGET-ID> # Get info about the target

echo "<TARGET-ID>" | acucli target setConfig # Set scan configuration defined on the .acucli.yaml file

cat targets.txt | acucli target add --gid=<TARGETGROUP-ID> # Add targets to a target group with given id
```

### Target Group

```bash
echo "TargetGroupName" | acucli targetGroup add # Create new target group

echo "<TARGETGROUP-ID>" | acucli targetGroup remove # Removed the target group with the given id

acucli targetGroup list # List the target groups

acucli targetGroup --id <TARGET-ID> # Get targets from target group

```

## Contributing

Contributions are welcome! If you'd like to contribute, please fork the repository and use a feature branch. Pull requests are warmly welcome.

## Links

- Project repository: https://github.com/tosbaa/acucli

## Licensing

The code in this project is licensed under MIT license.
