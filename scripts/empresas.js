/*$(document).ready(function(){ $("#enviar").validationEngine({promptPosition : "topRight", scroll: true}); });*/
function activateCancel(){ $("#cancelbtn").addClass("show") }
function deactivateCancel(){ $("#cancelbtn").removeClass("show") }
$(document).ready(function() {
 	$('#urlreq').hide()
	$('#loader').hide();
	$('#CveEnt').change(function(){
		$('#show_mun').fadeOut();
		$('#loader').show();
		getmsp()
		return false;
	});
	getmsp()

	$("#url1").blur(function() {
		if($('#quest1').attr('checked') && $('#url1').val()=='') { $('#urlreq').show(); } else {$('#urlreq').hide();}
		//alert($('#url1').val());
	});
	$("#enviar").submit(function() {
		if($('#quest1').attr('checked') && $('#url1').val()=='') { $('#urlreq').show(); return false; } else {$('#urlreq').hide(); return true;}
	});
	//$(".customDialog").easyconfirm({dialog: $("#question")});
	//$(".customDialog").easyconfirm({locale: $.easyconfirm.locales.esMX});
});
function getmsp() {
	$.get("/r/msp", {
		CveEnt: $('#CveEnt').val(),
		CveMun: $('#CveMun').val(),
	}, function(response){
		setTimeout("loadmun('show_mun', '"+escape(response)+"')", 400);
	});
}
function loadmun(id, response){
	$('#loader').hide();
	$('#'+id).html(unescape(response));
	$('#'+id).fadeIn();
} 
