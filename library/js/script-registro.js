(function() {
    var registrarse = function () {
        registros.registrarse();
        registros.actualizaregistro();
        registros.consultaregistro();
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
                   if(response.status == 'ok'){
                        $('#tusdatos').html(
                            '<h2>Gracias por registrarte en "El buen Fin"</h2>'
                            + '<p><strong>Estás a un paso de activar tu cuenta y comenzar con el registro de tu empresa</strong></p>'
                            + '<p><strong>En un momento recibiras un correo electrónico con tus datos de acceso</strong></p>'
                            + '<div class="row"><div class="span6">'
                            + '<dl class="dl-horizontal">'
                            + '<dt>Usuario:</dt><dd>' + response.Nombre + '</dd>'
                            + '<dt>Apellido:</dt><dd>' + response.Apellidos + '</dd>'
                            + '<dt>Correo:</dt><dd>' + response.Email + '</dd>'
                            + '<dt>Contraseña:</dt><dd>' + response.Pass + '</dd>'
                            + '</dl>'
                            + '</div>'
                            + '<div class="span6">'
                            + '<ul><li>Ingresa a tu cuenta de correo</li><li>Abre el correo que te enviamos</li><li>Haz click en la liga que aparece en él</li><li>Cpmienza el registro de tu(s) empresa(s)</li></ul>'
                            + '<p><strong>Es importante que actives antes de 15 días ya que la activación expira. Si expira tu código o no es válido puedes realizar tu registro nuevamente</strong></p>'
                            + '</div></div>'
                            );
                       Ajax.hidePreload($('#registro-enviado-block'));
                           //alert('registrado correctamente');
                   }
                   else if (response.status == 'alreadyOnSession') {
                        $('#tusdatos').html(
                                '<div class="alert alert-block">'
                                +'<h4>Alto!</h4>'
                                +'Ya estas logueado con tu cuenta'
                                +'</div>'
                            );
                        Ajax.hidePreload($('#registro-enviado-block'));
                   }
                   else if (response.status == 'alreadyRegistered') {
                        $('#tusdatos').html(
                                '<div class="alert alert-block">'
                                +'<h4>Alto!</h4>'
                                +'Ya estas registrado'
                                +'</div>'
                            );
                        Ajax.hidePreload($('#registro-enviado-block'));
                   }
           });
        })
    };    
    var actualizaregistro = function () {
        $('form#actualizaregistro').on('submit', function(event){
            event.preventDefault();
            var post = $(this).serialize();
            Ajax.post('/r/wsr/post', post, function(response){
                  // if(response.status == 'ok'){

                           //alert('registrado correctamente');
                 //  }
           });
        })
    };
    var consultaregistro = function() {
        Ajax.get('/r/wsr/get', function(response) {
            $('#actualizaregistro').formParams(response);

        });
    };
    return {
        registrarse: registrarse,
        actualizaregistro: actualizaregistro,
        consultaregistro: consultaregistro
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
