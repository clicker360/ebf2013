// JavaScript Document
function validate(){
		var regmail = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
		var regname = /^[a-zA-Z0-9.\.\-\_"'# ]+$/;
		var regtel = /^[0-9\-]+$/;
		var regpass = /^[a-zA-Z0-9.\.  \-@&#"'()\[\]?¿!¡]+$/;
		var regrfc = /^[a-zA-Z]{3}[0-9]{6}[a-zA-Z0-9]{2,}$/
		var regnum = /^[0-9]+$/
		var valid = true;
		
    	if( !regrfc.test(document.formempresa.rfc.value) ){
			document.formempresa.rfc.className += " invalid";
			valid = false;
    	}
		if( !regname.test(document.formempresa.nomcom.value) ){
			document.formempresa.nomcom.className += " invalid";
			valid = false;
    	}
		if( !regpass.test(document.formempresa.razon.value) ){
			document.formempresa.razon.className += " invalid";
			valid = false;
    	}
		if( !regname.test(document.formempresa.delmun.value) ){
			document.formempresa.delmun.className += " invalid";
			valid = false;
    	}
		if( !regname.test(document.formempresa.calle.value) ){
			document.formempresa.calle.className += " invalid";
			valid = false;
    	}
		if( !regnum.test(document.formempresa.cp.value) ){
			document.formempresa.cp.className += " invalid";
			valid = false;
    	}
		if( !regnum.test(document.formempresa.sucursales.value) ){
			document.formempresa.sucursales.className += " invalid";
			valid = false;
    	}
		if( document.formempresa.slist2.value == "ORG" ){
			document.formempresa.slist2.className += " invalid";
			valid = false;
    	}
		if( !regname.test(document.formempresa.registro.value) ){
			document.formempresa.registro.className += " invalid";
			valid = false;
    	}
    	
    	return valid;
	}