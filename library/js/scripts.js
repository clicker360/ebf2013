(function() {
  /* ---------------- Generales y utilerias -------------*/
    var municipios = function() {
       $(document).on('click','a.ver-sucursales', function(event) {
           event.preventDefault();
           var valmun = $(this).attr('value');
           sucursales.initSucursales(valmun);
       });
   };
   var registrarse = function () {
        registros.registrarse();
    };
    /* ---------------- Empresas -------------*/
   //carga automatica de empresas
   var initEmpresas = function() {
        empresas.initEmpresas(); //lista de empresas
    };
   // llena formulario de detalle de empresa
   var llenaformempresas = function() {
       $(document).on('click','a.editar-empresa', function(event) {
           event.preventDefault();
           var rel = $(this).attr('rel');
           empresas.empresaformdesdejson(rel);
       });
   };
   // abre formulario de nueva empresa
   var nuevaempresa = function() {
       $(document).on('click','a.nuevaempresa', function(event) {
           event.preventDefault();
           empresas.empresaformdesdejson();
       });
   };
  
   // Submit de datos de empresa ya sea PUt o POST
   var modificaempresa = function() {
       empresas.empresa_envia();
   };

   /* ---------------- Sucursales -------------*/
  var initSucursales = function() {
       $(document).on('click','a.ver-sucursales', function(event) {
           event.preventDefault();
           var rel = $(this).attr('rel');
           sucursales.initSucursales(rel);
       });
   };
   // llena formulario con datos de json 
   var llenaformsucursal = function() {
       $(document).on('click','a.editar-sucursal', function(event) {
           event.preventDefault();
           var rel = $(this).attr('rel');
           sucursales.sucursalformdesdejson(rel);
       });
   };
    // abre formulario de nueva sucursal
   var nuevasucursal = function() {
       $(document).on('click','a.nuevasucursal', function(event) {
           event.preventDefault();
           sucursales.sucursalformdesdejson();
       });
   };

    
    
    var execute = function() {
        $(document).ready(function() {
            //verSucusales2();
            registrarse();
            llenaformempresas();
            llenaformsucursal();
            nuevasucursal();
            initEmpresas();
            initSucursales();
            nuevaempresa();
            modificaempresa();
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
            $('.activo').removeClass('activo').addClass('inactivo');
            bloque.addClass('activo').removeClass('inactivo');   
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

var general = ( function() {
  var municipios = function (valmun) {
        var imprimeMunicipios = '',
                underMunicipiosID = $('#empresaTemplate').html(),
                underMunicipios = _.template(underMunicipiosID);
        Ajax.get('/r/wse/gets' + valmun, function(response) {
            imprimeMunicipios = underMunicipios({
                MunicipiosArray: response
            });
            $("#empresasBlock").html(imprimeMunicipios);
        });
  };
})();
var empresas = (function() {
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
    var empresaformdesdejson = function(codeRel) {
        if(codeRel) {
            $('#btn-empresa').html("Modificar");
            Ajax.get('/r/wse/get?IdEmp=' + codeRel, function(response) {
                $('#empresa-form').formParams(response, true);
                llenamuni(response.DirEnt, response.DirMun);
                llenaorganismos(response.OrgEmp);
                Ajax.hidePreload($('#empresas-detalle'));
            });
        } else {
            $('#btn-empresa').html("Crear");
            Ajax.hidePreload($('#empresas-detalle'));
        }
    };
    var empresa_envia = function () {
        $(document).on('submit','form#empresa-form', function(event){
            event.preventDefault();
            var post = $(this).serialize();
            if($('#btn-empresa').html() == 'Crear') {
                Ajax.post('/r/wse/put', post, function(response){
                       if(response.status=="ok"){
                               alert('registrado correctamente');
                       }
               });
            } else {
                Ajax.post('/r/wse/post', post, function(response){
                       if(response.status=="ok"){
                               alert('registrado correctamente');
                       }
               });
            }
        })
    };
    return {
        initEmpresas: initEmpresas,
        empresaformdesdejson: empresaformdesdejson,
        empresa_envia: empresa_envia
    };
})();


var sucursales = (function() {
    var initSucursales = function(codeRel) {
        //var imprimeTemplate = '',
        var sucursalesArray = {},
             imprimeTemplate = '',
                underTemplateID = $('#sucursalesTemplate').html(),
                underTemplate = _.template(underTemplateID);
        Ajax.get('/r/wss/gets?IdEmp=' + codeRel, function(response) {
            imprimeTemplate = underTemplate({sucursalesArray: response});
            $("#sucursalesBlock").html(imprimeTemplate);
            Ajax.hidePreload($('#sucursales-lista'));
        });
    };


  /*  var sucursalformdesdejson = function(codeRel) {
        Ajax.get('/r/wss/get?IdSuc=' + codeRel, function(response) {
            $('#sucursal-form').formParams(response, true);
              Ajax.hidePreload($('#sucursal-detalle'));
        });
    };*/
     var sucursalformdesdejson = function(codeRel) {
        if(codeRel) {
            $('#btn-sucursal').html("Modificar");
            Ajax.get('/r/wse/get?IdEmp=' + codeRel, function(response) {
                $('#empresa-form').formParams(response, true);
                llenamuni(response.DirEnt, response.DirMun);
                llenaorganismos(response.OrgEmp);
                Ajax.hidePreload($('#sucursal-detalle'));
            });
        } else {
            $('#btn-sucursal').html("Crear");
            Ajax.hidePreload($('#sucursal-detalle'));
        }
    };
    var sucursalnueva = function (codeEmpresa) {
        $('form#sucursal-form').on('submit', function(event){
            event.preventDefault();
            var post = $(this).serialize();
            Ajax.post('/r/wss/put?IdEmp=', post, function(response){
                   if(response.success){
                           alert('registrado correctamente');
                   }
           });
        })
    };
    return {
        initSucursales: initSucursales,
        sucursalformdesdejson: sucursalformdesdejson,
        sucursalnueva: sucursalnueva       
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
