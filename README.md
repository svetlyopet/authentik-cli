## authentik-cli

A lightweight command-line interface (CLI) tool written in Go to interact with the [Authentik](https://goauthentik.io) identity provider in a way that allows multi-tenancy and automation.

---

## ✨ Features

- Creates a tenant abstraction with RBAC controls over resources
- Manage users, groups, and applications
- Create, list, update, and delete Authentik resources
- Designed for scripting & automation (JSON/YAML output)

---

## 🔧 Configuration

`authentik-cli` reads its target instance and credentials from a config file.

### Generate config
```bash
authentik-cli config
```

### Config File example (`~/.authentik-cli`)

```yaml
url: https://authentik.example.com
token: your_api_token
```

---

## 🛠 Usage

Uses the syntax CMD [VERB] [NOUN] similar to tools like kubectl.

### Help

```bash
authentik-cli --help
```

### Basic Commands

```bash
authentik-cli create tenant example-tenant
```
```bash
authentik-cli create user example-user \
  --name John \
  --surname Doe \
  --email john.doe@example.com \
  --tenant-admin example-tenant
```
```bash
authentik-cli create app example-app \                                     
  --provider-type=oidc \ 
  --oidc-client-type=confidential \
  --oidc-redirect-uri-strict=http://localhost:8000 \
  --oidc-redirect-uri-regex='http://*.local:9000'
```

### Updating resources

`update` applies a partial update: only the attributes passed as flags are
changed, everything else keeps its current value.

```bash
authentik-cli update user example-user --email=jane.doe@example.com --active=false
```
```bash
authentik-cli update group example-group --name=new-group-name
```
```bash
authentik-cli update app example-app \
  --oidc-client-type=public \
  --oidc-redirect-uri-strict=https://app.example.com/callback
```

Each `--oidc-redirect-uri-*` flag replaces only the redirect URIs of its own
matching mode; pass it with an empty value (`--oidc-redirect-uri-regex=`) to
remove them.

---

## 📚 Command Overview

| Command   | Subcommands                       | Description                 |
|-----------|-----------------------------------|-----------------------------|
| `config`  |                                   | Set config                  |
| `create`  | `tenant`, `user`, `group`, `app`  | Create resources            |
| `get`     | `tenant`, `user`, `group`, `app`  | Get details for resources   |
| `update`  | `user`, `group`, `app`            | Partially update resources  |
| `delete`  | `tenant`, `user`, `group`, `app`  | Delete resources            |

Run `authentik-cli [command] --help` to see all flags and options.

---

## 🔄 Output Formats

All resource-fetching commands support output in multiple formats:

```bash
authentik-cli get user example --output json
authentik-cli get app example-app --output yaml
```

Supported formats: `json`(default), `yaml`

---

## 👷 Development

See [DEVELOPMENT.md](./DEVELOPMENT.md)

---

## 📝 License

Licensed under the [Apache-2.0 license](LICENSE).

---

## 📫 Support

- Issues: [GitHub Issues](https://github.com/svetlyopet/authentik-cli/issues)
- Authentik Docs: [https://goauthentik.io/docs](https://goauthentik.io/docs)
