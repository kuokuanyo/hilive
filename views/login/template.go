package login

// LoginTmpl 登入介面HTML語法
const LoginTmpl = `
<!DOCTYPE html>
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
    <title>{{.Title}}</title>
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
                <h2>{{.Title}}</h2>
                <div class="form-group">
                    <label for="phone" class="sr-only">phone</label>
                    <input type="text" class="form-control" id="phone" placeholder="電話號碼"
                           autocomplete="off">
                </div>
                <div class="form-group">
                    <label for="password" class="sr-only">Password</label>
                    <input type="password" class="form-control" id="password" placeholder="密碼"
                           autocomplete="off">
                </div>
                <div class="form-group">
                    <button class="btn btn-primary" onclick="submitData()">登入</button>
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
            url: '{{.URLPrefix}}/login',
            async: 'true',
            data: {
                'phone': $("#phone").val(),
                'password': $("#password").val()
            },
            success: function (data) {
                location.href = data.data.url
            },
            error: function (data) {
                alert(data.responseJSON.msg);
            }
        });
    }
</script>

</body>
</html>
`
