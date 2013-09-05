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
			micrositio.cargarmicrositio(empresaID);
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
			$('#añadir-suc-empresa').html('<a class="nuevasucursal span12 btn btn-success" rel="'+ rel+ '" ><i class="icon-plus"></i> Añadir nueva sucursal</a>');

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
	
	/* ---------------- Termina Sucursales ------------- */
	
	/* ---------------- Inicia Micrositios ------------- */
	
	var uploadImage = function() {
		$('#infile').on('change', function(event){
			$('#imgfile').val($('#infile').val());
			micrositio.enviarimagen();
		});
	}
	
	var registerMicrositio = function() {
		micrositio.enviarDatos();
	}
	
	/* ---------------- Termina Micrositios ------------ */
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
			url : url,
			dataType : 'json',
			success : function(response) {
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
			type : 'POST',
			url : url,
			dataType : 'json',
			data : data,
			success : function(response) {
				if (typeof callback === 'function') {
					callback.call(callback, response);
				}
			}
		});
	};
	
	// Registro de métodos públicos.
	return {
		hidePreload : hidePreload,
		get : get,
		post : post
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
				empresasArrayDash : response
			});
			imprimeTemplate = underTemplate({
				empresasArray : response
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
			$('#btn-empresa').html("Modificar");
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
			if ($('#btn-empresa').html() == 'Crear') {
				Ajax.post('/r/wse/put', post, function(response) {
					if (response.status == "ok") {
						alert('registrado correctamente');
						location.href = "/r/index";
					}
				});
			// De lo contrario, los datos se envían al método "post" de la API.
			} else {
				Ajax.post('/r/wse/post', post, function(response) {
					if (response.status == "ok") {
						alert('registrado correctamente');
						location.href = "/r/index";
					}
				});
			}
		})
	};
	
	// Registro de métodos públicos.
	return {
		initEmpresas : initEmpresas,
		empresaformdesdejson : empresaformdesdejson,
		empresa_envia : empresa_envia
	};
})();

/**
 * sucursales
 * Objeto que manipula la información de sucursales. 
 * @returns {object} sucursales - Regresa el registro de todos los métodos públicos para ser usados de la forma "sucursales.método".
 */
var sucursales = (function() {
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
				sucursalesArray : response
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
			if ($('#btn-sucursal').html() == 'Crear') {
				Ajax.post('/r/wss/put', post, function(response) {
					if (response.status == "ok") {
						alert('registrado correctamente');
						location.href = "/r/index";
					}
				});
			// de lo contrario, envía los datos al método "post" del API.
			} else {
				Ajax.post('/r/wss/post', post, function(response) {
					if (response.status == "ok") {
						alert('registrado correctamente');
						location.href = "/r/index";
					}
				});
			}
		})
	};
	
	// Registro de métodos públicos.
	return {
		initSucursales : initSucursales,
		sucursalformdesdejsonModifica : sucursalformdesdejsonModifica,
		sucursalformdesdejsonNueva : sucursalformdesdejsonNueva,
		sucursal_envia : sucursal_envia,
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
		})
	};
	
	// Registro de métodos públicos.
	return {
		registrarse : registrarse
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
		img = "<img  src = 'imgs/imageDefault.jpg' id='pic' width='258px' />"
		$('#urlimg').append(img);
	}
	
	/**
	 * _avoidCache
	 * Método privado que genera un identificador aleatorio para evitar el caché.
	 */
	var _avoidCache = function() {
		var numRam = Math.floor(Math.random() * 500);
		return numRam;
	}
	
	/**
	 * _updateimg
	 * Método privado que actualiza la imagen que se esta mostrando como logotipo, por la que se acaba de subir.
	 */
	var _updateimg = function(blob) {
		blobkey = blob; // set blobkey global
		$("#BlobKey").attr("value", blobkey);
		if (blob) {
			$('#pic').remove();
			var query = "id=" + blob + "&Avc=" + _avoidCache();
			img = "<img  src = '/extraimg?" + query + "' id='pic' width='256px' />"
			$('#urlimg').append(img);
		} else {
			_putDefault();
		}
	}
	
	/**
	 * cargarmicrositio
	 * Método que llena y muestra el formulario de micrositios.
	 * @param {string} empresaID - Identificador de la empresa.
	 */
	var cargarmicrositio = function(empresaID) {
		Ajax.get('/r/wsed/get?IdEmp=' + empresaID, function(response){
			if(response.status == 'ok'){
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
			dataType : 'json',
			beforeSend : function() {
				status.empty();
				var percentVal = '0%';
				bar.width(percentVal);
				percent.html(percentVal);
			},
			uploadProgress : function(event, position, total, percentComplete) {
				var percentVal = percentComplete + '%';
				bar.width(percentVal);
				// percent.html(percentVal);
			},
			success : function(data) {
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
			complete : function() {
				bar.width('100%');
			},
			error : function() {
				status.html("<p>Intente nuevamente con una imagen de menor peso, su imagen no puede ser integrada.</p>");
			}
		});
	};
	
	var enviarDatos = function() {
		$(document).on('submit', '#enviardata', function(event){
			event.preventDefault();
			var data = $(this).serialize();
			var action = $(this).attr('action');
			Ajax.post(action, data, function(response){
				if (response.status == "ok") {
					alert('Micrositio registrado correctamente');
					location.href = "/r/index";
				}
			});
		});
	}

	var extrasformulario = function() {
		$("#pic").error(function() {
			putDefault()
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
		cargarmicrositio : cargarmicrositio,
		enviarimagen : enviarimagen,
		extrasformulario : extrasformulario,
		enviarDatos:enviarDatos
	};

})();