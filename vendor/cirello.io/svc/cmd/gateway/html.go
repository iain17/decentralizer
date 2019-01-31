package main

const pkgHTML = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.FrontPackageDomain}}/{{ .RootPackage }} git https://github.com/{{.BaseGithubAccount}}/{{ .RootPackage }}">
<meta http-equiv="refresh" content="0; url=https://godoc.org/{{.FrontPackageDomain}}/{{ .Package }}">
</head>
<body>
Redirecting to docs at <a href="https://godoc.org/{{.FrontPackageDomain}}/{{ .Package }}">godoc.org/{{.FrontPackageDomain}}/{{ .Package }}</a>...
</body>
</html>
`

const ssoHTML = `
<!DOCTYPE html>
<html>
<head>
<meta name="google-signin-client_id" content="%s">
</head>
<body>
<div id="signin" class="g-signin2" data-onsuccess="onSignIn"></div>
<div id="msg"></div>
<script>
var signedout = false
function onSignIn(googleUser) {
	var auth = googleUser.getAuthResponse()
	var id_token = auth.id_token;
	var xhr = new XMLHttpRequest();
	xhr.open('POST', '/ssoLogin');
	xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
	xhr.onload = function() {
		document.getElementById("msg").innerHTML="redirecting in 5s"
		document.getElementById("signout").style=""
		document.getElementById("signin").style="display: none"
		setTimeout(function(){
			if (!signedout){
				location.reload()
			} else {
				setCookie('` + gatewayTokenCookie + `','',-1)
			}
		}, 5000)
	};
	xhr.send('id_token=' + id_token);
}
function setCookie(cname, cvalue, exdays) {
	var d = new Date();
	d.setTime(d.getTime() + (exdays*24*60*60*1000));
	var expires = "expires="+ d.toUTCString();
	document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}
</script>
<a id="signout" style="display: none; color: black; text-decoration: none" href="#" onclick="signOut();">Sign out</a>
<script>
function signOut() {
	signedout = true
	var auth2 = gapi.auth2.getAuthInstance();
	auth2.signOut().then(function () {
		document.getElementById('signin').style='display: none'
		document.getElementById('signout').style='display: none'
		document.getElementById('msg').innerHTML = 'logged out'
		setCookie('` + gatewayTokenCookie + `','',-1)
	});
}
</script>
<script src="https://apis.google.com/js/platform.js" async defer></script>
</body>
</html>`
