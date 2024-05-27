# SecureForward - Lightweight SSL/TLS Passthrough Proxy

SecureForward is a minimal, lightweight and efficient SSL/TLS passthrough proxy designed for secure forwarding of encrypted traffic. It provides features such as IP control and SNI filtering, making it suitable for various networking environments.

## Features
- SSL/TLS Passthrough: Securely forwards encrypted traffic without decrypting it.
- IP Control: Allows control over which IP addresses can access the proxy.
- SNI Filtering: Filters traffic based on Server Name Indication (SNI) to route it to the appropriate destination.

## Usage
SecureForward is easy to use and can be configured with simple command-line options. To get started, follow these steps:

1. Clone the repository:
   ```
   git clone https://github.com/a1383n/secureforward-proxy.git
   ```

2. Build the proxy:
   ```
   cd secureforward-proxy
   go mod download
   go build
   ```

3. Run the proxy:
   ```
   ./secureforward-proxy
   ```

## Configuration
SecureForward can be configured through command-line options or a configuration file. Common configurations include specifying allowed IP addresses and defining SNI filtering rules.

For detailed configuration options, refer to the documentation.

## Contributing
Contributions are welcome! If you have suggestions, feature requests, or bug reports, please open an issue or submit a pull request.

## License
This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
