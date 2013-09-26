/**
 * Métodos autoejecutables. Estos métodos se ejecutan al cargar la página.
 * @returns {function} execute - Retorna la ejecución de los métodos registrados en este método.
 */
(function() {
    /**
     * registrarse Método que ejecuta la acción del registro de usuario.
     */
    var registrarse = function() {
        registros.registrarse();
    };

    /* ---------------- Empresas ------------- */

    /**
     * initEmpresas Método que llena el listado de empresas del usuario.
     */
    var initEmpresas = function() {
        empresas.initEmpresas(); // lista de empresas
    };

    /**
     * micrositios Método que ejecuta la accion del registro de un micrositio.
     */
    var micrositios = function() {
        micrositio.extrasformulario();
    };

    /**
     * llenarformempresas Método que llena y muestra el formulario de modificar
     * empresa.
     */
    var llenaformempresas = function() {
        $(document).on('click', 'a.editar-empresa', function(event) {
            event.preventDefault();
            var empresaID = $(this).attr('rel');
            empresas.empresaformdesdejson(empresaID);
            
        });
    };

    /**
     * nuevaempresa Método que muestra el formulario de nueva empresa.
     */
    var nuevaempresa = function() {
        $(document).on('click', 'a.nuevaempresa', function(event) {
            event.preventDefault();
            empresas.empresaformdesdejson();
        });
    };

    /**
     * modificaempresa Método que envia la información de empresa, para ser
     * creada o actualizada.
     */
    var modificaempresa = function() {
        empresas.empresa_envia();
    };

    var eliminarEmpresa = function() {
        $(document).on('click', '.eliminarempresa', function(event) {
            event.preventDefault();
            var idEmpresa = $(this).attr('rel');
            if (confirm('¿Seguro que desea eliminar esta empresa?')) {
                empresas.eliminarEmpresa(idEmpresa);
            }
        });
    };
    /* ---------------- Fin Empresas ------------- */

    /* ---------------- Sucursales ------------- */
    /**
     * initSucursales
     * Método que muestra la lista de sucursales, segun la empresa seleccionada.
     */
    var initSucursales = function() {
        $(document).on('click', 'a.ver-sucursales', function(event) {
            event.preventDefault();
            var rel = $(this).attr('rel');
            sucursales.initSucursales(rel); // este
            $('#añadir-suc-empresa').html('<a class="nuevasucursal span12 btn btn-success" rel="' + rel + '" ><i class="icon-plus"></i> Añadir nueva sucursal</a>');

        });
    };

    /**
     * llenaformsucursal
     * Método que llena y muestra el formulario de modificar sucursal
     */
    var llenaformsucursal = function() {
        $(document).on('click', 'a.editar-sucursal', function(event) {
            event.preventDefault();
            var rel = $(this).attr('rel');
            sucursales.sucursalformdesdejsonModifica(rel);
        });
    };

    /**
     * nuevasucursal
     * Método que muestra el formulario de nueva sucursal
     */
    var nuevasucursal = function() {
        $(document).on('click', 'a.nuevasucursal', function(event) {
            event.preventDefault();
            var rel = $(this).attr('rel');
            sucursales.sucursalformdesdejsonNueva(rel);
        });
    };

    /**
     * modificasucursal
     * Método que envia los datos de sucursal, para ser creada o actualizada.
     */
    var modificasucursal = function() {
        sucursales.sucursal_envia();
    };

    var eliminarSucursal = function() {
        $(document).on('click', '.eliminarsucursal', function(event) {
            event.preventDefault();
            var idSucursal = $(this).attr('data-idsuc');
            var idEmpresa = $(this).attr('data-idemp');
            if (confirm('¿Seguro que desea eliminar esta sucursal?')) {
                sucursales.eliminarSucursal(idSucursal, idEmpresa);
            }
        });
    };
    /* ---------------- Termina Sucursales ------------- */

    /* ---------------- Inicia Micrositios ------------- */

    /**
     * uploadImage
     * Método que permite subir la imagen de logotipo del micrositio.
     */
    var uploadImage = function() {
        $('#submitImage').on('click', function(event) {
            event.preventDefault();
            $('#imgfile').val($('#infile').val());
            micrositio.enviarimagen();
        });
    };
    /**
     * registerMicrositio
     * Método que registra los datos del formulario de micrositios.
     */
    var registerMicrositio = function() {
        micrositio.enviarDatos();
    };

    /* ---------------- Termina Micrositios ------------ */

    /* ---------------- Inicia Ofertas ------------- */
    /**
     * initOfertas 
     * Método que llena el listado de ofertas del usuario x empresa.
     */
    var initOfertas = function() {
        $(document).on('click', 'a.ver-ofertas', function(event) {
            event.preventDefault();
            var rel = $(this).attr('rel');
            ofertas.initOfertas(rel); // este
            $('#añadir-oferta-empresa').html('<a class="nuevaoferta span12 btn btn-success" rel="' + rel + '" ><i class="icon-plus"></i> Añadir nueva oferta</a>');

        });
    };

    var showOfertaForm = function() {
        $(document).on('click', 'a.nuevaoferta, a.editar-oferta', function(event) {
            event.preventDefault();
            var rel = $(this).attr('rel');
            var esNueva = false;
            if ($(this).hasClass('nuevaoferta'))
                esNueva = true;
            ofertas.showOfertaForm(rel, esNueva);
        });
    };

    var registrarOfertaPaso1 = function() {
        ofertas.enviarDatosBasicos();
    };
    var addWords = function() {
        $('#AddWordButton').on('click', function(event) {
            var idOferta = $(this).attr('rel');
            ofertas.agregarPalabra(idOferta);
        });
    };

    var eliminarOferta = function() {
        $(document).on('click', '.eliminaroferta', function(event) {
            event.preventDefault();
            var idOferta = $(this).attr('data-idoft');
            var idEmpresa = $(this).attr('data-idemp');
            if (confirm('¿Seguro que desea eliminar esta oferta?'))
                ofertas.eliminarOferta(idOferta, idEmpresa);
        });
    };

    var eliminarPalabra = function() {
        $(document).on('click', '.eliminarpalabra', function(event) {
            event.preventDefault();
            var token = $(this).parent().find('span').html();
            var idOferta = $(this).parent().attr('data-idoft');
            ofertas.eliminarPalabra(idOferta, token);
        });
    };

    // var toggleSucursalOft = function() {
    // 	$(document).on('click','input[name="sucursales[]"]',function(event){
    // 		// event.preventDefault();
    // 		var idEmpresa = $(this).attr('data-idemp');
    // 		var idOferta = $(this).attr('data-idoft');
    // 		var idSucursal = $(this).val();
    // 		if($(this).is(':checked')){
    // 			ofertas.agregarSucursal(idEmpresa, idOferta, idSucursal);
    // 		}else{
    // 			ofertas.eliminarSucursal(idEmpresa, idOferta, idSucursal);
    // 		}
    // 	});
    // }

    var uploadImageOferta = function() {
        $('#submitImageOferta').on('click', function(event) {
            event.preventDefault();
            $('#imgfileoferta').val($('#infileoferta').val());
            ofertas.enviarimagen();
        });
    };
    /* ---------------- Termina Ofertas ------------ */
    /**
     * execute
     * Registro de métodos que se ejecutarán automáticamente cuando se cargue la página.
     */
    var execute = function() {
        $(document).ready(function() {
            registrarse();
            initEmpresas();
            initSucursales();
            llenaformempresas();
            llenaformsucursal();
            nuevaempresa();
            nuevasucursal();
            modificaempresa();
            modificasucursal();
            uploadImage();
            registerMicrositio();
            initOfertas();
            showOfertaForm();
            registrarOfertaPaso1();
            addWords();
            eliminarOferta();
            eliminarEmpresa();
            eliminarSucursal();
            eliminarPalabra();
            uploadImageOferta();
            // toggleSucursalOft();
        });
    };
    return execute();
})();

