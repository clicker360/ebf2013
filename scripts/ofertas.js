$(document).ready(function() {
	/* Este código comentado debe ir en el template html de lo contrario el 
	 * programa no puede planchar las variables de entorno
	 */
	/*$("#urlimg").attr('href', "{{with .FormDataOf}}{{.Url|js}}{{end}}");*/
	/*var idoft = "{{with .FormDataOf}}{{.IdOft|js}}{{end}}";*/
	/*var idblob = "{{with .FormDataOf}}{{.IdBlob|js}}{{end}}";*/

	if(idoft == 'none') {
		// se ocultan los campos que requieren IdOft
		putDefault();
		$('#imgform').hide();
		$('#modbtn').hide();
		$('#newbtn').show();
		$('#statuspub').attr("checked", true);
	} else {
		/* solo se actualizan estos datos si hay id de oferta */
		fillpcve(idoft, idemp);
		fillsucursales(idoft, idemp);
		$('#imgform').show();
		$('#modbtn').show();
		$('#newbtn').hide();	
	}

	if(idblob != "") {
		updateimg(idblob);
	} else {
		putDefault();
	}
	$('#loader').hide();
	$('#urlreq').hide();
	$('#placereq').hide();
	$('#tituloreq').hide();
	$('#enlinea').live('change', function() { 
		if($('#enlinea').attr('checked')) {
			$('#muestraurl').show();
		} else {
			$('#muestraurl').hide();
		}
	});

	if($('#enlinea').attr('checked')) {
		$('#muestraurl').show();
	} else {
		$('#muestraurl').hide();
	}

	$("#url").blur(function() {
		if($('#enlinea').attr('checked') && $('#url').val()=='') { $('#urlreq').show(); } else {$('#urlreq').hide();}
	});

	$("#oferta").blur(function() {
		var t = $('#oferta').val().toLowerCase();
		var ts = t.split(" ");
		var hit = 0;
		for (var i in ts) {
			console.log(ts[i])
			if(ts[i].replace(/^\s+|\s+$/g,"").indexOf("nueva") != -1) { hit++; }
			if(ts[i].replace(/^\s+|\s+$/g,"").indexOf("oferta") != -1) { hit++; }
			console.log(hit)
		}
		if(t == "" || hit > 1) { 
			$('#tituloreq').show(); console.log("nok");
			$('#tituloreq').goTo();
		} else { 
			$('#tituloreq').hide(); 
		}
	});

	$("#enviardata").submit(function() {
		/* 
		 * Manejo de sucursales
		 */
		var sucs = $("#listasuc").find("input");
		var pcves = $("#unpickpcve").find("a");
		var chain = ""; var sep = "";
		sucs.each(function() {
			if($(this).is(':checked')) {
				if($(this).attr('id') != "todassuc") {
					// Se arma una cadena de ID's
					chain += sep+$(this).attr('id');
					sep = " ";
				}
			}
		});
		$("#schain").val(chain);
		var chain = ""; var sep = "";
		pcves.each(function() { 
			chain += sep+$(this).text(); sep = " "; 
		});
		$("#pchain").val(chain);

		/* 
		 * Validaciones
		 */
		var ok = false;
		var sucs = $("#listasuc").find("input");
		if($("#enlinea").is(':checked')) {
			ok = true;
		} else {
			sucs.each(function() { if($(this).is(':checked')) { ok = true; } });
		}
		if(!ok) { $('#placereq').show(); return false; } 
		if($('#enlinea').attr('checked') && $('#url').val()=='') { $('#urlreq').show(); return false; } 
		$('#placereq').hide(); 
		$('#urlreq').hide(); 

		/* verifica el titulo de oferta */
		var t = $('#oferta').val().toLowerCase();
		var ts = t.split(" ");
		var hit = 0;
		for (var i in ts) {
			if(ts[i].replace(/^\s+|\s+$/g,"").indexOf("nueva") != -1) { hit++; }
			if(ts[i].replace(/^\s+|\s+$/g,"").indexOf("oferta") != -1) { hit++; }
		}
		if(t == "" || hit > 1) { 
			$('#tituloreq').show(); console.log("nok");
			$('#tituloreq').goTo();
			return false;
		} else { 
			$('#tituloreq').hide(); 
		}
		return true;
	});
	$("#enviar").validationEngine({promptPosition : "topRight", scroll: false});
	$("#enviardata").validationEngine({promptPosition : "topRight", scroll: false});
	var $pic = $("#pic");
	var $urlimg = $("#urimg");
	var max_size=400;
	

	/* 
	 * Ajax FORM para imagen de oferta
	 */
	var bar = $('.bar');
	var percent = $('.percent');
	var status = $('#status');
	var img;
	   
	$('#enviar').ajaxForm({
		dataType: 'json',
		beforeSend: function() {
			status.empty();
			var percentVal = '0%';
			bar.width(percentVal)
			percent.html(percentVal);
		},
		uploadProgress: function(event, position, total, percentComplete) {
			var percentVal = percentComplete + '%';
			bar.width(percentVal)
			percent.html(percentVal);
		},
		success: function(data) {
			var resp = "";
			switch (data.errstatus) {
				case "invalidUpload": resp = "<p>Intente nuevamente, su imagen no puede ser integrada.</p>";
				case "uploadSessionError": resp = "<p>Favor de refrescar la página para continuar.</p>";
				case "invalidId": resp = "<p>La oferta no existe.</p>";
				case "ok": 	resp = "<p>La imagen se integró exitosamente</p>";
							var uploadurl;
							(data.uploadurl.substr(0,5) != "https") ? uploadurl == data.uploadurl.replace("http","https") : uploadurl = data.uploadurl;
							$("#enviar").attr("action", uploadurl);
							setTimeout(function(){ updateimg(data.idblob); }, 1000); 	
							break;
				default:  resp = "<p>Intente nuevamente con una imagen de menor peso, su imagen no puede ser integrada.</p>";
			}
			status.html(resp);
		},
		error: function() {
			status.html("<p>Intente nuevamente con una imagen de menor peso, su imagen no puede ser integrada.</p>");
		}
	}); 

	$("#pic").error(function() { putDefault()});

	$('textarea[maxlength]').live('keyup blur', function() {
		var maxlength = $(this).attr('maxlength'); var val = $(this).val();
		if (val.length > maxlength) {
			$(this).val(val.slice(0, maxlength));
		}
	});

	$('oferta').live('keyup blur', function() {
		var maxlength = $(this).attr('maxlength'); var val = $(this).val();
		if (val.length > maxlength) {
			$(this).val(val.slice(0, maxlength));
		}
	});

	/* Palabras clave */
	$('#pcvepicker').on("click", "a", function(e){
		var token = $(this);
		if($(this).attr("value") == "0") {
			if($("#unpickpcve a").length < 5) {
				$.get("/r/addword", { token: ""+token.text()+"", id: ""+idoft+"" }, function(resp) {
					if(typeof(resp) != 'object') { resp = JSON.parse(resp); }
					if(resp.status=="ok") {
						token.attr("class", "wordselected");	
						token.attr("value", resp.id);	
						$('#unpickpcve').append(token);
					} else if(resp.status=="notFound") {
						/* No existe la oferta */
					} else if(resp.status=="invalidText") {
						/* El texto es inválido */
						alert("La palabra clave contiene carácteres no permitidos");
					} else if(resp.status=="writeErr") {
						alert("Hay problemas de conexión. Intente agregar de nuevo la palabra clave");
					}
				}, "json");
			} else {
				alert("Máximo 5 palabras");
			}
		} else if($(this).attr("value") == "E") {
			$.get("/r/rmword", { id: ""+idemp+"", token: ""+token.text()+"" }, function(resp) {
				if(typeof(resp) != 'object') { resp = JSON.parse(resp); }
				if(resp.status=="ok") {
					token.remove();
				}
			}, "json");
		} else {
			$.get("/r/delword", { id: ""+token.attr('value')+"", token: ""+token.text()+"" }, function(resp) {
				if(typeof(resp) != 'object') { resp = JSON.parse(resp); }
				if(resp.status=="ok") {
					token.attr("class", "sugestWord");	
					token.attr('value','0');
					$('#pickpcve').append(token);
				} else {
					alert("Hay problemas de conexión. Intente eliminar de nuevo la palabra clave");
				}
			}, "json");
		}
	});
	
	$("#nuevapcve").click(function(e) {
		var token = $("#tokenpcve");
		if(token.val().length >= 3) {
			if($("#unpickpcve a").length < 5) {
				$.get("/r/addword", { token: ""+token.val()+"", id: ""+idoft+"" }, function(resp) {
					if(typeof(resp) != 'object') { resp = JSON.parse(resp); }
					if(resp.status=="ok") {
						clearpcve();
						fillpcve(idoft,idemp);
					} else {
						alert("Hay problemas de conexión. Intente agregar de nuevo la palabra clave o recargue la página");
					}
				}, "json");
			} else {
				alert("Máximo 5 palabras");
			}
			$("#tokenpcve").val("");
		}
		$("#tokenpcve").trigger("blur");
	});

	/*
	 * Selecciona todas las sucursales
	 */
	$("#todassuc").click(function(e) {
		var sucs = $("#listasuc").find("input");
		sucs.each(function(ii, obj) {
			if($("#todassuc").is(':checked')) {
				if($(this).is('id')!="todassuc") $(this).attr('checked', true);
			} else {
				if($(this).is('id')!="todassuc") $(this).attr('checked', false);
			}
		});
   	});
	var elimActive= false;
	$("#elimilink").click(function(e){
		if (!elimActive) {
			$('.sugestWord').addClass("eliminateWord"); 
			$('.sugestWord').attr("value", "E"); 
			elimActive=true;

			$('#elimilink').addClass("button"); 
			$('#elimilink').addClass("red"); 
			$('#elimilink').addClass("small"); 
			
			$('#unpickpcve').addClass('hide');
			$('#tokenpcve').addClass('hide');
			$('#btnaddWords').addClass('hide');
			$('#elimilink').html("<span>Terminar </span>");
			$('#titleWord').html('Clic aquí para remover palabras que no necesites');  
		} else {
			$('.sugestWord').attr("value", "0"); 
			$('.eliminateWord').removeClass("eliminateWord"); 

			$('#elimilink').removeClass("button"); 
			$('#elimilink').removeClass("red"); 
			$('#elimilink').removeClass("small"); 

			$('#unpickpcve').removeClass('hide');
			$('#tokenpcve').removeClass('hide');
			$('#btnaddWords').removeClass('hide');
			$('#elimilink').html("Eliminar palabras capturadas");		
			$('#titleWord').html('Otras palabras que has capturado'); 
			elimActive=false;
		}
		return false;
	});
});/* termina onload */

