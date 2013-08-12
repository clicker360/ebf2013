$(document).ready(function() {
			var $pic = $("#pic");
			var $url = $("#url");
			var $size = $("#size");
			var img;
			function putDefault() {
				$('#pic').remove();
				img = "<img  src = 'imgs/imageDefault.jpg' id='pic' width='256px' />";
				$('#theImage').append(img);
			}
			function update() {
				$('#pic').remove();
				var query = "id="+id + "&Avc=" + avoidCache();
				img = "<img  src = '/simg?"+ query +"' id='pic' width='256px' />";
				$('#theImage').append(img);
			}
			update();	
		var bar = $('.bar');
		var percent = $('.percent');
		var status = $('#status');
		function avoidCache(){
				var numRam = Math.floor(Math.random() * 500);
				return numRam;
			}
	
		$('#enviar').ajaxForm({
			beforeSend: function() {
				status.empty();
				var percentVal = '0%';
				bar.width(percentVal)
				percent.html(percentVal);
			},
			uploadProgress: function(event, position, total, percentComplete) {
				$('#pic').remove();
				var percentVal = percentComplete + '%';
				bar.width(percentVal)
				percent.html(percentVal);
			},
			complete: function(xhr) {
				status.html(xhr.responseText);
				setTimeout(function(){
				 update(); }, 1000); 
			}
		}); 
		$("#pic").error(function() { putDefault()});
})