/**
 * Ajax
 * Objeto que ejecuta y procesa las llamadas Ajax.
 * @returns {object} Ajax - Regresa el registro de todos los métodos públicos para ser usados de la forma "Ajax.método".
 */
var Ajax = (function() {

    /**
     * _showPreload
     * Método privado que muestra la precarga al iniciar una llamada Ajax.
     */
    var _showPreload = function() {
        $('#preloader').fadeIn('fast');
    };

    /**
     * hidePreload
     * Método que oculta la precarga al terminar de procesar la respuesta Ajax, y muestra el objeto que se le pasa como parámetro.
     * @param {object} bloque - Bloque de html que se mostrará una vez procesada la respuesta.
     */
    var hidePreload = function(bloque) {
        $('#preloader').fadeOut('fast');
        if (typeof bloque !== 'undefined') {
            $('.activo').removeClass('activo').addClass('inactivo');
            bloque.addClass('activo').removeClass('inactivo');
        }
    };

    /**
     * get
     * Método que hace una llamada Ajax para obtener datos.
     * @param {string} url - Url a la cual se hará la petición.
     * @param {function} callback - Método que se ejecutará cuando se reciba la respuesta.
     */
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

    /**
     * post
     * Método que envía datos por medio de Ajax
     * @param {string} url - Url a la cual se enviarán los datos
     * @param {object} data - Datos a ser enviados.  Puede ser una serialización de un formulario o en formato json.
     * @param {function} callback - Método que se ejecutará al recibir la respuesta.
     */
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

    // Registro de métodos públicos.
    return {
        hidePreload: hidePreload,
        get: get,
        post: post
    };
})();

