/*$(document).ready(function(){ $("#enviar").validationEngine({promptPosition : "topRight", scroll: true}); });*/
function activateCancel(){ $("#cancelbtn").addClass("show") }
function deactivateCancel(){ $("#cancelbtn").removeClass("show") }
$(document).ready(function() {
	$('#loader').hide();
	$('#estado').change(function(){
		$('#show_mun').fadeOut();
		$('#loader').show();
		getmsp()
		return false;
	});
	getmsp()

});
function getmsp() {
	$.get("/r/msp", {
		CveEnt: $('#estado').val(),
		CveMun: $('#municipio').val(),
	}, function(response){
		setTimeout("loadmun('show_mun', '"+escape(response)+"')", 400);
	});
}
function loadmun(id, response){
	$('#loader').hide();
	$('#'+id).html(unescape(response));
	$('#'+id).fadeIn();
} 
