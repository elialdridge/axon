# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security seriously and appreciate your help in keeping Axon safe for everyone.

### How to Report

If you discover a security vulnerability, please report it privately by:

1. **Email**: Send details to [INSERT SECURITY EMAIL]
2. **GitHub Security Advisories**: Use GitHub's private vulnerability reporting feature
3. **Direct Message**: Contact maintainers through secure channels

### What to Include

Please include the following information in your report:

- **Description**: Clear description of the vulnerability
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Impact**: Potential security impact and attack scenarios
- **Affected Versions**: Which versions are affected
- **Suggested Fix**: If you have ideas for fixing the issue
- **Contact Information**: How we can reach you for follow-up

### Response Timeline

We aim to respond to security reports within:

- **Initial Response**: 24 hours
- **Assessment**: 3 business days
- **Fix Development**: 7 days for critical issues, 30 days for others
- **Public Disclosure**: After fix is available and deployed

### Security Update Process

1. **Verification**: We verify and reproduce the vulnerability
2. **Assessment**: We assess the severity and impact
3. **Fix Development**: We develop and test a fix
4. **Coordinated Disclosure**: We work with you on disclosure timing
5. **Release**: We release the security update
6. **Advisory**: We publish a security advisory

## Security Best Practices

### For Users

#### API Key Security

- **Never commit API keys** to version control
- **Use environment variables** for API key storage
- **Rotate keys regularly** as recommended by providers
- **Use minimal permissions** when possible

```bash
# Good: Environment variable
export OPENROUTER_API_KEY="your_key_here"

# Bad: Hardcoded in files
api_key = "sk-1234567890abcdef"
```

#### Configuration Security

- **Secure file permissions** on config files:
  ```bash
  chmod 600 ~/.axon/config.json
  ```
- **Validate configuration** before use
- **Use secure defaults** when possible

#### Network Security

- **Use HTTPS only** for API communications
- **Verify certificates** in production
- **Be cautious** with untrusted networks

### For Developers

#### Input Validation

```go
// Always validate and sanitize input
func sanitizeInput(input string) string {
    // Remove control characters
    input = regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(input, "")
    
    // Limit length
    if len(input) > 1000 {
        input = input[:1000]
    }
    
    return strings.TrimSpace(input)
}
```

#### Error Handling

```go
// Don't expose sensitive information in errors
func processAPIKey(key string) error {
    if !isValidAPIKey(key) {
        // Good: Generic error message
        return errors.New("invalid API key format")
        
        // Bad: Exposes key details
        // return fmt.Errorf("invalid key: %s", key)
    }
    return nil
}
```

#### File Operations

```go
// Use secure file permissions
func createSecureFile(path string) (*os.File, error) {
    return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
}

// Validate file paths to prevent directory traversal
func validatePath(path string) error {
    cleanPath := filepath.Clean(path)
    if strings.Contains(cleanPath, "../") {
        return errors.New("invalid path")
    }
    return nil
}
```

## Known Security Considerations

### AI API Communication

- **Data Transmission**: Game content is sent to AI providers
- **Data Retention**: AI providers may retain conversation data
- **Content Filtering**: Limited control over AI responses
- **Rate Limiting**: API abuse could affect service availability

### Local Data Storage

- **Save Files**: Stored in plaintext JSON format
- **Configuration**: May contain sensitive API keys
- **Logs**: May contain game content and system information
- **Temporary Files**: Created during normal operation

### Terminal Security

- **Escape Sequences**: Risk of terminal escape sequence injection
- **Screen Recording**: Terminal content may be recorded
- **Shared Systems**: Risk on multi-user systems

## Mitigation Strategies

### Implemented Protections

1. **Input Sanitization**: All user input is cleaned before processing
2. **Output Escaping**: Terminal output is properly escaped
3. **Path Validation**: File paths are validated to prevent traversal
4. **API Key Validation**: Keys are validated before use
5. **Error Handling**: Errors don't expose sensitive information
6. **Secure Defaults**: Safe default configurations

### User Responsibilities

1. **API Key Management**: Keep API keys secure and private
2. **System Security**: Maintain secure system configurations
3. **Network Security**: Use secure networks for API communications
4. **Update Management**: Keep Axon updated with latest security fixes

## Security Contact

For security-related questions or concerns:

- **Security Email**: [INSERT SECURITY EMAIL]
- **Maintainer Contact**: [INSERT MAINTAINER EMAIL]
- **GitHub Issues**: For non-sensitive security discussions

## Acknowledgments

We thank the security research community for helping keep Axon secure:

- [Security researchers who have contributed will be listed here]

## Legal

By reporting security vulnerabilities, you agree to:

- **Responsible Disclosure**: Allow reasonable time for fixes
- **No Harm**: Not exploit vulnerabilities maliciously
- **Cooperation**: Work with us to resolve issues

We commit to:

- **Credit**: Acknowledge your contribution (if desired)
- **Communication**: Keep you informed of progress
- **Fix Timeline**: Address issues in reasonable timeframes

---

**Security is a shared responsibility. Thank you for helping keep Axon safe!**