var catalogos = (function() {
})();

/**
 * empresas
 * Objeto que manipula la información de las empresas
 * @returns {object} empresas - Regresa el registro de todos los métodos públicos para ser usados de la forma "empresas.método".
 */
var empresas = (function() {
    /**
     * initEmpresas
     * Método que obtiene los datos para llenar el listado de empresas.
     */
    var initEmpresas = function() {
        var imprimeTemplateDash = '',
                underTemplateIDash = $('#empresaTemplateDash').html(),
                underTemplateDash = _.template(underTemplateIDash);
        var imprimeTemplate = '',
                underTemplateID = $('#empresaTemplate').html(),
                underTemplate = _.template(underTemplateID);

        Ajax.get('/r/wse/gets', function(response) {
            imprimeTemplateDash = underTemplateDash({
                empresasArrayDash: response
            });
            imprimeTemplate = underTemplate({
                empresasArray: response
            });
            $('#empresasBlockDash').html(imprimeTemplateDash);
            $("#empresasBlock").html(imprimeTemplate);
            Ajax.hidePreload();
        });
    };

    /**
     * empresaformdesdejson
     * Método llena y muestra el formulario de empresa.
     * @param {string} codeRel - Identificador de la empresa.
     */
    var empresaformdesdejson = function(codeRel) {
        // Si existe el parámetro codeRel, llena y muestra el formulario de modificar empresa.
        if (codeRel) {
            $('#btn-empresa').html("Siguiente paso");
            Ajax.get('/r/wse/get?IdEmp=' + codeRel, function(response) {
                $('#empresa-form').formParams(response, true);
                llenamuniEmpresa(response.DirEnt, response.DirMun);
                llenaorganismos(response.OrgEmp);
                Ajax.hidePreload($('#empresas-detalle'));
            });
            // De lo contrario, muestra el formulario de empresa nueva.
        } else {
            $('#btn-empresa').html("Crear");
            llenamuniEmpresa("01");
            llenaorganismos();
            Ajax.hidePreload($('#empresas-detalle'));
        }
    };

    /**
     * empresa_envia
     * Método que envia los datos del formulario de empresa, según sea para crear una nueva o para modificarla.
     */
    var empresa_envia = function() {
        $(document).on('submit', 'form#empresa-form', function(event) {
            event.preventDefault();
            var post = $(this).serialize();
            // Si el formulario es para crear, los datos se envían al método "put" de la API.
            if ($('#btn-empresa').html() === 'Crear') {
                Ajax.post('/r/wse/put', post, function(response) {
                    if (response.status === "ok") {
                        alert('registrado correctamente');
                        location.href = "/r/index";
                        
                    }
                });
                // De lo contrario, los datos se envían al método "post" de la API.
            } else {
                Ajax.post('/r/wse/post', post, function(response) {
                    if (response.status === "ok") {
                        alert('registrado correctamente');
                        //location.href = "/r/index";
                        micrositio.cargarmicrositio(response.IdEmp);
                    }
                });
            }
        });
    };

    var eliminarEmpresa = function(idEmpresa) {
        Ajax.get('/r/wse/del?IdEmp=' + idEmpresa, function(response) {
            if (response.status === 'ok') {
                alert('Empresa eliminada.');
                location.reload();
            }
        });
    };

    // Registro de métodos públicos.
    return {
        initEmpresas: initEmpresas,
        empresaformdesdejson: empresaformdesdejson,
        empresa_envia: empresa_envia,
        eliminarEmpresa: eliminarEmpresa
    };
})();

/**
 * sucursales
 * Objeto que manipula la información de sucursales. 
 * @returns {object} sucursales - Regresa el registro de todos los métodos públicos para ser usados de la forma "sucursales.método".
 */
