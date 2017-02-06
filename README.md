# CDS Contrib

Here you'll find extensions ready to use with [CDS](https://github.com/ovh/cds).

CDS support several kind of extensions:

- Actions
- Plugins
- Templates
- Secret Backends
- µServices

See [CDS documentation](https://github.com/ovh/cds) for more details.

## Actions

- [Docker Package](https://github.com/ovh/cds-contrib/tree/master/actions/cds/cds-docker-package.hcl)
- [Git clone](https://github.com/ovh/cds-contrib/tree/master/actions/cds/cds-git-clone.hcl)
- [Go Build](https://github.com/ovh/cds-contrib/tree/master/actions/cds/cds-go-build.hcl)

## Plugins

- [Kafka Publisher](https://github.com/ovh/cds-contrib/tree/master/plugins/plugin-kafka-publish)
- [Mesos/Marathon Deployment](https://github.com/ovh/cds-contrib/tree/master/plugins/plugin-marathon)
- [Tmpl](https://github.com/ovh/cds-contrib/tree/master/plugins/plugin-tmpl)
- [Mesos/Marathon Group-Tmpl](https://github.com/ovh/cds-contrib/tree/master/plugins/plugin-group-tmpl)

## Templates

- [Plain Template](https://github.com/ovh/cds-contrib/tree/master/templates/cds-template-plain)

## Secret Backends

- [Vault Secret Backend](https://github.com/ovh/cds-contrib/tree/master/secret-backends/secret-backend-vault)

## µServices

- [cds2xmpp](https://github.com/ovh/cds-contrib/tree/master/uservices/cds2xmpp)
- [cds2tat](https://github.com/ovh/cds-contrib/tree/master/uservices/cds2tat)

## Contributions

By convention, plugins must have prefix `plugin`, templates  must have prefix `templates`, secret backends must have prefix `secret-backend`.

If adding a package directly, don't forget to create a README.md inside with full description of your extension and a AUTHORS file.

Read [CONTRIBUTING guide](CONTRIBUTING.md) for more details about contributions.

## License

This work is under the BSD license, see the [LICENSE](LICENSE) file for details.
