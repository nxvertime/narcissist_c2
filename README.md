# narcissist_c2

**narcissist_c2** is an open-source Command and Control (C2) framework written in Go, designed exclusively for educational purposes. It offers a platform to manage multiple clients concurrently and broadcast commands to connected agents. 

⚠️ **Disclaimer**: This project is for educational purposes only. It should not be used in any malicious way, and the authors take no responsibility for any misuse.

## Features

- **Simultaneous Client Management**: Handle multiple clients at the same time, making it easy to control a fleet of connected machines.
- **Broadcast Commands**: Send commands to all connected clients simultaneously for efficient execution.

### Upcoming Features

- **Credential Stealing**: Capture and extract credentials from compromised machines.
- **Process Migration**: Move running processes between different contexts for stealth and persistence.
- **Privilege Escalation**: Utilize drivers to gain elevated privileges on compromised systems.
- **Network Analysis**: Perform deep network traffic analysis to understand target behavior.
- **DLL Hijacking**: Exploit vulnerabilities to hijack legitimate DLLs for malicious execution.
- **Persistence Mechanisms**: Implement methods to ensure the C2 agent remains on the system after reboot.

## Installation

1. Clone this repository:
    ```bash
    git clone https://github.com/yourusername/narcissist_c2.git
    ```
2. Navigate to the project directory:
    ```bash
    cd narcissist_c2
    ```
3. Build the project:
    ```bash
    go build -o narcissist_c2
    ```

## Usage

1. Start the C2 server:
    ```bash
    ./narcissist_c2
    ```
2. Connect clients using the provided agent code.
3. Execute commands either on individual clients or in broadcast mode.

## Educational Use Only

This project is strictly for educational use, to learn and understand how command and control frameworks operate. Do not use this for illegal purposes.

## Roadmap

- [ ] Implement credential stealing module
- [ ] Process migration functionality
- [ ] Privilege escalation using drivers
- [ ] Network traffic analysis tools
- [ ] DLL Hijacking exploits
- [ ] Persistence strategies for agent survival

## Documentation

Detailed usage guides and example configurations will be provided as the project progresses.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