var sucursales = (function() {
    var map, geocoder, marker, infowindow;  // Variables de localización en el mapa.

    /**
     * locateAddress
     * Método que localiza en el mapa la dirección indicada por el formulario.
     */
    var _locateAddress = function() {
        // Getting the address from the text input
        var dir = [];
        dir.push($('#DirMunSuc option:selected').text());
        dir.push($('#DirEntSuc option:selected').text());
        dir.push($('#calle').val());
        dir.push($('#colonia').val());
        //dir.push($('#cp').val());
        var address = '';
        var coma = '';
        $.each(dir, function(key, value) {
            if (value) {
                address = address + coma + value;
                coma = ', ';
            }
        });
        address = address + ", MEXICO";

        // Check to see if we already have a geocoded object. If not we create one
        if (!geocoder) {
            geocoder = new google.maps.Geocoder();
        }
        // Creating a GeocoderRequest object
        var geocoderRequest = {
            address: address
        };
        // Making the Geocode request
        geocoder.geocode(geocoderRequest, function(results, status) {
            // Check if status is OK before proceeding
            if (status === google.maps.GeocoderStatus.OK) {
                // Center the map on the returned location
                map.setCenter(results[0].geometry.location);
                map.setZoom(17);
                // Check to see if we've already got a Marker object
                if (!marker) {
                    // Creating a new marker and adding it to the map
                    marker = new google.maps.Marker({
                        map: map,
                        draggable: true
                    });
                }
                // Setting the position of the marker to the returned location
                marker.setPosition(results[0].geometry.location);

                document.getElementById('lat').value = results[0].geometry.location.lat();
                document.getElementById('lng').value = results[0].geometry.location.lng();
            }
        });
    };

    /**
     * mostrarMapa
     * Método que pinta el mapa segun se va llenando el formulario.
     */
    var mostrarMapa = function() {
        var zoom = 17;
        var lat = $('#lat').val();
        var lng = $('#lng').val();
        if (!lat) {
            lat = 19.434341;
            lng = -99.141483;
            zoom = 10;
            //console.log(lat+" : "+lng);
        }
        var center = new google.maps.LatLng(lat, lng);
        var options = {
            zoom: zoom,
            center: center,
            mapTypeId: google.maps.MapTypeId.ROADMAP,
            streetViewControl: false
        };
        map = new google.maps.Map(document.getElementById('mapas'), options);
        if (!marker) {
            // Creating a new marker and adding it to the map
            marker = new google.maps.Marker({
                map: map,
                draggable: true
            });
            marker.setPosition(center);
        }
        // Getting a reference to the HTML form
        $('#DirEntSuc').bind('change', function(e) {
            _locateAddress();
        });
        $('#calle').bind('change', function(e) {
            _locateAddress();
        });
        $('#colonia').bind('change', function(e) {
            _locateAddress();
        });
        /*$('#cp').bind('change',function(e){
         locateAddress();
         });*/
        $('#buscar').bind('keydown keyup mousedown', function(e) {
            _locateAddress();
        });
        google.maps.event.addListener(marker, 'dragend', function() {
            var tmppos = '' + this.getPosition();
            var latlng = tmppos.split(',');
            document.getElementById('lat').value = latlng[0].replace('(', '');
            document.getElementById('lng').value = latlng[1].replace(')', '');
            map.setCenter(this.getPosition());
        });
    };

    /**
     * initSucursales
     * Método que muestra la lista de sucursales de la empresa indicada.
     * @param {string} codeRel - Identificador de la empresa.
     */
    var initSucursales = function(codeRel) {
        var sucursalesArray = {}, imprimeTemplate = '',
                underTemplateID = $('#sucursalesTemplate').html(),
                underTemplate = _.template(underTemplateID);
        Ajax.get('/r/wss/gets?IdEmp=' + codeRel, function(response) {
            imprimeTemplate = underTemplate({
                sucursalesArray: response
            });
            $("#sucursalesBlock").html(imprimeTemplate);
            Ajax.hidePreload($('#sucursales-lista'));
        });
    };

    /**
     * sucursalformdesdejsonModifica
     * Método que llena y muestra el formulario de modificar sucursal
     * @param {string} codeRel - Identificador de la sucursal a modificar.
     */
    var sucursalformdesdejsonModifica = function(codeRel) {
        $('#btn-sucursal').html("Modificar");
        Ajax.get('/r/wss/get?IdSuc=' + codeRel, function(response) {
            $('#sucursal-form').formParams(response, true);
            llenamuniSucursal(response.DirEnt, response.DirMun);
            Ajax.hidePreload($('#sucursal-detalle'));
            mostrarMapa();
        });
    };

    /**
     * sucursalformdesdejsonNueva
     * Método que llena y muestra el formulario de nueva sucursal.
     * @param {string} codeRel - Identificador de la empresa.
     */
    var sucursalformdesdejsonNueva = function(codeRel) {
        $('#IdEmpSuc').val(codeRel);
        $('#btn-sucursal').html("Crear");
        llenamuniSucursal("01");
        Ajax.hidePreload($('#sucursal-detalle'));
        mostrarMapa();
    };

    /**
     * sucursal_envia
     * Método que envia la información de los formularios de sucursales, ya sea para nueva sucursal o modificar una existente.
     */
    var sucursal_envia = function() {
        $(document).on('submit', 'form#sucursal-form', function(event) {
            event.preventDefault();
            var post = $(this).serialize();
            // Si el formulario es de nueva sucursal, envia los datos al método "put" del API.
            if ($('#btn-sucursal').html() === 'Crear') {
                Ajax.post('/r/wss/put', post, function(response) {
                    if (response.status === "ok") {
                        alert('registrado correctamente');
                        location.href = "/r/index";
                    }
                });
                // de lo contrario, envía los datos al método "post" del API.
            } else {
                Ajax.post('/r/wss/post', post, function(response) {
                    if (response.status === "ok") {
                        alert('registrado correctamente');
                        location.href = "/r/index";
                    }
                });
            }
        });
    };
    var eliminarSucursal = function(idSucursal, idEmpresa) {
        Ajax.get('/r/wss/del?IdEmp=' + idEmpresa + '&IdSuc=' + idSucursal, function(response) {
            if (response.status === 'ok') {
                alert('Sucursal elmiminada.');
                initSucursales(idEmpresa);
            }
        });
    };

    // Registro de métodos públicos.
    return {
        initSucursales: initSucursales,
        sucursalformdesdejsonModifica: sucursalformdesdejsonModifica,
        sucursalformdesdejsonNueva: sucursalformdesdejsonNueva,
        sucursal_envia: sucursal_envia,
        mostrarMapa: mostrarMapa,
        eliminarSucursal: eliminarSucursal
    };
})();

