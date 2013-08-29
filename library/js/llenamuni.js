$(document).ready(function() {
    $('#DirEntEmp').on('change', function() {
        llenamuniEmpresa($(this).val());
    });
    $('#DirEntSuc').on('change', function() {
        llenamuniSucursal($(this).val());
    });
});
function llenamuniEmpresa(cvent, mun) {
    var imprimeTemplateMuni = '',
        underTemplateIMuni = $('#muniTemplateEmp').html(),
        underTemplateMuni = _.template(underTemplateIMuni);
    $.get('/r/wsu/municipios', { CveEnt: cvent }, function(response) {
        imprimeTemplateMuni = underTemplateMuni({
            MunicipiosArray: response.municipios
        });
        $('#DirMunEmp').html(imprimeTemplateMuni);
        $($('#DirMunEmp option[value='+mun+']')).attr('selected','selected');
    });
}
function llenamuniSucursal(cvent, mun) {
    var imprimeTemplateMuni = '',
        underTemplateIMuni = $('#muniTemplateSuc').html(),
        underTemplateMuni = _.template(underTemplateIMuni);
    $.get('/r/wsu/municipios', { CveEnt: cvent }, function(response) {
        imprimeTemplateMuni = underTemplateMuni({
            MunicipiosArray: response.municipios
        });
        $('#DirMunSuc').html(imprimeTemplateMuni);
        $($('#DirMunSuc option[value='+mun+']')).attr('selected','selected');
    });
}