function avoidCache(){
	var numRam = Math.floor(Math.random() * 500);
	return numRam;
}		

function putDefault() {
	$('#pic').remove();
	img = "<img  src = 'imgs/imageDefault.jpg' id='pic' width='258px' />" 
	$('#urlimg').append(img);
}

function updateimg(idblob) {
	if(idoft != 'none') {
		$('#pic').remove();
		var query = "id="+idblob + "&Avc=" + avoidCache();
		img = "<img  src = '/ofimg?"+ query +"' id='pic' width='256px' />" 
		$('#urlimg').append(img);
	} else {
		putDefault();
	}
}

/*
 * Llena palabras clave por oferta y empresa
 */
function fillpcve(idoft, idemp) {
	$.get("/r/wordsxo", { id: "" + idoft + ""})
	.success(function(data) {
		if(typeof(data) != 'object') { data = JSON.parse(data); }
		if($.isArray(data)) {
			$.each(data, function(i,item){
				var anchor = "<a href=\"#null\" class=\"wordselected\" id=\"pcve_"+item.token+"\" value=\""+item.id+"\">"+item.token+"</a>"
				$('#unpickpcve').append(anchor);
			});
		}
	})
	.error(function(){alert('Hay problemas de conexión, espere un momento y recargue la página');})
	.complete(function(){});

	$.get("/r/wordsxe", { id: "" + idemp + ""})
	.success(function(data) {
		if(typeof(data) != 'object') { data = JSON.parse(data); }
		if($.isArray(data)) {
			$.each(data, function(i,item){
				// Si en el ajax anterior no se añadio algo, aquí se añade como no seleccionado
				if($("#pcve_"+item.token).length == 0) {
					var anchor = "<a href=\"#null\" class=\"sugestWord pcve\" id=\"pcve_"+item.token+"\" value=\"0\">"+item.token+"</a>"
					$('#pickpcve').append(anchor);
				}
			});
		}
	})
	.error(function(){alert('Hay problemas de conexión, espere un momento y recargue la página');})
	.complete(function(){});
}

