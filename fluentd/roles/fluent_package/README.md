# fluent_package

An Ansible role to install fluent-package (Fluentd LTS) on RedHat-based systems (RHEL, CentOS, AlmaLinux, Rocky Linux).

## Features

- Automatically detects system architecture and OS version
- Downloads and installs fluent-release RPM repository
- Backs up existing unmanaged repo files to prevent conflicts
- Updates yum cache and installs fluent-package with all dependencies
- Idempotent and check-mode compatible
- Comprehensive error handling and validation

## Requirements

- Ansible 2.9+
- RedHat-based system (RHEL 8+, CentOS 8+, AlmaLinux, Rocky Linux, etc.)
- sudo/root access on target system
- Network connectivity to fluentd.cdn.cncf.io

## Role Variables

### Defaults (`defaults/main.yml`)

| Variable | Default | Description |
|----------|---------|-------------|
| `fluent_package_lts_version` | `6` | LTS version of fluent-package to install |
| `fluent_package_release_version` | `2025.9.29-1` | Release version of fluent-lts-release RPM |
| `fluent_package_cdn_url` | `https://fluentd.cdn.cncf.io` | CDN base URL for fluent releases |
| `fluent_package_backup_existing_repo` | `true` | Backup existing unmanaged repo files |
| `fluent_package_download_timeout` | `30` | Download timeout in seconds |
| `fluent_package_debug` | `false` | Enable debug output |

### Read-Only Variables (`vars/main.yml`)

These are automatically set based on system facts:

| Variable | Description |
|----------|-------------|
| `fluent_supported_distributions` | List of supported Linux distributions |
| `fluent_supported_versions` | List of supported major OS versions (8, 9) |
| `fluent_repo_file` | Repository configuration file location |
| `fluent_temp_dir` | Temporary directory for downloads |

### Discovered Facts

| Fact | Source | Description |
|------|--------|-------------|
| `fluent_arch` | `ansible_architecture` | System CPU architecture |
| `fluent_el_version` | `ansible_distribution_major_version` | RedHat major version (8, 9, etc.) |
| `fluent_release_rpm_url` | Computed | Full URL to fluent-release RPM |

## Dependencies

None

## Example Playbook

### Basic Installation

```yaml
---
- hosts: webservers
  become: yes
  roles:
    - role: fluent_package
```

### With Custom LTS Version

```yaml
---
- hosts: all
  roles:
    - role: fluent_package
      vars:
        fluent_package_lts_version: "6"
        fluent_package_release_version: "2025.9.29-1"
```

### In a Larger Playbook

```yaml
---
- name: Configure logging infrastructure
  hosts: production_servers
  become: yes
  tasks:
    - name: Install base packages
      yum:
        name:
          - curl
          - git
        state: present

    - name: Install fluent-package for log forwarding
      include_role:
        name: fluent_package
      vars:
        fluent_package_debug: true

    - name: Configure fluent-package
      template:
        src: fluent.conf.j2
        dest: /etc/fluent/fluent.conf
      notify: restart fluent
```

## Supported Platforms

- Red Hat Enterprise Linux (RHEL) 8, 9+
- CentOS 8, 9
- AlmaLinux 8, 9
- Rocky Linux 8, 9
- Fedora (recent versions)

## Behavior

### Task Flow

1. **Gather System Facts** — Collects OS and hardware information
2. **Validate System** — Ensures system is RedHat-based
3. **Check Prerequisites** — Verifies curl, rpm, and yum are available
4. **Build RPM URL** — Constructs download URL based on system facts
5. **Check Existing Repo** — Detects if fluent-package repo already exists
6. **Backup Old Repo** — Backs up unmanaged repo file if found (with `.rpmsave` suffix)
7. **Download RPM** — Fetches fluent-release RPM from CDN
8. **Install Release RPM** — Registers fluent-package repository
9. **Clean Up** — Removes temporary RPM file
10. **Update Cache** — Refreshes yum package metadata
11. **Install fluent-package** — Installs main fluent-package and dependencies

### Idempotency

This role is fully idempotent:
- **First run**: Downloads and installs fluent-package and all dependencies
- **Subsequent runs**: Verifies fluent-package is installed; no changes if already present
- **Check mode**: Dry-run compatible; shows what would be installed without making changes

### Error Handling

- **Missing prerequisites**: Role fails with clear message if curl, rpm, or yum are unavailable
- **Non-RedHat system**: Role fails with explicit message; only supports RedHat family
- **Network failure**: Role fails gracefully if CDN is unreachable
- **Existing repo conflicts**: Existing unmanaged repo files are backed up, not overwritten

## Troubleshooting

### Role fails with "This role only supports RedHat-based systems"

This role is designed for RedHat family distributions only. Ensure target hosts are:
- RHEL, CentOS, AlmaLinux, Rocky Linux, or similar
- Version 8 or 9 (or later supported versions)

**Solution**: Modify target hosts or adjust role conditions if supporting other distributions.

### "Failed to download RPM" error

The CDN may be experiencing issues or network connectivity is blocked.

**Solutions**:
- Verify network connectivity: `ping fluentd.cdn.cncf.io`
- Check firewall rules allow outbound HTTPS (port 443)
- Verify CDN URL is correct in `fluent_package_cdn_url`
- Check `fluent_package_release_version` matches available releases

### Role runs but fluent-package not installed

Check yum installation logs:
```bash
yum history list fluent-package
yum info fluent-package
```

Verify the fluent-package repository is registered:
```bash
yum repolist | grep -i fluent
```

## License

MIT

## Author Information

Created for automated fluent-package installation on RedHat-based systems.
