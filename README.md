# ghe-user-psycho-killer
Suspend GitHub Enterprise users who don't belong to a specific org

## Install
Downlad and unzip [the latest release](https://github.com/helaili/ghe-user-psycho-killer/releases/latest)

## Configure
Need the following env variable :

|Name|Description|Sample value|
|----|-----------|------------|
|GHE_SERVER|The fully qualified name of the GHE server|ghe.mydomain.com|
|GHE_WHITE_LIST_ORG|The name of the GitHub organization which contains all the users allowed to use the instance|Octocats|
|GHE_PERSONAL_ACCESS_TOKEN|The [token](https://help.github.com/articles/creating-an-access-token-for-command-line-use/) to use to connect to the server. This token needs the `admin:org` scope|aaaaaaaa000000007b8eba5a9dea532cf8980000|
|GHE_SKIP_VERIFY|Ignore SSL errors when working with untrusted certificate. *Optional - default to false*.|true/false|
