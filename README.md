# ttkeys
ttkeys helps with injecting secret keys securely into apps.

## how to install ttkeys
```
wget curl ttkeys.simplifiednetworks.co/ttkeys-v0.0.1.tar
tar -zxvf ttkeys-v0.0.1.tar
mv ttkeys-v0.0.1 /usr/bin/ttkeys
```

## how to use ttkeys
Once ttkeys is installed, you can use it to run a sample nodejs program like this:
```
ttkeys node index.js
```
ttkeys looks at the configuration file ttkeysconfig.yaml (discussed next) to determine where to fetch the secrets from. Once fetched successfully, ttkeys starts the application and injects the keys into the environment of the started application. This way the application can access the secrets from its environment variables.

## ttkeys config file
ttkeys uses a config file to determine where to get the secret keys from. Right now, ttkeys can fetch secrets from AWS SecretsManager.
To configure ttkeys to pull your secrets from AWS SecretManager, place a ttkeysconfig.yaml file in the root directory of your project.
The contents of the ttkeysconfig.yaml file should look like this:
```
secretStore: aws_sm
region: us-east-1
secretName: tt-test-secret
```

Possible values for AWS regions used in the ttkeysconfig.yam file are the following:
```
- ap-east-1      // Asia Pacific (Hong Kong).
- ap-northeast-1 // Asia Pacific (Tokyo).
- ap-northeast-2 // Asia Pacific (Seoul).
- ap-south-1     // Asia Pacific (Mumbai).
- ap-southeast-1 // Asia Pacific (Singapore).
- ap-southeast-2 // Asia Pacific (Sydney).
- ca-central-1   // Canada (Central).
- eu-central-1   // EU (Frankfurt).
- eu-north-1     // EU (Stockholm).
- eu-west-1      // EU (Ireland).
- eu-west-2      // EU (London).
- eu-west-3      // EU (Paris).
- sa-east-1      // South America (Sao Paulo).
- us-east-1      // US East (N. Virginia).
- us-east-2      // US East (Ohio).
- us-west-1      // US West (N. California).
- us-west-2      // US West (Oregon).
```

## ttkeys config file location
ttkeys' config file can be placed in any of the following locations:
- the root directory of your project. Most preferred option
- $HOME/.ttkeys
- /etc/ttkeys/. This can serve as a global configuration of ttkeys

## unrelated installations
### install node
```
cd /usr/local
sudo tar xvf ~/node-v10.16.0-linux-x64.tar.xz --strip=1
```

### add golang to environment path
```
export PATH=$PATH:/usr/local/go/bin
```

### tar up the ttkeys executable
```
tar -zcvf ttkeys.tar.gz ttkeys
```
