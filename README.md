# RemoteConfigServer

- A go server to test [Remote Configuration](https://sdk.collaboraonline.com/docs/installation/Configuration.html#remote-dynamic-configuration) feature in [Collabora Online](https://www.collaboraonline.com).

- It exposes static `http://localhost:8080/asset.json` for templates and fonts

```json
{
  "kind": "assetconfiguration",
  "server": "remoteserver",
  "templates": {
    "presentation": [
      {
        "uri": "http://localhost:8080/static/impress-template/template1.otp",
        "version": "1"
      },
      {
        "uri": "http://localhost:8080/static/impress-template/template2.otp",
        "version": "2"
      },
      {
        "uri": "http://localhost:8080/static/impress-template/template3.otp",
        "version": "3"
      }
    ]
  },
  "fonts": [
    {
      "uri": "http://localhost:8080/static/font/font1.ttf",
      "version": "1"
    },
    {
      "uri": "http://localhost:8080/static/font/font2.ttf",
      "version": "2"
    },
    {
      "uri": "http://localhost:8080/static/font/font3.ttf",
      "version": "3"
    }
  ]
}
```

- To use this remote server for your COOL setup. You need to define `remote_asset_config.url` in your `coolwsd.xml`. `COOLWSD` will fetch this json every 60s and if there are any changes like adding new template/deleting new template it will be carried out.

```xml
    <remote_asset_config>
      <url desc="URL of optional JSON file that lists fonts and impress template to be included in Online" type="string" default="">http://localhost:8080/asset.json</url>
    </remote_asset_config>
```

## Running locally

- You can run it locally using following command:

```sh
go run .
```

- It starts listening on "8080" by default but you can change it by passing `--port` flag. For example: `go run . --port=5000`
