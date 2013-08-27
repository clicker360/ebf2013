(function() {
    var initSucursales = function() {
        $('a.ver-sucursales').on('click', function(event) {
            event.preventDefault();
            var rel = $(this).attr('rel');
            sucursales.initSucursales(rel);
        });
    };
    var initEmpresas = function() {
        empresas.initEmpresas(); //lista de empresas
    };
    var execute = function() {
        $(document).ready(function() {
            //verSucusales2();
            initEmpresas();
            initSucursales();
        });
    };
    return execute();
})();
var Ajax = (function() {
    var _showPreload = function() {
        $('#preloader').fadeIn('fast');
    };
    var hidePreload = function() {
        $('#preloader').fadeOut('fast');
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

var empresas = (function() {
    var initEmpresas = function() {
        var imprimeTemplate = '',
                underTemplateID = $('#empresaTemplate').html(),
                underTemplate = _.template(underTemplateID);
        Ajax.get('/r/wse/gets', function(response) {
            imprimeTemplate = underTemplate({
                empresasArray: response
            });
            $("#empresasBlock").html(imprimeTemplate);
            Ajax.hidePreload();
        });
    };
    var empresaformdesdejson = function(codeRel) {
        Ajax.get('/r/wse/get?IdEmp=' + codeRel, function(response) {
            $('#empresa-form').formParams(response, true);
        });
    };
    return {
        initEmpresas: initEmpresas,
        empresaformdesdejson: empresaformdesdejson
    };
})();

var sucursales = (function() {
    var initSucursales = function(codeRel) {
        var imprimeTemplate = '',
                underTemplateID = $('#sucursalesTemplate').html(),
                underTemplate = _.template(underTemplateID);
        Ajax.get('/r/wss/gets?IdEmp=' + codeRel, function(response) {
            imprimeTemplate = underTemplate({sucursalesArray: response});
            $("#sucursalesBlock").html(imprimeTemplate);
            Ajax.hidePreload();
        });
    };

    var sucursalformdesdejson = function(codeRel) {
        Ajax.get('/r/wss/get?IdSuc=' + codeRel, function(response) {
            $('#sucursal-form').formParams(response, true);
        });
    };

    return {
        initSucursales: initSucursales,
        sucursalformdesdejson: sucursalformdesdejson
    };
})();