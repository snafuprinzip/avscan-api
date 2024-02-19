package main

import (
	"html/template"
	"net/http"
)

const indextemplate = `
{{define "index"}}
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <link href='http://fonts.googleapis.com/css?family=Roboto+Condensed:300' rel='stylesheet' type='text/css'>
  <link href='http://fonts.googleapis.com/css?family=Roboto+Mono' rel='stylesheet' type='text/css'>
  <style type="text/css">
      <!--
      body { font: 11pt/1.5 Roboto Condensed,Arial,Verdana,sans-serif; margin: 20px;}
      h1 { font-size: 15pt; font-weight:bold; color: #0073b3; margin-top: 0; margin-bottom: 1px;}
      h2 { font-size: 12pt; font-weight:bold; color: #0073b3; margin-top: 30px; margin-bottom: 10px;}
      table { font-size: 11pt; }
      table td { padding: 2px 6px; }
      pre {
          font: 10pt/1.5 Roboto Mono,Courier New;
          display: block;
          background-color: #eeeeee;
          padding: 5px;
      }
      .light {
          color: #888888;
          font-size: 9pt;
      }
      .cssButton {
          -moz-box-shadow:inset 0 1px 0 0 #7a8eb9;
          -webkit-box-shadow:inset 0 1px 0 0 #7a8eb9;
          box-shadow:inset 0 1px 0 0 #7a8eb9;
          background:-webkit-gradient(linear, left top, left bottom, color-stop(0.05, #637aad), color-stop(1, #5972a7));
          background:-moz-linear-gradient(top, #637aad 5%, #5972a7 100%);
          background:-webkit-linear-gradient(top, #637aad 5%, #5972a7 100%);
          background:-o-linear-gradient(top, #637aad 5%, #5972a7 100%);
          background:-ms-linear-gradient(top, #637aad 5%, #5972a7 100%);
          background:linear-gradient(to bottom, #637aad 5%, #5972a7 100%);
          filter:progid:DXImagetransform.Microsoft.gradient(startColorstr='#637aad', endColorstr='#5972a7',Gradienttype=0);
          background-color:#0073b3;
          border:1px solid #314179;
          display:inline-block;
          cursor:pointer;
          color:#ffffff;
          font-family:Arial;
          font-size:13px;
          font-weight:bold;
          padding:4px 10px;
          text-decoration:none;
      }
      .cssButton:hover {
          background:-webkit-gradient(linear, left top, left bottom, color-stop(0.05, #5972a7), color-stop(1, #637aad));
          background:-moz-linear-gradient(top, #5972a7 5%, #637aad 100%);
          background:-webkit-linear-gradient(top, #5972a7 5%, #637aad 100%);
          background:-o-linear-gradient(top, #5972a7 5%, #637aad 100%);
          background:-ms-linear-gradient(top, #5972a7 5%, #637aad 100%);
          background:linear-gradient(to bottom, #5972a7 5%, #637aad 100%);
          filter:progid:DXImagetransform.Microsoft.gradient(startColorstr='#5972a7', endColorstr='#637aad',Gradienttype=0);
          background-color:#5972a7;
      }
      .cssButton:active {
          position:relative;
          top:1px;
      }
      .warn {
          border:1px dotted #ff0000;
          color: #ff0000;
          padding: 5px;
          display: inline-block;
      }

      -->
  </style>
</head>

<body>
  <h1>{{ .Conf.Global.Title }}</h1>
  
  <p>Willkommen auf der &Uuml;bersichtsseite der AvScan API des Teams {{ .Conf.Global.Team }}. Dieser Kubernetes Microservice stellt eine <a href="https://de.wikipedia.org/wiki/Representational_State_transfer" target="_blank" rel="noopener">REST-API</a> f&uuml;r das Scannen nach Viren in hochgeladenen Dokumenten auf Basis von <a href="https://www.clamav.net/" target="_blank" rel="noopener">ClamAV</a> zur Verf&uuml;gung.</p>
  <p>Eine vollst&auml;ndige Dokumentation dieses Dienstes und der API finden Sie in zuk&uuml;nftig in Sharepoint.</p>
  <p class='warn'>
    Achtung !!! Diese API wird nicht zum Internet hin exponiert. Bitte sehen sie in Ihren Applikationen vor, diese API durch eine Komponente anzusteuern, die Zugriff auf das interne Netz hat.
  </p>
  
  <h2>Information:</h2>
  <table border=0>
    <tr>
      <td>Version:</td>
      <td>{{ .Conf.Global.Version }}</td>
    </tr>
    <tr>
      <td>Environment:</td>
      <td>{{ .Conf.Global.Environment }}</td>
    </tr>
    <tr>
      <td>URL:</td>
      <td>{{ .Conf.Global.URL }}</td>
    </tr>
    <tr>
      <td>Scanner:</td>
      <td>{{ .Conf.Scanner.Name }}</td>
    </tr>
    <tr>
      <td>Author:</td>
      <td>{{ .Conf.Global.Author }}</td>
    </tr>
    <tr>
      <td>Maintainer:</td>
      <td>{{ .Conf.Global.Maintainer }}</td>
    </tr>
    <tr>
      <td>Kontakt:</td>
      <td>{{ .Conf.Global.Email }}</td>
    </tr>
  </table>
  
  <h2>Endpoints:</h2>
  <table border=0>
    <tr>
      <td>Health check:</td>
      <td><A href='/api/v1/health'>/api/v1/health</A></td>
    </tr>
    <tr>
      <td>Display configuration:</td>
      <td><A href='/api/v1/config'>/api/v1/config</A></td>
    </tr>
    <tr>
      <td>Display scanner version:</td>
      <td><A href='/api/v1/version'>/api/v1/version</A></td>
    </tr>
    <tr>
      <td>Upload and scan file:</td>
      <td><A href='/api/v1/scan'>/api/v1/scan</A></td>
    </tr>
  </table>
  
  <h2>Response Codes (Scan):</h2>
  <table border=0>
    <tr>
      <td>200</td>
      <td>Check OK</td>
      <td>additional data in JSON-format</td>
    </tr>
    <tr>
      <td>400</td>
      <td>Bad request.</td>
      <td>one or more parameters are missing.</td>
    </tr>
    <tr>
      <td>403</td>
      <td>Forbidden </td>
      <td>IP Adress not in access list.</td>
    </tr>
    <tr>
      <td>502</td>
      <td>Infection found.</td>
      <td>additional data in JSON-format.</td>
    </tr>
    <tr>
      <td>413</td>
      <td>Document size exceeds maximum.</td>
      <td></td>
    </tr>
    <tr>
      <td>504</td>
      <td>Service unavailable.</td>
      <td></td>
    </tr>
  </table>
  
  
  <h2>cURL example request:</h2>
  <pre>
  curl -i -X POST \
    -F app-id=curl \
    -F correlation-id=curl1234 \
    -F 'file=@/etc/hosts' \
    http://{{ .Conf.Global.URL }}/api/v1/scan
  </pre>
  
  <h2>Example upload form (POST):</h2>
  <form action="/api/v1/scan" method="POST" enctype="multipart/form-data">
    <table border='0'>
      <tr><td>file:</td><td><input type="file" name="file" /></td></tr>
      <tr><td>app-id:</td><td> <input type="text" name="app-id" /> <span class="light">mandatory [a-z0-9\-]{1,36}</span></td></tr>
      <tr><td>correlation-id:</td><td><input type="text" NAME="correlation-id" /> <span class="light">mandatory [a-z0-9\-]{1,36}</span></td></tr>
    </table>
    <br/>
    <input class="cssButton" type="submit" value="  SCAN ME  "/></br>
</form>

</body>
</html>

{{ end }}
`

func RenderIndex(w http.ResponseWriter, r *http.Request) {
	if !isAccessGrantedByIP(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	idxtpl, err := template.New("index").Parse(indextemplate)
	if err != nil {
		logf("crit", "render", http.StatusInternalServerError, "Unable to render index template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	data := struct {
		Conf ConfigStruct
	}{
		Conf: *Config,
	}

	idxtpl.Execute(w, data)
}