/**
 * registros
 * Objeto que maneja el registro de usuarios
 * @returns {object} registros - Regresa el registro de todos los métodos públicos para ser usados de la forma "registros.método".
 */
var registros = (function() {
    /**
     * registrarse
     * Método que envia la información para hacer el registro de un nuevo usuario.
     */
    var registrarse = function() {
        $('form#registroform').on('submit', function(event) {
            event.preventDefault();
            var post = $(this).serialize();
            Ajax.post('/r/wsr/put', post, function(response) {
                if (response.success) {
                    alert('registrado correctamente');
                }
            });
        });
    };

    // Registro de métodos públicos.
    return {
        registrarse: registrarse
    };
})();

/**
 * micrositio
 * Objeto que manipula la información de micrositios.
 */
var micrositio = (function() {

    /**
     * _putDefault
     * Metodo privado que muestra una imagen por defecto en caso que no haya imagen de logotipo.
     */
    var _putDefault = function() {
        $('#pic').remove();
        img = "<img  src = 'imgs/imageDefault.jpg' id='pic' width='258px' />";
        $('#urlimg').append(img);
    };

    /**
     * _avoidCache
     * Método privado que genera un identificador aleatorio para evitar el caché.
     */
    var _avoidCache = function() {
        var numRam = Math.floor(Math.random() * 500);
        return numRam;
    };

    /**
     * _updateimg
     * Método privado que actualiza la imagen que se esta mostrando como logotipo, por la que se acaba de subir.
     * @param {string} blob Identificador de la imagen actual
     */
    var _updateimg = function(blob) {
        blobkey = blob; // set blobkey global
        $("#BlobKey").attr("value", blobkey);
        if (blob) {
            $('#pic').remove();
            var query = "id=" + blob + "&Avc=" + _avoidCache();
            img = "<img  src = '/extraimg?" + query + "' id='pic' width='256px' />";
            $('#urlimg').append(img);
        } else {
            _putDefault();
        }
    };

    /**
     * cargarmicrositio
     * Método que llena y muestra el formulario de micrositios.
     * @param {string} empresaID - Identificador de la empresa.
     */
    var cargarmicrositio = function(empresaID) {
        Ajax.get('/r/wsed/get?IdEmp=' + empresaID, function(response) {
            if (response.status === 'ok') {
                $('div#micrositio-detalle').removeClass('inactivo');
                var blobkey = response.BlobKey;
                var uploadurl = response.UploadUrl;
                $("#IdEmp").attr("value", response.IdEmp);
                $("#enviar").attr('action', uploadurl);
                $("#BlobKey").attr("value", blobkey);
                $("#uploadimg_id").attr('value', empresaID);
                $("#empresa").html(response.Empresa);
                $("#descripcion").val(response.Desc);
                $("#facebook").val(response.Facebook);
                $("#twitter").val(response.Twitter);
                if (blobkey) {
                    _updateimg(blobkey);
                } else {
                    _putDefault();
                }
            } else {
                _putDefault();
            }
        });
    };

    /**
     * enviarimagen
     * Método que realiza el upload de la imagen.
     */
    var enviarimagen = function() {
        var bar = $('.bar');
        var percent = $('.percent');
        var status = $('#status');
        $('#enviar').ajaxSubmit({
            dataType: 'json',
            iframe: true,
            beforeSend: function() {
                status.empty();
                var percentVal = '0%';
                bar.width(percentVal);
                percent.html(percentVal);
            },
            uploadProgress: function(event, position, total, percentComplete) {
                var percentVal = percentComplete + '%';
                bar.width(percentVal);
                // percent.html(percentVal);
            },
            success: function(data) {
                //console.log(data);
                var resp = "";
                switch (data.status) {
                    case "invalidContentType":
                        resp = "<p>El archivo que intenta subir no es una imagen válida (sólo se aceptan con formato jpg o png).</p>";
                        break;
                    case "invalidUpload":
                        resp = "<p>Intente nuevamente, su imagen no puede ser integrada.</p>";
                        break;
                    case "uploadSessionError":
                        resp = "<p>Favor de refrescar la página para continuar.</p>";
                        break;
                    case "notFound":
                        resp = "<p>La oferta no existe.</p>";
                        break;
                    case "ok":
                        resp = "<p>La imagen se integró exitosamente</p>";
                        var uploadurl;
                        uploadurl = data.UploadUrl;
                        $("#enviar").attr("action", uploadurl);
                        setTimeout(function() {
                            _updateimg(data.BlobKey);
                        }, 1000);
                        break;
                    default:
                        resp = "<p>Intente nuevamente con una imagen de menor peso, su imagen no puede ser integrada.</p>";
                }
                status.html(resp);
            },
            complete: function() {
                bar.width('100%');
            },
            error: function() {
                status.html("<p>Intente nuevamente con una imagen de menor peso, su imagen no puede ser integrada.</p>");
            }
        });
    };

    /**
     * enviarDatos
     * Método que envia los datos del formulario de micrositios, via Ajax.
     */
    var enviarDatos = function() {
        $(document).on('submit', '#enviardata', function(event) {
            event.preventDefault();
            var data = $(this).serialize();
            var action = $(this).attr('action');
            Ajax.post(action, data, function(response) {
                if (response.status === "ok") {
                    alert('Micrositio registrado correctamente');
                    location.href = "/r/index";
                }
            });
        });
    };

    /**
     * extrasformulario
     * ?? Método desconocido
     * @todo Averiguar que hace.
     */
    var extrasformulario = function() {
        $("#pic").error(function() {
            putDefault();
        });

        $('textarea[maxlength]').live('keyup blur', function() {
            var maxlength = $(this).attr('maxlength');
            var val = $(this).val();
            if (val.length > maxlength) {
                $(this).val(val.slice(0, maxlength));
            }
        });

        $('input[maxlength]').live('keyup blur', function() {
            var maxlength = $(this).attr('maxlength');
            var val = $(this).val();
            if (val.length > maxlength) {
                $(this).val(val.slice(0, maxlength));
            }
        });
    };

    return {
        cargarmicrositio: cargarmicrositio,
        enviarimagen: enviarimagen,
        extrasformulario: extrasformulario,
        enviarDatos: enviarDatos
    };

})();