function clearpcve() {
	$("#unpickpcve").empty();
	$("#pickpcve").empty();
}

/*
 * Llena las sucursales
 */
function fillsucursales(idoft, idemp) {
	$.get("/r/ofsuc", { idoft: "" + idoft + "", idemp: "" + idemp + ""})
	.success(function(data) {
		if(typeof(data) != 'object') { data = JSON.parse(data); }
		$.each(data, function(i,item){
			var div = "<div class=\"gridsubRow bg-Gry2\"><label class=\"col-5 marg-L10pix\">"+item.sucursal+"</label><input name=\""+item.idsuc+"\" type=\"checkbox\" class=\"last marg-U5pix marg-R10pix\" id=\""+item.idsuc+"\"/></div>";
			$('#listasuc').append(div);
			if(item.idoft!="") {
				$("#"+item.idsuc).attr('checked', true);
			} else {
				$("#"+item.idsuc).attr('checked', false);
			}
		});
	})
	.error(function(){alert('Hay problemas de conexión, espere un momento y recargue la página');})
	.complete(function(){});

	/* scroll to */
	(function($) {
	    $.fn.goTo = function() {
			$('html, body').animate({
				scrollTop: $(this).offset().top + 'px'
			}, 'fast');
			return this; 
		}
	})(jQuery);
}

function activateCancel(){ $("#cancelbtn").addClass("show") }
function deactivateCancel(){ $("#cancelbtn").removeClass("show") }

/*
$("#"+item.idsuc).change(function() { 
	if($(this).is(':checked')) {
		$.get("/r/addofsuc", { idoft: "" + idoft + "", idemp: "" + idemp + "", idsuc: "" + item.idsuc + ""}, function(data) { })
		if(typeof(data) != 'object') { data = JSON.parse(data); }
		.success(function(){})
		.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
	} else {
		$.get("/r/delofsuc", { idoft: "" + idoft + "", idsuc: "" + item.idsuc + ""}, function(data) { })
		if(typeof(data) != 'object') { data = JSON.parse(data); }
		.success(function(){})
		.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
	}
});
*/
