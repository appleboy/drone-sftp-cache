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
* **ignore_branch** - boolean flag to ignore commit branch name on hash value

The following secret values can be set to configure the plugin.

* **SFTP_CACHE_SERVER** - corresponds to **server**
* **SFTP_CACHE_PATH** - corresponds to **path**
* **SFTP_CACHE_USERNAME** - corresponds to **username**
* **SFTP_CACHE_PASSWORD** - corresponds to **password**
* **SFTP_CACHE_PRIVATE_KEY** - corresponds to **key**

See [secrets](http://readme.drone.io/usage/secret-guide/) for additional
information on secrets

## Example

The following is a sample configuration in your .drone.yml file:

```yaml
pipeline:
  restore_cache:
    image: applebot/drone-sftp-cache
    path: /var/cache/drone
    restore: true
    mount:
      - node_modules

  build:
    image: node:latest
    commands:
      - npm install

  rebuild_cache:
    image: applebot/drone-sftp-cache
    path: /var/cache/drone
    rebuild: true
    mount:
      - node_modules
```
