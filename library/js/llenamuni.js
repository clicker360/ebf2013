$(document).ready(function() {
    $('#DirEnt').on('change', function() {
        llenamuni($(this).val());
    });
});
function llenamuni(cvent, mun) {
    var imprimeTemplateMuni = '',
        underTemplateIMuni = $('#muniTemplate').html(),
        underTemplateMuni = _.template(underTemplateIMuni);
    $.get('/r/wsu/municipios', { CveEnt: cvent }, function(response) {
        imprimeTemplateMuni = underTemplateMuni({
            MunicipiosArray: response.municipios
        });
        $('#DirMunId').html(imprimeTemplateMuni);
        $($('#DirMunId option[value='+mun+']')).attr('selected','selected');
    });
}

