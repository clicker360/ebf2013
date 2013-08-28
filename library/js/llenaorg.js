$(document).ready(function() {
    $('#OrgEmp').on('change', function() {
        llenaorganismos($(this).val());
    });
});
function llenaorganismos(siglas) {
    var imprimeTemplateOrg = '',
        underTemplateIOrg = $('#orgTemplate').html(),
        underTemplateOrg = _.template(underTemplateIOrg);
    $.get('/r/wsu/organismos', function(response) {
        imprimeTemplateOrg = underTemplateOrg({
            OrganismosArray: response.organismos
        });
        $('#OrgEmpId').html(imprimeTemplateOrg);
        $($('#OrgEmpId option[value='+siglas+']')).attr('selected','selected');
    });
}

