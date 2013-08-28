(function() {
    var registrarse = function () {
        registros.registrarse();
    };
    var execute = function() {
        $(document).ready(function() {
            registrarse();
        });
    };
    return execute();
})();
var Ajax = (function() {
    var _showPreload = function() {
        $('#preloader').fadeIn('fast');
    };
    var hidePreload = function(bloque) {
        $('#preloader').fadeOut('fast');
        if(typeof bloque !== 'undefined'){
            $('.active').removeClass('active').addClass('inactive');
            bloque.addClass('active').removeClass('inactive');   
        }
        

    };
    var get = function(url, callback) {
        _showPreload();
        $.ajax({
            url: url,
            dataType: 'json',
            success: function(response) {
                if (typeof callback === 'function') {
                    callback.call(callback, response);
                }
            }
        });
    };
    var post = function(url, data, callback) {
        _showPreload();
        $.ajax({
            type: 'POST',
            url: url,
            dataType: 'json',
            data: data,
            success: function(response) {
                if (typeof callback === 'function') {
                    callback.call(callback, response);
                }
            }
        });
    };
    return {
        hidePreload: hidePreload,
        get: get,
        post: post
    };
})();


var registros = (function() {
    var registrarse = function () {
        $('form#registroform').on('submit', function(event){
            event.preventDefault();
            var post = $(this).serialize();
            Ajax.post('/r/wsr/put', post, function(response){
                   if(response.success){
                           alert('registrado correctamente');
                   }
           });
        })
    };
    return {
        registrarse: registrarse
    };
})();

// Función para validar email alternativo, permite meter varios email separados por coma
function validateEmail(field, rules, i, options) {
	var err=0;
	$.each( field.val().split(","), function(i,candidate) { 
		if($.trim(candidate) != "") {
			if(!$.trim(candidate).match(options.allrules.email.regex)) err++;
		}
	});
	if(err) return options.allrules.email.alertText+". Puede introducir varios correos separados por coma";
}
