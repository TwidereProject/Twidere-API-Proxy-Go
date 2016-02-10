# Twidere-API-Proxy-Go
Yet another Twitter API proxy written in Go language


----

##Deploy within one minute


###OpenShift

1. Create an application on OpenShift, Choose **Go Language**
2. Click **Create Application**
3. git clone your **openshift repo** to local
4. git clone **this repo** to local
5. move all files of this repo to the openshift repo
6. delete web.go in openshift repo
7. git push your openshift repo
8. Done!


###Heroku

Just click [![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy) to deploy **within seconds!**

###Any devices with Golang support
1. Clone the repo;
2. Run ```go build apiproxy.go```;
3. Run ```PORT=8080 ./apiproxy```;
4. Done!
5. You can use Apache or Nginx to act as a reverse proxy to enable TLS encryption.

##Support my work

**Donation methods**

* Me: mariotaku.lee[AT]gmail.com

PayPal & Alipay accepted.

Bitcoin: 1FHAVAzge7cj1LfCTMfnLL49DgA3mVUCuW
