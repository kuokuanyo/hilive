package views

// TemplateList 放置所有模板
var TemplateList = map[string]string{"head": `{{define "head"}}
	<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<title>晶橙資訊</title>
		<!-- Tell the browser to be responsive to screen width -->
		<meta content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" name="viewport">

		<!--[if lt IE 9]>
		<script src="{{ "/admin/assets/dist/js/html5shiv.min.js"}}"></script>
		<script src="{{ "/admin/assets/dist/js/respond.min.js"}}"></script>
		<![endif]-->

		<script src="/admin/assets/dist/js/all.min.8425540791.js"></script>
		<link rel="stylesheet" href="/admin/assets/dist/css/all.min.7bcc29700e.css">
	</head>
{{end}}`, "header": `{{define "header"}}
	<header class="main-header">
		<a href={{.URLRoute.IndexURL}} class="logo">
			<span class="logo-mini">{{.Config.MiniLogo}}</span>
			<span class="logo-lg">{{.Config.Logo}}</span>
		</a>
		<nav class="navbar navbar-static-top">
			<div id="firstnav">
				<a href="#" class="sidebar-toggle" data-toggle="offcanvas" role="button">
					<span class="sr-only">Toggle navigation</span>
				</a>
				<div style="float: left;">
					<ul class="nav navbar-nav">
						<li class="navbar-nav-btn-left" style="display: none;">
							<a href="javascript:;" style="border-left: none;border-right: solid 1px #dedede;">
								<i class="fa fa-angle-double-left"></i>
							</a>
						</li>
					</ul>
				</div>
				<div class="nav-tabs-content">
					<ul class="nav nav-tabs nav-addtabs">
					</ul>
				</div>
				<div style="float: left;">
					<ul class="nav navbar-nav">
						<li class="navbar-nav-btn-right" style="display: none;">
							<a href="javascript:;" style="border-left: solid 1px #dedede;border-right: none;">
								<i class="fa fa-angle-double-right"></i>
							</a>
						</li>
					</ul>
				</div>
				{{ template "admin_panel" . }}
			</div>
		</nav>
	</header>
{{end}}`, "admin_panel": `{{define "admin_panel"}}
	<div class="navbar-custom-menu">
		<ul class="nav navbar-nav">
			<li title="{{"刷新"}}">
				<a href="javascript:void(0);" class="container-refresh">
					<i class="fa fa-refresh"></i>
				</a>
			</li>
			<li class="dropdown user user-menu">
				<a href="#" class="dropdown-toggle" data-toggle="dropdown">
					{{if eq .User.Picture ""}}
						<img src="/admin/assets/dist/img/avatar04.png" class="user-image" alt="User Image">
					{{else}}
						<img src="{{.User.Picture}}" class="user-image" alt="User Image">
					{{end}}
					<span class="hidden-xs">{{.User.UserName}}</span>
				</a>
				<ul class="dropdown-menu">
					<li class="user-header">
						{{if eq .User.Picture ""}}
							<img src="/admin/assets/dist/img/avatar04.png" class="img-circle"
								alt="User Image">
						{{else}}
							<img src="{{.User.Picture}}" class="img-circle" alt="User Image">
						{{end}}
						<p>
							{{.User.UserName}} -{{.User.LevelName}}
						</p>
					</li>
					<li class="user-footer">
						<div class="pull-right">
							<a href="{{.URLRoute.URLPrefix}}/logout"
							class="no-pjax btn btn-default btn-flat">{{"登出"}}</a>
						</div>
					</li>
				</ul>
			</li>
		</ul>
	</div>
{{end}}`, "sidebar": `{{define "sidebar"}}
<aside class="main-sidebar">
	<section class="sidebar" style="height: auto;">
		<ul class="sidebar-menu" data-widget="tree">
			{{$URLPrefix := .URLRoute.URLPrefix}}
			{{range $key, $list := .Menu.List }}
				{{if eq (len $list.ChildrenList) 0}}
					{{if $list.Header}}
						<li class="header" data-rel="external">{{$list.Header}}</li>
					{{end}}
					<li class='{{$list.Active}}'>
						{{if eq $list.URL "/"}}
							<a href='{{$URLPrefix}}'>
						{{else if isLinkURL $list.URL}}
							<a href='{{$list.URL}}'>
						{{else}}
							<a href='{{$URLPrefix}}{{$list.URL}}'>
						{{end}}
							<i class="fa {{$list.Icon}}"></i><span> {{$list.Name}}</span>
							<span class="pull-right-container"><!-- <small class="label pull-right bg-green">new</small> --></span>
						</a>
					</li>
				{{else}}
					<li class="treeview {{$list.Active}}">
						<a href="#">
							<i class="fa {{$list.Icon}}"></i><span> {{$list.Name}}</span>
							<span class="pull-right-container">
							<i class="fa fa-angle-left pull-right"></i>
						</span>
						</a>
						<ul class="treeview-menu">
							{{range $key2, $item := $list.ChildrenList}}
								{{if eq (len $item.ChildrenList) 0}}
								<li>
									{{if eq $item.URL "/"}}
										<a href='{{$URLPrefix}}'>
									{{else if isLinkURL $item.URL}}
										<a href='{{$item.URL}}'>
									{{else}}
										<a href='{{$URLPrefix}}{{$item.URL}}'>
									{{end}}                            
										<i class="fa {{$item.Icon}}"></i> {{$item.Name}}
									</a>
								</li>
								{{else}}
									<li class="treeview {{$item.Active}}">
										<a href="#">
											<i class="fa {{$item.Icon}}"></i><span> {{$item.Name}}</span>
											<span class="pull-right-container">
												<i class="fa fa-angle-left pull-right"></i>
											</span>
										</a>
										<ul class="treeview-menu">
											{{range $key3, $subItem := $item.ChildrenList}}
												<li>
													{{if eq $subItem.URL "/"}}
														<a href='{{$URLPrefix}}'>
													{{else if isLinkURL $subItem.URL}}
														<a href='{{$subItem.URL}}'>
													{{else}}
														<a href='{{$URLPrefix}}{{$subItem.URL}}'>
													{{end}}                                             
														<i class="fa {{$subItem.Icon}}"></i> {{$subItem.Name}}
													</a>
												</li>
											{{end}}
										</ul>
									</li>
								{{end}}
							{{end}}
						</ul>
					</li>
				{{end}}
			{{end}}
		</ul>
	</section>
</aside>
{{end}}`, "layout": `{{define "layout"}}
 <!DOCTYPE html>
    <html>
		{{ template "head" . }}
		<body class="hold-transition skin-black sidebar-mini">
			<div class="wrapper">

				{{ template "header" . }}

				{{ template "sidebar" . }}


				<div class="content-wrapper" id="pjax-container">
					{{if eq .TmplName "menu"}}
						{{ template "menu_content" . }}
					{{else if eq .TmplName "info"}}
						{{ template "info_content" . }}
					{{else if eq .TmplName "form"}}
						{{ template "form_content" . }}
					{{else if eq .TmplName "alert"}}
						{{ template "alert_content" . }}
					{{end}}

				</div>

			</div>
			<script src="/admin/assets/dist/js/all_2.min.38a2a946b0.js"></script>
		</body>
    </html>
{{end}}`, "alert_content": `{{define "alert_content"}}
	<script src="/admin/assets/dist/js/datatable.min.581cdc109b.js"></script>
	<script src="/admin/assets/dist/js/form.min.f8678914e9.js"></script>
	<script src="/admin/assets/dist/js/treeview.min.7780d3bb0f.js"></script>
	<script src="/admin/assets/dist/js/tree.min.e1faf8b7de.js"></script>
	<section class="content-header">
		<h1>
			錯誤
			<small>發生錯誤</small>
		</h1>
		<ol class="breadcrumb" style="margin-right: 30px;">
			<li><a href={{.URLRoute.IndexURL}}><i class="fa fa-dashboard"></i> 首頁</a></li>
		</ol>
	</section> 
	<section class="content">
		{{if ne .AlertContent ""}}
		<div class="alert alert-warning alert-dismissible">
			<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>
			{{ .AlertContent}}
		</div>
		{{end}}
	</section>
{{end}}`, "form_content": `{{define "form_content"}}
	<script src="/admin/assets/dist/js/datatable.min.581cdc109b.js"></script>
	<script src="/admin/assets/dist/js/form.min.f8678914e9.js"></script>
	<script src="/admin/assets/dist/js/treeview.min.7780d3bb0f.js"></script>
	<script src="/admin/assets/dist/js/tree.min.e1faf8b7de.js"></script>
	<section class="content-header">
		<h1>
			{{.FormInfo.Title}}
			<small>{{.FormInfo.Description}}</small>
		</h1>
		<ol class="breadcrumb" style="margin-right: 30px;">
			<li><a href={{.URLRoute.IndexURL}}><i class="fa fa-dashboard"></i> 首頁</a></li>
		</ol>
	</section> 
	<section class="content">
		{{if ne .AlertContent ""}}
			<div class="alert alert-warning alert-dismissible">
				<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>
				<h4>發生錯誤</h4>
				{{ .AlertContent}}
			</div>
		{{end}}
		{{ template "form" . }}
	<section>
{{end}}`, "info_content": `{{define "info_content"}}
	<script src="/admin/assets/dist/js/datatable.min.581cdc109b.js"></script>
	<script src="/admin/assets/dist/js/form.min.f8678914e9.js"></script>
	<script src="/admin/assets/dist/js/treeview.min.7780d3bb0f.js"></script>
	<script src="/admin/assets/dist/js/tree.min.e1faf8b7de.js"></script>
	<section class="content-header">
		<h1>
			{{.PanelInfo.Title}}
			<small>{{.PanelInfo.Description}}</small>
		</h1>
		<ol class="breadcrumb" style="margin-right: 30px;">
			<li><a href={{.URLRoute.IndexURL}}><i class="fa fa-dashboard"></i> 首頁</a></li>
		</ol>
	</section> 
	<section class="content">
		{{if ne .AlertContent ""}}
			<div class="alert alert-warning alert-dismissible">
				<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>
				<h4>發生錯誤</h4>
				{{ .AlertContent}}
			</div>
		{{end}}
		<div>
			<div class="box box-" >
				<div class="box-header with-border">
					<div class="pull-right">
						<div class="dropdown pull-right column-selector" style="margin-right: 10px">
							<button type="button" class="btn btn-sm btn-instagram dropdown-toggle" data-toggle="dropdown">
								<i class="fa fa-table"></i>
								&nbsp;
								<span class="caret"></span>
							</button>
							<ul class="dropdown-menu" role="menu" style="padding: 10px;max-height: 400px;overflow: scroll;">
								<li>
									<ul style="padding: 0;">
										{{range $key, $head := .PanelInfo.FieldList}}
											<li class="checkbox icheck" style="margin: 0;">
												<label style="width: 100%;padding: 3px;">
													<input type="checkbox" class="column-select-item" data-id="{{$head.Field}}"
														style="position: absolute; opacity: 0;">&nbsp;&nbsp;&nbsp;{{$head.Header}}
												</label>
											</li>
										{{end}}
									</ul>
								</li>
								<li class="divider">
								</li>
								<li class="text-right">
									<button class="btn btn-sm btn-default column-select-all">{{"全選"}}</button>&nbsp;&nbsp;
									<button class="btn btn-sm btn-primary column-select-submit">{{"提交"}}</button>
								</li>
							</ul>
						</div>
						<div class="btn-group pull-right" style="margin-right: 10px">
							<a href="javascript:;" class="btn btn-sm btn-primary" id="filter-btn"><i
										class="fa fa-filter"></i>&nbsp;&nbsp;{{"篩選"}}</a>
						</div>
						<script>
							$("#filter-btn").click(function () {
								$('.filter-area').toggle();
							});
						</script>
						<div class="btn-group pull-right" style="margin-right: 10px">
						{{if .URLRoute.NewURL}}
							<a href="{{.URLRoute.NewURL}}" class="btn btn-sm btn-success">
								<i class="fa fa-plus"></i>&nbsp;&nbsp;{{"創建"}}
							</a>
						{{end}}
						</div>
					</div>
					<span>
					<a class="btn btn-sm btn-primary grid-refresh">
						<i class="fa fa-refresh"></i> {{"刷新"}}
					</a>
					</span>
					<script>
					let toastMsg = '{{"刷新成功"}} !';
					$('.grid-refresh').unbind('click').on('click', function () {
						$.pjax.reload('#pjax-container');
						toastr.success(toastMsg);
					});
					</script>
				</div>
				<div class="box-header filter-area " style="display: none;">
					<form id={{.FormID}} action="{{.URLRoute.InfoURL}}" method="get" accept-charset="UTF-8" class="form-horizontal" pjax-container style="background-color: white;">
						<div class="box-body">
							<div class="box-body">
								<div class="fields-group">
								{{range $key, $data := .PanelInfo.FilterFormData}}
									<div class="form-group" >
										{{if ne $data.Header ""}}
											<label for="{{$data.Field}}"
												class="col-sm-2 control-label">{{$data.Header}}</label>
										{{end}}
										<div class="col-sm-10 ">
											{{if eq $data.FormType.String "default"}}
												<div class="box box-solid box-default no-margin">
												<div class="box-body" style="min-height: 40px;">
													{{$data.Value}}
												</div>
												</div>
												<input type="hidden" class="{{$data.Field}}" name="{{$data.Field}}" value='{{$data.Value}}'>
											{{else if eq $data.FormType.String "text"}}
												{{if $data.Editable}}
													<div class="input-group">
														<span class="input-group-addon"><i class="fa fa-pencil fa-fw"></i></span>
														<input {{if $data.Must}}required="1"{{end}} type="text" name="{{$data.Field}}" value='{{$data.Value}}'
															class="form-control {{$data.Field}}" placeholder="{{$data.Placeholder}}">
													</div>
												{{end}}
											{{else if eq $data.FormType.String "datetime"}}
												{{if not .Editable}}
													<div class="box box-solid box-default no-margin">
														<div class="box-body" style="min-height: 40px;">
															{{.Value}}
														</div>
														<input type="hidden" class="{{.Field}}" name="{{.Field}}" value='{{.Value}}'>
													</div>
												{{else}}
													<div class="input-group">
														<span class="input-group-addon"><i class="fa fa-calendar fa-fw"></i></span>
														<input {{if .Must}}required="1"{{end}} style="width: 170px" type="text"
															name="{{.Field}}"
															value="{{.Value}}"
															class="form-control {{.Field}}" placeholder="{{.Placeholder}}">
													</div>
													<script>
														$(function () {
															$('input.{{.Field}}').parent().datetimepicker({{.OptionExt}});
														});
													</script>
												{{end}}
											{{else if eq $data.FormType.String "select"}}
												<select class="form-control {{.Field}} select2-hidden-accessible" style="width: 100%;" name="{{.Field}}[]"
												multiple="" data-placeholder="{{.Placeholder}}" tabindex="-1" aria-hidden="true"
												{{if not .Editable}}disabled="disabled"{{end}}>
													{{range $key, $v := .FieldOptions }}
														<option value='{{$v.Value}}' {{$v.SelectedLabel}}>{{if ne $v.TextHTML ""}}{{$v.TextHTML}}{{else}}{{$v.Text}}{{end}}</option>
													{{end}}
												</select>
													<script>
														$("select.{{.Field}}").select2({{.OptionExt}});
													</script>
											{{end}}	
										</div>
									</div>       
								{{end}}
								</div>
							</div>
						</div>
						<div class="box-footer">            
							<div class="col-md-2 ">
							</div>
							<div class="col-md-8 ">
								<div class="btn-group pull-left" >
									<button type="submit" class="btn btn-sm btn-primary" data-loading-text="&nbsp;search">
										<i class="icon fa fa-search"></i>&nbsp;&nbsp;search
									</button>
								</div>
								<div class="btn-group pull-left" style="margin-left:12px;">
									<a href= {{ .URLRoute.InfoURL}} type="reset" class="btn btn-sm btn-default" data-loading-text="&nbsp;Save">
										<i class="icon fa fa-undo"></i>&nbsp;&nbsp;reset
									</a>
								</div>
							</div>
						</div>
					</form>
				</div>
				<div class="box-body" style="overflow-x: scroll;overflow-y: hidden;padding:0;">
					<table class="table table-hover" style="min-width: 1000px;table-layout: auto;">
						<tbody>
							{{$SortUrlParam := .URLRoute.SortURL}}
							<tr>
								<th style="text-align: center;">
									<input type="checkbox" class="grid-select-all" style="position: absolute; opacity: 0;">
								</th>
								{{range $key, $head := .PanelInfo.FieldList}}
									{{if eq $head.Hide false}}
										<th>
											{{$head.Header}}
											{{if $head.Sortable}}
												<a class="fa fa-fw fa-sort" id="sort-{{$head.Field}}"
													href="?__sort={{$head.Field}}&__sort_type=desc"></a>
											{{end}}
										</th>
									{{end}}
								{{end}}
								<th style="text-align: center;">操作</th>
							</tr>
							{{$FieldList := .PanelInfo.FieldList}}
							{{$EditUrl := .URLRoute.EditURL}}
							{{$DeleteUrl := .URLRoute.DeleteURL}}
							{{$PrimaryKey := .PanelInfo.PrimaryKey}}
							{{range $key1, $info := .PanelInfo.InfoList}}
							<tr>
								<td style="text-align: center;">
									<input type="checkbox" class="grid-row-checkbox"
											data-id="{{(index $info $PrimaryKey).Content}}"
											style="position: absolute; opacity: 0;">
								</td>
								{{range $key2, $head2 := $FieldList}}
									{{if eq $head2.Hide false}}
										{{if $head2.Editable}}
											<td>
												<a href="#" class="editable-td-"
													data-pk="{{(index $info $PrimaryKey).Content}}"
													data-source=''

													data-value="{{(index $info $head2.Field).Value}}"
													data-name="{{$head2.Field}}"
													data-title="Enter {{$head2.Header}}">{{(index $info $head2.Field).Content}}</a>
											</td>
										{{else}}
											<td>{{(index $info $head2.Field).Content}}</td>
										{{end}}
									{{end}}
								{{end}}
								<td style="text-align: center;">
									{{if $EditUrl}}
										<a href='{{$EditUrl}}&__edit_pk={{(index $info $PrimaryKey).Content}}'><i
													class="fa fa-edit"></i></a>
									{{end}}
									{{if $DeleteUrl}}
										<a href="javascript:void(0);" data-id='{{(index $info $PrimaryKey).Content}}'
											class="grid-row-delete"><i class="fa fa-trash"></i></a>
									{{end}}
								</td>
							</tr>
							{{end}}
						</tbody>
					</table>
					<script>
						window.selectedRows = function () {
							let selected = [];
							$('.grid-row-checkbox:checked').each(function () {
								selected.push($(this).data('id'));
							});
							return selected;
						};

						const selectedAllFieldsRows = function () {
							let selected = [];
							$('.column-select-item:checked').each(function () {
								selected.push($(this).data('id'));
							});
							return selected;
						};

						const pjaxContainer = "#pjax-container";
						const noAnimation = "__no_animation_";

						function iCheck(el) {
							el.iCheck({checkboxClass: 'icheckbox_minimal-blue'}).on('ifChanged', function () {
								if (this.checked) {
									$(this).closest('tr').css('background-color', "#ffffd5");
								} else {
									$(this).closest('tr').css('background-color', '');
								}
							});
						}

						$(function () {

							$('.grid-select-all').iCheck({checkboxClass: 'icheckbox_minimal-blue'}).on('ifChanged', function (event) {
								if (this.checked) {
									$('.grid-row-checkbox').iCheck('check');
								} else {
									$('.grid-row-checkbox').iCheck('uncheck');
								}
							});
							let items = $('.column-select-item');
							iCheck(items);
							iCheck($('.grid-row-checkbox'));
							let columns = getQueryVariable("__columns");
							if (columns === -1) {
								items.iCheck('check');
							} else {
								let columnsArr = columns.split(",");
								for (let i = 0; i < columnsArr.length; i++) {
									for (let j = 0; j < items.length; j++) {
										if (decodeURI(columnsArr[i]) === $(items[j]).attr("data-id")) {
											$(items[j]).iCheck('check');
										}
									}
								}
							}

							$('.filter-area').hide();


							let lastTd = $("table tr:last td:last div");
							if (lastTd.hasClass("dropdown")) {
								let popUpHeight = $("table tr:last td:last div ul").height();

								let trs = $("table tr");
								let totalHeight = 0;
								for (let i = 1; i < trs.length - 1; i++) {
									totalHeight += $(trs[i]).height();
								}
								if (popUpHeight > totalHeight) {
									let h = popUpHeight + 16;
									$("table tbody").append("<tr style='height:" + h + "px;'></tr>");
								}

								trs = $("table tr");
								for (let i = trs.length - 1; i > 1; i--) {
									let td = $(trs[i]).find("td:last-child div");
									let combineHeight = $(trs[i]).height() / 2 - 20;
									for (let j = i + 1; j < trs.length; j++) {
										combineHeight += $(trs[j]).height();
									}
									if (combineHeight < popUpHeight) {
										td.removeClass("dropdown");
										td.addClass("dropup");
									}
								}
							}

							let sort = getQueryVariable("__sort");
							let sort_type = getQueryVariable("__sort_type");

							if (sort !== -1 && sort_type !== -1) {
								let sortFa = $('#sort-' + sort);
								if (sort_type === 'asc') {
									sortFa.attr('href', '?__sort=' + sort + "&__sort_type=desc" + decodeURIComponent("{{.URLRoute.SortURL}}"))
								} else {
									sortFa.attr('href', '?__sort=' + sort + "&__sort_type=asc" + decodeURIComponent("{{.URLRoute.SortURL}}"))
								}
								sortFa.removeClass('fa-sort');
								sortFa.addClass('fa-sort-amount-' + sort_type);
							} else {
								let sortParam = decodeURIComponent("{{.URLRoute.SortURL}}");
								let sortHeads = $(".fa.fa-fw.fa-sort");
								for (let i = 0; i < sortHeads.length; i++) {
									$(sortHeads[i]).attr('href', $(sortHeads[i]).attr('href') + sortParam)
								}
							}
						});


						$('.column-select-all').on('click', function () {
							if ($(this).data('check') === '') {
								$('.column-select-item').iCheck('check');
								$(this).data('check', 'true')
							} else {
								$('.column-select-item').iCheck('uncheck');
								$(this).data('check', '')
							}
						});

						$('.column-select-submit').on('click', function () {

							let param = new Map();
							param.set('__columns', selectedAllFieldsRows().join(','));
							param.set(noAnimation, 'true');

							$.pjax({
								url: addParameterToURL(param),
								container: pjaxContainer
							});

							toastr.success('{{"刷新成功"}} !');
						});




						$('.grid-row-delete').click(function () {
							DeletePost($(this).data('id'))
						});

						$('.grid-batch-0').on('click', function () {
							let rows = selectedRows();
							if (rows.length > 0) {
								DeletePost(rows.join())
							}
						});

						function DeletePost(id) {
							swal({
									title: {{"確定要刪除資料嗎?"}},
									type: "warning",
									showCancelButton: true,
									confirmButtonColor: "#DD6B55",
									confirmButtonText: {{"確定"}},
									closeOnConfirm: false,
									cancelButtonText: {{"取消"}},
								},
								function () {
									$.ajax({
										method: 'post',
										url: {{.URLRoute.DeleteURL}},
										data: {
											id: id
										},
										success: function (data) {
											let param = new Map();
											param.set(noAnimation, "true");
											$.pjax({
												url: addParameterToURL(param),
												container: pjaxContainer
											});
											if (typeof (data) === "string") {
												data = JSON.parse(data);
											}
											if (data.code === 200) {
												$('#_TOKEN').val(data.data);
												let lastTd = $("table tr:last td:last div");
												if (lastTd.hasClass("dropdown")) {
													let popUpHeight = $("table tr:last td:last div ul").height();

													let trs = $("table tr");
													let totalHeight = 0;
													for (let i = 1; i < trs.length - 1; i++) {
														totalHeight += $(trs[i]).height();
													}
													if (popUpHeight > totalHeight) {
														let h = popUpHeight + 16;
														$("table tbody").append("<tr style='height:" + h + "px;'></tr>");
													}
												}
												swal(data.msg, '', 'success');
											} else {
												swal(data.msg, '', 'error');
											}
										},
										error: function (data) {
											if (data.responseText !== "") {
												swal(data.responseJSON.msg, '', 'error');
											} else {
												swal("{{"錯誤"}}", '', 'error');
											}
										},
									});
								});
						}


						function getQueryVariable(variable) {
							let query = window.location.search.substring(1);
							let vars = query.split("&");
							for (let i = 0; i < vars.length; i++) {
								let pair = vars[i].split("=");
								if (pair[0] === variable) {
									return pair[1];
								}
							}
							return -1;
						}

						function addParameterToURL(params) {
							let newUrl = location.href.replace("#", "");

							for (let [field, value] of params) {
								if (getQueryVariable(field) !== -1) {
									newUrl = replaceParamVal(newUrl, field, value);
								} else {
									if (newUrl.indexOf("?") > 0) {
										newUrl = newUrl + "&" + field + "=" + value;
									} else {
										newUrl = newUrl + "?" + field + "=" + value;
									}
								}
							}

							return newUrl
						}


						function replaceParamVal(oUrl, paramName, replaceWith) {
							let re = eval('/(' + paramName + '=)([^&]*)/gi');
							return oUrl.replace(re, paramName + '=' + replaceWith);
						}

						$(function () {

							$('.editable-td-select').editable({
								"type": "select",
								"emptytext": "<i class=\"fa fa-pencil\"><\/i>"
							});
							$('.editable-td-text').editable({
								emptytext: "<i class=\"fa fa-pencil\"><\/i>",
								type: "text"
							});
							$('.editable-td-datetime').editable({
								"type": "combodate",
								"emptytext": "<i class=\"fa fa-pencil\"><\/i>",
								"format": "YYYY-MM-DD HH:mm:ss",
								"viewformat": "YYYY-MM-DD HH:mm:ss",
								"template": "YYYY-MM-DD HH:mm:ss",
								"combodate": {"maxYear": 2035}
							});
							$('.editable-td-date').editable({
								"type": "combodate",
								"emptytext": "<i class=\"fa fa-pencil\"><\/i>",
								"format": "YYYY-MM-DD",
								"viewformat": "YYYY-MM-DD",
								"template": "YYYY-MM-DD",
								"combodate": {"maxYear": 2035}
							});
							$('.editable-td-year').editable({
								"type": "combodate",
								"emptytext": "<i class=\"fa fa-pencil\"><\/i>",
								"format": "YYYY",
								"viewformat": "YYYY",
								"template": "YYYY",
								"combodate": {"maxYear": 2035}
							});
							$('.editable-td-month').editable({
								"type": "combodate",
								"emptytext": "<i class=\"fa fa-pencil\"><\/i>",
								"format": "MM",
								"viewformat": "MM",
								"template": "MM",
								"combodate": {"maxYear": 2035}
							});
							$('.editable-td-day').editable({
								"type": "combodate",
								"emptytext": "<i class=\"fa fa-pencil\"><\/i>",
								"format": "DD",
								"viewformat": "DD",
								"template": "DD",
								"combodate": {"maxYear": 2035}
							});
							$('.editable-td-textarea').editable({
								"type": "textarea",
								"rows": 10,
								"emptytext": "<i class=\"fa fa-pencil\"><\/i>"
							});
							$(".info_edit_switch").bootstrapSwitch({
								onSwitchChange: function (event, state) {
									let obejct = $(event.target);
									let val = "";
									if (state) {
										val = obejct.closest('.bootstrap-switch').next().val();
									} else {
										val = obejct.closest('.bootstrap-switch').next().next().val()
									}
								}
							})
						});
					</script>
					<style>
						table tbody tr td {
							word-wrap: break-word;
							word-break: break-all;
						}
					</style>
				</div>
				<div class="box-footer clearfix">
					<ul class="pagination pagination-sm no-margin pull-right">
						<li class="page-item {{.PanelInfo.Paginator.PreviousClass}}">
							{{if eq .PanelInfo.Paginator.PreviousClass "disabled"}}
								<span class="page-link">«</span>
							{{else}}
								<a class="page-link" href='{{.PanelInfo.Paginator.PreviousURL}}' rel="next">«</a>
							{{end}}
						</li>
						{{range $key, $page := .PanelInfo.Paginator.Pages}}
							{{if eq (index $page "isSplit") "0"}}
								{{if eq (index $page "active") "active"}}
									<li class="page-item active"><span class="page-link">{{index $page "page"}}</span></li>
								{{else}}
									<li class="page-item"><a class="page-link" href='{{index $page "url"}}'>{{index $page "page"}}</a>
									</li>
								{{end}}
							{{else}}
								<li class="page-item disabled"><span class="page-link">...</span></li>
							{{end}}
						{{end}}
						<li class='page-item {{.PanelInfo.Paginator.NextClass}}'>
							{{if eq .PanelInfo.Paginator.NextClass "disabled"}}
								<span class="page-link">»</span>
							{{else}}
								<a class="page-link" href='{{.PanelInfo.Paginator.NextURL}}' rel="next">»</a>
							{{end}}
						</li>
					</ul>
					<label class="control-label pull-right" style="margin-right: 10px; font-weight: 100;">
						<small>{{"show"}}</small>&nbsp;
						{{$option := .PanelInfo.Paginator.Option}}
						{{$url := .PanelInfo.Paginator.URL}}
						<select class="input-sm grid-per-pager" name="per-page">
							{{range $key, $pageSize := .PanelInfo.Paginator.PageSizeList}}
								<option value="{{$url}}&__pageSize={{$pageSize}}" {{index $option $pageSize}}>
									{{$pageSize}}
								</option>
							{{end}}
						</select>
						<small>{{"entries"}}</small>
					</label>
					<script>
						let gridPerPaper = $('.grid-per-pager');
						gridPerPaper.on('change', function () {
							$.pjax({url: this.value, container: '#pjax-container'});
						});
					</script>
				</div>
			</div>
		</div>
	</section>
{{end}}`, "menu_content": `{{define "menu_content"}}
	<script src="/admin/assets/dist/js/datatable.min.581cdc109b.js"></script>
	<script src="/admin/assets/dist/js/form.min.f8678914e9.js"></script>
	<script src="/admin/assets/dist/js/treeview.min.7780d3bb0f.js"></script>
	<script src="/admin/assets/dist/js/tree.min.e1faf8b7de.js"></script>
	<section class="content-header">
		<h1>
		{{.PanelInfo.Title}}
		<small>{{.PanelInfo.Description}}</small>
		</h1>
		<ol class="breadcrumb" style="margin-right: 30px;">
			<li><a href={{.URLRoute.IndexURL}}><i class="fa fa-dashboard"></i> 首頁</a></li>
		</ol>
	</section>
	<section class="content">
		<div>
			{{if ne .AlertContent ""}}
				<div class="alert alert-warning alert-dismissible">
					<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>
					<h4>發生錯誤</h4>
					{{ .AlertContent}}
				</div>
			{{end}}
			<div class="row">
				<div class="col-md-6">
					<div class="box box-"}>
						<div class="box-header ">
							<div class="btn-group">
								<a class="btn btn-warning btn-sm tree-model-refresh"><i class="fa fa-refresh"></i>&nbsp;{{"刷新頁面"}}</a>
							</div>
							<div class="btn-group"></div>
						</div>
						<div class="box-body">
							<div class="dd" id="tree-model">
								{{$EditURL := .URLRoute.EditURL}}
								{{$URLPrefix := .URLRoute.URLPrefix}}
								<ol class="dd-list">
									{{range $key, $list := .Menu.List}}
										<li class="dd-item" data-id='{{$list.ID}}'>
											<div class="dd-handle">
												{{if eq $list.URL ""}}
													<i class="fa {{$list.Icon}}"></i>&nbsp;<strong>{{$list.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
															href="{{$list.URL}}" class="dd-nodrag">{{$list.URL}}</a>
												{{else if eq $list.URL "/"}}
													<i class="fa {{$list.Icon}}"></i>&nbsp;<strong>{{$list.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
															href="{{$URLPrefix}}" class="dd-nodrag">{{$URLPrefix}}</a>
												{{else if (isLinkURL $list.URL)}}
													<i class="fa {{$list.Icon}}"></i>&nbsp;<strong>{{$list.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
															href="{{$list.URL}}" class="dd-nodrag">{{$list.URL}}</a>
												{{else}}
													<i class="fa {{$list.Icon}}"></i>&nbsp;<strong>{{$list.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
															href="{{$URLPrefix}}{{$list.URL}}" class="dd-nodrag">{{$URLPrefix}}{{$list.URL}}</a>
												{{end}}
												<span class="pull-right dd-nodrag">
													<a href="{{$EditURL}}?id={{$list.ID}}"><i class="fa fa-edit"></i></a>
													<a href="javascript:void(0);" data-id="{{$list.ID}}" class="tree_branch_delete"><i
																class="fa fa-trash"></i></a>
												</span>
											</div>
											{{if gt (len $list.ChildrenList) 0}}
												<ol class="dd-list">
													{{range $key, $item := $list.ChildrenList}}
														<li class="dd-item" data-id='{{$item.ID}}'>
															<div class="dd-handle">
																{{if eq $item.URL ""}}
																	<i class="fa {{$item.Icon}}"></i>&nbsp;
																	<strong>{{$item.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
																			href="{{$item.URL}}" class="dd-nodrag">{{$item.URL}}</a>
																{{else if eq $item.URL "/"}}
																	<i class="fa {{$item.Icon}}"></i>&nbsp;
																	<strong>{{$item.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
																			href="{{$URLPrefix}}" class="dd-nodrag">{{$URLPrefix}}</a>
																{{else if (isLinkURL $item.URL)}}
																	<i class="fa {{$item.Icon}}"></i>&nbsp;
																	<strong>{{$item.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
																			href="{{$item.URL}}" class="dd-nodrag">{{$item.URL}}</a>
																{{else}}
																	<i class="fa {{$item.Icon}}"></i>&nbsp;
																	<strong>{{$item.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
																			href="{{$URLPrefix}}{{$item.URL}}"
																			class="dd-nodrag">{{$URLPrefix}}{{$item.URL}}</a>
																{{end}}
																<span class="pull-right dd-nodrag">
																	<a href="{{$EditURL}}?id={{$item.ID}}"><i class="fa fa-edit"></i></a>
																	<a href="javascript:void(0);" data-id="{{$item.ID}}"
																	class="tree_branch_delete"><i class="fa fa-trash"></i></a>
																</span>
															</div>
															{{if gt (len $item.ChildrenList) 0}}
																<ol class="dd-list">
																	{{range $key2, $subItem := $item.ChildrenList}}
																		<li class="dd-item" data-id='{{$subItem.ID}}'>
																			<div class="dd-handle">
																				{{if eq $subItem.URL ""}}
																					<i class="fa {{$subItem.Icon}}"></i>&nbsp;
																					<strong>{{$subItem.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
																							href="{{$subItem.URL}}"
																							class="dd-nodrag">{{$subItem.URL}}</a>
																				{{else if eq $subItem.URL "/"}}
																					<i class="fa {{$subItem.Icon}}"></i>&nbsp;
																					<strong>{{$subItem.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
																							href="{{$URLPrefix}}"
																							class="dd-nodrag">{{$URLPrefix}}</a>
																				{{else if (isLinkURL $subItem.URL)}}
																					<i class="fa {{$subItem.Icon}}"></i>&nbsp;
																					<strong>{{$subItem.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
																							href="{{$subItem.URL}}"
																							class="dd-nodrag">{{$subItem.URL}}</a>
																				{{else}}
																					<i class="fa {{$subItem.Icon}}"></i>&nbsp;
																					<strong>{{$subItem.Name}}</strong>&nbsp;&nbsp;&nbsp;<a
																							href="{{$URLPrefix}}{{$subItem.URL}}"
																							class="dd-nodrag">{{$URLPrefix}}{{$subItem.URL}}</a>
																				{{end}}
																				<span class="pull-right dd-nodrag">
																					<a href="{{$EditURL}}?id={{$subItem.ID}}"><i
																								class="fa fa-edit"></i></a>
																					<a href="javascript:void(0);" data-id="{{$subItem.ID}}"
																					class="tree_branch_delete"><i
																								class="fa fa-trash"></i></a>
																				</span>
																			</div>
																		</li>
																	{{end}}
																</ol>
															{{end}}
														</li>
													{{end}}
												</ol>
											{{end}}
										</li>
									{{end}}
								</ol>
							</div>
							<script data-exec-on-popstate="">
							$(function () {
								$('#tree-model').nestable([]);
								$('.tree_branch_delete').click(function () {
									let id = $(this).data('id');
									swal({
											title: {{"確定要刪除嗎?"}},
											type: "warning",
											showCancelButton: true,
											confirmButtonColor: "#DD6B55",
											confirmButtonText: {{"確定"}},
											closeOnConfirm: false,
											cancelButtonText: {{"取消"}}
										},
										function () {
											$.ajax({
												method: 'post',
												url: {{.URLRoute.DeleteURL}} +'?id=' + id,
												data: {},
												success: function (data) {
													$.pjax.reload('#pjax-container');
													if (data.code === 200) {
														swal(data.msg, '', "success");
													} else {
														swal(data.msg, '', "error");
													}
												},
												error: function (data) {
													if (data.responseText !== "") {
														swal(data.responseJSON.msg, '', 'error');
													} else {
														swal("{{"錯誤"}}", '', 'error');
													}
												},
											});
										});
								});
								$('.tree-model-save').click(function () {
									let serialize = $('#tree-model').nestable('serialize');
									$.post("", {
											_order: JSON.stringify(serialize)
										},
										function (data) {
											$.pjax.reload('#pjax-container');
											toastr.success('Save succeeded !');
										});
								});
								$('.tree-model-refresh').click(function () {
									$.pjax.reload('#pjax-container');
									toastr.success(toastMsg);
								});
								$('.tree-model-tree-tools').on('click', function (e) {
									let target = $(e.target),
										action = target.data('action');
									if (action === 'expand') {
										$('.dd').nestable('expandAll');
									}
									if (action === 'collapse') {
										$('.dd').nestable('collapseAll');
									}
								});
								$(".parent_id").select2({"allowClear": true, "placeholder": "Parent"});
								$(".roles").select2({"allowClear": true, "placeholder": "Roles"});
							});
							</script>
						</div>
					</div>
				</div>
				<div class="col-md-6">
				{{ template "form" . }}
				</div>
			</div>
		</div>
	</section>
{{end}}`, "form": `{{define "form"}}
	<div class="box box-"}>
		<div class="box-header with-border"> 
			<h3 class="box-title">請完整填寫表格</h3>
			
			{{if not .FormInfo.HideBackButton}}
				<div class="box-tools">
					<div class="btn-group pull-right" style="margin-right: 10px">
						<a href={{.URLRoute.PreviousURL}} class="btn btn-sm btn-default form-history-back"><i class="fa fa-arrow-left"></i>&nbsp;返回</a>
					</div>
				</div>
			{{end}}
		</div>
		<div class="box-body" style=" ">
			<form id={{.FormID}} action="{{.URLRoute.InfoURL}}" method="post" accept-charset="UTF-8" class="form-horizontal" pjax-container style="background-color: white;">
				<div class="box-body">

					<div class="box-body">
						<div class="fields-group">
							{{range $key, $data := .FormInfo.FieldList}}
								{{if $data.Hide}}
									<input type="hidden" name="{{$data.Field}}" value='{{$data.Value}}'>
								{{else}}
									<div class="form-group">
										{{if ne $data.Header ""}}
											<label for="{{$data.Field}}"
												class="col-sm-2 {{if $data.Must}}asterisk{{end}} control-label">{{$data.Header}}</label>
										{{end}}
										<div class="col-sm-8">
											{{if eq $data.FormType.String "default"}}
												<div class="box box-solid box-default no-margin">
												<div class="box-body" style="min-height: 40px;">
													{{$data.Value}}
												</div>
												</div>
												<input type="hidden" class="{{$data.Field}}" name="{{$data.Field}}" value='{{$data.Value}}'>
											{{else if eq $data.FormType.String "text"}}
												{{if $data.Editable}}
													<div class="input-group">
														<span class="input-group-addon"><i class="fa fa-pencil fa-fw"></i></span>
														<input {{if $data.Must}}required="1"{{end}} type="text" name="{{$data.Field}}" value='{{$data.Value}}'
															class="form-control {{$data.Field}}" placeholder="{{$data.Placeholder}}">
													</div>
												{{end}}
											{{else if eq $data.FormType.String "select"}}
												<select class="form-control {{.Field}} select2-hidden-accessible" style="width: 100%;" name="{{.Field}}[]"
												multiple="" data-placeholder="{{.Placeholder}}" tabindex="-1" aria-hidden="true"
												{{if not .Editable}}disabled="disabled"{{end}}>
													{{range $key, $v := .FieldOptions }}
														<option value='{{$v.Value}}' {{$v.SelectedLabel}}>{{if ne $v.TextHTML ""}}{{$v.TextHTML}}{{else}}{{$v.Text}}{{end}}</option>
													{{end}}
												</select>
													<script>
														$("select.{{.Field}}").select2({{.OptionExt}});
													</script>
											{{else if eq $data.FormType.String "select_single"}}
												<select class="form-control {{.Field}} select2-hidden-accessible" style="width: 100%;" name="{{.Field}}"
												data-multiple="false" data-placeholder="{{.Placeholder}}" tabindex="-1" aria-hidden="true"
													{{if not .Editable}}disabled="disabled"{{end}}>
														<option></option>
														{{range $key, $v := .FieldOptions }}
															<option value='{{$v.Value}}' {{$v.SelectedLabel}}>{{if ne $v.TextHTML ""}}{{$v.TextHTML}}{{else}}{{$v.Text}}{{end}}</option>
														{{end}}
												</select>
												<script>
													$("select.{{.Field}}").select2({{.OptionExt}});
												</script>
											{{else if eq $data.FormType.String "iconpicker"}}
												<div class="input-group">
													<span class="input-group-addon"><i class="fa"></i></span>
														{{if eq $data.Value ""}}
															<input style="width: 140px" type="text" name="{{$data.Field}}" value="fa-bars"
																class="form-control {{.Field}}"
																placeholder="{{"Input Icon"}}">
														{{else}}
															<input style="width: 140px" type="text" name="{{$data.Field}}" value="{{$data.Value}}"
																class="form-control {{.Field}}"
																placeholder="{{"Input Icon"}}">
														{{end}}
												</div>
												<script>
													$('.{{.Field}}').iconpicker({placement: 'bottomLeft'});
												</script>
											{{else if eq $data.FormType.String "password"}}
												{{if .Editable}}
													<div class="input-group">
														<span class="input-group-addon"><i class="fa fa-eye-slash"></i></span>
														<input {{if .Must}}required="1"{{end}} type="password" name="{{$data.Field}}"
															value="{{$data.Value}}"
															class="form-control {{.Field}}" placeholder="{{.Placeholder}}">
													</div>
												{{else}}
													<div class="box box-solid box-default no-margin">
														<div class="box-body">********</div>
													</div>
												{{end}}
											{{else if eq .FormType.String "selectbox"}}
												<select class="form-control {{.Field}}" style="width: 100%;" name="{{.Field}}[]" multiple="multiple"
												data-placeholder="Input {{.Header}}" {{if not .Editable}}disabled="disabled"{{end}}>
												{{range  $key, $v := .FieldOptions }}
													<option value='{{$v.Value}}' {{$v.SelectedLabel}}>{{if ne $v.TextHTML ""}}{{$v.TextHTML}}{{else}}{{$v.Text}}{{end}}</option>
												{{end}}
												</select>
												<script>
													$("select.{{.Field}}").bootstrapDualListbox({
														"infoText": "Showing all {0}",
														"infoTextEmpty": "Empty list",
														"infoTextFiltered": "{0} \/ {1}",
														"filterTextClear": "Show all",
														"filterPlaceHolder": "Filter"
													});
												</script>
											{{else if eq .FormType.String "textarea"}}
												<textarea {{if .Must}}required="1"{{end}} name="{{.Field}}" class="form-control" rows="5"
												placeholder="{{.Placeholder}}"
												{{if not .Editable}}disabled="disabled"{{end}}>{{.Value}}</textarea>
											{{else if eq $data.FormType.String "datetime"}}
												{{if not .Editable}}
													<div class="box box-solid box-default no-margin">
														<div class="box-body" style="min-height: 40px;">
															{{.Value}}
														</div>
														<input type="hidden" class="{{.Field}}" name="{{.Field}}" value='{{.Value}}'>
													</div>
												{{else}}
													<div class="input-group">
														<span class="input-group-addon"><i class="fa fa-calendar fa-fw"></i></span>
														<input {{if .Must}}required="1"{{end}} style="width: 170px" type="text"
															name="{{.Field}}"
															value="{{.Value}}"
															class="form-control {{.Field}}" placeholder="{{.Placeholder}}">
													</div>
													<script>
														$(function () {
															$('input.{{.Field}}').parent().datetimepicker({{.OptionExt}});
														});
													</script>
												{{end}}
											{{else if eq $data.FormType.String "datetime_range"}}
												{{if .Editable}}
													<div class="input-group">
														<span class="input-group-addon"><i class="fa fa-calendar fa-fw"></i></span>
														<input type="text" id="{{.Field}}_start" name="{{.Field}}_start" value="{{.Value}}"
															class="form-control {{.Field}}_start" placeholder="{{.Placeholder}}">
														<span class="input-group-addon" style="border-left: 0; border-right: 0;">-</span>
														<input type="text" id="{{.Field}}_end" name="{{.Field}}_end" value="{{.Value2}}"
															class="form-control {{.Field}}_end" placeholder="{{.Placeholder}}">
													</div>
													<script>
														$(function () {
															$('input.{{.Field}}_start').datetimepicker({{.OptionExt}});
															$('input.{{.Field}}_end').datetimepicker({{.OptionExt2}});
															$('input.{{.Field}}_start').on("dp.change", function (e) {
																$('input.{{.Field}}_end').data("DateTimePicker").minDate(e.date);
															});
															$('input.{{.Field}}_end').on("dp.change", function (e) {
																$('input.{{.Field}}_start').data("DateTimePicker").maxDate(e.date);
															});
														});
													</script>
												{{else}}
													<div class="box box-solid box-default no-margin">
														<div class="box-body">{{.Value}}</div>
													</div>
													<input type="hidden" class="{{.Field}}" name="{{.Field}}" value='{{.Value}}'>
												{{end}}
											{{end}}
											{{if ne .HelpMsg ""}}
												<span class="help-block">
													<i class="fa fa-info-circle"></i>&nbsp;{{.HelpMsg}}
												</span>
											{{end}}
										</div>
									</div>	           
								{{end}}
							{{end}}						
						</div>
					</div>
				</div>
				<div class="box-footer">
					<div class="col-md-2 "></div>
					<div class="col-md-8 ">
						<div class="btn-group pull-right" >
							<button type="submit" class="btn  btn-primary" data-loading-text="&nbsp;Save">新增
							</button>	
						</div>
						<div class="btn-group pull-left" >
							<button type="reset" class="btn  btn-warning" data-loading-text="&nbsp;Save">重置
							</button>
						</div>
					</div>
				</div>
				<input type="hidden" name="__previous_" value="{{.URLRoute.PreviousURL}}">
				<input type="hidden" name="__token_" value="{{.Token}}">
			</form>
		</div>
	</div>
{{end}}`,
}
