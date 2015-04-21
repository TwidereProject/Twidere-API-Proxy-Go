# Twidere-API-Proxy-Go
Yet another Twitter API proxy written in Go language

----

##Deploy within one minute


###OpenShift

1. Create an application on OpenShift, Choose **Go Language**
2. Copy URL of this repo
3. Paste URL into **Source Code** in OpenShift configuration page
4. Click **Create Application**
5. Done! You don't need any code skill to deploy, yay!


###Heroku

Deployment to heroku is a bit more complicated than OpenShift, but still easier than any other api proxies.

1. Fork this repo
2. Create an application in **Heroku Dashboard**, you will be redirected to **Settings** of your newly created application.
3. Find **Config Variables** in **Settings** segment, click **Reveal Config Vars**, then press **Edit**
4. Add a new variable, the **key** is ````BUILDPACK_URL````, and the **value** is ````https://github.com/heroku/heroku-buildpack-go````, click **Save**.
5. Find **Connect to Github** in **Deploy** segment, gives Heroku your Github access, then type **the repo name that you've forked**, click **Connect**.
6. Find **Manual deploy**, click **Deploy Branch**.
7. Done! Why deploy to heroku? I don't know. Just provide one more choice for you ;)

##Support my work

**Donation methods**

* Me: mariotaku.lee[AT]gmail.com

* Our designer: pay[AT]uucky.me

PayPal & Alipay accepted.

Bitcoin: 1FHAVAzge7cj1LfCTMfnLL49DgA3mVUCuW