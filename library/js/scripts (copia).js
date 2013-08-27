(function(){

/*
  var verSucusales = function(){
    $('.ver-sucursales').on('click', function(event){
      //event.preventDefault();
      //event.stopPropagation();
      var rel = $(this).attr('rel');
      sucursales.initSucursales(rel);
    });
  }

   var verSucusales2 = function() {
    $('.ver-sucursales').click(function(){
      var rel = $(this).attr('rel');
      sucursales.initSucursales(rel);
    });
   }

*/

 var initSucursales = function(){
    //sucursales.initSucursales(); // lista de sucursales
    sucursales.sucursalformdesdejson(); // llenar formulario de sucursales
  }

  var initEmpresas = function(){
    //empresas.initEmpresas(); //lista de empresas
    empresas.empresaformdesdejson(); // llenar formulario de empresas
  }
  var execute = function() {
    $(document).ready(function(){
        //verSucusales2();
        initEmpresas();
        initSucursales();
    });
  }
  return execute();
})();


var empresas = (function(){
  //listas desde
  var initEmpresas = function(){
    var empresasArray = {},
        imprimeTemplate = '',
        underTemplateID = $('#empresaTemplate').html(),
        underTemplate = _.template(underTemplateID);        
        $.ajax({
          url: '/r/wse/gets',
          async: false,
          dataType: 'json',
          success: function(json) {
            empresasArray = json;
          }
        });
        imprimeTemplate = underTemplate({ empresasArray:empresasArray });
        $("#empresasBlock").html(imprimeTemplate);
  }
  //llenar formulario desde el json
  var empresaformdesdejson  = function(){
   //var forma = $('#sucursal-form');
    var datosform = {}
        $.ajax({
          url: '/r/wse/get?IdEmp=ecusboweznoeuqemhrna',
          async: false,
          dataType: 'json',
          success: function(json) {
            datosform = json;
          }
        });
    $('#empresa-form').formParams(datosform, true);
  }
  return {
    initEmpresas:initEmpresas,
    empresaformdesdejson:empresaformdesdejson
  }
})();

var sucursales = (function(){
  //llenar las listas desde el json
  var initSucursales = function(){
    var sucursalesArray = {},
        imprimeTemplate = '',
        underTemplateID = $('#sucursalesTemplate').html(),
        underTemplate = _.template(underTemplateID);        
        $.ajax({
          url: '/r/wss/gets?IdEmp=ghjkhxiwknbpvftwdrsm', 

          async: false,
          dataType: 'json',
          success: function(json) {
            sucursalesArray = json;
          }
        });
        imprimeTemplate = underTemplate({ sucursalesArray:sucursalesArray });
        $("#sucursalesBlock").html(imprimeTemplate);
  }
  // llenar el formulario desde el json
  var sucursalformdesdejson  = function(){ 
   //var forma = $('#sucursal-form');
    var datosform = {}
        $.ajax({
          url: '/r/wss/get?IdSuc=fbpchxospboklwmfiult',
          async: false,
          dataType: 'json',
          success: function(json) {
            datosform = json;
          }
        });
    $('#sucursal-form').formParams(datosform, true);
  }
  return {
    initSucursales:initSucursales,
    sucursalformdesdejson:sucursalformdesdejson
  }
})();