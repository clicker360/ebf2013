var empresas = (function(){
	var thepage = 1;
	var jalalasmalditasrecetas = function(){
		var tag=[];
		jQuery('.receta-option input').each( function() {
			if(jQuery(this).attr('checked')) {
				tag.push(jQuery(this).attr('name'));
			}
		} );
		jQuery("#animacion").show();
		//var ajaxurl = 'http://localhost:26936//ajax?tags=';
		var ajaxurl = 'http://dev.clicker360.com/daily_salad/ajax?tags=';
		jQuery.ajax({   
			//url: 'http://localhost:26936//ajax?tags='+tag+'&pg='+thepage,   
			url: 'http://dev.clicker360.com/daily_salad/ajax?tags='+tag+'&pg='+thepage,   
			type:'GET',   
			success: function(response){
				jQuery('#lasrecetas').append(response);
				jQuery('#animacion').hide();
				jQuery('#mas-recetas')!.show();
			}
		})
		thepage++;
	}
	
	return {
		jalalasmalditasrecetas:jalalasmalditasrecetas
	}
})();


window.onload = initPage;
      function initPage() {
        // Declare all the variables.
        var exampleValues = {},
          parsedTemplate = "",
          // Get the template from the script block.
          templateText = $('#profileTemplate').html(),
          // Then we use Underscore's template function to compile it 
          // into a stand alone function that we can feed values to and 
          // get HTML output.
          demoTemplate = _.template(templateText);
        // Here, we grab the data, and put the results into exampleValues. 
        // I'm using .ajax instead of .getJSON so I can set 'async: false' 
        // for the demo; in production you'd be checking for async calls 
        // to return, but I wanted to keep the script simple.
        $.ajax({
        	url: "http://localhost:8080/r/wse/gets",
        	async: false,
        	dataType: "json",
        	success: function(json) {
        		exampleValues = json;
        	}
        });
        // Finally, we call the demoTemplate function we created earlier, 
        // passing in the data we just retrieved, and then put the  
        // resulting HTML into the empty div.
        parsedTemplate = demoTemplate(exampleValues);
        $("#profileBlock").html(parsedTemplate);
      }