var ofertas = (function() {

    /**
     * _putDefault
     * Metodo privado que muestra una imagen por defecto en caso que no haya imagen de logotipo.
     */
    var _putDefault = function() {
        $('#OfertaPic').remove();
        img = "<img  src = 'imgs/imageDefault.jpg' id='OfertaPic' width='258px' />";
        $('#ofertaurlimg').append(img);
    };

    /**
     * _avoidCache
     * Método privado que genera un identificador aleatorio para evitar el caché.
     */
    var _avoidCache = function() {
        var numRam = Math.floor(Math.random() * 500);
        return numRam;
    };

    /**
     * _updateimg
     * Método privado que actualiza la imagen que se esta mostrando como logotipo, por la que se acaba de subir.
     * @param {string} blob Identificador de la imagen actual.
     */
    var _updateimg = function(blob) {
        blobkey = blob; // set blobkey global
        $("#OfertaBlobKey").attr("value", blobkey);
        if (blob) {
            $('#OfertaPic').remove();
            var query = "id=" + blob + "&Avc=" + _avoidCache();
            img = "<img  src = '/extraimg?" + query + "' id='OfertaPic' width='256px' />";
            $('#ofertaurlimg').append(img);
        } else {
            _putDefault();
        }
    };


    /**
     * initOfertas
     * Método que llena y muestra la lista de ofertas.
     * @param {string} rel - Identificador de la empresa.
     */
    var initOfertas = function(rel) {
        var imprimeTemplate = '',
                underTemplateID = $('#ofertasTemplate').html(),
                underTemplate = _.template(underTemplateID);
        Ajax.get('/r/wso/gets?IdEmp=' + rel, function(response) {
            imprimeTemplate = underTemplate({
                ofertasArray: response.ofertas
            });
            $("#ofertasLista").html(imprimeTemplate);
            Ajax.hidePreload($('#ofertas-lista'));
        });

    };
    /**
     * showOfertaForm
     * Método que muestra el formulario de la oferta, y si es para actualizar, lo llena.
     * @param {string} rel - Identificador de la oferta.
     * @param {boolean} esNueva Identifica si la oferta que se envia es nueva o es una actualización.
     */
    var showOfertaForm = function(rel, esNueva) {
        if (esNueva) {
            $('#imgform').hide();
            $('#boton-enviar-oferta').html('Siguiente Paso');
            $('.OfertaIdEmp').val(rel);
            Ajax.get('/r/wss/gets?IdEmp=' + rel, function(resp) {
                $('#ofertas-lista-sucursales tbody').empty();
                for (var a in resp) {
                    $('#ofertas-lista-sucursales tbody').append('<tr><td>' + resp[a].Nombre + '</td><td><input type="checkbox" name="sucursales[]" value="' + resp[a].IdSuc + '"></td></tr>');
                }
            });
        } else {
            $('#oferta-paso-1').attr('action', '/r/wso/post');
            $('#oferta-paso-1').append('<input type="hidden" name="IdOft" value="' + rel + '" />');
            $('#boton-enviar-oferta').html('Editar Oferta');
            $('#AddWordButton').attr('rel', rel);
            $('#oferta-paso-2').parent().removeClass('inactivo');
            $('#oferta-paso-3').parent().removeClass('inactivo');
            $('#oferta-paso-4').parent().removeClass('inactivo');
            Ajax.get('/r/wso/get?IdOft=' + rel, function(response) {
                if (response.status === 'ok') {
                    var d = new Date(Date.parse(response.FechaPub));
                    var blobkey = response.BlobKey;
                    Ajax.get('/r/wss/gets?IdEmp=' + response.IdEmp, function(resp) {
                        $('#ofertas-lista-sucursales tbody').empty();
                        for (var a in resp) {
                            $('#ofertas-lista-sucursales tbody').append('<tr><td>' + resp[a].Nombre + '</td><td><input type="checkbox" name="sucursales[]" value="' + resp[a].IdSuc + '"></td></tr>');
                        }
                    });
                    Ajax.get('/r/ofsuc?idemp=' + response.IdEmp + '&idoft=' + rel, function(resp) {
                        for (var a in resp) {
                            if (resp[a].idoft === rel) {
                                $('input[value="' + resp[a].idsuc + '"]').attr('checked', true);
                            }
                        }
                    });
                    $('.OfertaIdEmp').val(response.IdEmp);
                    $('#oferta-paso-1 #oferta').val(response.Oferta);
                    $('#oferta-paso-1 #descripcion').val(response.Descripcion);
                    $('#oferta-paso-1 #categoria').val(response.IdCat);
                    $('#oferta-paso-1 #date1').val(d.getUTCDate() + " Nov");
                    $('#oferta-paso-1 #ofertaDescuento').val(response.Descuento);
                    $('#oferta-paso-1 #ofertaMeses').val(response.Meses);
                    $('#oferta-paso-1 #ofertaPromocion').val(response.Promocion);
                    $('#oferta-paso-1 #ofertaPrecio').val(response.Precio);
                    $('#oferta-paso-1 #url').val(response.Url);
                    $('#oferta-paso-4').attr('action', response.UploadUrl);
                    $('#oferta-paso-4 #oferta_uploadimg_id').val(rel);
                    if (blobkey && blobkey !== 'none') {
                        _updateimg(blobkey);
                    } else {
                        _putDefault();
                    }
                }
            });
        }
        Ajax.hidePreload($('#oferta-detalle'));
    };

    var enviarDatosBasicos = function() {
        $(document).on('submit', '#oferta-paso-1', function(event) {
            event.preventDefault();
            $('input[name="sucursales[]"]').each(function() {
                if ($(this).is(':checked')) {
                    $('#schain').val($('#schain').val() + ' ' + $(this).val());
                }
            });
            var data = $(this).serialize();
            var action = $(this).attr('action');

            Ajax.post(action, data, function(response) {
                if (response.status === 'ok') {
                    alert('Datos enviados exitosamente.');
                    // $('#AddWordButton').attr('rel',response.IdOft);
                    // $('#oferta-paso-2').parent().removeClass('inactivo');
                    // $('#oferta-paso-3').parent().removeClass('inactivo');
                    Ajax.get('/r/wso/get?IdOft=' + response.IdOft, function(resp) {
                        $('#oferta-paso-4').attr('action', resp.UploadUrl);
                        $('#oferta-paso-4 #oferta_uploadimg_id').val(resp.IdOft);
                        $('#oferta-paso-4').parent().removeClass('inactivo');
                        if (resp.blobkey) {
                            _updateimg(blobkey);
                        } else {
                            _putDefault();
                        }
                    });

                    // initOfertas(idEmp);
                }
            });
        });
    };

    var agregarPalabra = function(idOferta) {
        var data = {
            token: $('#AddWordInput').val(),
            id: idOferta
        };
        Ajax.post('/r/addword', data, function(response) {
            var dataResp = {
                id: idOferta
            };
            if (response.status === 'ok') {
                Ajax.post('/r/wordsxo', dataResp, function(resp) {
                    $('#WordsAdded').empty();
                    $('#WordsAdded').append('<div class="palabraAgregada" data-idoft="' + idOferta + '"><span class="label-bf">' + response.token + '</span><a class="close-bf eliminarpalabra" href="">&times;</a></div>');
                    if (resp.status !== 'notFound') {
                        for (var a in resp) {
                            $('#WordsAdded').append('<div class="palabraAgregada" data-idoft="' + idOferta + '"><span class="label-bf">' + resp[a].token + '</span><a class="close-bf eliminarpalabra" href="">&times;</a></div>');
                        }
                    }
                });
            }

        });
    };

    var eliminarPalabra = function(idOferta, palabra) {
        data = {
            token: palabra,
            id: idOferta
        };
        Ajax.post('/r/delword', data, function(response) {
            if (response.status === 'ok') {
                var dataResp = {
                    id: idOferta
                };
                Ajax.post('/r/wordsxo', dataResp, function(resp) {
                    $('#WordsAdded').empty();
                    if (resp.status !== 'notFound') {
                        for (var a in resp) {
                            $('#WordsAdded').append('<div class="palabraAgregada" data-idoft="' + idOferta + '"><span class="label-bf">' + resp[a].token + '</span><a class="close-bf eliminarpalabra" href="">&times;</a></div>');
                        }
                    }
                });
            }
        });
    };

    // var agregarSucursal = function(idEmpresa, idSucursal, idOferta) {
    // 	Ajax.get('/r/addofsuc?idemp='+idEmpresa+'&idsuc='+idSucursal+'&idoft='+idOferta, function(response){
    // 		alert('Sucursal Agregada.');
    // 	});
    // }

    // var eliminarSucursal = function(idEmpresa, idSucursal, idOferta) {
    // 	Ajax.get('/r/delofsuc?idemp='+idEmpresa+'&idsuc='+idSucursal+'&idoft='+idOferta, function(response){
    // 		alert('Sucursal Eliminada.');
    // 	});
    // }

    var eliminarOferta = function(idOferta, idEmpresa) {
        Ajax.get('/r/wso/del?IdEmp=' + idEmpresa + '&IdOft=' + idOferta, function(response) {
            if (response.status === 'ok') {
                alert('Oferta elmiminada.');
                initOfertas(idEmpresa);
            }
        });
    };

    /**
     * enviarimagen
     * Método que realiza el upload de la imagen.
     */
    var enviarimagen = function() {
        var bar = $('.bar');
        var percent = $('.percent');
        var status = $('#ofertaimgstatus');
        $('#oferta-paso-4').ajaxSubmit({
            dataType: 'json',
            iframe: true,
            beforeSend: function() {
                status.empty();
                var percentVal = '0%';
                bar.width(percentVal);
                percent.html(percentVal);
            },
            uploadProgress: function(event, position, total, percentComplete) {
                var percentVal = percentComplete + '%';
                bar.width(percentVal);
                // percent.html(percentVal);
            },
            success: function(data) {
                // console.log(data);
                var resp = "";
                switch (data.status) {
                    case "invalidContentType":
                        resp = "<p>El archivo que intenta subir no es una imagen válida (sólo se aceptan con formato jpg o png).</p>";
                        break;
                    case "invalidUpload":
                        resp = "<p>Intente nuevamente, su imagen no puede ser integrada.</p>";
                        break;
                    case "uploadSessionError":
                        resp = "<p>Favor de refrescar la página para continuar.</p>";
                        break;
                    case "notFound":
                        resp = "<p>La oferta no existe.</p>";
                        break;
                    case "ok":
                        resp = "<p>La imagen se integró exitosamente</p>";
                        var uploadurl;
                        uploadurl = data.UploadUrl;
                        $("#oferta-paso-4").attr("action", uploadurl);
                        setTimeout(function() {
                            _updateimg(data.blobkey);
                        }, 1000);
                        break;
                    default:
                        resp = "<p>Intente nuevamente con una imagen de menor peso, su imagen no puede ser integrada.</p>";
                }
                status.html(resp);
            },
            complete: function() {
                bar.width('100%');
            },
            error: function() {
                status.html("<p>Intente nuevamente con una imagen de menor peso, su imagen no puede ser integrada.</p>");
            }
        });
    };


    // Registro de métodos públicos.
    return {
        initOfertas: initOfertas,
        showOfertaForm: showOfertaForm,
        enviarDatosBasicos: enviarDatosBasicos,
        agregarPalabra: agregarPalabra,
        eliminarOferta: eliminarOferta,
        eliminarPalabra: eliminarPalabra,
        enviarimagen: enviarimagen
                // agregarSucursal:agregarSucursal,
                // eliminarSucursal:eliminarSucursal
    };
})();