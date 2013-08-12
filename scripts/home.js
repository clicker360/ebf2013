$(document).ready(function(){
		tgle1();
		$("#enviar").validationEngine({promptPosition : "topRight", scroll: true});
		$("#acceso").validationEngine({promptPosition : "topRight", scroll: false});
		if(getURLParameter("login") == '1') activeLghtbx();
});
function tgle1(){ $( "#toggler1" ).delay(800).animate({ height: 284, opacity: 1 }, 1000, "linear", function() { tgle2(); }); }
function tgle2(){ $( "#toggler2" ).animate({ height: 284, opacity: 1 }, 800, "linear", function() { tgle3(); }); }
function tgle3(){ $( "#toggler3" ).animate({ height: 284, opacity: 1 }, 600, "linear", function() { tgle3(); }); }
function activeLghtbx(){ $('#light').removeClass("hidden"); $('#box').removeClass("hidden"); }
function deactiveLghtbx(){ $('#light').addClass("hidden"); $('#box').addClass("hidden"); }
function validateEmail(field, rules, i, options) {
	var err=0;
	$.each( field.val().split(","), function(i,candidate) { 
		if($.trim(candidate) != "") {
			if(!$.trim(candidate).match(options.allrules.email.regex)) err++;
		}
	});
	if(err) return options.allrules.email.alertText+". Puede introducir varios correos separados por coma";
}

function getURLParameter(name) {
		return decodeURIComponent((new RegExp('[?|&]' + name + '=' + '([^&;]+?)(&|#|;|$)').exec(location.search)||[,""])[1].replace(/\+/g, '%20'))||null;
}

