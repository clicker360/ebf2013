$(document).ready(function() {
    var blobkey;
    var idemp = "qlixmfnbmngqarnccsky";
    var uploadurl;
	/* Este código comentado debe ir en el template html de lo contrario el 
	 * programa no puede planchar las variables de entorno
	 */
    $.get("/r/wsed/get", { IdEmp: ""+idemp+""}, function(resp) {
        console.log(resp);
        if(typeof(resp) != 'object') { resp = JSON.parse(resp); }
        if(resp.status=="ok") {
            blobkey = resp.BlobKey;
            uploadurl = resp.UploadUrl;
            $("#IdEmp").attr("value", resp.IdEmp);
            $("#enviar").attr('action', uploadurl);
            $("#BlobKey").attr("value", blobkey);
            $("#uploadimg_id").attr('value', idemp);
            $("#empresa").html(resp.Empresa);
            $("#descripcion").val(resp.Desc);
            $("#facebook").val(resp.Facebook);
            $("#twitter").val(resp.Twitter);
            if(blobkey) {
                updateimg(blobkey);
            } else {
                putDefault();
            }
        } else {
            putDefault();
        }

        $('#loader').hide();
    }, "json");

	$("#enviardata").submit(function() {
		return true;
	});

	$("#enviar").validationEngine({promptPosition : "topRight", scroll: false});
	$("#enviardata").validationEngine({promptPosition : "topRight", scroll: false});
	var $pic = $("#pic");
	

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
            console.log(data);
			var resp = "";
			switch (data.status) {
				case "invalidUpload": 
                    resp = "<p>Intente nuevamente, su imagen no puede ser integrada.</p>";
				case "uploadSessionError": 
                    resp = "<p>Favor de refrescar la página para continuar.</p>";
				case "notFound": 
                    resp = "<p>La oferta no existe.</p>";
				case "ok": 	
                    resp = "<p>La imagen se integró exitosamente</p>";
                    var uploadurl;
                    uploadurl = data.uploadurl;
                    $("#enviar").attr("action", uploadurl);
                    setTimeout(function(){ updateimg(data.blobkey); }, 1000); 	
                    break;
				default:  
                    resp = "<p>Intente nuevamente con una imagen de menor peso, su imagen no puede ser integrada.</p>";
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

    $('input[maxlength]').live('keyup blur', function() {
		var maxlength = $(this).attr('maxlength'); var val = $(this).val();
		if (val.length > maxlength) {
			$(this).val(val.slice(0, maxlength));
		}
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

function updateimg(blob) {
    blobkey = blob; // set blobkey global
    $("#BlobKey").attr("value", blobkey);
	if(blob) {
		$('#pic').remove();
		var query = "id="+blob + "&Avc=" + avoidCache();
		img = "<img  src = '/extraimg?"+ query +"' id='pic' width='256px' />" 
		$('#urlimg').append(img);
	} else {
		putDefault();
	}
}
