# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.2.x   | :white_check_mark: |
| < 0.2   | :x:                |

## Reporting a Vulnerability

We take the security of `skube` seriously. If you discover a security vulnerability, please follow these steps:

### 1. **Do Not** Open a Public Issue

Please do not report security vulnerabilities through public GitHub issues.

### 2. Report Privately

Send an email to: **[your-email@example.com]** (replace with your actual email)

Include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### 3. Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Timeline**: Depends on severity
  - Critical: 1-7 days
  - High: 7-14 days
  - Medium: 14-30 days
  - Low: 30-90 days

## Security Measures

`skube` implements several security measures:

### Input Sanitization
- All user inputs are sanitized before being passed to `kubectl`
- Dangerous shell characters are blocked: `;`, `&`, `|`, `$`, `` ` ``, `\`, newlines
- Command substitution patterns are prevented: `$(...)`, `${...}`
- Flag injection is prevented by prepending `./` to inputs starting with `-`
- Input length is limited to 1024 characters to prevent DoS

### Command Execution
- Uses Go's `exec.Command()` which safely handles arguments (no shell interpretation)
- Arguments are passed as a slice, not concatenated into a shell string
- No use of `sh -c` or similar shell invocation

### kubectl Dependency
- `skube` validates that `kubectl` is installed on startup
- All Kubernetes operations are delegated to `kubectl`
- No direct API calls to Kubernetes (relies on kubectl's security)

## Known Limitations

1. **kubectl Trust**: `skube` trusts `kubectl` to be secure. If `kubectl` is compromised, `skube` is also compromised.
2. **Cluster Access**: `skube` inherits the permissions of the current kubectl context. Users should ensure their kubeconfig is properly secured.
3. **Secrets Display**: When using `skube get secrets`, secrets are displayed in plain text (same as `kubectl get secrets -o yaml`). Use with caution.

## Best Practices for Users

1. **Protect your kubeconfig**: Ensure `~/.kube/config` has proper file permissions (600)
2. **Use RBAC**: Configure Kubernetes RBAC to limit what `skube` can do
3. **Audit logs**: Enable Kubernetes audit logging to track `skube` operations
4. **Verify kubectl**: Ensure your `kubectl` binary is from a trusted source
5. **Update regularly**: Keep `skube` updated to get the latest security patches

## Security Audit History

| Date       | Auditor | Findings | Status |
|------------|---------|----------|--------|
| 2025-11-25 | Internal | Input sanitization review | ✅ Implemented |

## Acknowledgments

We appreciate the security research community's efforts in responsibly disclosing vulnerabilities. Contributors who report valid security issues will be acknowledged (with permission) in our release notes.

---

**Note**: This security policy applies to the `skube` CLI tool itself. For Kubernetes security, refer to the [Kubernetes Security Documentation](https://kubernetes.io/docs/concepts/security/).
