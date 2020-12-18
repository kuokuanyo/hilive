package signup

// SignupTmpl 註冊用戶模板
const SignupTmpl = `<!DOCTYPE html>
<!--[if lt IE 7]>
<html class="no-js lt-ie9 lt-ie8 lt-ie7">
<![endif]-->
<!--[if IE 7]>
<html class="no-js lt-ie9 lt-ie8">
<![endif]-->
<!--[if IE 8]>
<html class="no-js lt-ie9">
<![endif]-->
<!--[if gt IE 8]><!-->
<html class="no-js">
<!--<![endif]-->
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>註冊新用戶</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">

	<link rel="stylesheet" href="{{link .CdnURL .URLPrefix "/assets/login/dist/all.min.css"}}">

	<!--[if lt IE 9]>
	<script src="{{link .CdnURL .URLPrefix "/assets/login/dist/respond.min.js"}}"></script>
	<![endif]-->

</head>
<body>

<div class="container">
	<div class="row" style="margin-top: 80px;">
		<div class="col-md-4 col-md-offset-4">
			<form action="##" onsubmit="return false" method="post" id="sign-up-form" class="fh5co-form animate-box"
				  data-animate-effect="fadeIn">
				<h2>註冊新用戶</h2>
				<div class="form-group">
					<tr>
						<label for="username" class="sr-only">用戶名稱</label>
						<td><input type="text" class="form-control" id="username" autocomplete="off" required="required" placeholder="輸入用戶名稱" ></td>
					<tr/>
				</div>
				<div class="form-group">
					<tr>
						<label for="phone" class="sr-only">phone</label>
						<td><input type="tel" class="form-control" id="phone" autocomplete="off" required="required" placeholder="輸入電話號碼" ></td>
					<tr/>
				</div>
				<div class="form-group">
					<tr>
						<label for="password" class="sr-only">password</label>
						<td><input type="password" class="form-control" id="password" autocomplete="off" required="required" placeholder="輸入密碼"></td>
					</tr>
				</div>
				<div class="form-group">
					<tr>
						<label for="checkPassword" class="sr-only">checkPassword</label>
						<td><input type="password" class="form-control" id="checkPassword" autocomplete="off" required="required" placeholder="再次輸入密碼"></td>
					</tr>
				</div>
				<div class="form-group">
					<tr>
						<label for="email" class="sr-only">email</label>
						<td><input type="email" class="form-control" id="email" autocomplete="off" required="required" placeholder="輸入電子郵件(gmail)"></td>
					</tr>
				</div>
				<div class="form-group">
					<button class="btn btn-primary" onclick="submitData()">註冊</button>
				</div>
			</form>
		</div>
	</div>
	<div class="row" style="padding-top: 60px; clear: both;">
		<div class="col-md-12 text-center"></div>
	</div>
</div>

<div id="particles-js">
	<canvas class="particles-js-canvas-el" width="1606" height="1862" style="width: 100%; height: 100%;"></canvas>
</div>

<script src="{{link .CdnURL .URLPrefix "/assets/login/dist/all.min.js"}}"></script>

<script>
	function submitData() {
		$.ajax({
			dataType: 'json',
			type: 'POST',
			url: '{{.URLPrefix}}/signup',
			async: 'true',
			data: {
				'username': $("#username").val(),
				'phone': $("#phone").val(),
				'password': $("#password").val(),
				'checkPassword': $("#checkPassword").val(),
				'email': $("#email").val(),
			},
			success: function (data) {
				alert('註冊新用戶成功');
				location.href = data.data.url
			},
			error: function (data) {
				alert(data.responseJSON.msg);
			}
		});
	}
	
	function getQueryVariable(variable){
		var query = window.location.search.substring(1);
		var vars = query.split("&");
		for (var i=0;i<vars.length;i++) {
				var pair = vars[i].split("=");
				if(pair[0] == variable){return pair[1];}
		}
		   return(false);
	}

	function getCharFromUtf8(str) {  
		var cstr = "";  
		var nOffset = 0;  
		if (str == "")  
		return "";  
			str = str.toLowerCase();  
			nOffset = str.indexOf("%e");  
		if (nOffset == -1)  
		return str;  
		while (nOffset != -1) {  
				cstr += str.substr(0, nOffset);  
				str = str.substr(nOffset, str.length - nOffset);  
		if (str == "" || str.length < 9)  
		return cstr;  
				cstr += utf8ToChar(str.substr(0, 9));  
				str = str.substr(9, str.length - 9);  
				nOffset = str.indexOf("%e");  
			}  
		return cstr + str;  
	} 

	function utf8ToChar(str) {  
		var iCode, iCode1, iCode2;  
			iCode = parseInt("0x" + str.substr(1, 2));  
			iCode1 = parseInt("0x" + str.substr(4, 2));  
			iCode2 = parseInt("0x" + str.substr(7, 2));  
		return String.fromCharCode(((iCode & 0x0F) << 12) | ((iCode1 & 0x3F) << 6) | (iCode2 & 0x3F));  
		} 

</script>

</body>
</html>`
