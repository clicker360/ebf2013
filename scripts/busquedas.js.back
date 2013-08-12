$(document).ready(function(){
	getcarrousel();
	pagina = 0;
	var queryVars = getVars();
	var hasVars = false;
	//if(queryVars.hasOwnProperty('word') || queryVars.hasOwnProperty('catMenu') || queryVars.hasOwnProperty('estadoMenu') || queryVars.hasOwnProperty('tipoMenu') )
	if(queryVars.hasOwnProperty('word')){
		hasVars = true;
		$("input[name=word]").val(queryVars.word);
	}
	if(queryVars.hasOwnProperty('catMenu')){
		hasVars = true;
		$("select[name=catMenu]").val(queryVars.catMenu);
	}
	if(queryVars.hasOwnProperty('estadoMenu')){
		hasVars = true;
		$("select[name=estadoMenu]").val(queryVars.estadoMenu);
	}
	if(queryVars.hasOwnProperty('tipoMenu')){
		hasVars = true;
		$("select[name=tipoMenu]").val(queryVars.tipoMenu);
	}

	searchOfertas();
	$("#buscarOferta").click(function(){                
		pagina = 0;    
		searchOfertas();
		return false;
	});
	cargaOfertas = false;
        inSearch = false;
	$(window).scroll(function(){
		if(cargaOfertas && !inSearch){
			if(($(document).scrollTop() + $(window).height()) >= ($(document).height() - 10)){
				searchOfertas();
			}
		}
	})

	$("#CloserLight").click(function() {
                $('#mapCont').hide();
                $('#imgOft img').show();
		$('#cuerpo').removeClass('noscroll');
		$('#lightback').addClass("hide"); 
		$('#lightfront').addClass("hide");
				  return false;
	});

	/*$(".lighter").click(function() {
		$('#cuerpo').addClass('noscroll');//importante ese impide que el fondo scrolle mientras la oferta si lo hace
		$('#lightback').removeClass("hide"); 
		$('#lightfront').removeClass("hide"); 
				  return false;
	});*/
	});
	function getVars() {
	var delimiter = "?"; // using '#' here is great for AJAX apps.
	var separator = "&";
	var url = location.href;
	var get_exists = (url.indexOf(delimiter) > -1) ? true : false;
	if (get_exists) {
		var url_get = {};
		var params = url.substring(url.indexOf(delimiter)+1);
		var params_array = params.split(separator);
		for (var i=0; i < params_array.length; i++) {
			var param_name = params_array[i].substring(0,params_array[i].indexOf('='));
			var param_value = params_array[i].substring(params_array[i].indexOf('=')+1);
			url_get[param_name] = param_value;
		}
	  return url_get;
	}
	return false;
	}
	function lighterAjax(){
            $(".lighter").click(function() {
                    var id = $(this).parent().attr('id');
                    $.get('wsdetalle',{id:id},function(data){
                        if(typeof(data) != 'object')
                            data = JSON.parse(data);
                            //console.log(data);
						if(data.hasOwnProperty('srvurl') && data.srvurl != '') {
							var imgurl = data.srvurl;
						} else {
                                                        var imgAnt = $("#"+data.idoft+" span.imgcont img").attr('src');
                                                        if(imgAnt.substring(0,8) != 'simg?id=')
                                                            var imgurl = imgAnt
                                                        else
                                                            var imgurl = (data.hasOwnProperty('imgurl') && data.imgurl != 'none') ? 'ofimg?id='+data.imgurl : false;
						}
                        var titOft = (data.hasOwnProperty('oferta')) ? data.oferta : '';
                        var desOft = (data.hasOwnProperty('descripcion')) ? data.descripcion :'';
                        var nomEmp = (data.hasOwnProperty('empresa')) ? data.empresa : '';
                        var idEmp = (data.hasOwnProperty('idemp')) ? data.idemp : '';
                        var enLinea = (data.hasOwnProperty('enlinea')) ? (data.enlinea) ? data.url : false : false;
                        if(enLinea){
                            $("#enLinea").html('<div class="col-12 bgRd marg-B10px marg-T10px padd-R10px marg-L5px" ><h4 class="typ-Wh"> El Buen Fin en Línea</h4></div><div class="col-13 marg-B10px marg-T10px padd-R10px marg-L10px"><a target="_blank" href="'+enLinea+'" class="first" >'+enLinea+'</a></div>')
                        }else{
                            $("#enLinea").html('')
                        }
                        var urlOferta = 'http://www.elbuenfin.org/detalleoferta.html?id='+data.idoft;
                        var mtOft = '<a onClick="window.open(\'mailto:?subject=Conoce esta oferta&body=Conoce esta oferta de El Buen Fin ' + urlOferta +'\', this.target, \'width=600,height=400\'); return false;" href="'+urlOferta+'">'
			mtOft += '<img src="/imgs/ofrtTemp/mtShare.jpg" alt="Enviar por correo electrónico" />'
			mtOft += '</a>'
                        var fbOft = '<a onClick="window.open(this.href, this.target, \'width=600,height=400\'); return false;" href="http://www.facebook.com/sharer.php?s=100&p[url]=' + urlOferta + '&p[images][0]=' + imgurl + '&p[title]= ' + titOft +' &p[summary] = '+desOft+'">';
			fbOft += '<img src="/imgs/ofrtTemp/fbShare.jpg" alt="Compartir en Facebook" />';
			fbOft += '</a>';
                        var twOft = '<a onClick="window.open(\'https://twitter.com/intent/tweet?text='+nomEmp+ ' ' + titOft +'  \' + this.href, this.target, \'width=600,height=400\'); return false" href="' + urlOferta +'" class="btwitter" title="Compartelo en Twitter">';
			twOft += '<img src="/imgs/ofrtTemp/twShare.jpg" alt="Compartir en Twitter" />';
			twOft += '</a>';
                        $(".sucList").html('');
                        if(idEmp){
                            var imgEmp = 'simg?id='+idEmp;
                            $(".logoOferta img").attr('src',imgEmp);
                            $(".logoOferta img").show();
							$(".logoOferta img").error(function() {
                                                            $(this).attr("src", "imgs/imageDefault.jpg");
								//$(this).hide();
							});
                            $.get('wssucs',{id:idEmp},function(sucursales){
                                if(typeof(sucursales) != 'object')
                                    sucursales = JSON.parse(sucursales);
                                if(sucursales.length >= 1){
                                    for(var i in sucursales){
                                        $(".sucList").append('<li><a href="#null" onClick="showMap(\''+sucursales[i].lat+'\',\''+sucursales[i].lng+'\',\'map\'); return false;">'+sucursales[i].sucursal+'</a></li>')
                                    }
                                }
                            });
                        }
                        $("#msEmp a").attr('href','/micrositio.html?id='+idEmp)
                        if(imgurl)
                            $("#imgOftLb").html('<img class="img" src="'+imgurl+'" width="430"/>');
                        else
                            $("#imgOftLb").html('');
                        $("#titOft h3").html(titOft);
                        $("#desOft p").html(desOft);
                        $("#nomEmp h4").html(nomEmp);
                        $("#mtOft").html(mtOft);
                        $("#fbOft").html(fbOft);
                        $("#twOft").html(twOft);
                        $('#cuerpo').addClass('noscroll');//importante ese impide que el fondo scrolle mientras la oferta si lo hace
                        $('#lightback').removeClass("hide");
                        $('#lightfront').removeClass("hide");

						$("#imgOftLb img").error(function() {
							$(this).attr("src", "imgs/imageDefault.jpg");
						});
                    });
              return false;
            });
	}
	function getcarrousel() {
	$.get("carr", "", function(response){
		$('#logo1').html(response);
	});
	}
	function searchOfertas(){

            pagina ++;
            if(pagina == '1')
		  $(".ofertCont").html('')
                 $(".ofertCont").append('<div class="col-98PR first Bg-ky padd-5px" id="cargador"><div class="cargando"><h4>CARGANDO</h4></div><div>');
            inSearch = true;
            var keywords = ($("input[name=word]").val() == '¿Qué buscas?') ? '' : $("input[name=word]").val();
            var categoria = $("select[name=catMenu]").attr("value");
            var estado = $("select[name=estadoMenu]").attr("value");
            var tipo = $("select[name=tipoMenu]").attr("value");
            $.get("http://movil.ebfmxorg.appspot.com/search?callback=?",{pagina:pagina, keywords:keywords, categoria:categoria, estado:estado , tipo:tipo, kind: 'Oferta'},function(data){
		var ofertas = data;
                $(".cargando").remove();
		if(ofertas.length >= 1){
		  cargaOfertas = true;
		  for(var i in ofertas){
			if(ofertas[i].Oferta == "Nueva oferta") continue;
			var logo;
			//(ofertas[i].Logo != "/ofimg?id=none") ? logo = ofertas[i].Logo : logo =  'simg?id='+ofertas[i].IdEmp;
			(ofertas[i].Logo != "") ? logo = ofertas[i].Logo : logo =  'simg?id='+ofertas[i].IdEmp;
			urlOferta = 'http://www.elbuenfin.org/detalleoferta.html?id='+ofertas[i].IdOft;
			addOferta = '<div class="oferta bgWh pagina_'+pagina+'" id="'+ofertas[i].IdOft+'">'
			addOferta += '<a href="#" class="lighter">'
			addOferta += '<span class="imgcont">'
			addOferta += '<img src="'+logo+'" width="212" alt="'+ofertas[i].Oferta+'" title="'+ofertas[i].Oferta+'" id="'+ofertas[i].IdOft+'"/>'
			addOferta += '</span>'
			addOferta += '<h3>'+ofertas[i].Oferta+'</h3>'
			addOferta += '</a>'
			addOferta += '<div class="col-30PR first" style="">'
			addOferta += '<a onClick="window.open(\'mailto:?subject=Conoce esta oferta&body=Conoce esta oferta de El Buen Fin ' + urlOferta +'\', this.target, \'width=600,height=400\'); return false;" href="'+urlOferta+'">'
			addOferta += '<img src="/imgs/ofrtTemp/mtShare.jpg" alt="Enviar por correo electrónico" />'
			addOferta += '</a>'
			addOferta += '</div>'
			addOferta += '<div class="col-40PR first" style="margin-top:5px;">'
			if(logo.substring(0,4) == 'http') {
				addOferta += '<a onClick="window.open(this.href, this.target, \'width=600,height=400\'); return false;" href="http://www.facebook.com/sharer.php?s=100&p[url]=' + urlOferta + '&p[images][0]=' + logo + '&p[title]= ' + ofertas[i].Oferta +'&p[summary] = '+ofertas[i].Descripcion+'">'
			} else {
				addOferta += '<a onClick="window.open(this.href, this.target, \'width=600,height=400\'); return false;" href="http://www.facebook.com/sharer.php?s=100&p[url]=' + urlOferta + '&p[images][0]=http://www.elbuenfin.org' + logo + '&p[title]= ' + ofertas[i].Oferta +'&p[summary] = '+ofertas[i].Descripcion+'">'
			}
			addOferta += '<img src="/imgs/ofrtTemp/fbShare.jpg" alt="Compartir en Facebook" />'
			addOferta += '</a>'
			addOferta += '</div>'
			addOferta += '<div class="col-30PR first">'
			addOferta += '<a onClick="window.open(\'https://twitter.com/intent/tweet?text='+ofertas[i].Oferta+ ' \' + this.href, this.target, \'width=600,height=400\'); return false" href="' + urlOferta +'" class="btwitter" title="Compartelo en Twitter">'
			addOferta += '<img src="/imgs/ofrtTemp/twShare.jpg" alt="Compartir en Twitter" />'
			addOferta += '</a>'
			addOferta += '</div>'
			addOferta += '</div>';
			$(".ofertCont").append(addOferta);
			$("img").error(function() {
				img = "<img  src = 'imgs/imageDefault.jpg' id='pic' width='212'/>";
				$(this).replaceWith(img);
			});
		  }
		lighterAjax();
		}
                if(ofertas.length < 12){
                    cargaOfertas = false;
                    $(".ofertCont").append('<h4 style="float:left; width:100%; text-align:center;">No hay más ofertas para esta búsqueda.</h4>');
                }
        $(".cargando").remove();
        $("#cargador").remove();
        inSearch = false;

	},'jsonp');
}
function showMap(lat, lng, div){
        $('#imgOft img').hide();
        $('#mapCont #closeMap').remove();
        $("#mapCont").prepend('<a href="#null" id="closeMap" class="button-2 marg-R20px"><span class="col-9 Dsblock alineCenter marg-L20px" onClick="$(\'#mapCont\').slideToggle(\'slow\'); $(\'#imgOft img\').show();">CERRAR MAPA [X]</span></a>');
        if($("#mapCont").is(":visible")){
            $('#mapCont').slideToggle('slow', function() {                
                $('#'+div).html('');
            });
        }else
            $('#map').html('');
        $('#mapCont').slideToggle('slow', function() {
            sucMap(lat, lng, div);
        });
        //setTimeout('sucMap("'+lat+'" , "'+lng+'" , "'+div+'")',1000);
}
function sucMap(lat, lng, div){
    var zoom = 17;
    var marker = false;
	if(lat == '0') {
		lat = 22.770856;
		lng = -102.583243;
		zoom = 4;
	}
	var center = new google.maps.LatLng(lat,lng);
	var options = {
		zoom: zoom,
		center: center,
		mapTypeId: google.maps.MapTypeId.ROADMAP,
		streetViewControl: false
	};
	map = new google.maps.Map(document.getElementById(div), options);
	if (!marker) {
		// Creating a new marker and adding it to the map
		marker = new google.maps.Marker({
			map: map,
			draggable: true
		});
		marker.setPosition(center);
	}
	google.maps.event.addListener(marker, 'dragend', function() {
		var tmppos = ''+this.getPosition();
		var latlng = tmppos.split(',');
		lat = latlng[0].replace('(','');
		lng = latlng[1].replace(')','')
		map.setCenter(this.getPosition());
	});
        // Creating a new map
}

