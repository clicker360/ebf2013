// Defining some global variables
var map, geocoder, marker, infowindow;
$(document).ready(function() {
	// Creating a new map
	var zoom = 17;
	var lat = $('#lat').val();
	var lng = $('#lng').val();
	if(!lat) { 
		lat = 19.434341;
		lng = -99.141483; 
		zoom = 10;
        console.log(lat+" : "+lng);
	}
	var center = new google.maps.LatLng(lat,lng);
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
	$('#DirEntSuc').bind('change',function(e){
		locateAddress();
	});
	$('#calle').bind('change',function(e){
		locateAddress();
	});
	$('#colonia').bind('change',function(e){
		locateAddress();
	});
	/*$('#cp').bind('change',function(e){
		locateAddress();
	});*/
	$('#buscar').bind('keydown keyup mousedown',function(e){
		locateAddress();
	});
	google.maps.event.addListener(marker, 'dragend', function() {
		var tmppos = ''+this.getPosition();
		var latlng = tmppos.split(',');
		document.getElementById('lat').value = latlng[0].replace('(','');
		document.getElementById('lng').value = latlng[1].replace(')','')
		map.setCenter(this.getPosition());
	});
});

function locateAddress() {
	// Getting the address from the text input
	var dir = [];
	dir.push($('#DirMunSuc option:selected').text());
	dir.push($('#DirEntSuc option:selected').text());
	dir.push($('#calle').val());
	dir.push($('#colonia').val());
	//dir.push($('#cp').val());
	var address = '';
	var coma = '';
	$.each(dir, function(key, value) { if(value) { address = address+coma+value; coma = ', '; } });
	address = address+", MEXICO";
	
	// Check to see if we already have a geocoded object. If not we create one
	if(!geocoder) {
		geocoder = new google.maps.Geocoder();
	}
	// Creating a GeocoderRequest object
	var geocoderRequest = {
		address: address
	}
	// Making the Geocode request
	geocoder.geocode(geocoderRequest, function(results, status) {
		// Check if status is OK before proceeding
		if (status == google.maps.GeocoderStatus.OK) {
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
}

