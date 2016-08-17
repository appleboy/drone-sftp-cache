Use this plugin for caching build artifacts to speed up your build times. This
plugin can create and restore caches of any folders.

## Config

The following parameters are used to configure the plugin:

* **server** - host of the sftp server
* **username** - authenticate with this user against sftp server
* **password** - authenticate with this password against sftp server
* **key** - private key for ssh authentication
* **path** - root path on the sftp server
* **mount** - one or an array of folders to cache
* **rebuild** - boolean flag to trigger a rebuild
* **restore** - boolean flag to trigger a restore

The following secret values can be set to configure the plugin.

* **SFTP_CACHE_SERVER** - corresponds to **server**
* **SFTP_CACHE_PATH** - corresponds to **path**
* **SFTP_CACHE_USERNAME** - corresponds to **username**
* **SFTP_CACHE_PASSWORD** - corresponds to **password**
* **SFTP_CACHE_PRIVATE_KEY** - corresponds to **key**

It is highly recommended to put the **SFTP_CACHE_USERNAME** and
**SFTP_CACHE_PASSWORD** or **SFTP_CACHE_PRIVATE_KEY** into a secret so it is
not exposed to users. This can be done using the drone-cli.

```bash
drone secret add --image=sftp-cache \
    octocat/hello-world SFTP_CACHE_USERNAME octocat

drone secret add --image=sftp-cache \
    octocat/hello-world SFTP_CACHE_PASSWORD pa55word

drone secret add --image=sftp-cache \
    octocat/hello-world SFTP_CACHE_PRIVATE_KEY @path/to/private/key
```

Then sign the YAML file after all secrets are added.

```bash
drone sign octocat/hello-world
```

See [secrets](http://readme.drone.io/0.5/usage/secrets/) for additional
information on secrets

## Example

The following is a sample configuration in your .drone.yml file:

```yaml
pipeline:
  sftp_cache:
    restore: true
  	mount:
  	  - node_modules

  build:
    image: node:latest
    commands:
      - npm install

  sftp_cache:
    rebuild: true
  	mount:
  	  - node_modules
```
