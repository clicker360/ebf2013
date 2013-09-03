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
    var micrositios = function(){
      micrositio.enviarimagen();
      micrositio.extrasformulario();

    };
   // llena formulario de detalle de empresa
   var llenaformempresas = function() {
       $(document).on('click','a.editar-empresa', function(event) {
           event.preventDefault();
           var rel = $(this).attr('rel');
           var empresaID = $(this).attr('rel');
           empresas.empresaformdesdejson(rel);
           micrositio.cargarmicrositio(empresaID);
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
           sucursales.initSucursales(rel);  // este
            $('#añadir-suc-empresa').html('<a class="nuevasucursal span12 btn btn-success" rel="'+rel+'" ><i class="icon-plus"></i> Añadir nueva sucursal</a>' );


       });
   };
   // llena formulario con datos de json 
   var llenaformsucursal = function() {
       $(document).on('click','a.editar-sucursal', function(event) {
           event.preventDefault();
           var rel = $(this).attr('rel');
           sucursales.sucursalformdesdejsonModifica(rel);
       });
   };
    // abre formulario de nueva sucursal
   var nuevasucursal = function() {
       $(document).on('click','a.nuevasucursal', function(event) {
           event.preventDefault();
           var rel = $(this).attr('rel');
           sucursales.sucursalformdesdejsonNueva(rel);
       });
   };

   // Submit de datos de empresa ya sea PUt o POST
   var modificasucursal = function() {
       sucursales.sucursal_envia();
   };

    var execute = function() {
        $(document).ready(function() {
            //verSucusales2();
            registrarse();
            initEmpresas();
            initSucursales();
            llenaformempresas();
            llenaformsucursal();
            nuevaempresa();
            nuevasucursal();
            modificaempresa();
            modificasucursal();
        });
    };
    return execute();
})();
var Ajax = (function() {
    var _showPreload = function() {
        $('.preloader').fadeIn('fast');
    };
    var hidePreload = function(bloque) {
        $('.preloader').fadeOut('fast');
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
                llenamuniEmpresa(response.DirEnt, response.DirMun);
                llenaorganismos(response.OrgEmp);
                Ajax.hidePreload($('#empresas-detalle'));
            });
        } else {
            $('#btn-empresa').html("Crear");
            llenamuniEmpresa("01");
            llenaorganismos();
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
                           location.href = "/r/index";
                       }
               });
            } else {
                Ajax.post('/r/wse/post', post, function(response){
                       if(response.status=="ok"){
                          alert('registrado correctamente');
                           location.href = "/r/index";
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

    // codeRel trae sucursal
     var sucursalformdesdejsonModifica = function(codeRel) {
            $('#btn-sucursal').html("Modificar");
            Ajax.get('/r/wss/get?IdSuc=' + codeRel, function(response) {
                $('#sucursal-form').formParams(response, true);
                llenamuniSucursal(response.DirEnt, response.DirMun);
                Ajax.hidePreload($('#sucursal-detalle'));
            });
    };

    // codeRel trae ID Empresa
     var sucursalformdesdejsonNueva = function(codeRel) {
            $('#IdEmpSuc').val(codeRel);
            $('#btn-sucursal').html("Crear");
            llenamuniSucursal("01");
            Ajax.hidePreload($('#sucursal-detalle'));
    };

    var sucursal_envia = function () {
        $(document).on('submit','form#sucursal-form', function(event){
            event.preventDefault();
            var post = $(this).serialize();
            if($('#btn-sucursal').html() == 'Crear') {
                Ajax.post('/r/wss/put', post, function(response){
                       if(response.status=="ok"){
                           alert('registrado correctamente');
                           location.href = "/r/index";
                       }
               });
            } else {
                Ajax.post('/r/wss/post', post, function(response){
                       if(response.status=="ok"){
                           alert('registrado correctamente');
                           location.href = "/r/index";
                       }
               });
            }
        })
    };
    
    return {
        initSucursales: initSucursales,
        sucursalformdesdejsonModifica: sucursalformdesdejsonModifica,
        sucursalformdesdejsonNueva: sucursalformdesdejsonNueva,
        sucursal_envia: sucursal_envia,
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


var micrositio = (function(){
    var blobkey;
    var uploadurl;

    var cargarmicrositio = function (empresaID){
      $.get('/r/wsed/get?IdEmp='+empresaID, function(resp) {
        console.log(resp);
        if(typeof(resp) != 'object') { resp = JSON.parse(resp); }
        if(resp.status=="ok") {
            blobkey = resp.BlobKey;
            uploadurl = resp.UploadUrl;
            $("#IdEmp").attr("value", resp.IdEmp);
            $("#enviar").attr('action', uploadurl);
            $("#BlobKey").attr("value", blobkey);
            $("#uploadimg_id").attr('value', empresaID);
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
    };

    var $pic = $("#pic");
   /* 
   * Ajax FORM para imagen de oferta
   */
    var bar = $('.bar');
    var percent = $('.percent');
    var status = $('#status');
    var img;

    var enviarimagen = function () {
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
    };

    var extrasformulario = function () {
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
    };


    return {
        cargarmicrositio: cargarmicrositio,
        enviarimagen: enviarimagen,
        extrasformulario: extrasformulario
    };

})();

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