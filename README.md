Polymer Gopher
==============

This project is a sample web application hosted on [App Engine](1) composed of
two [App Engine modules](2):

- a frontend written using [Polymer](3)
- a backend written in [Go](4)

## Running locally

To run this application locally install the [Go App Engine SDK](7) and then execute:

```
$ goapp serve dispatch.yaml frontend/app.yaml backend/step6/app.yaml
```

## Deploying the app on the cloud

And to deploy it:

- Create a new Google Cloud project on the [Google Cloud Console](8)

- Write your project id in every single `yaml` file:

```yaml
	application: your-application-id
```

- Then execute

```
$ goapp deploy backend/step6/app.yaml
$ goapp deploy frontend/app.yaml
$ appcfg.py update_dispatch .
```

Then visiting http://your-project-id.appspot.com should show you the application
running on the cloud.

[1]: https://cloud.google.com/appengine/docs
[2]: https://cloud.google.com/appengine/docs/go/modules
[3]: https://www.polymer-project.org
[4]: https://golang.org
[7]: https://cloud.google.com/appengine/downloads
[9]: https://console.developers.google.com


### Disclaimer

This is not an official Google product (experimental or otherwise), it is just
code that happens to be owned by Google.
