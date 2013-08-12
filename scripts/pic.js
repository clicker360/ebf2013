$(document).ready(function() {
		var id = "{{.|js}}";
		var $pic = $("#pic");
		var $url = $("#url");
		var $share = $("#share");
		var x = 0;
		var y = 0;
		function update() {
			var query = "id="+id+"&x="+x+"&y="+y+
				"&s="+$("#size").val()+
				"&d="+$("#droop").val();
			$pic.attr("src", "simg?"+query);
			$url.text("simg?"+query);
			$share.attr("href", "/share?"+query);
		}
		$pic.click(function(e) {
			x = e.pageX - this.offsetLeft;
			y = e.pageY - this.offsetTop;
			update();
		});
		$("#size, #droop").bind("mouseup", update);
		update();
	})

